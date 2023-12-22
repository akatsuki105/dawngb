package renderer

import (
	"image/color"

	"github.com/akatsuki105/dugb/core/gb/video/renderer/software"
)

type Renderer interface {
	DrawScanline(y int, scanline []color.RGBA)
}

func New(id string) Renderer {
	switch id {
	case "software":
		return software.New()
	default:
		panic("invalid renderer id. valid renderer id is {software}")
	}
}
