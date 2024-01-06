package main

import (
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"os"

	"github.com/akatsuki105/dugb/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/oto"
)

var music = true
var samples = make([]byte, 4096)

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

type Emu struct {
	c            core.Core
	active       bool
	sampleBuffer *bytes.Buffer
	context      *oto.Context
	music        *oto.Player
}

func createEmu() *Emu {
	e := &Emu{
		sampleBuffer: bytes.NewBuffer(make([]byte, 0)),
	}
	e.c = core.New("GB", e.sampleBuffer)

	if music {
		context, err := oto.NewContext(44100, 2, 2, 4096)
		if err != nil {
			panic("oto.NewContext failed: " + err.Error())
		}
		e.context = context
		e.music = context.NewPlayer()
	}
	return e
}

func (e *Emu) Title() string {
	if !e.active {
		return "DuGB"
	}
	return fmt.Sprintf("DuGB - %s", e.c.Title())
}

func (e *Emu) LoadROMFromPath(path string) error {
	if path == "" {
		return fmt.Errorf("rom path is not specified")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return e.LoadROM(data)
}

func (e *Emu) LoadROM(data []byte) error {
	err := e.c.LoadROM(data)
	e.active = err == nil
	return err
}

func (e *Emu) Update() error {
	if e.active {
		e.pollInput()
		e.c.RunFrame()
		if e.music != nil {
			for i := 0; i < len(samples); i++ {
				samples[i] = 0
			}
			n, _ := e.sampleBuffer.Read(samples)
			e.music.Write(samples[:n])
		}
	} else {
		file := ebiten.DroppedFiles()
		if file != nil {
			entries, err := fs.ReadDir(file, ".")
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				if len(name) < 4 {
					continue
				}
				data, err := fs.ReadFile(file, name)
				if err != nil {
					return err
				}
				return e.LoadROM(data)
			}
		}
	}
	return nil
}

func (e *Emu) Draw(screen *ebiten.Image) {
	if e.active {
		data := e.c.Screen()
		w, h := e.c.Resolution()
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.Set(x, y, data[y*w+x])
			}
		}
		screen.DrawImage(ebiten.NewImageFromImage(img), nil)
	}
}

func (e *Emu) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.c.Resolution()
}

func (e *Emu) pollInput() {
	for key, input := range keyMap {
		e.c.SetKeyInput(input, ebiten.IsKeyPressed(key))
	}
}
