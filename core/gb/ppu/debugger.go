package ppu

import "unsafe"

func (p *PPU) Frame() uint64 {
	return p.FrameCounter
}

func (p *PPU) GetTileImage(ppuID, tileID int, buffer unsafe.Pointer, bpp, paletteID uint8) {
}

func (p *PPU) GetTilemap(ppuID int, buffer unsafe.Pointer, w, h, n int) {
	p.r.GetTilemap(buffer, n)
}

func (p *PPU) GetPalette(paletteID uint8) unsafe.Pointer {
	return nil
}
