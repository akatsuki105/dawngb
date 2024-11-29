package gb

func (g *GB) Read(addr uint16) uint8 {
	switch {
	case addr < 0x8000: // ROM
		return g.cartridge.Read(addr)
	case addr < 0xA000: // VRAM
		return g.ppu.Read(addr)
	case addr < 0xC000: // SRAM
		return g.cartridge.Read(addr)
	case addr < 0xD000: // WRAM0
		return g.wram[addr&0xFFF]
	case addr < 0xE000: // WRAMn
		return g.wram[(uint(g.wramBank)<<12)|uint(addr&0xFFF)]
	case addr < 0xF000: // mirror of WRAM0
		return g.wram[addr&0xFFF]
	case addr < 0xFE00: // mirror of WRAMn
		return g.wram[(uint(g.wramBank)<<12)|uint(addr&0xFFF)]
	case addr < 0xFEA0: // OAM
		return g.ppu.Read(addr)
	case addr < 0xFF00: // unused
		return 0xFF
	}

	switch addr {
	case 0xFF00, 0xFF01, 0xFF02, 0xFF04, 0xFF05, 0xFF06, 0xFF07, 0xFF0F: // CPU
		return g.cpu.ReadIO(addr)
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26, 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F: // APU
		return g.apu.Read(addr)
	case 0xFF40, 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45, 0xFF46, 0xFF47, 0xFF48, 0xFF49, 0xFF4A, 0xFF4B: // PPU
		return g.ppu.Read(addr)
	case 0xFF4C: // KEY0
		if g.IsCGB() {
			return g.cpu.ReadIO(addr)
		}
		return 1
	case 0xFF4D, 0xFF50, 0xFF55, 0xFF72, 0xFF73, 0xFF74: // CPU(CGB only)
		if g.IsCGB() {
			return g.cpu.ReadIO(addr)
		}
	case 0xFF4F, 0xFF68, 0xFF69, 0xFF6A, 0xFF6B: // PPU(CGB only)
		if g.IsCGB() {
			return g.ppu.Read(addr)
		}
	case 0xFF56: // RP
		return 0x02 // TODO: infrared
	case 0xFF70:
		if g.IsCGB() {
			return g.wramBank
		}
		return 1
	case 0xFFFF: // IE
		return g.cpu.IE
	}

	return 0
}

func (g *GB) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x8000: // ROM
		g.cartridge.Write(addr, val)
		return
	case addr < 0xA000: // VRAM
		g.ppu.Write(addr, val)
		return
	case addr < 0xC000: // SRAM
		g.cartridge.Write(addr, val)
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
		g.ppu.Write(addr, val)
		return
	}

	switch addr {
	case 0xFF00, 0xFF01, 0xFF02, 0xFF03, 0xFF04, 0xFF05, 0xFF06, 0xFF07, 0xFF0F: // CPU
		g.cpu.WriteIO(addr, val)
		return
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26, 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F: // APU
		g.apu.Write(addr, val)
		return
	case 0xFF40, 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45, 0xFF47, 0xFF48, 0xFF49, 0xFF4A, 0xFF4B: // PPU
		g.ppu.Write(addr, val)
		return
	case 0xFF46: // OAM DMA
		g.ppu.TriggerDMA(uint16(val)<<8, g.cpu.Clock)
		return
	case 0xFF4C, 0xFF4D, 0xFF50, 0xFF51, 0xFF52, 0xFF53, 0xFF54, 0xFF55, 0xFF72, 0xFF73, 0xFF74: // CPU(CGB only)
		if g.IsCGB() {
			g.cpu.WriteIO(addr, val)
		}
		return
	case 0xFF4F, 0xFF68, 0xFF69, 0xFF6A, 0xFF6B: // PPU(CGB only)
		if g.IsCGB() {
			g.ppu.Write(addr, val)
		}
		return
	case 0xFF70:
		if g.IsCGB() {
			g.wramBank = (val & 0b111)
			if g.wramBank == 0 {
				g.wramBank = 1
			}
		} else {
			g.wramBank = 1
		}
		return
	case 0xFFFF:
		g.cpu.IE = val
		return
	}
}
