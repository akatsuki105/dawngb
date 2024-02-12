package video

import (
	"github.com/akatsuki105/dawngb/util"
)

func (v *Video) Read(addr uint16) uint8 {
	if addr >= 0xFE00 && addr <= 0xFE9F {
		return v.oam[addr&0xFF]
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		if !v.canAccessVRAM() {
			return 0xFF
		}
		return v.ram.data[(v.ram.bank<<13)|uint(addr&0x1FFF)]
	}

	v.CatchUp()

	switch addr {
	case 0xFF40:
		return v.lcdc
	case 0xFF41:
		if !util.Bit(v.lcdc, 7) {
			return 0x80
		}
		return v.stat
	case 0xFF44:
		return uint8(v.ly)
	case 0xFF45:
		return v.lyc
	case 0xFF4F:
		val := uint8(0xFE)
		val = util.SetBit(val, 0, v.ram.bank == 1)
		return val
	case 0xFF69:
		// ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
		return v.r.GetBGPD()
	case 0xFF6B:
		// ゲームによってはパレットの値を読み取ることがある(ロックマンX1など)
		return v.r.GetOBPD()
	default:
		if addr >= 0xFF40 && addr < 0xFF70 {
			return v.ioreg[addr-0xFF40]
		}
	}
	return 0xFF
}

func (v *Video) Write(addr uint16, val uint8) {
	v.CatchUp() // OAMやVRAMへの変更前に今の内容で描画をしておいてもらう

	if addr >= 0xFE00 && addr <= 0xFE9F {
		v.oam[addr&0xFF] = val
		return
	}

	switch addr >> 12 {
	case 0x8, 0x9:
		v.ram.data[(v.ram.bank<<13)|uint(addr&0x1FFF)] = val
		return
	}

	switch addr {
	case 0xFF40:
		wasEnabled := util.Bit(v.lcdc, 7)
		v.lcdc = val
		v.r.SetLCDC(val)
		if wasEnabled != util.Bit(v.lcdc, 7) { // Toggle
			v.stat = (v.stat & 0xFC)
			v.lx, v.ly = 0, 0
			if util.Bit(v.lcdc, 7) {
				v.enableLatch = true
			}
		}
	case 0xFF41:
		oldStat := v.stat
		v.stat = (v.stat & 0x7) | (val & 0x78)
		if !statIRQAsserted(oldStat) && statIRQAsserted(v.stat) {
			v.onInterrupt(1)
		}
	case 0xFF42:
		v.r.SetSCY(val)
	case 0xFF43:
		v.r.SetSCX(val)
	case 0xFF44:
		v.ly = 0
		v.compareLYC()
	case 0xFF45:
		v.lyc = val
		v.compareLYC()
	case 0xFF47:
		v.r.SetBGP(val)
	case 0xFF48:
		v.r.SetOBP0(val)
	case 0xFF49:
		v.r.SetOBP1(val)
	case 0xFF4A:
		v.r.SetWY(val)
	case 0xFF4B:
		v.r.SetWX(val)
	case 0xFF4F:
		v.ram.bank = uint(val & 0b1)
	case 0xFF68:
		v.r.SetBGPI(val)
	case 0xFF69:
		v.ioreg[0x28] = v.r.SetBGPD(val)
	case 0xFF6A:
		v.r.SetOBPI(val)
	case 0xFF6B:
		v.ioreg[0x2A] = v.r.SetOBPD(val)
	}
	if addr >= 0xFF40 && addr < 0xFF70 {
		v.ioreg[addr-0xFF40] = val
	}
}

func (v *Video) canAccessVRAM() bool {
	if util.Bit(v.lcdc, 7) {
		mode := v.stat & 0b11
		switch mode {
		case 2:
			return ((v.lx >> 2) != 20)
		case 3:
			return false
		}
	}
	return true
}
