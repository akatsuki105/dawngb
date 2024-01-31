package core

import (
	"image"
	"image/color"
	"io"

	"github.com/akatsuki105/dugb/core/gb"
)

type ID = string

const (
	GB ID = "GB"
)

type Core interface {
	// Get Core ID
	ID() ID

	// LoadROM loads game rom
	//
	// It assumes an environment with enough memory, so it is necessary to pass the complete ROM data in advance.
	//
	// NOTE: romData is mutable(not copied).
	LoadROM(romData []byte) error

	LoadSRAM(sramData []byte) error

	SRAM() []byte

	// RunFrame runs emulator until a next frame
	RunFrame()

	// Get display resolution
	Resolution() (w int, h int)

	Screen() []color.RGBA

	SetKeyInput(key string, press bool)

	Title() string

	// debug
	DebugVRAM() image.Image
}

func New(id ID, audioBuffer io.Writer) Core {
	switch id {
	case GB:
		return gb.New(audioBuffer)
	default:
		panic("invalid core id. valid core id is {GB}")
	}
}
