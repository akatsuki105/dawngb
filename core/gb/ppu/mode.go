package ppu

// Mode 0
func (p *PPU) hblank() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC)
	if ((p.lcdc & (1 << 7)) != 0) && !p.enableLatch {
		p.r.DrawScanline(p.ly, p.screen[p.ly*160:(p.ly+1)*160])
	}
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
		p.StatIRQ.Triggered = true
		p.StatIRQ.Mode, p.StatIRQ.Lx, p.StatIRQ.Ly = 0, uint8(p.lx), uint8(p.ly)
	}
	p.cpu.HBlank()
}

// Mode 1
func (p *PPU) vblank() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC) | 1
	p.cpu.IRQ(0)

	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
		p.StatIRQ.Triggered = true
		p.StatIRQ.Mode, p.StatIRQ.Lx, p.StatIRQ.Ly = 1, uint8(p.lx), uint8(p.ly)
	}
}

// Mode 2
func (p *PPU) scanOAM() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC) | 2
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
		p.StatIRQ.Triggered = true
		p.StatIRQ.Mode, p.StatIRQ.Lx, p.StatIRQ.Ly = 2, uint8(p.lx), uint8(p.ly)
	}
}

// Mode 3
func (p *PPU) drawing() {
	p.stat = (p.stat & 0xFC) | 3

	// Count scanline objects
	h := 8
	if (p.lcdc & (1 << 2)) != 0 {
		h = 16
	}
	o := uint8(0)
	for i := 0; i < 40; i++ {
		y := int(p.OAM[i*4]) - 16
		if y <= p.ly && p.ly < y+h {
			o++
		}
	}
	p.objCount = o
}
