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
	return 0
}

func (g *GB) GetChunk(chunkID uint8) []uint8 {
	return nil
}
