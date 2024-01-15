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

var standardButtonToString = map[ebiten.StandardGamepadButton]string{
	ebiten.StandardGamepadButtonRightTop:    "B",
	ebiten.StandardGamepadButtonRightRight:  "A",
	ebiten.StandardGamepadButtonCenterLeft:  "SELECT",
	ebiten.StandardGamepadButtonCenterRight: "START",
	ebiten.StandardGamepadButtonLeftBottom:  "DOWN",
	ebiten.StandardGamepadButtonLeftRight:   "RIGHT",
	ebiten.StandardGamepadButtonLeftLeft:    "LEFT",
	ebiten.StandardGamepadButtonLeftTop:     "UP",
}

var inputMap = map[string]bool{
	"A":      false,
	"B":      false,
	"START":  false,
	"SELECT": false,
	"UP":     false,
	"DOWN":   false,
	"LEFT":   false,
	"RIGHT":  false,
}

type Emu struct {
	c            core.Core
	active       bool
	sampleBuffer *bytes.Buffer
	context      *oto.Context
	music        *oto.Player
	isDebugMode  bool
	turbo        int
}

func createEmu(isDebugMode bool) *Emu {
	e := &Emu{
		sampleBuffer: bytes.NewBuffer(make([]byte, 0)),
		isDebugMode:  isDebugMode,
		turbo:        1,
	}
	e.c = core.New("GB", e.sampleBuffer)

	if e.turbo > 1 {
		music = false
	}

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

	err = e.LoadROM(data)
	if err != nil {
		return err
	}

	// xxxx.gb or xxxx.gbc -> xxxx.sav
	if len(path) > 4 && path[len(path)-4:] == ".gbc" {
		path = path[:len(path)-4]
	} else if len(path) > 3 && path[len(path)-3:] == ".gb" {
		path = path[:len(path)-3]
	}
	savPath := path + ".sav"
	if _, err := os.Stat(savPath); err == nil {
		savData, err := os.ReadFile(savPath)
		if err == nil {
			e.c.LoadSRAM(savData)
		}
	}

	return nil
}

func (e *Emu) LoadROM(data []byte) error {
	err := e.c.LoadROM(data)
	e.active = err == nil
	return err
}

func (e *Emu) Update() error {
	if e.active {
		e.pollInput()
		for i := 0; i < e.turbo; i++ {
			e.c.RunFrame()
		}
		if music && e.music != nil {
			for i := 0; i < len(samples); i++ {
				samples[i] = 0
			}
			n, err := e.sampleBuffer.Read(samples)
			if err != nil && n > 0 {
				e.music.Write(samples[:n])
			}
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
		if e.isDebugMode {
			screen.DrawImage(ebiten.NewImageFromImage(e.c.DebugVRAM()), nil)
		} else {
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
}

func (e *Emu) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if e.isDebugMode {
		return 256, 256
	}
	return e.c.Resolution()
}

func (e *Emu) pollInput() {
	for key := range inputMap {
		inputMap[key] = false
	}

	e.pollKeyInput()
	e.pollGamepadInput()

	for key, input := range inputMap {
		e.c.SetKeyInput(key, input)
	}
}

func (e *Emu) pollKeyInput() {
	for key, input := range keyMap {
		if _, ok := inputMap[input]; ok {
			if ebiten.IsKeyPressed(key) {
				inputMap[input] = true
			}
		}
	}
}

func (e *Emu) pollGamepadInput() {
	ids := ebiten.AppendGamepadIDs(nil)
	for _, id := range ids {

		for b, input := range standardButtonToString {
			switch {
			case ebiten.IsStandardGamepadButtonPressed(id, b):
				if _, ok := inputMap[input]; ok {
					inputMap[input] = true
				}
			}
		}

		switch ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) {
		case 1:
			inputMap["RIGHT"] = true
		case -1:
			inputMap["LEFT"] = true
		}
		switch ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) {
		case 1:
			inputMap["DOWN"] = true
		case -1:
			inputMap["UP"] = true
		}
	}
}
