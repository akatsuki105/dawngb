package cartridge

type mbc2 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint8 // 0..15
}

func newMBC2(c *Cartridge) *mbc2 {
	return &mbc2{
		c:       c,
		romBank: 1,
	}
}

func (m *mbc2) reset() {
	m.ramEnabled = false
	m.romBank = 1
}

func (m *mbc2) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3: // ROMバンク0
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7: // ROMバンク1..15
		return m.c.ROM[(uint(m.romBank)<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB: // RAM
		if m.ramEnabled {
			return 0xF0 | (m.c.ram[addr&0x1FF] & 0x0F)
		}
	}
	return 0xFF
}

func (m *mbc2) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3: // RAM有効化 or ROMバンク切り替え
		mode := (addr >> 8) & 0x1
		if mode == 0 {
			m.ramEnabled = val == 0x0A
		} else {
			m.romBank = val & 0x0F
			if m.romBank == 0 {
				m.romBank = 1
			}
		}
	case 0xA, 0xB: // RAM書き込み
		if m.ramEnabled {
			m.c.ram[addr&0x1FF] = val & 0x0F
		}
	}
}
