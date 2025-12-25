package ppu

import "github.com/akatsuki105/dawngb/core/gb/internal"

func (p *PPU) Read(addr uint16) uint8 {
	if addr >= 0xFE00 && addr <= 0xFE9F {
		return p.OAM[addr&0xFF]
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		if !p.canAccessVRAM() {
			return 0xFF
		}
		return p.RAM.Data[(uint(p.RAM.Bank)<<13)|uint(addr&0x1FFF)]
	}

	switch addr {
	case 0xFF40:
		return p.LCDC
	case 0xFF41:
		if (p.LCDC & (1 << 7)) == 0 {
			return 0x80
		}
		return p.STAT | 0x80
	case 0xFF44:
		return uint8(p.Ly)
	case 0xFF45:
		return p.LYC
	case 0xFF4F:
		return 0xFE | (p.RAM.Bank & 1)
	case 0xFF69:
		return internal.Byte(p.Palette[((p.BGPI&0x3F)>>1)], p.BGPI&1) // ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
	case 0xFF6B:
		return internal.Byte(p.Palette[32+((p.OBPI&0x3F)>>1)], p.OBPI&1) // ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
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
		p.RAM.Data[(uint(p.RAM.Bank)<<13)|uint(addr&0x1FFF)] = val
		return
	}

	switch addr {
	case 0xFF40: // LCDC
		wasEnabled := (p.LCDC & (1 << 7)) != 0
		p.LCDC = val
		p.r.SetLCDC(val)
		enabled := (val & (1 << 7)) != 0
		if wasEnabled != enabled { // Toggle
			p.STAT = (p.STAT & 0xFC)
			p.Lx, p.Ly = 0, 0
			if enabled { // Turn on
				p.enableLatch = true
			}
		}
	case 0xFF41:
		oldStat := p.STAT
		p.STAT = (p.STAT & 0x7) | (val & 0x78)
		if !statIRQAsserted(oldStat) && statIRQAsserted(p.STAT) {
			p.cpu.IRQ(1)
		}
	case 0xFF42:
		p.r.SetSCY(val)
	case 0xFF43:
		p.r.SetSCX(val)
	case 0xFF44:
		p.Ly = 0
		p.compareLYC()
	case 0xFF45:
		p.LYC = val
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
		p.RAM.Bank = val & 0b1
		if !p.cpu.IsCGBMode() {
			p.RAM.Bank = 0
		}
	case 0xFF68:
		p.BGPI = val
	case 0xFF69:
		p.setBGPD(val)
	case 0xFF6A:
		p.OBPI = val
	case 0xFF6B:
		p.setOBPD(val)
	}
	if addr >= 0xFF40 && addr < 0xFF70 {
		p.ioreg[addr-0xFF40] = val
	}
}

func (p *PPU) canAccessVRAM() bool {
	if (p.LCDC & (1 << 7)) != 0 {
		mode := p.STAT & 0b11
		switch mode {
		case 2:
			return ((p.Lx >> 2) != 20)
		case 3:
			return false
		}
	}
	return true
}

func (p *PPU) setBGPD(val uint8) {
	palID := int((p.BGPI & 0x3F) / 8)
	colorID := int(p.BGPI&7) >> 1
	idx := ((palID * 4) + colorID) & 0x1F
	p.Palette[idx] = internal.SetByte(p.Palette[idx], p.BGPI&1, val)

	if (p.BGPI & (1 << 7)) != 0 {
		bgpi := (p.BGPI + 1) & 0x3F
		p.BGPI &= 0xC0
		p.BGPI |= bgpi
	}
}

func (p *PPU) setOBPD(val uint8) {
	palID := int((p.OBPI & 0x3F) / 8)
	colorID := int(p.OBPI&7) >> 1
	idx := 32 | ((palID*4 + colorID) & 0x1F)
	p.Palette[idx] = internal.SetByte(p.Palette[idx], p.OBPI&1, val)

	if (p.OBPI & (1 << 7)) != 0 {
		obpi := (p.OBPI + 1) & 0x3F
		p.OBPI &= 0xC0
		p.OBPI |= obpi
	}
}
