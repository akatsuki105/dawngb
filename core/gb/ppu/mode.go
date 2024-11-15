package ppu

import (
	"github.com/akatsuki105/dawngb/util"
)

// Mode 0
func (p *PPU) hblank() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC)
	if util.Bit(p.lcdc, 7) && !p.enableLatch {
		p.r.DrawScanline(p.ly, p.screen[p.ly*160:(p.ly+1)*160])
	}
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.irq(1)
	}
	if p.onHBlank != nil {
		p.onHBlank()
	}
}

// Mode 1
func (p *PPU) vblank() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC) | 1
	p.irq(0)

	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.irq(1)
	}
}

// Mode 2
func (p *PPU) scanOAM() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC) | 2
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.irq(1)
	}
}

// Mode 3
func (p *PPU) drawing() {
	p.stat = (p.stat & 0xFC) | 3

	// Count scanline objects
	h := 8
	if util.Bit(p.lcdc, 2) {
		h = 16
	}
	o := 0
	for i := 0; i < 40; i++ {
		y := int(p.oam[i*4]) - 16
		if y <= p.ly && p.ly < y+h {
			o++
		}
	}
	p.objCount = o
}
