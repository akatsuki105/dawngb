package cartridge

type mbc0 struct {
	c *Cartridge
}

func newMBC0(c *Cartridge) mbc {
	return &mbc0{c: c}
}

func (m *mbc0) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		return m.c.rom[addr]
	case 0xA, 0xB:
		return m.c.ram[addr&0x1FFF]
	}
	return 0xFF
}

func (m *mbc0) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0xA, 0xB:
		m.c.ram[addr&0x1FFF] = val
	}
}
