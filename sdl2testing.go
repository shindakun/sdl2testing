package main

import (
	"math/rand"
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
		renderer.Present()
		clear(pixels)
		elpasedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elpasedTime < 7 {
			sdl.Delay(7 - uint32(elpasedTime))
			elpasedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
