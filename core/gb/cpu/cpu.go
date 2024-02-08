package cpu

import (
	"fmt"
)

type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type Cpu struct {
	r    Registers
	m    Memory
	inst struct {
		opcode uint8
		addr   uint16
		cb     bool
	}
	IME        bool
	halt, stop func()
	_tick      func(mastercycles int64)
	Cycle      int64 // 8 or 4
}

func New(m Memory, halt, stop func(), tick func(int64)) *Cpu {
	return &Cpu{
		m:     m,
		halt:  halt,
		stop:  stop,
		_tick: tick,
	}
}

func (c *Cpu) Reset(hasBIOS bool) {
	c.r = Registers{}
	c.IME = false
	c.Cycle = 8
	if !hasBIOS {
		c.skipBIOS()
	}
}

func (c *Cpu) skipBIOS() {
	c.r.a = 0x11
	c.r.f.unpack(0x80)
	c.r.bc.unpack(0x0000)
	c.r.de.unpack(0xFF56)
	c.r.hl.unpack(0x000D)

	c.r.sp = 0xFFFE
	c.r.pc = 0x100
}

func (c *Cpu) Step() {
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

	c.tick(opCycles[opcode] * c.Cycle)
}

func (c *Cpu) fetch() uint8 {
	pc := c.r.pc
	c.r.pc++
	return c.m.Read(pc)
}

func (c *Cpu) tick(mastercycles int64) {
	if c._tick != nil {
		c._tick(mastercycles)
	}
}
