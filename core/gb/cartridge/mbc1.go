package cartridge

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
		romBank := uint(m.romBank)
		if m.mode == 0 {
			if len(m.c.rom) >= int(1*MB) {
				romBank |= (uint(m.ramBank) << 5)
			}
		}
		return m.c.rom[(romBank<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB:
		if m.ramEnabled {
			ramBank := uint(0)
			if m.mode == 1 {
				if len(m.c.ram) >= int(32*KB) {
					ramBank = uint(m.ramBank)
				}
			}
			return m.c.ram[(ramBank<<13)|uint(addr&0x1FFF)]
		}
	}
	return 0xFF
}

func (m *mbc1) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.ramEnabled = (val&0x0F == 0x0A)
	case 0x2, 0x3:
		m.romBank = val & 0b11111
		if m.romBank == 0 {
			m.romBank = 1
		}
	case 0x4, 0x5:
		m.ramBank = val & 0b11
	case 0x6, 0x7:
		m.mode = val & 0b1
	case 0xA, 0xB:
		if m.ramEnabled {
			ramBank := uint(0)
			if m.mode == 1 {
				ramBank = uint(m.ramBank)
			}
			bank := m.c.ram[(8*KB)*ramBank:]
			addr &= 0x1FFF
			if len(bank) > int(addr) {
				bank[addr] = val
			}
		}
	}
}
