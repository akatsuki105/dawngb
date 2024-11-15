package sm83

import (
	"github.com/akatsuki105/dawngb/util"
)

type opcode = func(c *SM83)

// GameBoy Opcode Table
var opTable = [256]opcode{
	/* 0x00 */ op00, op01, op02, op03, op04, op05, op06, op07, op08, op09, op0A, op0B, op0C, op0D, op0E, op0F,
	/* 0x10 */ op10, op11, op12, op13, op14, op15, op16, op17, op18, op19, op1A, op1B, op1C, op1D, op1E, op1F,
	/* 0x20 */ op20, op21, op22, op23, op24, op25, op26, op27, op28, op29, op2A, op2B, op2C, op2D, op2E, op2F,
	/* 0x30 */ op30, op31, op32, op33, op34, op35, op36, op37, op38, op39, op3A, op3B, op3C, op3D, op3E, op3F,
	/* 0x40 */ op40, op41, op42, op43, op44, op45, op46, op47, op48, op49, op4A, op4B, op4C, op4D, op4E, op4F,
	/* 0x50 */ op50, op51, op52, op53, op54, op55, op56, op57, op58, op59, op5A, op5B, op5C, op5D, op5E, op5F,
	/* 0x60 */ op60, op61, op62, op63, op64, op65, op66, op67, op68, op69, op6A, op6B, op6C, op6D, op6E, op6F,
	/* 0x70 */ op70, op71, op72, op73, op74, op75, op76, op77, op78, op79, op7A, op7B, op7C, op7D, op7E, op7F,
	/* 0x80 */ op80, op81, op82, op83, op84, op85, op86, op87, op88, op89, op8A, op8B, op8C, op8D, op8E, op8F,
	/* 0x90 */ op90, op91, op92, op93, op94, op95, op96, op97, op98, op99, op9A, op9B, op9C, op9D, op9E, op9F,
	/* 0xA0 */ opA0, opA1, opA2, opA3, opA4, opA5, opA6, opA7, opA8, opA9, opAA, opAB, opAC, opAD, opAE, opAF,
	/* 0xB0 */ opB0, opB1, opB2, opB3, opB4, opB5, opB6, opB7, opB8, opB9, opBA, opBB, opBC, opBD, opBE, opBF,
	/* 0xC0 */ opC0, opC1, opC2, opC3, opC4, opC5, opC6, opC7, opC8, opC9, opCA, opCB, opCC, opCD, opCE, opCF,
	/* 0xD0 */ opD0, opD1, opD2, todo, opD4, opD5, opD6, opD7, opD8, opD9, opDA, todo, opDC, todo, opDE, opDF,
	/* 0xE0 */ opE0, opE1, opE2, todo, todo, opE5, opE6, opE7, opE8, opE9, opEA, todo, todo, todo, opEE, opEF,
	/* 0xF0 */ opF0, opF1, opF2, opF3, todo, opF5, opF6, opF7, opF8, opF9, opFA, opFB, todo, todo, opFE, opFF,
}

var opCycles = [256]int64{
	/* 0x00 */ 1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1,
	/* 0x10 */ 1, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1,
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
	/* 0xC0 */ 2, 1, 3, 3, 3, 2, 2, 1, 2, 1, 3, 0, 3, 3, 2, 1,
	/* 0xD0 */ 2, 1, 3, 0, 3, 2, 2, 1, 2, 1, 3, 0, 3, 0, 2, 1,
	/* 0xE0 */ 3, 1, 2, 0, 0, 2, 2, 1, 4, 0, 4, 0, 0, 0, 2, 1,
	/* 0xF0 */ 3, 1, 2, 1, 0, 2, 2, 1, 3, 2, 4, 1, 0, 0, 2, 1,
}

func op00(c *SM83) { /* nop */ }

func op01(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.bc.unpack((hi << 8) | lo)
}

func op02(c *SM83) { c.bus.Write(c.r.bc.pack(), c.r.a) }

func op03(c *SM83) { c.r.bc.unpack(c.r.bc.pack() + 1) }

func op04(c *SM83) {
	c.r.bc.hi++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.bc.hi == 0), false, (c.r.bc.hi&0x0F == 0x00)
}

