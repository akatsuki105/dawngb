package cpu

import (
	"errors"

	"github.com/akatsuki105/dawngb/core/gb/cpu/sm83"
	"github.com/akatsuki105/dawngb/internal/debugger"
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
	isCGB  bool  // ハードがCGBかどうか
	Cycles int64 // 8MHzのマスターサイクル単位
	*sm83.SM83
	bus              sm83.Bus
	Clock            int64 // 8(1x) or 4(2x)
	timer            *timer
	DMA              *DMA
	joypad           *joypad
	serial           *serial
	bios             BIOS
	HRAM             [0x7F]uint8
	halted           bool
	IE, IF           uint8
	key0             uint8 // FF4C
	key1             uint8 // FF4D
	ff72, ff73, ff74 uint8
	Usage            uint32 // フレーム中のCPU使用率(haltしてないときのみカウントしたサイクル数)
	Debugger         debugger.Debugger
}

// a.k.a. Boot ROM
type BIOS struct {
	ff50 bool
	data []uint8
}

func New(isCGB bool, bus sm83.Bus) *CPU {
	c := &CPU{
		isCGB: isCGB,
		bus:   bus,
	}
	c.SM83 = sm83.New(c, c.halt, c.stop, c.wait)
	c.timer = newTimer(c.IRQ, &c.Clock)
	c.joypad = newJoypad(c.IRQ)
	c.DMA = newDMA(c)
	c.serial = newSerial(c.IRQ)
	return c
}

func (c *CPU) Reset() {
	c.Cycles = 0
	c.SM83.Reset()
	c.Clock = 8
	c.timer.reset()
	c.DMA.reset()
	c.joypad.reset()
	c.serial.reset()
	clear(c.HRAM[:])
	c.halted = false
	c.IE, c.IF = 0, 0
	c.bios.ff50 = true
	c.key0, c.key1 = 0, 0
	c.ff72, c.ff73, c.ff74 = 0, 0, 0
}

func (c *CPU) SkipBIOS() {
	c.bios.ff50 = false
	c.timer.tac = 0xF8
	c.joypad.write(0x30)
	c.joypad.write(0xCF)
	c.DMA.skipBIOS()

	if c.isCGB {
		c.R.A = 0x11
		c.R.F.Unpack(0x80)
		c.R.BC.Unpack(0x0000)
		c.R.DE.Unpack(0xFF56)
		c.R.HL.Unpack(0x000D)
	} else {
		c.R.A = 0x01
		c.R.F.Unpack(0x80)
		c.R.BC.Unpack(0x0013)
		c.R.DE.Unpack(0x00D8)
		c.R.HL.Unpack(0x014D)
	}
	c.R.SP, c.R.PC = 0xFFFE, 0x0100
}

func (c *CPU) wait(n int64) {
	c.Cycles += n * c.Clock
	c.Usage += uint32(n * c.Clock)
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

func (c *CPU) HBlank() {
	if c.isCGB {
		c.DMA.startHDMA()
	}
}

func (c *CPU) Step() int64 {
	cycles := c.step()
	c.timer.run(cycles)
	c.serial.run(cycles)
	return cycles
}

func (c *CPU) step() int64 {
	prev := c.Cycles
	if c.DMA.doHDMA {
		c.DMA.doHDMA = false
		c.DMA.runHDMA()
		c.Cycles += 64
		return c.Cycles - prev
	}

	irqID := c.checkInterrupt()
	if irqID >= 0 {
		c.halted = false
		if c.IME {
			c.IF &^= uint8(1 << irqID)
			c.Interrupt(irqID)
		} else {
			c.instruction()
		}
	} else if c.halted {
		c.Cycles++
	} else {
		c.instruction()
	}

	return c.Cycles - prev
}

func (c *CPU) instruction() {
	if c.Debugger != nil {
		c.Debugger.InstructionHook(0, uint64(c.R.PC))
	}
	c.SM83.Step()
}

// IRQ id: 0: VBLANK, 1: LCDSTAT, 2: TIMER, 3: SERIAL, 4: JOYPAD
func (c *CPU) IRQ(id int) { c.IF |= (1 << id) }

func (c *CPU) SendInputs(inputs uint8) {
	c.joypad.inputs = inputs
}

func (c *CPU) checkInterrupt() int {
	irq := c.IE & c.IF
	for i := 0; i < 5; i++ {
		if (irq & (1 << i)) != 0 {
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

func (c *CPU) IsCGBMode() bool {
	return c.isCGB && c.key0 != 4
}
