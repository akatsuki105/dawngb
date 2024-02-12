package main

import (
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/akatsuki105/dawngb/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/oto"
)

var emu *Emu

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

var inputMapWeb = map[string]bool{
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
	c      core.Core
	active bool
	paused bool

	// Audio
	soundEnabled bool
	samples      []byte
	sampleBuffer *bytes.Buffer
	context      *oto.Context
	music        *oto.Player

	turbo     int
	taskQueue []func() // Run at the start of the frame, so safe to access the core
}

func createEmu() *Emu {
	if emu != nil {
		return emu
	}
	e := &Emu{
		samples:      make([]byte, 4096),
		sampleBuffer: bytes.NewBuffer(make([]byte, 0)),
		turbo:        1,
		taskQueue:    make([]func(), 0, 10),
	}
	e.c = core.New("GB", e.sampleBuffer)

	// init Audio
	context, err := oto.NewContext(44100, 2, 2, len(e.samples))
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	e.context = context
	e.music = context.NewPlayer()

	emu = e
	return e
}

func (e *Emu) title() string {
	if !e.active {
		return "DawnGB"
	}
	return fmt.Sprintf("DawnGB - %s", e.c.Title())
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

	// Load Save Data
	ext := filepath.Ext(path)
	if ext == ".gbc" || ext == ".gb" {
		savPath := strings.ReplaceAll(path, ext, ".sav")
		if _, err := os.Stat(savPath); err == nil {
			if savData, err := os.ReadFile(savPath); err == nil {
				err := e.c.LoadSRAM(savData)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (e *Emu) LoadROM(data []byte) error {
	err := e.c.LoadROM(data)
	e.active = err == nil
	if e.active {
		ebiten.SetWindowTitle(e.title())
	}
	return err
}

func (e *Emu) Update() error {
	if len(e.taskQueue) > 0 {
		for _, task := range e.taskQueue {
			task()
		}
		e.taskQueue = e.taskQueue[:0]
	}

	if e.active && !e.paused {
		e.pollInput()
		for i := 0; i < e.turbo; i++ {
			e.c.RunFrame()
		}

		e.playSound()
	}

	err := e.handleDropFile()
	if err != nil {
		return err
	}

	return nil
}

func (e *Emu) Draw(screen *ebiten.Image) {
	if e.active && !e.paused {
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
	for key := range inputMap {
		inputMap[key] = inputMapWeb[key]
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

func (e *Emu) playSound() {
	for i := 0; i < len(e.samples); i++ {
		e.samples[i] = 0
	}
	n, err := e.sampleBuffer.Read(e.samples)
	if e.soundEnabled && err == nil && n > 0 {
		e.music.Write(e.samples[:n])
	}
}

func (e *Emu) setTurbo(speed int) {
	e.queueTask(func() {
		e.turbo = speed
	})
}

func (e *Emu) enableSound(enabled bool) {
	e.queueTask(func() {
		e.soundEnabled = enabled
	})
}

func (e *Emu) queueTask(f func()) {
	e.taskQueue = append(e.taskQueue, f)
}

func (e *Emu) handleDropFile() error {
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
			ext := filepath.Ext(name)
			data, err := fs.ReadFile(file, name)
			if err != nil {
				return err
			}
			switch ext {
			case ".gb", ".gbc":
				return e.LoadROM(data)
			case ".sav":
				return e.c.LoadSRAM(data)
			}
		}
	}
	return nil
}

func (e *Emu) setPaused(paused bool) {
	e.queueTask(func() {
		if e.active {
			e.paused = paused
		}
	})
}
