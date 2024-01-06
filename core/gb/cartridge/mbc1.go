package cartridge

type mbc1 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint8
}

func newMBC1(c *Cartridge) mbc {
	return &mbc1{
		c:       c,
		romBank: 1,
	}
}

func (m *mbc1) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.rom[addr]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.rom[(uint32(m.romBank)<<14)|(uint32(addr&0x3FFF))]
	case 0xA, 0xB:
		if m.ramEnabled {
			return m.c.ram[addr&0x1FFF]
		}
	}
	return 0xFF
}

func (m *mbc1) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		if val&0b1111 == 0x0A {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}
	case 0x2, 0x3:
		m.romBank = val & 0b11111
	}
}
