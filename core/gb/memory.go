package gb

import (
	"github.com/akatsuki105/dugb/util"
	. "github.com/akatsuki105/dugb/util/datasize"
)

type Memory struct {
	gb       *GB
	wram     [4 * KB * 8]uint8
	wramBank uint
	hram     [0x7F]uint8
}

func newMemory(gb *GB) *Memory {
	return &Memory{
		gb:       gb,
		wramBank: 1,
	}
}

func (m *Memory) Read(addr uint16) byte {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0xA, 0xB:
		return m.gb.cartridge.Read(addr)
	case 0x8, 0x9:
		return (m.gb.video.VRAM())[addr&0x1FFF]
	case 0xC, 0xE:
		return m.wram[addr&0xFFF]
	case 0xD:
		bank := m.wramBank * (4 * KB)
		return m.wram[bank+uint(addr&0xFFF)]
	case 0xF:
		if addr >= 0xFE00 && addr <= 0xFE9F {
			return m.gb.video.OAM[addr&0xFF]
		}
		switch addr {
		case 0xFF00:
			return m.gb.input.ReadIO(addr)
		case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
			return m.gb.timer.ReadIO(addr)
		case 0xFF0F:
			val := uint8(0)
			for i := 0; i < 5; i++ {
				val |= (uint8(util.Btoi(m.gb.interrupt[i])) << i)
			}
			return val
		case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26:
			return m.gb.audio.Read(addr)
		case 0xFF40, 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45, 0xFF47, 0xFF48, 0xFF49, 0xFF4A, 0xFF4B, 0xFF4F, 0xFF68, 0xFF69, 0xFF6A, 0xFF6B:
			return m.gb.video.ReadIO(addr)
		case 0xFF4D:
			is2x := m.gb.cpu.Cycle == 4
			val := uint8(0x7E)
			val = util.SetBit(val, 7, is2x)
			val = util.SetBit(val, 0, m.gb.key1)
			return val
		case 0xFF70:
			return uint8(m.wramBank)
		case 0xFF80, 0xFF81, 0xFF82, 0xFF83, 0xFF84, 0xFF85, 0xFF86, 0xFF87, 0xFF88, 0xFF89, 0xFF8A, 0xFF8B, 0xFF8C, 0xFF8D, 0xFF8E, 0xFF8F, 0xFF90, 0xFF91, 0xFF92, 0xFF93, 0xFF94, 0xFF95, 0xFF96, 0xFF97, 0xFF98, 0xFF99, 0xFF9A, 0xFF9B, 0xFF9C, 0xFF9D, 0xFF9E, 0xFF9F, 0xFFA0, 0xFFA1, 0xFFA2, 0xFFA3, 0xFFA4, 0xFFA5, 0xFFA6, 0xFFA7, 0xFFA8, 0xFFA9, 0xFFAA, 0xFFAB, 0xFFAC, 0xFFAD, 0xFFAE, 0xFFAF, 0xFFB0, 0xFFB1, 0xFFB2, 0xFFB3, 0xFFB4, 0xFFB5, 0xFFB6, 0xFFB7, 0xFFB8, 0xFFB9, 0xFFBA, 0xFFBB, 0xFFBC, 0xFFBD, 0xFFBE, 0xFFBF, 0xFFC0, 0xFFC1, 0xFFC2, 0xFFC3, 0xFFC4, 0xFFC5, 0xFFC6, 0xFFC7, 0xFFC8, 0xFFC9, 0xFFCA, 0xFFCB, 0xFFCC, 0xFFCD, 0xFFCE, 0xFFCF, 0xFFD0, 0xFFD1, 0xFFD2, 0xFFD3, 0xFFD4, 0xFFD5, 0xFFD6, 0xFFD7, 0xFFD8, 0xFFD9, 0xFFDA, 0xFFDB, 0xFFDC, 0xFFDD, 0xFFDE, 0xFFDF, 0xFFE0, 0xFFE1, 0xFFE2, 0xFFE3, 0xFFE4, 0xFFE5, 0xFFE6, 0xFFE7, 0xFFE8, 0xFFE9, 0xFFEA, 0xFFEB, 0xFFEC, 0xFFED, 0xFFEE, 0xFFEF, 0xFFF0, 0xFFF1, 0xFFF2, 0xFFF3, 0xFFF4, 0xFFF5, 0xFFF6, 0xFFF7, 0xFFF8, 0xFFF9, 0xFFFA, 0xFFFB, 0xFFFC, 0xFFFD, 0xFFFE:
			return m.hram[addr&0x7F]
		case 0xFFFF:
			return m.gb.ie
		}
	}
	return 0
}

