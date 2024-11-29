package cartridge

type mbc5 struct {
	c          *Cartridge
	hasRam     bool
	ramEnabled bool
	romBank    uint16 // 0..511
	ramBank    uint8  // 0..15
}

func newMBC5(c *Cartridge) mbc {
	hasRam := c.ROM[0x147] != 25 && c.ROM[0x147] != 28
	return &mbc5{
		c:       c,
		hasRam:  hasRam,
		romBank: 1,
	}
}

func (m *mbc5) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[(uint32(m.romBank)<<14)|(uint32(addr&0x3FFF))]
	case 0xA, 0xB:
		if m.hasRam && m.ramEnabled {
			n := int((uint(m.ramBank) << 13) | uint(addr&0x1FFF))
			if n >= len(m.c.ram) {
				n &= len(m.c.ram) - 1
			}
			return m.c.ram[n]
		}
	}
	return 0xFF
}

func (m *mbc5) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.ramEnabled = (val&0x0F == 0x0A)
	case 0x2:
		m.romBank &= 0x100
		m.romBank |= uint16(val)
	case 0x3:
		m.romBank &= 0xFF
		m.romBank |= uint16(val&0b1) << 8
	case 0x4, 0x5:
		m.ramBank = (val & 0b1111)
	case 0xA, 0xB:
		if m.hasRam && m.ramEnabled {
			n := int((uint(m.ramBank) << 13) | uint(addr&0x1FFF))
			if n >= len(m.c.ram) {
				n &= len(m.c.ram) - 1
			}
			m.c.ram[n] = val
		}
	}
}