func op05(c *SM83) {
	c.r.bc.hi--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.bc.hi == 0), true, (c.r.bc.hi&0x0F == 0x0F)
}

func op06(c *SM83) { c.r.bc.hi = c.fetch() }

// rlca
func op07(c *SM83) {
	msb := util.Bit(c.r.a, 7)
	c.r.f.c = msb
	c.r.a = (c.r.a << 1) | (util.Btou8(msb))
	c.r.f.z, c.r.f.n, c.r.f.h = false, false, false
}

func op08(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.bus.Write(addr, uint8(c.r.sp))
	c.bus.Write(addr+1, uint8(c.r.sp>>8))
}

func op09(c *SM83) {
	hl := c.r.hl.pack()
	bc := c.r.bc.pack()
	c.r.hl.unpack(hl + bc)
	c.r.f.n, c.r.f.h, c.r.f.c = false, ((hl&0x0FFF)+(bc&0x0FFF) > 0x0FFF), (uint(hl)+uint(bc) > 0xFFFF)
}

func op0A(c *SM83) { c.r.a = c.bus.Read(c.r.bc.pack()) }

func op0B(c *SM83) { c.r.bc.unpack(c.r.bc.pack() - 1) }

func op0C(c *SM83) {
	c.r.bc.lo++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.bc.lo == 0), false, (c.r.bc.lo&0x0F == 0x00)
}

func op0D(c *SM83) {
	c.r.bc.lo--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.bc.lo == 0), true, (c.r.bc.lo&0x0F == 0x0F)
}

func op0E(c *SM83) { c.r.bc.lo = c.fetch() }

// rrca
func op0F(c *SM83) {
	lsb := util.Bit(c.r.a, 0)
	c.r.f.c = lsb
	c.r.a = (c.r.a >> 1) | (util.Btou8(lsb) << 7)
	c.r.f.z, c.r.f.n, c.r.f.h = false, false, false
}

func op10(c *SM83) {
	c.r.pc++ // NOTE: 遊戯王DM4はこれをしっかりしないと動かない
	c.stop()
}

func op11(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.de.unpack((hi << 8) | lo)
}

func op12(c *SM83) { c.bus.Write(c.r.de.pack(), c.r.a) }

func op13(c *SM83) { c.r.de.unpack(c.r.de.pack() + 1) }

func op14(c *SM83) {
	c.r.de.hi++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.de.hi == 0), false, (c.r.de.hi&0x0F == 0x00)
}

func op15(c *SM83) {
	c.r.de.hi--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.de.hi == 0), true, (c.r.de.hi&0x0F == 0x0F)
}

func op16(c *SM83) { c.r.de.hi = c.fetch() }

// rla
func op17(c *SM83) {
	carry := util.Btou8(c.r.f.c)
	c.r.f.c = util.Bit(c.r.a, 7)
	c.r.a = (c.r.a << 1) | carry
	c.r.f.z, c.r.f.n, c.r.f.h = false, false, false
}

func op18(c *SM83) {
	rel := int8(c.fetch())
	c.branch(c.r.pc + uint16(rel))
}

func op19(c *SM83) {
	hl := c.r.hl.pack()
	de := c.r.de.pack()
	c.r.hl.unpack(hl + de)
	c.r.f.n = false
	c.r.f.h = (hl&0x0FFF)+(de&0x0FFF) > 0x0FFF
	c.r.f.c = uint(hl)+uint(de) > 0xFFFF
}

func op1A(c *SM83) { c.r.a = c.bus.Read(c.r.de.pack()) }

func op1B(c *SM83) { c.r.de.unpack(c.r.de.pack() - 1) }

func op1C(c *SM83) {
	c.r.de.lo++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.de.lo == 0), false, (c.r.de.lo&0x0F == 0x00)
}

func op1D(c *SM83) {
	c.r.de.lo--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.de.lo == 0), true, (c.r.de.lo&0x0F == 0x0F)
}

func op1E(c *SM83) { c.r.de.lo = c.fetch() }

// rra
func op1F(c *SM83) {
	carry := util.Btou8(c.r.f.c)
	c.r.f.c = util.Bit(c.r.a, 0)
	c.r.a = (c.r.a >> 1) | (carry << 7)
	c.r.f.z, c.r.f.n, c.r.f.h = false, false, false
}

