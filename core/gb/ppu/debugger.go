package ppu

import "unsafe"

func (p *PPU) Frame() uint64 {
	return p.FrameCounter
}

func (p *PPU) GetTilemap(id int, buffer unsafe.Pointer, w, h, n int) {
	p.r.GetTilemap(id, buffer, n)
}

func (p *PPU) GetPalette(paletteID uint8) unsafe.Pointer {
	return unsafe.Pointer(&p.Palette[0])
}
