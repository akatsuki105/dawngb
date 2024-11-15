package ppu

import (
	"github.com/akatsuki105/dawngb/util"
)

func (p *PPU) Read(addr uint16) uint8 {
	if addr >= 0xFE00 && addr <= 0xFE9F {
		return p.oam[addr&0xFF]
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		if !p.canAccessVRAM() {
			return 0xFF
		}
		return p.ram.data[(p.ram.bank<<13)|uint(addr&0x1FFF)]
	}

	switch addr {
	case 0xFF40:
		return p.lcdc
	case 0xFF41:
		if !util.Bit(p.lcdc, 7) {
			return 0x80
		}
		return p.stat
	case 0xFF44:
		return uint8(p.ly)
	case 0xFF45:
		return p.lyc
	case 0xFF4F:
		val := uint8(0xFE)
		val = util.SetBit(val, 0, p.ram.bank == 1)
		return val
	case 0xFF69:
		// ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
		return p.r.GetBGPD()
	case 0xFF6B:
		// ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
		return p.r.GetOBPD()
	default:
		if addr >= 0xFF40 && addr < 0xFF70 {
			return p.ioreg[addr-0xFF40]
		}
	}
	return 0xFF
}

func (p *PPU) Write(addr uint16, val uint8) {
	if addr >= 0xFE00 && addr <= 0xFE9F {
		p.oam[addr&0xFF] = val
		return
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		p.ram.data[(p.ram.bank<<13)|uint(addr&0x1FFF)] = val
		return
	}

	switch addr {
	case 0xFF40:
		wasEnabled := util.Bit(p.lcdc, 7)
		p.lcdc = val
		p.r.SetLCDC(val)
		if wasEnabled != util.Bit(p.lcdc, 7) { // Toggle
			p.stat = (p.stat & 0xFC)
			p.lx, p.ly = 0, 0
			if util.Bit(p.lcdc, 7) {
				p.enableLatch = true
			}
		}
	case 0xFF41:
		oldStat := p.stat
		p.stat = (p.stat & 0x7) | (val & 0x78)
		if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
			p.irq(1)
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
		p.r.SetOBP0(val)
	case 0xFF49:
		p.r.SetOBP1(val)
	case 0xFF4A:
		p.r.SetWY(val)
	case 0xFF4B:
		p.r.SetWX(val)
	case 0xFF4F:
		p.ram.bank = uint(val & 0b1)
	case 0xFF68:
		p.r.SetBGPI(val)
	case 0xFF69:
		p.ioreg[0x28] = p.r.SetBGPD(val)
	case 0xFF6A:
		p.r.SetOBPI(val)
	case 0xFF6B:
		p.ioreg[0x2A] = p.r.SetOBPD(val)
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
