package cartridge

import . "github.com/akatsuki105/dugb/util/datasize"

type mbc3 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint8
	ramBank    uint8
}

func newMBC3(c *Cartridge) mbc {
	return &mbc3{c: c}
}

func (m *mbc3) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.rom[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		bank := m.c.rom[(16*KB)*uint(m.romBank):]
		return bank[addr&0x3FFF]
	case 0xA, 0xB:
		if m.ramEnabled {
			return m.c.ram[(uint32(m.ramBank)<<13)|(uint32(addr&0x1FFF))]
		}
	}
	return 0xFF
}

func (m *mbc3) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		if val&0xF == 0x0A {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}
	case 0x2, 0x3:
		m.romBank = val & 0b111_1111
		if m.romBank&0x1F == 0 {
			m.romBank |= 0x1
		}
	case 0x4, 0x5:
		m.ramBank = val & 0b11
	case 0xA, 0xB:
		if m.ramEnabled {
			m.c.ram[(uint32(m.ramBank)<<13)|(uint32(addr&0x1FFF))] = val
		}
	}
}