func op20(c *SM83) {
	rel := int8(c.fetch())
	if !c.r.f.z {
		c.branch(c.r.pc + uint16(rel))
	}
}

func op21(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.hl.unpack((hi << 8) | lo)
}

func op22(c *SM83) {
	c.bus.Write(c.r.hl.pack(), c.r.a)
	c.r.hl.unpack(c.r.hl.pack() + 1)
}

func op23(c *SM83) { c.r.hl.unpack(c.r.hl.pack() + 1) }

func op24(c *SM83) {
	c.r.hl.hi++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.hl.hi == 0), false, (c.r.hl.hi&0x0F == 0x00)
}

func op25(c *SM83) {
	c.r.hl.hi--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.hl.hi == 0), true, (c.r.hl.hi&0x0F == 0x0F)
}

func op26(c *SM83) { c.r.hl.hi = c.fetch() }

// daa
func op27(c *SM83) {
	carry := c.r.f.c
	if !c.r.f.n {
		if carry || c.r.a > 0x99 {
			c.r.a += 0x60
			c.r.f.c = true
		}
		if c.r.f.h || (c.r.a&0xF) > 0x09 {
			c.r.a += 0x06
		}
	} else {
		if carry {
			c.r.a -= 0x60
		}
		if c.r.f.h {
			c.r.a -= 0x06
		}
	}
	c.r.f.z, c.r.f.h = (c.r.a == 0), false
}

func op28(c *SM83) {
	rel := int8(c.fetch())
	if c.r.f.z {
		c.branch(c.r.pc + uint16(rel))
	}
}

// add hl, hl
func op29(c *SM83) {
	hl := c.r.hl.pack()
	result := uint32(hl) + uint32(hl)
	c.r.hl.unpack(uint16(result))
	c.r.f.n, c.r.f.h, c.r.f.c = false, ((hl&0x0FFF)+(hl&0x0FFF) > 0x0FFF), (result > 0xFFFF)
}

func op2A(c *SM83) {
	c.r.a = c.bus.Read(c.r.hl.pack())
	c.r.hl.unpack(c.r.hl.pack() + 1)
}

func op2B(c *SM83) { c.r.hl.unpack(c.r.hl.pack() - 1) }

func op2C(c *SM83) {
	c.r.hl.lo++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.hl.lo == 0), false, (c.r.hl.lo&0x0F == 0x00)
}

// dec l
func op2D(c *SM83) {
	c.r.hl.lo--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.hl.lo == 0), true, (c.r.hl.lo&0x0F == 0x0F)
}

func op2E(c *SM83) { c.r.hl.lo = c.fetch() }

func op2F(c *SM83) {
	c.r.a = ^c.r.a
	c.r.f.n, c.r.f.h = true, true
}

func op30(c *SM83) {
	rel := int8(c.fetch())
	if !c.r.f.c {
		c.branch(c.r.pc + uint16(rel))
	}
}

func op31(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.r.sp = (hi << 8) | lo
}

func op32(c *SM83) {
	c.bus.Write(c.r.hl.pack(), c.r.a)
	c.r.hl.unpack(c.r.hl.pack() - 1)
}

func op33(c *SM83) { c.r.sp++ }

func op34(c *SM83) {
	hl := c.r.hl.pack()
	val := c.bus.Read(hl)
	val++
	c.bus.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h = (val == 0), false, (val&0x0F == 0x00)
}

func op35(c *SM83) {
	hl := c.r.hl.pack()
	val := c.bus.Read(hl)
	val--
	c.bus.Write(hl, val)
	c.r.f.z, c.r.f.n, c.r.f.h = (val == 0), true, (val&0x0F == 0x0F)
}

func op36(c *SM83) {
	hl := c.r.hl.pack()
	val := c.fetch()
	c.bus.Write(hl, val)
}

func op37(c *SM83) { c.r.f.n, c.r.f.h, c.r.f.c = false, false, true }

func op38(c *SM83) {
	rel := int8(c.fetch())
	if c.r.f.c {
		c.branch(c.r.pc + uint16(rel))
	}
}

