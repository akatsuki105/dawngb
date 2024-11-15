package sm83

import (
	"fmt"
)

type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type SM83 struct {
	r    Registers
	bus  Bus
	inst struct {
		opcode uint8
		addr   uint16
		cb     bool
	}
	IME        bool
	halt, stop func()
	_tick      func(mastercycles int64)
}

func New(bus Bus, halt, stop func(), tick func(int64)) *SM83 {
	return &SM83{
		bus:   bus,
		halt:  halt,
		stop:  stop,
		_tick: tick,
	}
}

func (c *SM83) Reset(hasBIOS bool) {
	c.r = Registers{}
	c.IME = false
	if !hasBIOS {
		c.skipBIOS()
	}
}

func (c *SM83) skipBIOS() {
	c.r.a = 0x11
	c.r.f.unpack(0x80)
	c.r.bc.unpack(0x0000)
	c.r.de.unpack(0xFF56)
	c.r.hl.unpack(0x000D)

	c.r.sp = 0xFFFE
	c.r.pc = 0x100
}

func (c *SM83) Step() {
	pc := c.r.pc
	c.inst.addr = pc
	opcode := c.fetch()
	c.inst.opcode = opcode
	c.inst.cb = false

	fn := opTable[opcode]
	if fn != nil {
		// fmt.Printf("0x%02X in 0x%04X\n", opcode, pc)
		fn(c)
	} else {
		panic(fmt.Sprintf("illegal opcode: 0x%02X in 0x%04X", opcode, pc))
	}

	c.tick(opCycles[opcode])
}

func (c *SM83) fetch() uint8 {
	pc := c.r.pc
	c.r.pc++
	return c.bus.Read(pc)
}

func (c *SM83) tick(mastercycles int64) {
	if c._tick != nil {
		c._tick(mastercycles)
	}
}
