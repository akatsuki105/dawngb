package ppu

import "unsafe"

func (p *PPU) Frame() uint64 {
	return p.FrameCounter
}

func (p *PPU) GetTilemap(id int, buffer unsafe.Pointer, w, h, n int) {
	p.r.GetTilemap(id, buffer, n)
}
