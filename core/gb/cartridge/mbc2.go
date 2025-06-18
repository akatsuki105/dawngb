package cartridge

type MBC2 struct {
	c          *Cartridge
	RAMEnabled bool
	ROMBank    uint8 // 0..15
}

func newMBC2(c *Cartridge) *MBC2 {
	return &MBC2{
		c:       c,
		ROMBank: 1,
	}
}

func (m *MBC2) reset() {
	m.RAMEnabled = false
	m.ROMBank = 1
}

func (m *MBC2) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3: // ROMバンク0
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7: // ROMバンク1..15
		return m.c.ROM[(uint(m.ROMBank)<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB: // RAM
		if m.RAMEnabled {
			return 0xF0 | (m.c.RAM[addr&0x1FF] & 0x0F)
		}
	}
	return 0xFF
}

func (m *MBC2) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3: // RAM有効化 or ROMバンク切り替え
		mode := (addr >> 8) & 0x1
		if mode == 0 {
			m.RAMEnabled = val == 0x0A
		} else {
			m.ROMBank = val & 0x0F
			if m.ROMBank == 0 {
				m.ROMBank = 1
			}
		}
	case 0xA, 0xB: // RAM書き込み
		if m.RAMEnabled {
			m.c.RAM[addr&0x1FF] = val & 0x0F
		}
	}
}

type MBC2Snapshot struct {
	Header     uint64
	RAMEnabled bool
	ROMBank    uint8
	Reserved   [14]uint8
}

func (m *MBC2) CreateSnapshot() MBC2Snapshot {
	return MBC2Snapshot{
		RAMEnabled: m.RAMEnabled,
		ROMBank:    m.ROMBank,
	}
}

func (m *MBC2) RestoreSnapshot(snap MBC2Snapshot) bool {
	m.RAMEnabled = snap.RAMEnabled
	m.ROMBank = snap.ROMBank
	return true
}
