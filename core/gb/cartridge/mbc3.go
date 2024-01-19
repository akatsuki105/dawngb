package cartridge

import . "github.com/akatsuki105/dugb/util/datasize"

type rtc struct {
	latch uint8
}

type mbc3 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint
	ramBank    uint
	rtc        rtc
	ramBankMax uint
}

func newMBC3(c *Cartridge) mbc {
	m := &mbc3{
		c:          c,
		ramBankMax: 4,
	}
	if m.isMBC30() {
		m.ramBankMax = 8
	}
	return m
}

func (m *mbc3) isMBC30() bool {
	return len(m.c.ram) == int(64*KB)
}

func (m *mbc3) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.rom[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.rom[(m.romBank<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB:
		if m.ramEnabled {
			if m.ramBank < m.ramBankMax {
				return m.c.ram[(m.ramBank<<13)|uint(addr&0x1FFF)]
			}

			// RTC
			return m.rtc.latch
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
		m.romBank = uint(val & 0b111_1111)
		if m.romBank == 0 {
			m.romBank |= 0x1
		}
	case 0x4, 0x5:
		if val <= 0x0C {
			m.ramBank = uint(val)
		}
	case 0xA, 0xB:
		if m.ramEnabled {
			if m.ramBank < m.ramBankMax {
				m.c.ram[(m.ramBank<<13)|uint(addr&0x1FFF)] = val
			}
		}
	}
}
