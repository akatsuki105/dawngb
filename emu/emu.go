package emu

import (
	"fmt"
	"image"
	"os"

	"github.com/akatsuki105/dugb/core"
	"github.com/akatsuki105/dugb/core/gb"
	"github.com/hajimehoshi/ebiten/v2"
)

var keyMap = map[ebiten.Key]string{
	ebiten.KeyX:          "A",
	ebiten.KeyZ:          "B",
	ebiten.KeyBackspace:  "SELECT",
	ebiten.KeyEnter:      "START",
	ebiten.KeyArrowUp:    "UP",
	ebiten.KeyArrowDown:  "DOWN",
	ebiten.KeyArrowLeft:  "LEFT",
	ebiten.KeyArrowRight: "RIGHT",
}

var btnMap = map[ebiten.GamepadButton]string{
	ebiten.GamepadButton2:  "A",
	ebiten.GamepadButton1:  "B",
	ebiten.GamepadButton3:  "B",
	ebiten.GamepadButton4:  "L",
	ebiten.GamepadButton5:  "R",
	ebiten.GamepadButton8:  "SELECT",
	ebiten.GamepadButton9:  "START",
	ebiten.GamepadButton15: "UP",
	ebiten.GamepadButton16: "RIGHT",
	ebiten.GamepadButton17: "DOWN",
	ebiten.GamepadButton18: "LEFT",
}

type Emu struct {
	c core.Core
}

func New() *Emu {
	return &Emu{
		c: gb.New(),
	}
}

func (e *Emu) Title() string {
	return "TODO"
}

func (e *Emu) LoadROM(path string) error {
	if path == "" {
		return fmt.Errorf("rom path is not specified")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return e.c.LoadROM(data)
}

func (e *Emu) Update() error {
	e.pollInput()
	e.c.RunFrame()
	return nil
}

func (e *Emu) Draw(screen *ebiten.Image) {
	data := e.c.FrameBuffer()
	w, h := e.c.Resolution()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, data[y*w+x])
		}
	}
	screen.DrawImage(ebiten.NewImageFromImage(img), nil)
}

func (e *Emu) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.c.Resolution()
}

func (e *Emu) pollInput() {
	for key, input := range keyMap {
		e.c.SetKeyInput(input, ebiten.IsKeyPressed(key))
	}
}
