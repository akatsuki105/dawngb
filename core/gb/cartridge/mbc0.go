package cartridge

type mbc0 struct {
	c      *Cartridge
	hasRam bool
}

func newMBC0(c *Cartridge) mbc {
	hasRam := c.ROM[0x147] != 0
	return &mbc0{
		c:      c,
		hasRam: hasRam,
	}
}

func (m *mbc0) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[addr]
	case 0xA, 0xB:
		if m.hasRam {
			return m.c.ram[addr&0x1FFF]
		}
	}
	return 0xFF
}

func (m *mbc0) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0xA, 0xB:
		if m.hasRam {
			m.c.ram[addr&0x1FFF] = val
		}
	}
}
