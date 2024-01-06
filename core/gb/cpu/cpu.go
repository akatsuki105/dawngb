package cpu

import (
	"fmt"

	"github.com/akatsuki105/dugb/util/sched"
)

type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type Cpu struct {
	r       Registers
	m       Memory
	Blocked bool
	inst    struct {
		opcode uint8
		addr   uint16
		cb     bool
	}
	s          *sched.Sched
	IME        bool
	halt, stop func()
	Cycle      int64 // 8 or 4
}

func New(s *sched.Sched, m Memory, halt, stop func()) *Cpu {
	return &Cpu{
		s:     s,
		m:     m,
		halt:  halt,
		stop:  stop,
		Cycle: 8,
	}
}

func (c *Cpu) Reset(hasBIOS bool) {
	if !hasBIOS {
		c.skipBIOS()
	}
}

func (c *Cpu) skipBIOS() {
	c.r.a = 0x11
	c.r.f.unpack(0x80)
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

	c.s.Add(opCycles[opcode] * c.Cycle)
}

func (c *Cpu) fetch() uint8 {
	pc := c.r.pc
	c.r.pc++
	return c.m.Read(pc)
}
