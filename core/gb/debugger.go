package gb

import "github.com/akatsuki105/dawngb/internal/unsafeslice"

func (g *GB) VideoUnit() any {
	return g.PPU
}

// GetValue returns the value of the specified state.
func (g *GB) GetValue(which uint64) uint64 {
	category := which >> 56
	switch category {
	case 0: // CPU
		subcategory := (which >> 48) & 0xFF
		switch subcategory {
		case 0: // Registers
			target := which & 0xFF
			switch target {
			case 0: // Register (AF)
				val := uint64(g.CPU.R.A)
				val |= uint64(g.CPU.R.F.Pack()) << 8
				return val
			case 1: // BC
				return uint64(g.CPU.R.BC.Pack())
			case 2: // DE
				return uint64(g.CPU.R.DE.Pack())
			case 3: // HL
				return uint64(g.CPU.R.HL.Pack())
			case 4: // SP
				return uint64(g.CPU.R.SP)
			case 5: // PC
				return uint64(g.CPU.R.PC)
			}
		case 1: // Usage
			return uint64(g.CPU.Usage)
		}
	case 1: // PPU
		subcategory := (which >> 48) & 0xFF
		switch subcategory {
		case 0: // LCDSTAT IRQ info
			target := which & 0xFF
			switch target {
			case 0: // IRQ is triggered
				if !g.PPU.StatIRQ.Triggered {
					return 0
				}
				return 1
			case 1: // IRQ mode
				return uint64(g.PPU.StatIRQ.Mode)
			case 2: // IRQ line X
				return uint64(g.PPU.StatIRQ.Lx)
			case 3: // IRQ line Y
				return uint64(g.PPU.StatIRQ.Ly)
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

func (g *GB) GetChunk(chunkID uint64) []uint8 {
	memory := (chunkID >> 56) & 0xFF
	switch memory {
	case 0: // CPUアドレス空間全体 (0x0000..FFFF)
		return nil
	case 1: // ROM(バンクも含む全部)
		return g.cartridge.ROM[:]
	case 2: // VRAM(バンクも含む全部)
		return g.PPU.RAM.Data[:]
	case 3: // WRAM(バンクも含む全部)
		return g.wram[:]
	case 4: // Palette
		target := chunkID & 0xFF
		if target < 16 {
			return unsafeslice.ByteSliceFromUint16Slice(g.PPU.Palette[target*4 : (target+1)*4])
		}
		switch target {
		case 0xFF:
			return unsafeslice.ByteSliceFromUint16Slice(g.PPU.Palette[:])
		}
	case 5: // OAM
		return g.PPU.OAM[:]
	}
	return nil
}