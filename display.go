package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	Window       *sdl.Window
	Surface      *sdl.Surface
	Instructions chan func()
}

func (d *Display) InitDisplay() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	var err error

	defer sdl.Quit()

	d.Window, err = sdl.CreateWindow("CHIP 8 Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	defer d.Window.Destroy()

	d.Surface, err = d.Window.GetSurface()
	if err != nil {
		panic(err)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

	}
}

// Input of x and y should be the original pixel coordinate, and will be scaled by a factor of 10
func (d *Display) TurnOnPixel(x int32, y int32) {
	rect := sdl.Rect{x * 10, y * 10, 10, 10}
	color := sdl.Color{255, 255, 255, 255}
	pixel := sdl.MapRGBA(d.Surface.Format, color.R, color.G, color.B, color.A)
	d.Surface.FillRect(&rect, pixel)
	d.Window.UpdateSurface()
}

// Input of x and y should be the original pixel coordinate, and will be scaled by a factor of 10
func (d *Display) TurnOffPixel(x int32, y int32) {
	rect := sdl.Rect{x * 10, y * 10, 10, 10}
	color := sdl.Color{0, 0, 0, 0}
	pixel := sdl.MapRGBA(d.Surface.Format, color.R, color.G, color.B, color.A)
	d.Surface.FillRect(&rect, pixel)
	d.Window.UpdateSurface()
}

func (d *Display) IsPixelOn(x int32, y int32) bool {
	color := d.Surface.At(int(x), int(y))
	r, g, b, a := color.RGBA()
	if (r == 0xFF) && (g == 0xFF) && (b == 0xFF) && (a == 0xFF) {
		return true
	} else {
		return false
	}
}

func (d *Display) TogglePixel(x int32, y int32) {
	if d.IsPixelOn(x, y) {
		d.TurnOffPixel(x, y)
	} else {
		d.TurnOnPixel(x, y)
	}
}

func (d *Display) ClearScreen() {
	d.Surface.FillRect(nil, 0)
	d.Window.UpdateSurface()
}

func (d *Display) Run() {
	running := true
	for running {
		sdl.PollEvent()
	}
}