func op39(c *SM83) {
	sp := c.r.sp
	hl := c.r.hl.pack()
	result := uint32(sp) + uint32(hl)
	c.r.hl.unpack(uint16(result))
	c.r.f.n, c.r.f.h, c.r.f.c = false, ((sp&0x0FFF)+(hl&0x0FFF) > 0x0FFF), (result > 0xFFFF)
}

func op3A(c *SM83) {
	c.r.a = c.bus.Read(c.r.hl.pack())
	c.r.hl.unpack(c.r.hl.pack() - 1)
}

func op3B(c *SM83) { c.r.sp-- }

func op3C(c *SM83) {
	c.r.a++
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.a == 0), false, (c.r.a&0x0F == 0x00)
}

func op3D(c *SM83) {
	c.r.a--
	c.r.f.z, c.r.f.n, c.r.f.h = (c.r.a == 0), true, (c.r.a&0x0F == 0x0F)
}

func op3E(c *SM83) { c.r.a = c.fetch() }

func op3F(c *SM83) { c.r.f.n, c.r.f.h, c.r.f.c = false, false, !c.r.f.c }

func op40(c *SM83) { /* ld b, b */ }

func op41(c *SM83) { c.r.bc.hi = c.r.bc.lo }

func op42(c *SM83) { c.r.bc.hi = c.r.de.hi }

func op43(c *SM83) { c.r.bc.hi = c.r.de.lo }

func op44(c *SM83) { c.r.bc.hi = c.r.hl.hi }

func op45(c *SM83) { c.r.bc.hi = c.r.hl.lo }

func op46(c *SM83) { c.r.bc.hi = c.bus.Read(c.r.hl.pack()) }

func op47(c *SM83) { c.r.bc.hi = c.r.a }

func op48(c *SM83) { c.r.bc.lo = c.r.bc.hi }

func op49(c *SM83) { /* ld c, c */ }

func op4A(c *SM83) { c.r.bc.lo = c.r.de.hi }

func op4B(c *SM83) { c.r.bc.lo = c.r.de.lo }

func op4C(c *SM83) { c.r.bc.lo = c.r.hl.hi }

func op4D(c *SM83) { c.r.bc.lo = c.r.hl.lo }

func op4E(c *SM83) { c.r.bc.lo = c.bus.Read(c.r.hl.pack()) }

func op4F(c *SM83) { c.r.bc.lo = c.r.a }

func op50(c *SM83) { c.r.de.hi = c.r.bc.hi }

func op51(c *SM83) { c.r.de.hi = c.r.bc.lo }

func op52(c *SM83) { /* ld d, d */ }

func op53(c *SM83) { c.r.de.hi = c.r.de.lo }

func op54(c *SM83) { c.r.de.hi = c.r.hl.hi }

func op55(c *SM83) { c.r.de.hi = c.r.hl.lo }

func op56(c *SM83) { c.r.de.hi = c.bus.Read(c.r.hl.pack()) }

func op57(c *SM83) { c.r.de.hi = c.r.a }

func op58(c *SM83) { c.r.de.lo = c.r.bc.hi }

func op59(c *SM83) { c.r.de.lo = c.r.bc.lo }

// ld e, d
func op5A(c *SM83) { c.r.de.lo = c.r.de.hi }

func op5B(c *SM83) { /* ld e, e */ }

func op5C(c *SM83) { c.r.de.lo = c.r.hl.hi }

func op5D(c *SM83) { c.r.de.lo = c.r.hl.lo }

func op5E(c *SM83) { c.r.de.lo = c.bus.Read(c.r.hl.pack()) }

func op5F(c *SM83) { c.r.de.lo = c.r.a }

func op60(c *SM83) { c.r.hl.hi = c.r.bc.hi }

func op61(c *SM83) { c.r.hl.hi = c.r.bc.lo }

func op62(c *SM83) { c.r.hl.hi = c.r.de.hi }

func op63(c *SM83) { c.r.hl.hi = c.r.de.lo }

func op64(c *SM83) { /* ld h, h */ }

func op65(c *SM83) { c.r.hl.hi = c.r.hl.lo }

func op66(c *SM83) { c.r.hl.hi = c.bus.Read(c.r.hl.pack()) }

