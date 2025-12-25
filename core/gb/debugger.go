package gb

import (
	"github.com/akatsuki105/dawngb/internal/debugger"
	"github.com/akatsuki105/dawngb/internal/unsafeslice"
)

func (g *GB) AttachDebugger(d debugger.Debugger) {
	g.Debugger = d
	g.CPU.Debugger = d
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
	case 2: // Cartridge
		target := which & 0xFF
		switch target {
		case 0: // ROM Bank number
			return uint64(g.Cart.ROMBankNumber())
		}
	case 0xFF: // Misc
		target := which & 0xFF
		switch target {
		case 0: // SRAM size
			return uint64(g.Cart.RAMSize())
		}
	}
	return 0
}

func (g *GB) ViewMemory(_ int, addr uint32, width int) uint64 {
	addr16 := uint16(addr)
	switch width {
	case 1:
		return uint64(g.read(addr16, true))
	case 2:
		return uint64(g.read(addr16, true)) | (uint64(g.read(addr16+1, true)) << 8)
	case 4:
		return uint64(g.read(addr16, true)) | (uint64(g.read(addr16+1, true)) << 8) | (uint64(g.read(addr16+2, true)) << 16) | (uint64(g.read(addr16+3, true)) << 24)
	default:
		return 0
	}
}

func (g *GB) PokeMemory(memID int, addr uint32, width int, data uint32) bool {
	switch width {
	case 2:
		ok := g.PokeMemory(memID, addr, 1, uint32(data&0xFF))
		ok = g.PokeMemory(memID, addr+1, 1, uint32(data>>8)) && ok
		return ok
	case 4:
		ok := g.PokeMemory(memID, addr, 2, uint32(data&0xFFFF))
		ok = g.PokeMemory(memID, addr+2, 2, uint32(data>>16)) && ok
		return ok
	}

	addr16 := uint16(addr)
	switch {
	case addr16 < 0x8000: // ROM
		return false
	case addr16 < 0xA000: // VRAM
		return false
	case addr16 < 0xC000: // SRAM
		return false
	case addr16 < 0xD000: // WRAM0
		g.WRAM.Data[addr16&0xFFF] = uint8(data)
		return true
	case addr16 < 0xE000: // WRAMn
		g.WRAM.Data[(uint(g.WRAM.Bank)<<12)|uint(addr16&0xFFF)] = uint8(data)
		return true
	case addr16 < 0xF000: // mirror of WRAM0
		g.WRAM.Data[addr16&0xFFF] = uint8(data)
		return true
	case addr16 < 0xFE00: // mirror of WRAMn
		g.WRAM.Data[(uint(g.WRAM.Bank)<<12)|uint(addr16&0xFFF)] = uint8(data)
		return true
	case addr16 < 0xFEA0: // OAM
		return false
	case addr16 < 0xFF00: // unused
		return false
	case addr16 < 0xFF80: // I/O
		return false
	case addr16 < 0xFFFF: // HRAM
		g.CPU.HRAM[addr16&0x7F] = uint8(data)
		return true
	case addr16 == 0xFFFF: // IE
		g.CPU.IE = uint8(data)
	}
	return false
}

func (g *GB) GetChunk(chunkID uint64) []uint8 {
	memory := (chunkID >> 56) & 0xFF
	switch memory {
	case 0: // CPUアドレス空間全体 (0x0000..FFFF)
		return nil
	case 1: // ROM(バンクも含む全部)
		return g.Cart.ROM[:]
	case 2: // VRAM(バンクも含む全部)
		return g.PPU.RAM.Data[:]
	case 3: // WRAM(バンクも含む全部)
		return g.WRAM.Data[:]
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
	case 6: // HRAM
		return g.CPU.HRAM[:]
	}
	return nil
}
