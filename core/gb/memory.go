package gb

func (g *GB) Read(addr uint16) uint8 {
	return g.read(addr, false)
}

func (g *GB) read(addr uint16, peek bool) uint8 {
	if peek && addr >= 0xFF80 && addr < 0xFFFF { // High RAM
		return g.CPU.HRAM[addr&0x7F]
	}

	switch {
	case addr < 0x8000: // ROM
		return g.Cart.Read(addr)
	case addr < 0xA000: // VRAM
		return g.PPU.Read(addr)
	case addr < 0xC000: // SRAM
		return g.Cart.Read(addr)
	case addr < 0xD000: // WRAM0
		return g.wram[addr&0xFFF]
	case addr < 0xE000: // WRAMn
		return g.wram[(uint(g.wramBank)<<12)|uint(addr&0xFFF)]
	case addr < 0xF000: // mirror of WRAM0
		return g.wram[addr&0xFFF]
	case addr < 0xFE00: // mirror of WRAMn
		return g.wram[(uint(g.wramBank)<<12)|uint(addr&0xFFF)]
	case addr < 0xFEA0: // OAM
		return g.PPU.Read(addr)
	case addr < 0xFF00: // unused
		return 0xFF
	}

	switch addr {
	case 0xFF00, 0xFF01, 0xFF02, 0xFF04, 0xFF05, 0xFF06, 0xFF07, 0xFF0F: // CPU
		return g.CPU.ReadIO(addr)
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26, 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F: // APU
		return g.APU.Read(addr, peek)
	case 0xFF40, 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45, 0xFF46, 0xFF47, 0xFF48, 0xFF49, 0xFF4A, 0xFF4B: // PPU
		return g.PPU.Read(addr)
	case 0xFF4C: // KEY0
		if g.IsColor() {
			return g.CPU.ReadIO(addr)
		}
		return 1
	case 0xFF4D, 0xFF50, 0xFF55, 0xFF72, 0xFF73, 0xFF74: // CPU(CGB only)
		if g.IsColor() {
			return g.CPU.ReadIO(addr)
		}
	case 0xFF4F, 0xFF68, 0xFF69, 0xFF6A, 0xFF6B: // PPU(CGB only)
		if g.IsColor() {
			return g.PPU.Read(addr)
		}
	case 0xFF56: // RP
		return 0x02 // TODO: infrared
	case 0xFF70:
		if g.IsColor() {
			return g.wramBank
		}
		return 1
	case 0xFFFF: // IE
		return g.CPU.IE
	}
	return 0
}

func (g *GB) Write(addr uint16, val uint8) {
	g.write(addr, val, false)
}

func (g *GB) write(addr uint16, val uint8, override bool) {
	switch {
	case addr < 0x8000: // ROM
		g.Cart.Write(addr, val)
		return
	case addr < 0xA000: // VRAM
		g.PPU.Write(addr, val)
		return
	case addr < 0xC000: // SRAM
		g.Cart.Write(addr, val)
		return
	case addr < 0xD000: // WRAM0
		g.wram[addr&0xFFF] = val
		return
	case addr < 0xE000: // WRAMn
		g.wram[(uint(g.wramBank)<<12)|uint(addr&0xFFF)] = val
		return
	case addr < 0xF000: // mirror of WRAM0
		g.wram[addr&0xFFF] = val
		return
	case addr < 0xFE00: // mirror of WRAMn
		g.wram[(uint(g.wramBank)<<12)|uint(addr&0xFFF)] = val
		return
	case addr < 0xFEA0: // OAM
		g.PPU.Write(addr, val)
		return
	}

	switch addr {
	case 0xFF00, 0xFF01, 0xFF02, 0xFF03, 0xFF04, 0xFF05, 0xFF06, 0xFF07, 0xFF0F: // CPU
		g.CPU.WriteIO(addr, val)
		return
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26, 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F: // APU
		g.APU.Write(addr, val)
		return
	case 0xFF40, 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45, 0xFF47, 0xFF48, 0xFF49, 0xFF4A, 0xFF4B: // PPU
		g.PPU.Write(addr, val)
		return
	case 0xFF46: // OAM DMA
		g.PPU.TriggerDMA(uint16(val)<<8, g.CPU.Clock)
		return
	case 0xFF4C, 0xFF4D, 0xFF50, 0xFF51, 0xFF52, 0xFF53, 0xFF54, 0xFF55, 0xFF72, 0xFF73, 0xFF74: // CPU(CGB only)
		if g.IsColor() {
			g.CPU.WriteIO(addr, val)
		}
		return
	case 0xFF4F, 0xFF68, 0xFF69, 0xFF6A, 0xFF6B: // PPU(CGB only)
		if g.IsColor() {
			g.PPU.Write(addr, val)
		}
		return
	case 0xFF70:
		if g.IsColor() {
			g.wramBank = (val & 0b111)
			if g.wramBank == 0 {
				g.wramBank = 1
			}
		} else {
			g.wramBank = 1
		}
		return
	case 0xFFFF:
		g.CPU.IE = val
		return
	}
}
