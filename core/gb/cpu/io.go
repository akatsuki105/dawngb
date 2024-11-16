package cpu

import "github.com/akatsuki105/dawngb/util"

func (c *CPU) ReadIO(addr uint16) uint8 {
	switch addr {
	case 0xFF00:
		return c.joypad.read()
	case 0xFF01:
		return c.serial.sb
	case 0xFF02:
		return c.serial.sc
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		return c.timer.Read(addr)
	case 0xFF0F:
		val := uint8(0)
		for i := 0; i < 5; i++ {
			val |= (util.Btou8(c.interrupt[i]) << i)
		}
		return val
	case 0xFF4D:
		key1 := c.key1 | 0x7E
		if c.Clock == 4 { // 2x
			key1 |= 1 << 7
		}
		return key1
	case 0xFF51, 0xFF52, 0xFF53, 0xFF54, 0xFF55:
		return c.dma.Read(addr)
	case 0xFF72:
		return c.ff72
	case 0xFF73:
		return c.ff73
	case 0xFF74:
		return c.ff74
	}
	return 0
}

func (c *CPU) WriteIO(addr uint16, val uint8) {
	switch addr {
	case 0xFF00:
		c.joypad.write(val)
	case 0xFF01:
		c.serial.sb = val
	case 0xFF02:
		c.serial.setSC(val)
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		c.timer.Write(addr, val)
	case 0xFF0F:
		for i := 0; i < 5; i++ {
			c.interrupt[i] = (val & (1 << i)) != 0
		}
	case 0xFF4D:
		c.key1 = (c.key1 & 0x80) | (val & 0x01)
	case 0xFF50: // BANK
		c.bios.ff50 = false
	case 0xFF51, 0xFF52, 0xFF53, 0xFF54:
		c.dma.Write(addr, val)
	case 0xFF55:
		cycles := c.dma.Write(addr, val)
		c.Cycles += cycles
	case 0xFF72:
		c.ff72 = val
	case 0xFF73:
		c.ff73 = val
	case 0xFF74:
		c.ff74 = val
	}
}
