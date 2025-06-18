package cartridge

type MBC5 struct {
	c          *Cartridge
	hasRam     bool
	RAMEnabled bool
	ROMBank    uint16 // 0..511
	RAMBank    uint8  // 0..15
}

func newMBC5(c *Cartridge) *MBC5 {
	hasRam := c.ROM[0x147] != 25 && c.ROM[0x147] != 28
	return &MBC5{
		c:       c,
		hasRam:  hasRam,
		ROMBank: 1,
	}
}

func (m *MBC5) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[(uint32(m.ROMBank)<<14)|(uint32(addr&0x3FFF))]
	case 0xA, 0xB:
		if m.hasRam && m.RAMEnabled {
			n := int((uint(m.RAMBank) << 13) | uint(addr&0x1FFF))
			if n >= len(m.c.RAM) {
				n &= len(m.c.RAM) - 1
			}
			return m.c.RAM[n]
		}
	}
	return 0xFF
}

func (m *MBC5) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.RAMEnabled = (val&0x0F == 0x0A)
	case 0x2:
		m.ROMBank &= 0x100
		m.ROMBank |= uint16(val)
	case 0x3:
		m.ROMBank &= 0xFF
		m.ROMBank |= uint16(val&0b1) << 8
	case 0x4, 0x5:
		m.RAMBank = (val & 0b1111)
	case 0xA, 0xB:
		if m.hasRam && m.RAMEnabled {
			n := int((uint(m.RAMBank) << 13) | uint(addr&0x1FFF))
			if n >= len(m.c.RAM) {
				n &= len(m.c.RAM) - 1
			}
			m.c.RAM[n] = val
		}
	}
}

type MBC5Snapshot struct {
	Header     uint64
	RAMEnabled bool
	ROMBank    uint16
	RAMBank    uint8
}

func (m *MBC5) CreateSnapshot() MBC5Snapshot {
	return MBC5Snapshot{
		RAMEnabled: m.RAMEnabled,
		ROMBank:    m.ROMBank,
		RAMBank:    m.RAMBank,
	}
}

func (m *MBC5) RestoreSnapshot(snap *MBC5Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	m.RAMEnabled = snap.RAMEnabled
	m.ROMBank, m.RAMBank = snap.ROMBank, snap.RAMBank
	return nil
}
