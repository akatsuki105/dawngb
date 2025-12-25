package cpu

func (c *CPU) ReadIO(addr uint16) uint8 {
	switch addr {
	case 0xFF00:
		return c.Joypad.read()
	case 0xFF01:
		return c.Serial.SB
	case 0xFF02:
		return c.Serial.SC
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		return c.Timer.Read(addr)
	case 0xFF0F:
		return c.IF & 0x1F
	case 0xFF4C:
		return c.Key0
	case 0xFF4D:
		if c.isCGB {
			key1 := c.Key1 | 0x7E
			if c.Clock == 4 { // 2x
				key1 |= 1 << 7
			}
			return key1
		}
	case 0xFF50:
		return 1
	case 0xFF55:
		if c.isCGB {
			return c.DMA.Read(addr)
		}
	case 0xFF72:
		return c.FF72
	case 0xFF73:
		return c.FF73
	case 0xFF74:
		return c.FF74
	}
	return 0
}

func (c *CPU) WriteIO(addr uint16, val uint8) {
	switch addr {
	case 0xFF00:
		c.Joypad.write(val)
	case 0xFF01:
		c.Serial.SB = val
	case 0xFF02:
		c.Serial.setSC(val)
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		c.Timer.Write(addr, val)
	case 0xFF0F:
		c.IF = val & 0x1F
	case 0xFF4C:
		if c.Key0 == 0 {
			c.Key0 = val
		}
	case 0xFF4D:
		if c.isCGB {
			c.Key1 = (c.Key1 & 0x80) | (val & 0x01)
		}
	case 0xFF50: // BANK
		c.BIOS.FF50 = false
	case 0xFF51, 0xFF52, 0xFF53, 0xFF54:
		if c.isCGB {
			c.DMA.Write(addr, val)
		}
	case 0xFF55:
		if c.isCGB {
			cycles := c.DMA.Write(addr, val)
			c.Cycles += cycles
		}
	case 0xFF72:
		c.FF72 = val
	case 0xFF73:
		c.FF73 = val
	case 0xFF74:
		c.FF74 = val
	}
}
