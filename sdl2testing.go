package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"time"

	"github.com/shindakun/vec3"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	wWidth  = 1280
	wHeight = 720
)

type rgb struct {
	r byte
	g byte
	b byte
}

type star struct {
	pos   vec3.Vector3
	dir   vec3.Vector3
	color rgb
}

type starField struct {
	stars []star
}

type letter struct {
	letter rune
	x      int
	y      int
	h      int
	w      int
}

type letters struct {
	letters []letter
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)
	return tex
}

func pngToTexture(renderer *sdl.Renderer, filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, wWidth*wHeight*4)
	pIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[pIndex] = byte(r / 256)
			pIndex++
			pixels[pIndex] = byte(g / 256)
			pIndex++
			pixels[pIndex] = byte(b / 256)
			pIndex++
			pixels[pIndex] = byte(a / 256)
			pIndex++
		}
	}

	tex := pixelsToTexture(renderer, pixels, w, h)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func setPix(pixels []byte, x, y int, color rgb) {
	count := (y*wWidth + x) * 4

	if count < len(pixels)-4 && count >= 0 {
		pixels[count] = color.r
		pixels[count+1] = color.g
		pixels[count+2] = color.b
	}
}

func (s *starField) update() {
	for i := 0; i < len(s.stars); i++ {
		newX := vec3.Add(s.stars[i].pos, s.stars[i].dir)
		s.stars[i].pos.X = newX.X
	}
}

func (s *starField) draw(pixels []byte) {
	for i := 0; i < len(s.stars); i++ {
		setPix(pixels, int(s.stars[i].pos.X), int(s.stars[i].pos.Y), s.stars[i].color)
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func letterIn(letter rune, letters []letter) letter {
	for i := 0; i < len(letters); i++ {
		if letter == letters[i].letter {
			return letters[i]
		}
	}
	return letters[0]
}

func printFont(renderer *sdl.Renderer, font *sdl.Texture, letters []letter) {
	hello := "HELLO WORLD"

	for i := 0; i < len(hello); i++ {
		a := letterIn(rune(hello[i]), letters)
		renderer.Copy(font, &sdl.Rect{int32(a.x), int32(a.y), int32(a.w), int32(a.h)}, &sdl.Rect{int32(i * 32), 0, 32, 32})
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Test", 300, 300, wWidth, wHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(wWidth), int32(wHeight))
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()

	pixels := make([]byte, wWidth*wHeight*4)

	var sf starField

	s := make([]star, 400)
	for i := 0; i < len(s); i++ {
		s[i].pos.X = rand.Float32() * wWidth
		s[i].pos.Y = rand.Float32() * wHeight
		s[i].dir.X = 1
		s[i].color.r = 255
		s[i].color.b = 255
		s[i].color.g = 255

		sf.stars = append(sf.stars, s[i])
	}

	var tf starField
	t := make([]star, 400)
	for i := 0; i < len(s); i++ {
		t[i].pos.X = rand.Float32() * wWidth
		t[i].pos.Y = rand.Float32() * wHeight
		t[i].dir.X = .7
		t[i].color.r = 170
		t[i].color.b = 170
		t[i].color.g = 170

		tf.stars = append(tf.stars, t[i])
	}

	var zf starField
	z := make([]star, 400)
	for i := 0; i < len(s); i++ {
		z[i].pos.X = rand.Float32() * wWidth
		z[i].pos.Y = rand.Float32() * wHeight
		z[i].dir.X = .5
		z[i].color.r = 70
		z[i].color.b = 70
		z[i].color.g = 70

		zf.stars = append(zf.stars, z[i])
	}

	font := pngToTexture(renderer, "font2.png")

	startingChar := 32
	index := 0
	letters := make([]letter, 60)
	for y := 0; y < 6; y++ {
		for x := 0; x < 10; x++ {
			letters[index].letter = rune(startingChar + index)
			letters[index].h = 32
			letters[index].w = 32
			letters[index].x = x * 32
			letters[index].y = y * 32
			index++
		}
	}
	// fmt.Println(letters)
	fmt.Println(letterIn('!', letters))

	var elpasedTime float32
	for {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		//moveStars(t, pixels)
		//moveStars(s, pixels)
		sf.update()
		tf.update()
		zf.update()
		zf.draw(pixels)
		tf.draw(pixels)
		sf.draw(pixels)
		tex.Update(nil, pixels, wWidth*4)
		renderer.Copy(tex, nil, nil)
		printFont(renderer, font, letters)
		renderer.Present()
		clear(pixels)
		elpasedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elpasedTime < 7 {
			sdl.Delay(7 - uint32(elpasedTime))
			elpasedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
