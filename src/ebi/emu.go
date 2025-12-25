package main

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/akatsuki105/dawngb/core/gb"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/constraints"
)

type Emu struct {
	Core    *gb.GB
	Paused  bool
	Reset   bool
	HasBIOS bool
	active  bool

	Snapshot struct {
		Enabled bool
		Data    bytes.Buffer
	}
}

func createEmu[V constraints.Integer](model V) *Emu {
	return &Emu{
		Core:  gb.New(gb.Model(model), App.Audio),
		Reset: true,
	}
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
			err := e.LoadSave(savData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Emu) LoadROM(data []uint8) error {
	err := e.Core.Load(gb.LOAD_ROM, data)
	if err != nil {
		return err
	}
	e.active = true
	e.Reset = true
	return nil
}

func (e *Emu) Update() error {
	if !e.Paused && e.active {
		if e.Reset {
			e.Reset = false
			e.Core.Reset()
			if !e.HasBIOS || !App.Config.GB.Intro {
				e.Core.DirectBoot()
			}
		}

		for key, input := range Inputs {
			e.Core.SetKeyInput(key, input)
		}
		e.Core.RunFrame()
	}
	return nil
}

func (e *Emu) Draw(screen *ebiten.Image) {
	if !e.Paused && e.active {
		data := e.Core.Screen()
		w, h := 160, 144
		img := image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetNRGBA(x, y, data[y*w+x])
			}
		}
		screen.DrawImage(ebiten.NewImageFromImage(img), nil)
	}
}

func (e *Emu) LoadSave(data []uint8) error {
	err := e.Core.Load(gb.LOAD_SAVE, data)
	if err != nil {
		return err
	}
	e.Reset = true
	return nil
}

func (e *Emu) LoadBIOS(data []uint8) error {
	size := len(data)
	if size == 256 || size == 2048 || size == 2048+256 {
		err := e.Core.Load(gb.LOAD_BIOS, data)
		if err != nil {
			return err
		}
		e.HasBIOS, e.Reset = true, true
	}
	return nil
}

func (e *Emu) SaveState() bool {
	fmt.Println("State save")
	ok := e.Core.Serialize(&e.Snapshot.Data)
	if ok {
		fmt.Println("State save failed")
		e.Snapshot.Enabled = true
	}
	return ok
}

func (e *Emu) LoadState() bool {
	fmt.Println("State load")
	if e.Snapshot.Enabled {
		ok := e.Core.Deserialize(&e.Snapshot.Data)
		if !ok {
			fmt.Println("State load failed")
			return false
		}
	}
	return true
}
