package cartridge

import (
	"github.com/akatsuki105/dawngb/core/gb/internal"
)

type Time struct {
	Sec, Min, Hour uint8
	DayCarry       bool
	Day            uint16
}

type RTC struct {
	Enabled     bool
	Time, Latch Time
}

type MBC3 struct {
	c                *Cartridge
	RAMEnabled       bool
	ROMBank, RAMBank uint8
	RAMBankMax       uint8 // 4(Normal) or 8(MBC30)
	RTC              RTC
}

func newMBC3(c *Cartridge) *MBC3 {
	m := &MBC3{
		c:          c,
		ROMBank:    1,
		RAMBankMax: 4,
	}
	if m.isMBC30() {
		m.RAMBankMax = 8
	}
	return m
}

// ポケモンクリスタルなどは、MBC30と呼ばれる特殊なMBC3を使っている
// これを見分ける方法は今のところ、カートリッジヘッダのROMサイズとRAMサイズを見るしかない
func (m *MBC3) isMBC30() bool {
	return (len(m.c.ROM) > int(2*MB)) || (len(m.c.RAM) > int(32*KB))
}

func (m *MBC3) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[(uint(m.ROMBank)<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB:
		if m.RAMEnabled {
			if m.RAMBank < m.RAMBankMax {
				return m.c.RAM[(uint(m.RAMBank)<<13)|uint(addr&0x1FFF)]
			}

			// RTC
			switch m.RAMBank {
			case 0x8:
				return m.RTC.Latch.Sec
			case 0x9:
				return m.RTC.Latch.Min
			case 0xA:
				return m.RTC.Latch.Hour
			case 0xB:
				return uint8(m.RTC.Latch.Day & 0xFF)
			case 0xC:
				val := uint8(0x0)
				val = internal.SetBit(val, 0, m.RTC.Latch.Day >= 0x100)
				val = internal.SetBit(val, 6, !m.RTC.Enabled)
				val = internal.SetBit(val, 7, m.RTC.Latch.DayCarry)
				return val
			}
		}
	}
	return 0xFF
}

func (m *MBC3) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.RAMEnabled = (val&0x0F == 0x0A)
	case 0x2, 0x3:
		m.ROMBank = (val & 0b111_1111)
		if m.isMBC30() {
			m.ROMBank = val
		}
		if m.ROMBank == 0 {
			m.ROMBank = 1
		}
	case 0x4, 0x5:
		if val <= 0x0C {
			m.RAMBank = val
		}
	case 0x6, 0x7:
		m.RTC.Latch = m.RTC.Time // NOTE: 任天堂のドキュメントはここに0と1を書き込むことでラッチすると書いてあるが、実際には何を書き込んでもすぐにラッチされる
	case 0xA, 0xB:
		if m.RAMEnabled {
			if m.RAMBank < m.RAMBankMax {
				m.c.RAM[(uint(m.RAMBank)<<13)|uint(addr&0x1FFF)] = val
			} else {
				switch m.RAMBank {
				case 0xC:
					m.RTC.Time.Day &= 0xFF
					m.RTC.Time.Day |= uint16(val&0x1) << 8
					m.RTC.Enabled = !internal.Bit(val, 6)
					m.RTC.Time.DayCarry = internal.Bit(val, 7)
				}
			}
		}
	}
}

type MBC3Snapshot struct {
	Header           uint64
	RAMEnabled       bool
	ROMBank, RAMBank uint8
	RAMBankMax       uint8 // 4(Normal) or 8(MBC30)
	RTC              RTC
}

func (m *MBC3) CreateSnapshot() MBC3Snapshot {
	return MBC3Snapshot{
		RAMEnabled: m.RAMEnabled,
		ROMBank:    m.ROMBank,
		RAMBank:    m.RAMBank,
		RAMBankMax: m.RAMBankMax,
		RTC:        m.RTC,
	}
}

func (m *MBC3) RestoreSnapshot(snap *MBC3Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	m.RAMEnabled = snap.RAMEnabled
	m.ROMBank, m.RAMBank = snap.ROMBank, snap.RAMBank
	m.RAMBankMax = snap.RAMBankMax
	m.RTC = snap.RTC
	return nil
}