func (m *Memory) Write(addr uint16, val byte) {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0xA, 0xB:
		m.gb.cartridge.Write(addr, val)
	case 0x8, 0x9:
		(m.gb.video.VRAM())[addr&0x1FFF] = val
	case 0xC, 0xE:
		m.wram[addr&0xFFF] = val
	case 0xD:
		bank := m.wramBank * (4 * KB)
		m.wram[bank+uint(addr&0xFFF)] = val
	case 0xF:
		if addr >= 0xFE00 && addr <= 0xFE9F {
			m.gb.video.OAM[addr&0xFF] = val
			return
		}

		switch addr {
		case 0xFF00:
			m.gb.input.WriteIO(addr, val)
		case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
			m.gb.timer.WriteIO(addr, val)
		case 0xFF0F:
			for i := 0; i < 5; i++ {
				m.gb.interrupt[i] = util.Bit(val, i)
			}
		case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26:
			m.gb.audio.Write(addr, val)
		case 0xFF40, 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45, 0xFF4A, 0xFF4B, 0xFF4F:
			m.gb.video.WriteIO(addr, val)
		case 0xFF47, 0xFF48, 0xFF49:
			if !m.gb.cartridge.IsCGB() {
				m.gb.video.WriteIO(addr, val)
			}
		case 0xFF68, 0xFF69, 0xFF6A, 0xFF6B:
			if m.gb.cartridge.IsCGB() {
				m.gb.video.WriteIO(addr, val)
			}
		case 0xFF46:
			m.gb.triggerGDMA(uint16(val) << 8)
		case 0xFF4D:
			m.gb.key1 = util.Bit(val, 0)
		case 0xFF70:
			m.wramBank = uint(val & 0b111)
			if m.wramBank == 0 {
				m.wramBank = 1
			}
		case 0xFF80, 0xFF81, 0xFF82, 0xFF83, 0xFF84, 0xFF85, 0xFF86, 0xFF87, 0xFF88, 0xFF89, 0xFF8A, 0xFF8B, 0xFF8C, 0xFF8D, 0xFF8E, 0xFF8F, 0xFF90, 0xFF91, 0xFF92, 0xFF93, 0xFF94, 0xFF95, 0xFF96, 0xFF97, 0xFF98, 0xFF99, 0xFF9A, 0xFF9B, 0xFF9C, 0xFF9D, 0xFF9E, 0xFF9F, 0xFFA0, 0xFFA1, 0xFFA2, 0xFFA3, 0xFFA4, 0xFFA5, 0xFFA6, 0xFFA7, 0xFFA8, 0xFFA9, 0xFFAA, 0xFFAB, 0xFFAC, 0xFFAD, 0xFFAE, 0xFFAF, 0xFFB0, 0xFFB1, 0xFFB2, 0xFFB3, 0xFFB4, 0xFFB5, 0xFFB6, 0xFFB7, 0xFFB8, 0xFFB9, 0xFFBA, 0xFFBB, 0xFFBC, 0xFFBD, 0xFFBE, 0xFFBF, 0xFFC0, 0xFFC1, 0xFFC2, 0xFFC3, 0xFFC4, 0xFFC5, 0xFFC6, 0xFFC7, 0xFFC8, 0xFFC9, 0xFFCA, 0xFFCB, 0xFFCC, 0xFFCD, 0xFFCE, 0xFFCF, 0xFFD0, 0xFFD1, 0xFFD2, 0xFFD3, 0xFFD4, 0xFFD5, 0xFFD6, 0xFFD7, 0xFFD8, 0xFFD9, 0xFFDA, 0xFFDB, 0xFFDC, 0xFFDD, 0xFFDE, 0xFFDF, 0xFFE0, 0xFFE1, 0xFFE2, 0xFFE3, 0xFFE4, 0xFFE5, 0xFFE6, 0xFFE7, 0xFFE8, 0xFFE9, 0xFFEA, 0xFFEB, 0xFFEC, 0xFFED, 0xFFEE, 0xFFEF, 0xFFF0, 0xFFF1, 0xFFF2, 0xFFF3, 0xFFF4, 0xFFF5, 0xFFF6, 0xFFF7, 0xFFF8, 0xFFF9, 0xFFFA, 0xFFFB, 0xFFFC, 0xFFFD, 0xFFFE:
			m.hram[addr&0x7F] = val
		case 0xFFFF:
			m.gb.ie = val
		}
	}
}