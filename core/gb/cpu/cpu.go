package cpu

import (
	"errors"

	"github.com/akatsuki105/dawngb/core/gb/cpu/sm83"
	"github.com/akatsuki105/dawngb/util"
)

const (
	IRQ_VBLANK = iota
	IRQ_LCDSTAT
	IRQ_TIMER
	IRQ_SERIAL
	IRQ_JOYPAD
)

// DMG-CPU, CGB-CPU
type CPU struct {
	Cycles int64 // 8MHzのマスターサイクル単位
	*sm83.SM83
	bus              sm83.Bus
	Clock            int64 // 8(1x) or 4(2x)
	dma              *DMA
	joypad           *joypad
	bios             BIOS
	HRAM             [0x7F]uint8
	halted           bool
	IE               uint8
	interrupt        [5]bool // IF
	key1             uint8   // FF4D
	ff72, ff73, ff74 uint8
}

// a.k.a. Boot ROM
type BIOS struct {
	ff50 bool
	data []uint8
}

func New(bus sm83.Bus) *CPU {
	c := &CPU{
		bus: bus,
	}
	c.SM83 = sm83.New(c, c.halt, c.stop, c.wait)
	c.joypad = newJoypad(c.IRQ)
	c.dma = newDMA(c)
	return c
}

func (c *CPU) Reset(hasBIOS bool) {
	c.Cycles = 0
	c.SM83.Reset(hasBIOS)
	c.Clock = 8
	c.dma.Reset(hasBIOS)
	c.joypad.reset(hasBIOS)
	clear(c.HRAM[:])
	c.halted = false
	c.IE, c.interrupt = 0, [5]bool{}
	c.key1 = 0
	c.ff72, c.ff73, c.ff74 = 0, 0, 0
}

func (c *CPU) wait(n int64) {
	c.Cycles += n * c.Clock
}

func (c *CPU) LoadBIOS(bios []uint8) error {
	c.bios.ff50 = false

	switch len(bios) {
	case 256: // DMG, MGB, SGB
		c.bios.data = make([]uint8, 256)
		copy(c.bios.data[:], bios)
	case 2048: // CGB, AGB
		c.bios.data = make([]uint8, 2048)
		copy(c.bios.data[:], bios)
	case 2048 + 256: // CGB, AGB (0x100..200 is padded)
		c.bios.data = make([]uint8, 2048)
		copy(c.bios.data[:256], bios[:256]) // 0x000..100
		copy(c.bios.data[256:], bios[512:]) // 0x200..900
	default:
		return errors.New("invalid BIOS size")
	}

	c.bios.ff50 = true
	return nil
}

func (c *CPU) stop() {
	if c.key1&(1<<0) != 0 {
		if c.Clock == 4 {
			c.Clock = 8
		} else {
			c.Clock = 4
		}
		c.key1 &^= 1 << 0
	}
}

func (c *CPU) StartHDMA() {
	c.dma.startHDMA()
}

func (c *CPU) Step() int64 {
	prev := c.Cycles
	if c.dma.doHDMA {
		c.dma.doHDMA = false
		c.dma.runHDMA()
		c.Cycles += 64
		return c.Cycles - prev
	}

	irqID := c.checkInterrupt()
	if irqID >= 0 {
		c.halted = false
		if c.IME {
			c.interrupt[irqID] = false
			c.Interrupt(irqID)
		} else {
			c.SM83.Step()
		}
	} else if c.halted {
		c.Cycles++
	} else {
		c.SM83.Step()
	}

	return c.Cycles - prev
}

func (c *CPU) ReadIO(addr uint16) uint8 {
	switch addr {
	case 0xFF00:
		return c.joypad.read()
	case 0xFF0F:
		val := uint8(0)
		for i := 0; i < 5; i++ {
			val |= (uint8(util.Btoi(c.interrupt[i])) << i)
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
	case 0xFF0F:
		for i := 0; i < 5; i++ {
			c.interrupt[i] = util.Bit(val, i)
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

func (c *CPU) IRQ(id int) { c.interrupt[id] = true }

func (c *CPU) SendInputs(inputs uint8) {
	c.joypad.inputs = inputs
}

func (c *CPU) checkInterrupt() int {
	for i := 0; i < 5; i++ {
		if util.Bit(c.IE, i) && c.interrupt[i] {
			return i
		}
	}
	return -1
}

func (c *CPU) halt() {
	if c.IME {
		c.halted = true
	} else {
		if c.checkInterrupt() < 0 {
			c.halted = true
		}
	}
}
