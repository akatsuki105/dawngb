package cartridge

type MBC0 struct {
	c *Cartridge
}

func newMBC0(c *Cartridge) *MBC0 {
	return &MBC0{
		c: c,
	}
}

func (m *MBC0) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[addr]
	case 0xA, 0xB:
		addr &= 0x1FFF
		if len(m.c.RAM) > int(addr) {
			return m.c.RAM[addr]
		}
	}
	return 0xFF
}

func (m *MBC0) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0xA, 0xB:
		addr &= 0x1FFF
		if len(m.c.RAM) > int(addr) {
			m.c.RAM[addr] = val
		}
	}
}
