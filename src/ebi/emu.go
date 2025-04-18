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

var stateSaveData bytes.Buffer
var saved bool

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

func createEmu(model uint8) *Emu {
	if emu != nil {
		return emu
	}
	e := &Emu{
		sampleBuffer: newSampleBuffer(make([]uint8, 0, 8192)),
		turbo:        1,
		volume:       1,
		taskQueue:    make([]func(), 0, 10),
	}
	e.c = gb.New(gb.Model(model), e.sampleBuffer)

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
	return "DawnGB"
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
	e.c.Reset()
	e.c.DirectBoot()

	// Load Save Data
	ext := filepath.Ext(path)
	if ext == ".gbc" || ext == ".gb" {
		var savData []uint8

		savExt := []string{".sav", ".srm"}
		for _, sav := range savExt {
			savPath := strings.ReplaceAll(path, ext, sav)
			if _, err := os.Stat(savPath); err == nil {
				data, err := os.ReadFile(savPath)
				if err == nil {
					savData = data
					break
				}
			}
		}

		if len(savData) > 0 {
			err := e.c.Load(gb.LOAD_SAVE, savData)
			if err != nil {
				return err
			}
			e.c.Reset()
			e.c.DirectBoot()
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

func (e *Emu) setTurbo(speed int) {
	e.queueTask(func() {
		e.turbo = speed
	})
}

func (e *Emu) enableSound(enabled bool) {
	e.queueTask(func() {
		prev := e.soundEnabled
		e.soundEnabled = enabled
		if prev != enabled {
			e.sampleBuffer.Reset()
		}
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
				e.c.Reset()
				e.c.DirectBoot()

			case ".sav", ".srm": // Save Data
				err := e.c.Load(gb.LOAD_SAVE, data)
				if err != nil {
					return err
				}
				e.c.Reset()
				e.c.DirectBoot()

			case ".bin": // BIOS
				size := len(data)
				if size == 256 || size == 2048 || size == 2048+256 {
					err := e.c.Load(gb.LOAD_BIOS, data)
					if err != nil {
						return err
					}
					e.c.Reset()
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
