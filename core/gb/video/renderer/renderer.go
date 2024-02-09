package renderer

import (
	"image/color"

	"github.com/akatsuki105/dawngb/core/gb/video/renderer/dummy"
	"github.com/akatsuki105/dawngb/core/gb/video/renderer/software"
)

type Renderer interface {
	DrawScanline(y int, scanline []color.RGBA)
	SetLCDC(val uint8)
	SetBGP(val uint8)
	SetOBP0(val uint8)
	SetOBP1(val uint8)
	SetSCX(val uint8)
	SetSCY(val uint8)
	SetWX(val uint8)
	SetWY(val uint8)

	SetBGPI(val uint8)
	GetBGPD() uint8
	SetBGPD(val uint8) uint8
	SetOBPI(val uint8)
	GetOBPD() uint8
	SetOBPD(val uint8) uint8
}

func New(id string, vram, oam []uint8, model int) Renderer {
	switch id {
	case "software":
		return software.New(vram, oam, model)
	case "dummy":
		return dummy.New()
	default:
		panic("invalid renderer id. valid renderer id is {software, dummy}")
	}
}
