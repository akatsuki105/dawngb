package main

import (
	"bytes"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/akatsuki105/dawngb/core/gb"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
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
	c      *gb.GB
	active bool
	paused bool

	// Audio
	soundEnabled bool
	volume       float64
	sampleBuffer *sampleBuffer
	music        *oto.Player

	turbo     int
	taskQueue []func() // Run at the start of the frame, so safe to access the core
}

func createEmu() *Emu {
	if emu != nil {
		return emu
	}
	e := &Emu{
		sampleBuffer: newSampleBuffer(make([]uint8, 0, 8192)),
		turbo:        1,
		volume:       0.5,
		taskQueue:    make([]func(), 0, 10),
	}
	e.c = gb.New(e.sampleBuffer)

	// init Audio
	op := oto.NewContextOptions{
		SampleRate:   32768,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE, // RetroArch はこれを使っているので合わせると楽
	}
	context, readyChan, err := oto.NewContext(&op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	<-readyChan
	e.music = context.NewPlayer(e.sampleBuffer)
	e.music.SetVolume(e.volume)
	e.music.SetBufferSize(8192)

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
	e.c.Reset(false)

	// Load Save Data
	ext := filepath.Ext(path)
	if ext == ".gbc" || ext == ".gb" {
		savPath := strings.ReplaceAll(path, ext, ".sav")
		if _, err := os.Stat(savPath); err == nil {
			if savData, err := os.ReadFile(savPath); err == nil {
				err := e.c.Load(gb.LOAD_SAVE, savData)
				if err != nil {
					return err
				}
				e.c.Reset(false)
			}
		}
	}

	return nil
}

func (e *Emu) LoadROM(data []uint8) error {
	err := e.c.Load(gb.LOAD_ROM, data)
	if err != nil {
		e.active = false
		return err
	}

	e.active = true
	ebiten.SetWindowTitle(e.title())
	return nil
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

		if e.soundEnabled {
			e.music.Play()
		}
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
		img := image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetNRGBA(x, y, data[y*w+x])
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
			case ".gb", ".gbc": // ROM
				err := e.LoadROM(data)
				if err != nil {
					return err
				}
				e.c.Reset(false)

			case ".sav": // Save Data
				err := e.c.Load(gb.LOAD_SAVE, data)
				if err != nil {
					return err
				}
				e.c.Reset(false)

			case ".bin": // BIOS
				size := len(data)
				if size == 256 || size == 2048 || size == 2048+256 {
					err := e.c.Load(gb.LOAD_BIOS, data)
					if err != nil {
						return err
					}
					e.c.Reset(true)
				}
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

// Read で n == 0 のときに EOF を返すと音が途切れるので、 nil を返すようにしただけ
type sampleBuffer struct {
	*bytes.Buffer
}

func newSampleBuffer(buf []uint8) *sampleBuffer {
	return &sampleBuffer{bytes.NewBuffer(buf)}
}

func (s *sampleBuffer) Read(p []uint8) (int, error) {
	n, _ := s.Buffer.Read(p)
	if n == 0 {
		return 0, nil // EOF を返すと音が途切れるので、 nil を返す
	}
	return n, nil
}
