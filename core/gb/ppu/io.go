package ppu

import (
	"github.com/akatsuki105/dawngb/util"
)

func (p *PPU) Read(addr uint16) uint8 {
	if addr >= 0xFE00 && addr <= 0xFE9F {
		return p.OAM[addr&0xFF]
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		if !p.canAccessVRAM() {
			return 0xFF
		}
		return p.ram.data[(uint(p.ram.bank)<<13)|uint(addr&0x1FFF)]
	}

	switch addr {
	case 0xFF40:
		return p.lcdc
	case 0xFF41:
		if !util.Bit(p.lcdc, 7) {
			return 0x80
		}
		return p.stat | 0x80
	case 0xFF44:
		return uint8(p.ly)
	case 0xFF45:
		return p.lyc
	case 0xFF4F:
		return 0xFE | (p.ram.bank & 1)
	case 0xFF69:
		// ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
		return uint8(p.Palette[(p.bgpi>>1)] >> ((p.bgpi & 1) * 8))
	case 0xFF6B:
		// ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
		return uint8(p.Palette[32+(p.obpi>>1)] >> ((p.obpi & 1) * 8))
	default:
		if addr >= 0xFF40 && addr < 0xFF70 {
			return p.ioreg[addr-0xFF40]
		}
	}
	return 0xFF
}

func (p *PPU) Write(addr uint16, val uint8) {
	if addr >= 0xFE00 && addr <= 0xFE9F {
		p.OAM[addr&0xFF] = val
		return
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		p.ram.data[(uint(p.ram.bank)<<13)|uint(addr&0x1FFF)] = val
		return
	}

	switch addr {
	case 0xFF40: // LCDC
		wasEnabled := (p.lcdc & (1 << 7)) != 0
		p.lcdc = val
		p.r.SetLCDC(val)
		enabled := (val & (1 << 7)) != 0
		if wasEnabled != enabled { // Toggle
			p.stat = (p.stat & 0xFC)
			p.lx, p.ly = 0, 0
			if enabled { // Turn on
				p.enableLatch = true
			}
		}
	case 0xFF41:
		oldStat := p.stat
		p.stat = (p.stat & 0x7) | (val & 0x78)
		if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
			p.cpu.IRQ(1)
		}
	case 0xFF42:
		p.r.SetSCY(val)
	case 0xFF43:
		p.r.SetSCX(val)
	case 0xFF44:
		p.ly = 0
		p.compareLYC()
	case 0xFF45:
		p.lyc = val
		p.compareLYC()
	case 0xFF47:
		p.r.SetBGP(val)
	case 0xFF48:
		p.r.SetOBP(0, val)
	case 0xFF49:
		p.r.SetOBP(1, val)
	case 0xFF4A:
		p.r.SetWY(val)
	case 0xFF4B:
		p.r.SetWX(val)
	case 0xFF4F:
		p.ram.bank = val & 0b1
	case 0xFF68:
		p.bgpi = val
	case 0xFF69:
		p.setBGPD(val)
	case 0xFF6A:
		p.obpi = val
	case 0xFF6B:
		p.setOBPD(val)
	}
	if addr >= 0xFF40 && addr < 0xFF70 {
		p.ioreg[addr-0xFF40] = val
	}
}

func (p *PPU) canAccessVRAM() bool {
	if util.Bit(p.lcdc, 7) {
		mode := p.stat & 0b11
		switch mode {
		case 2:
			return ((p.lx >> 2) != 20)
		case 3:
			return false
		}
	}
	return true
}

func (p *PPU) setBGPD(val uint8) {
	palID := int((p.bgpi & 0x3F) / 8)
	colorID := int(p.bgpi&7) >> 1
	idx := ((palID * 4) + colorID) & 0x1F
	rgb555 := p.Palette[idx]
	isHi := (p.bgpi & 1) == 1
	if isHi {
		rgb555 = (rgb555 & 0x00FF) | (uint16(val) << 8)
	} else {
		rgb555 = (rgb555 & 0xFF00) | (uint16(val) << 0)
	}
	p.Palette[idx] = rgb555

	if (p.bgpi & (1 << 7)) != 0 {
		bgpi := (p.bgpi + 1) & 0x3F
		p.bgpi &= 0xC0
		p.bgpi |= bgpi
	}
}

func (p *PPU) setOBPD(val uint8) {
	palID := int((p.obpi & 0x3F) / 8)
	colorID := int(p.obpi&7) >> 1
	idx := 32 | ((palID*4 + colorID) & 0x1F)
	rgb555 := p.Palette[idx]
	isHi := (p.obpi & 1) == 1
	if isHi {
		rgb555 = (rgb555 & 0x00FF) | (uint16(val) << 8)
	} else {
		rgb555 = (rgb555 & 0xFF00) | uint16(val)
	}
	p.Palette[idx] = rgb555

	if (p.obpi & (1 << 7)) != 0 {
		obpi := (p.obpi + 1) & 0x3F
		p.obpi &= 0xC0
		p.obpi |= obpi
	}
}
