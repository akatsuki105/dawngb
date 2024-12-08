package gb

func (g *GB) GetCPU(_ int) any {
	return g.cpu
}

func (g *GB) VideoUnit() any {
	return g.ppu
}

func (g *GB) GetValue(which uint64) uint64 {
	return 0
}

func (g *GB) ViewMemory(memID int, addr uint32, width int, flag uint32) uint32 {
	addr16 := uint16(addr)
	switch width {
	case 1:
		return uint32(g.cpu.Read(addr16))
	case 2:
		return uint32(g.cpu.Read(addr16)) | (uint32(g.cpu.Read(addr16+1)) << 8)
	case 4:
		return uint32(g.cpu.Read(addr16)) | (uint32(g.cpu.Read(addr16+1)) << 8) | (uint32(g.cpu.Read(addr16+2)) << 16) | (uint32(g.cpu.Read(addr16+3)) << 24)
	default:
		return 0
	}
}

func (g *GB) GetChunk(chunkID uint8) []uint8 {
	return nil
}
