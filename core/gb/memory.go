package gb

import "github.com/akatsuki105/dugb/util"

type Memory struct {
	gb *GB
}

func newMemory(gb *GB) *Memory {
	return &Memory{
		gb: gb,
	}
}

func (m *Memory) Read(addr uint16) byte {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.gb.cartridge.rom[addr]
	case 0x4, 0x5, 0x6, 0x7:
		return m.gb.cartridge.rom[addr]
	case 0xF:
		switch addr {
		case 0xFF44:
			return m.gb.video.Scanline()
		case 0xFFFF:
			return byte(util.Btoi(m.gb.cpu.IME))
		}
	}
	return 0
}

func (m *Memory) Write(addr uint16, val byte) {}
