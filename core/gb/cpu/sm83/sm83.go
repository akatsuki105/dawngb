package sm83

import (
	"fmt"
)

type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type SM83 struct {
	R    Registers
	bus  Bus
	inst struct {
		opcode uint8
		addr   uint16
		cb     bool
	}
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
	pc := c.R.PC
	c.R.PC++
	return c.bus.Read(pc)
}
