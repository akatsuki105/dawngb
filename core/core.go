package core

import (
	"image/color"
	"io"

	"github.com/akatsuki105/dawngb/core/gb"
)

type Core interface {
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

func NewGB(audioBuffer io.Writer) *gb.GB {
	return gb.New(audioBuffer)
}
