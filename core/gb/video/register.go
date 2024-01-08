package video

func (v *Video) ReadIO(addr uint16) uint8 {
	switch addr {
	case 0xFF40:
		return v.lcdc
	case 0xFF41:
		return v.stat
	case 0xFF42:
		return v.scy
	case 0xFF43:
		return v.scx
	case 0xFF44:
		return uint8(v.ly)
	case 0xFF45:
		return v.lyc
	case 0xFF47:
		return v.bgp
	case 0xFF48:
		return v.obp0
	case 0xFF49:
		return v.obp1
	case 0xFF4A:
		return v.wy
	case 0xFF4B:
		return v.wx
	case 0xFF4F:
		return uint8(v.ram.bank)
	default:
		return v.ioreg[addr-0xFF40]
	}
}

func (v *Video) WriteIO(addr uint16, val uint8) {
	switch addr {
	case 0xFF40:
		v.lcdc = val
		v.r.SetLCDC(val)
	case 0xFF41:
		v.stat = (v.stat & 0x7) | (val & 0x78)
	case 0xFF42:
		v.scy = val
		v.r.SetSCY(val)
	case 0xFF43:
		v.scx = val
		v.r.SetSCX(val)
	case 0xFF45:
		v.lyc = val
	case 0xFF47:
		v.bgp = val
		v.r.SetBGP(val)
	case 0xFF48:
		v.obp0 = val
		v.r.SetOBP0(val)
	case 0xFF49:
		v.obp1 = val
		v.r.SetOBP1(val)
	case 0xFF4A:
		v.wy = val
		v.r.SetWY(val)
	case 0xFF4B:
		v.wx = val
		v.r.SetWX(val)
	case 0xFF4F:
		v.ram.bank = int(val & 0b1)
	case 0xFF68:
		v.r.SetBGPI(val)
	case 0xFF69:
		v.r.SetBGPD(val)
	case 0xFF6A:
		v.r.SetOBPI(val)
	case 0xFF6B:
		v.r.SetOBPD(val)
	default:
		v.ioreg[addr-0xFF40] = val
	}
}