func op67(c *SM83) { c.r.hl.hi = c.r.a }

func op68(c *SM83) { c.r.hl.lo = c.r.bc.hi }

func op69(c *SM83) { c.r.hl.lo = c.r.bc.lo }

func op6A(c *SM83) { c.r.hl.lo = c.r.de.hi }

func op6B(c *SM83) { c.r.hl.lo = c.r.de.lo }

func op6C(c *SM83) { c.r.hl.lo = c.r.hl.hi }

func op6D(c *SM83) { /* ld l, l */ }

func op6E(c *SM83) { c.r.hl.lo = c.bus.Read(c.r.hl.pack()) }

func op6F(c *SM83) { c.r.hl.lo = c.r.a }

func op70(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.bc.hi) }

func op71(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.bc.lo) }

func op72(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.de.hi) }

func op73(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.de.lo) }

func op74(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.hl.hi) }

func op75(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.hl.lo) }

func op76(c *SM83) { c.halt() }

func op77(c *SM83) { c.bus.Write(c.r.hl.pack(), c.r.a) }

func op78(c *SM83) { c.r.a = c.r.bc.hi }

func op79(c *SM83) { c.r.a = c.r.bc.lo }

func op7A(c *SM83) { c.r.a = c.r.de.hi }

func op7B(c *SM83) { c.r.a = c.r.de.lo }

func op7C(c *SM83) { c.r.a = c.r.hl.hi }

func op7D(c *SM83) { c.r.a = c.r.hl.lo }

func op7E(c *SM83) { c.r.a = c.bus.Read(c.r.hl.pack()) }

func op7F(c *SM83) { /* ld a, a */ }

func op80(c *SM83) { c.add(c.r.bc.hi, false) }

func op81(c *SM83) { c.add(c.r.bc.lo, false) }

func op82(c *SM83) { c.add(c.r.de.hi, false) }

func op83(c *SM83) { c.add(c.r.de.lo, false) }

func op84(c *SM83) { c.add(c.r.hl.hi, false) }

func op85(c *SM83) { c.add(c.r.hl.lo, false) }

func op86(c *SM83) { c.add(c.bus.Read(c.r.hl.pack()), false) }

func op87(c *SM83) { c.add(c.r.a, false) }

func op88(c *SM83) { c.add(c.r.bc.hi, c.r.f.c) }

func op89(c *SM83) { c.add(c.r.bc.lo, c.r.f.c) }

func op8A(c *SM83) { c.add(c.r.de.hi, c.r.f.c) }

func op8B(c *SM83) { c.add(c.r.de.lo, c.r.f.c) }

func op8C(c *SM83) { c.add(c.r.hl.hi, c.r.f.c) }

func op8D(c *SM83) { c.add(c.r.hl.lo, c.r.f.c) }

func op8E(c *SM83) { c.add(c.bus.Read(c.r.hl.pack()), c.r.f.c) }

func op8F(c *SM83) { c.add(c.r.a, c.r.f.c) }

func op90(c *SM83) { c.sub(c.r.bc.hi, false) }

func op91(c *SM83) { c.sub(c.r.bc.lo, false) }

func op92(c *SM83) { c.sub(c.r.de.hi, false) }

func op93(c *SM83) { c.sub(c.r.de.lo, false) } // sub a, e

func op94(c *SM83) { c.sub(c.r.hl.hi, false) }

func op95(c *SM83) { c.sub(c.r.hl.lo, false) }

func op96(c *SM83) { c.sub(c.bus.Read(c.r.hl.pack()), false) }

func op97(c *SM83) { c.sub(c.r.a, false) }

func op98(c *SM83) { c.sub(c.r.bc.hi, c.r.f.c) }

func op99(c *SM83) { c.sub(c.r.bc.lo, c.r.f.c) }

func op9A(c *SM83) { c.sub(c.r.de.hi, c.r.f.c) } // sbc a, d

func op9B(c *SM83) { c.sub(c.r.de.lo, c.r.f.c) }

func op9C(c *SM83) { c.sub(c.r.hl.hi, c.r.f.c) }

func op9D(c *SM83) { c.sub(c.r.hl.lo, c.r.f.c) }

