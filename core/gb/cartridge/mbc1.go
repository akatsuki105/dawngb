package cartridge

type MBC1 struct {
	c                *Cartridge
	ramEnabled       bool
	romBank, ramBank uint8
	Mode             uint8
}

func newMBC1(c *Cartridge) *MBC1 {
	return &MBC1{
		c:       c,
		romBank: 1,
	}
}

func (m *MBC1) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		romBank := uint(m.romBank)
		if m.Mode == 0 {
			if len(m.c.ROM) >= int(1*MB) {
				romBank |= (uint(m.ramBank) << 5)
			}
		}
		return m.c.ROM[(romBank<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB:
		if m.ramEnabled {
			ramBank := uint(0)
			if m.Mode == 1 {
				if len(m.c.RAM) >= int(32*KB) {
					ramBank = uint(m.ramBank)
				}
			}
			return m.c.RAM[(ramBank<<13)|uint(addr&0x1FFF)]
		}
	}
	return 0xFF
}

func (m *MBC1) write(addr uint16, val uint8) {
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
		m.Mode = val & 0b1
	case 0xA, 0xB:
		if m.ramEnabled {
			ramBank := uint(0)
			if m.Mode == 1 {
				ramBank = uint(m.ramBank)
			}
			bank := m.c.RAM[(8*KB)*ramBank:]
			addr &= 0x1FFF
			if len(bank) > int(addr) {
				bank[addr] = val
			}
		}
	}
}

type MBC1Snapshot struct {
	Header           uint64
	RamEnabled       bool
	ROMBank, RAMBank uint8
	Mode             uint8
	Reserved         [16]uint8
}

func (m *MBC1) CreateSnapshot() MBC1Snapshot {
	return MBC1Snapshot{
		RamEnabled: m.ramEnabled,
		ROMBank:    m.romBank,
		RAMBank:    m.ramBank,
		Mode:       m.Mode,
	}
}

func (m *MBC1) RestoreSnapshot(snap MBC1Snapshot) bool {
	m.ramEnabled = snap.RamEnabled
	m.romBank, m.ramBank = snap.ROMBank, snap.RAMBank
	m.Mode = snap.Mode
	return true
}
