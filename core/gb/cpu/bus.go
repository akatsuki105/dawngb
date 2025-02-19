package cpu

func (c *CPU) Read(addr uint16) uint8 {
	if c.Debugger != nil {
		c.Debugger.ReadMemoryHook(0, uint64(addr), 1)
	}

	if c.bios.ff50 {
		if addr < 0x100 {
			return c.bios.data[addr]
		}
		if len(c.bios.data) == 2048 && (addr >= 0x200 && addr < 0x900) {
			return c.bios.data[addr-0x100]
		}
	}
	if addr >= 0xFF80 && addr <= 0xFFFE { // HRAM
		return c.HRAM[addr&0x7F]
	}
	return c.bus.Read(addr)
}

func (c *CPU) Write(addr uint16, val uint8) {
	if c.Debugger != nil {
		c.Debugger.WriteMemoryHook(0, uint64(addr), 1, uint64(val))
	}

	if addr >= 0xFF80 && addr <= 0xFFFE { // HRAM
		c.HRAM[addr&0x7F] = val
		return
	}
	c.bus.Write(addr, val)
}