func op9E(c *SM83) { c.sub(c.bus.Read(c.r.hl.pack()), c.r.f.c) }

func op9F(c *SM83) { c.sub(c.r.a, c.r.f.c) }

func opA0(c *SM83) {
	c.r.a &= c.r.bc.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA1(c *SM83) {
	c.r.a &= c.r.bc.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA2(c *SM83) {
	c.r.a &= c.r.de.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA3(c *SM83) {
	c.r.a &= c.r.de.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA4(c *SM83) {
	c.r.a &= c.r.hl.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA5(c *SM83) {
	c.r.a &= c.r.hl.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA6(c *SM83) {
	c.r.a &= c.bus.Read(c.r.hl.pack())
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA7(c *SM83) {
	c.r.a &= c.r.a
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opA8(c *SM83) {
	c.r.a ^= c.r.bc.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opA9(c *SM83) {
	c.r.a ^= c.r.bc.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opAA(c *SM83) {
	c.r.a ^= c.r.de.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opAB(c *SM83) {
	c.r.a ^= c.r.de.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opAC(c *SM83) {
	c.r.a ^= c.r.hl.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opAD(c *SM83) {
	c.r.a ^= c.r.hl.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opAE(c *SM83) {
	c.r.a ^= c.bus.Read(c.r.hl.pack())
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opAF(c *SM83) {
	c.r.a ^= c.r.a
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB0(c *SM83) {
	c.r.a |= c.r.bc.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB1(c *SM83) {
	c.r.a |= c.r.bc.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB2(c *SM83) {
	c.r.a |= c.r.de.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB3(c *SM83) {
	c.r.a |= c.r.de.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB4(c *SM83) {
	c.r.a |= c.r.hl.hi
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB5(c *SM83) {
	c.r.a |= c.r.hl.lo
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB6(c *SM83) {
	c.r.a |= c.bus.Read(c.r.hl.pack())
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB7(c *SM83) {
	c.r.a |= c.r.a
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opB8(c *SM83) { c.cp(c.r.bc.hi) }

func opB9(c *SM83) { c.cp(c.r.bc.lo) }

func opBA(c *SM83) { c.cp(c.r.de.hi) }

func opBB(c *SM83) { c.cp(c.r.de.lo) }

func opBC(c *SM83) { c.cp(c.r.hl.hi) }

func opBD(c *SM83) { c.cp(c.r.hl.lo) }

func opBE(c *SM83) { c.cp(c.bus.Read(c.r.hl.pack())) }

func opBF(c *SM83) { c.cp(c.r.a) }

func opC0(c *SM83) {
	if !c.r.f.z {
		c.ret()
	}
}

func opC1(c *SM83) { c.r.bc.unpack(c.pop16()) }

func opC2(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.r.f.z {
		c.branch((hi << 8) | lo)
	}
}

func opC3(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.branch((hi << 8) | lo)
}

func opC4(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.r.f.z {
		c.call((hi << 8) | lo)
	}
}

func opC5(c *SM83) { c.push16(c.r.bc.pack()) }

func opC6(c *SM83) {
	a := c.r.a
	val := c.fetch()
	result := uint16(a) + uint16(val)
	c.r.a = a + val
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, ((a&0x0F)+(val&0x0F) > 0x0F), (result > 0xFF)
}

func opC7(c *SM83) { c.call(0x00) }

func opC8(c *SM83) {
	if c.r.f.z {
		c.ret()
	}
}

func opC9(c *SM83) { c.ret() }

func opCA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.r.f.z {
		c.branch((hi << 8) | lo)
	}
}

func opCB(c *SM83) {
	opcode := c.fetch()
	c.inst.opcode = opcode
	c.inst.cb = true

	(cbTable[opcode])(c)
	c.tick(cbCycles[opcode])
}

func opCC(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.r.f.z {
		c.call((hi << 8) | lo)
	}
}

func opCD(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.call((hi << 8) | lo)
}

func opCE(c *SM83) { c.add(c.fetch(), c.r.f.c) } // adc a, u8

func opCF(c *SM83) { c.call(0x08) }

func opD0(c *SM83) {
	if !c.r.f.c {
		c.ret()
	}
}

func opD1(c *SM83) { c.r.de.unpack(c.pop16()) }

func opD2(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.r.f.c {
		c.branch((hi << 8) | lo)
	}
}

func opD4(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.r.f.c {
		c.call((hi << 8) | lo)
	}
}

func opD5(c *SM83) { c.push16(c.r.de.pack()) }

// sub a, u8
func opD6(c *SM83) {
	a, val := c.r.a, c.fetch()
	diff := int(a) - int(val)
	c.r.a = uint8(diff)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), true, (int(a&0x0F)-int(val&0x0F) < 0), (diff < 0)
}

func opD7(c *SM83) { c.call(0x10) }

func opD8(c *SM83) {
	if c.r.f.c {
		c.ret()
	}
}

// reti
func opD9(c *SM83) {
	c.IME = true
	c.ret()
}

func opDA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.r.f.c {
		c.branch((hi << 8) | lo)
	}
}

// call c, u16
func opDC(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.r.f.c {
		c.call((hi << 8) | lo)
	}
}

// sbc a, u8
func opDE(c *SM83) {
	carry := util.Btou8(c.r.f.c)
	a, val := c.r.a, c.fetch()
	result := int(a) - int(val) - int(carry)
	c.r.a = uint8(result)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), true, (int(a&0x0F)-int(val&0x0F)-int(carry) < 0), (result < 0)
}

func opDF(c *SM83) { c.call(0x18) }

func opE0(c *SM83) {
	addr := 0xFF00 | uint16(c.fetch())
	c.bus.Write(addr, c.r.a)
}

func opE1(c *SM83) { c.r.hl.unpack(c.pop16()) } // pop hl

func opE2(c *SM83) {
	addr := 0xFF00 | uint16(c.r.bc.lo)
	c.bus.Write(addr, c.r.a)
}

func opE5(c *SM83) { c.push16(c.r.hl.pack()) } // push hl

// and a, u8
func opE6(c *SM83) {
	c.r.a &= c.fetch()
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, true, false
}

func opE7(c *SM83) { c.call(0x20) }

func opE8(c *SM83) {
	sp := c.r.sp
	rel := int8(c.fetch())
	val := sp + uint16(rel)
	c.r.sp = val
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = false, false, ((sp&0x0F)+(uint16(rel)&0x0F) > 0x0F), ((val & 0xFF) < (sp & 0xFF))
}

func opE9(c *SM83) { c.branch(c.r.hl.pack()) }

func opEA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.bus.Write(addr, c.r.a)
}

// xor a, u8
func opEE(c *SM83) {
	c.r.a ^= c.fetch()
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opEF(c *SM83) { c.call(0x28) }

func opF0(c *SM83) {
	addr := 0xFF00 | uint16(c.fetch())
	c.r.a = c.bus.Read(addr)
}

func opF1(c *SM83) {
	af := c.pop16()
	c.r.a = uint8(af >> 8)
	c.r.f.unpack(uint8(af))
}

func opF2(c *SM83) {
	addr := 0xFF00 | uint16(c.r.bc.lo)
	c.r.a = c.bus.Read(addr)
}

func opF3(c *SM83) { c.IME = false }

// push af
func opF5(c *SM83) {
	a := uint16(c.r.a)
	f := uint16(c.r.f.pack())
	af := (a << 8) | f
	c.push16(af)
}

func opF6(c *SM83) {
	c.r.a |= c.fetch()
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, false, false
}

func opF7(c *SM83) { c.call(0x30) }

func opF8(c *SM83) {
	rel := int8(c.fetch())
	val := c.r.sp + uint16(rel)
	c.r.hl.unpack(val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = false, false, ((c.r.sp&0x0F)+(uint16(rel)&0x0F) > 0x0F), ((int(c.r.sp)&0xFF)+int(rel)&0xFF) > 0xFF
}

func opF9(c *SM83) { c.r.sp = c.r.hl.pack() }

func opFA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.r.a = c.bus.Read(addr)
}

func opFB(c *SM83) { c.IME = true }

func opFE(c *SM83) {
	val := c.fetch()
	diff := int(c.r.a) - int(val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (diff == 0), true, ((c.r.a & 0x0F) < (val & 0x0F)), (diff < 0)
}

func opFF(c *SM83) { c.call(0x38) }
