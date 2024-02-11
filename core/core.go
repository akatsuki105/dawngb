package core

import (
	"image/color"
	"io"

	"github.com/akatsuki105/dawngb/core/gb"
)

type ID = string

const (
	GB ID = "GB"
)

type Core interface {
	ID() ID // Get Core ID
	Reset(hasBIOS bool)
	LoadROM(romData []byte) error // romData is mutable(not copied).
	SRAM() []byte
	LoadSRAM(sramData []byte) error
	RunFrame()
	Resolution() (w int, h int) // Get display resolution
	Screen() []color.RGBA
	SetKeyInput(key string, press bool)
	Title() string // Get title of the game

	// Serialize
	Serialize(state io.Writer)
	Deserialize(state io.Reader)
}

func New(id ID, audioBuffer io.Writer) Core {
	switch id {
	case GB:
		return gb.New(audioBuffer)
	default:
		panic("invalid core id. valid core id is {GB}")
	}
}
