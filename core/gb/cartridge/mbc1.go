package cartridge

import . "github.com/akatsuki105/dugb/util/datasize"

type mbc1 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint8
	ramBank    uint8
	mode       uint8
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
		return m.c.rom[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		bank := m.c.rom[(16*KB)*uint(m.romBank):]
		return bank[addr&0x3FFF]
	case 0xA, 0xB:
		if m.ramEnabled {
			bank := m.c.ram[(8*KB)*uint(m.ramBank):]
			return bank[addr&0x1FFF]
		}
	}
	return 0xFF
}

func (m *mbc1) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		if val&0xF == 0x0A {
			m.ramEnabled = true
		} else {
			m.ramEnabled = false
		}
	case 0x2, 0x3:
		m.romBank &= 0b110_0000
		m.romBank |= val & 0b11111
		if m.romBank&0x1F == 0 {
			m.romBank |= 0x1
		}
	case 0x4, 0x5:
		if m.mode == 0 {
			m.romBank &= 0b001_1111
			m.romBank |= (val & 0b11) << 5
		} else {
			m.ramBank = val & 0b11
		}
	case 0x6, 0x7:
		m.mode = val & 0b1
	case 0xA, 0xB:
		if m.ramEnabled {
			bank := m.c.ram[(8*KB)*uint(m.ramBank):]
			addr &= 0x1FFF
			if addr < uint16(len(bank)) {
				bank[addr] = val
			}
		}
	}
}
