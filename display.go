package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	Window        *sdl.Window
	Surface       *sdl.Surface
	Instructions  chan Instruction
	PixelStatuses [][]bool
}

type Pixel struct {
	X    int32
	Y    int32
	IsOn bool
}

type Instruction struct {
	Fn   func(args ...any)
	args []any
}

var key *Key
var onColor sdl.Color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
var offColor sdl.Color = sdl.Color{R: 0, G: 0, B: 0, A: 0}

func (d *Display) InitDisplay(k *Key) {
	key = k
	d.Instructions = make(chan Instruction, 1000)
	go d.PerformInstructions()

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
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if e.Type == uint32(sdl.KEYDOWN) {
					key.KeyboardEvent = e.Keysym.Sym
					key.SetMappedKey()
					key.IsPressed = true

				} else {
					key.IsPressed = false
				}

			case *sdl.QuitEvent:
				fmt.Println("QUIT")
				running = false

			}
		}
	}
}

// Input of x and y should be the original pixel coordinate, and will be scaled by a factor of 10
func (d *Display) TurnOnPixel(x int32, y int32) {
	rect := sdl.Rect{X: x * 10, Y: y * 10, W: 10, H: 10}
	pixel := sdl.MapRGBA(d.Surface.Format, onColor.R, onColor.G, onColor.B, onColor.A)
	d.Surface.FillRect(&rect, pixel)
	// d.Window.UpdateSurface()
}

// Input of x and y should be the original pixel coordinate, and will be scaled by a factor of 10
func (d *Display) TurnOffPixel(x int32, y int32) {
	rect := sdl.Rect{X: x * 10, Y: y * 10, W: 10, H: 10}
	pixel := sdl.MapRGBA(d.Surface.Format, offColor.R, offColor.G, offColor.B, offColor.A)
	d.Surface.FillRect(&rect, pixel)
	// d.Window.UpdateSurface()
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

func (d *Display) SendInstruction(i Instruction) {
	d.Instructions <- i
}

func (d *Display) PerformInstructions() {
	for i := range d.Instructions {
		i.Fn(i.args...)
	}
}
