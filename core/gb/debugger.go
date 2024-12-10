package gb

func (g *GB) GetCPU(_ int) any {
	return g.cpu
}

func (g *GB) VideoUnit() any {
	return g.ppu
}

// GetValue returns the value of the specified state.
func (g *GB) GetValue(which uint64) uint64 {
	category := which >> 56
	switch category {
	case 0: // CPU
		target := which & 0xFF
		switch target {
		case 0: // Register (AF)
			val := uint64(g.cpu.R.A)
			val |= uint64(g.cpu.R.F.Pack()) << 8
			return val
		case 1: // BC
			return uint64(g.cpu.R.BC.Pack())
		case 2: // DE
			return uint64(g.cpu.R.DE.Pack())
		case 3: // HL
			return uint64(g.cpu.R.HL.Pack())
		case 4: // SP
			return uint64(g.cpu.R.SP)
		case 5: // PC
			return uint64(g.cpu.R.PC)
		}
	case 1: // PPU
		subcategory := (which >> 48) & 0xFF
		switch subcategory {
		case 0: // LCDSTAT IRQ info
			target := which & 0xFF
			switch target {
			case 0: // IRQ is triggered
				if !g.ppu.StatIRQ.Triggered {
					return 0
				}
				return 1
			case 1: // IRQ mode
				return uint64(g.ppu.StatIRQ.Mode)
			case 2: // IRQ line X
				return uint64(g.ppu.StatIRQ.Lx)
			case 3: // IRQ line Y
				return uint64(g.ppu.StatIRQ.Ly)
			}
		}
	}
	return 0
}

func (g *GB) ViewMemory(memID int, addr uint32, width int, flag uint32) uint32 {
	addr16 := uint16(addr)
	switch width {
	case 1:
		return uint32(g.read(addr16, true))
	case 2:
		return uint32(g.read(addr16, true)) | (uint32(g.read(addr16+1, true)) << 8)
	case 4:
		return uint32(g.read(addr16, true)) | (uint32(g.read(addr16+1, true)) << 8) | (uint32(g.read(addr16+2, true)) << 16) | (uint32(g.read(addr16+3, true)) << 24)
	default:
		return 0
	}
}

func (g *GB) GetChunk(chunkID uint8) []uint8 {
	return nil
}
