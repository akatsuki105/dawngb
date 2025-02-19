package sm83

import (
	"fmt"
)

type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type Context struct {
	Opcode uint8
	Addr   uint16
	CB     bool
}

type SM83 struct {
	R          Registers
	bus        Bus
	inst       Context
	IME        bool
	halt, stop func()
	tick       func(clockCycles int64)
}

func New(bus Bus, halt, stop func(), tick func(int64)) *SM83 {
	if tick == nil {
		panic("tick function is required")
	}
	return &SM83{
		bus:  bus,
		halt: halt,
		stop: stop,
		tick: tick,
	}
}

func (c *SM83) Reset() {
	c.R.reset()
	c.IME = false
}

func (c *SM83) Step() {
	pc := c.R.PC
	c.inst.Addr = pc
	opcode := c.fetch()
	c.inst.Opcode = opcode
	c.inst.CB = false

	fn := opTable[opcode]
	if fn != nil {
		fn(c)
	} else {
		panic(fmt.Sprintf("illegal opcode: 0x%02X in 0x%04X", opcode, pc))
	}

	c.tick(opCycles[opcode])
}

func (c *SM83) fetch() uint8 {
	pc := c.R.PC
	c.R.PC++
	return c.bus.Read(pc)
}

func btou8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
