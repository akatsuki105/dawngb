package cpu

import "fmt"

type opcode = func(c *Cpu)

var opTable = [256]opcode{
	/* 0x00 */ op00, op01, op02, op03, op04, op05, op06, todo, op08, op09, op0A, op0B, op0C, op0D, op0E, todo,
	/* 0x10 */ todo, op11, op12, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x20 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x30 */ todo, op31, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x40 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x50 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x60 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x70 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x80 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0x90 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0xA0 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0xB0 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0xC0 */ todo, todo, todo, opC3, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0xD0 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0xE0 */ todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo,
	/* 0xF0 */ opF0, todo, todo, opF3, todo, todo, todo, todo, todo, todo, todo, todo, todo, todo, opFE, todo,
}

var opCycles = [256]int64{
	/* 0x00 */ 1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1,
	/* 0x10 */ 0, 3, 2, 2, 1, 1, 2, 1, 3, 2, 2, 2, 1, 1, 2, 1,
	/* 0x20 */ 2, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1,
	/* 0x30 */ 2, 3, 2, 2, 3, 3, 3, 1, 2, 2, 2, 2, 1, 1, 2, 1,
	/* 0x40 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0x50 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0x60 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0x70 */ 2, 2, 2, 2, 2, 2, 0, 2, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0x80 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0x90 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0xA0 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0xB0 */ 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	/* 0xC0 */ 2, 3, 3, 4, 3, 4, 2, 4, 2, 4, 3, 0, 3, 6, 2, 4,
	/* 0xD0 */ 2, 3, 3, 0, 3, 4, 2, 4, 2, 4, 3, 0, 3, 0, 2, 4,
	/* 0xE0 */ 3, 3, 2, 0, 0, 4, 2, 4, 4, 1, 4, 0, 0, 0, 2, 4,
	/* 0xF0 */ 3, 3, 2, 1, 0, 4, 2, 4, 3, 2, 4, 1, 0, 0, 2, 4,
}

func todo(c *Cpu) {
	panic(fmt.Sprintf("todo opcode: 0x%02X in 0x%04X", c.inst.opcode, c.inst.addr))
}

func op00(c *Cpu) { /* nop */ }

func op01(c *Cpu) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.bc.unpack((hi << 8) | lo)
}

func op02(c *Cpu) { c.m.Write(c.r.bc.pack(), c.r.a) }

func op03(c *Cpu) { c.r.bc.unpack(c.r.bc.pack() + 1) }

func op04(c *Cpu) {
	c.r.bc.hi++
	c.r.f.z = c.r.bc.hi == 0
	c.r.f.n = false
	c.r.f.h = c.r.bc.hi&0x0F == 0x00
}

func op05(c *Cpu) {
	c.r.bc.hi--
	c.r.f.z = c.r.bc.hi == 0
	c.r.f.n = true
	c.r.f.h = c.r.bc.hi&0x0F == 0x0F
}

func op06(c *Cpu) { c.r.bc.hi = c.fetch() }

func op08(c *Cpu) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.m.Write(addr, uint8(c.r.sp))
	c.m.Write(addr+1, uint8(c.r.sp>>8))
}

func op09(c *Cpu) {
	hl := c.r.hl.pack()
	bc := c.r.bc.pack()
	c.r.hl.unpack(hl + bc)
	c.r.f.n = false
	c.r.f.h = (hl&0x0FFF)+(bc&0x0FFF) > 0x0FFF
	c.r.f.c = hl+bc > 0xFFFF
}

func op0A(c *Cpu) { c.r.a = c.m.Read(c.r.bc.pack()) }

func op0B(c *Cpu) { c.r.bc.unpack(c.r.bc.pack() - 1) }

func op0C(c *Cpu) {
	c.r.bc.lo++
	c.r.f.z = c.r.bc.lo == 0
	c.r.f.n = false
	c.r.f.h = c.r.bc.lo&0x0F == 0x00
}

func op0D(c *Cpu) {
	c.r.bc.lo--
	c.r.f.z = c.r.bc.lo == 0
	c.r.f.n = true
	c.r.f.h = c.r.bc.lo&0x0F == 0x0F
}

func op0E(c *Cpu) { c.r.bc.lo = c.fetch() }

func op11(c *Cpu) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.de.unpack((hi << 8) | lo)
}

func op12(c *Cpu) { c.m.Write(c.r.de.pack(), c.r.a) }

func op31(c *Cpu) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.sp = (hi << 8) | lo
}

func opC3(c *Cpu) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.pc = (hi << 8) | lo
}

func opF0(c *Cpu) {
	lo := uint16(c.fetch())
	addr := 0xFF00 | lo
	c.r.a = c.m.Read(addr)
}

func opF3(c *Cpu) { c.IME = false }

func opFE(c *Cpu) {
	val := c.fetch()
	c.r.f.z = c.r.a == val
	c.r.f.n = true
	c.r.f.h = (c.r.a & 0x0F) < (val & 0x0F)
}
