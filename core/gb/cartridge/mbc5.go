package cartridge

type mbc5 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint
	ramBank    uint8
}

func newMBC5(c *Cartridge) mbc {
	return &mbc5{c: c}
}

func (m *mbc5) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.rom[addr]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.rom[(uint32(m.romBank)<<14)|(uint32(addr&0x3FFF))]
	case 0xA, 0xB:
		if m.ramEnabled {
			return m.c.ram[(uint32(m.ramBank)<<13)|(uint32(addr&0x1FFF))]
		}
	}
	return 0xFF
}

func (m *mbc5) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		if val&0b1111 == 0x0A {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}
	case 0x2:
		m.romBank &= 0x100
		m.romBank |= uint(val)
	case 0x3:
		m.romBank &= 0xFF
		m.romBank |= uint(val&0b1) << 8
	case 0x4, 0x5:
		m.ramBank = val & 0b1111
	case 0xA, 0xB:
		if m.ramEnabled {
			m.c.ram[(uint32(m.ramBank)<<13)|(uint32(addr&0x1FFF))] = val
		}
	}
}
