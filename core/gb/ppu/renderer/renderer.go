package renderer

import (
	"image/color"
	"unsafe"
)

type Renderer interface {
	DrawScanline(y int, scanline []color.NRGBA)
	SetLCDC(val uint8)
	SetBGP(val uint8)
	SetOBP(bank, val uint8)
	SetSCX(val uint8)
	SetSCY(val uint8)
	SetWX(val uint8)
	SetWY(val uint8)

	// Debug
	GetTilemap(id int, buffer unsafe.Pointer, n int)
}
