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
	c.R.BC.Unpack((hi << 8) | lo)
}

func op02(c *SM83) { c.bus.Write(c.R.BC.Pack(), c.R.A) }

func op03(c *SM83) { c.R.BC.Unpack(c.R.BC.Pack() + 1) }

func op04(c *SM83) {
	c.R.BC.Hi++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.BC.Hi == 0), false, (c.R.BC.Hi&0x0F == 0x00)
}

func op05(c *SM83) {
	c.R.BC.Hi--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.BC.Hi == 0), true, (c.R.BC.Hi&0x0F == 0x0F)
}

func op06(c *SM83) { c.R.BC.Hi = c.fetch() }

// rlca
func op07(c *SM83) {
	msb := (c.R.A >> 7) & 1
	c.R.F.c = msb != 0
	c.R.A = (c.R.A << 1) | msb
	c.R.F.z, c.R.F.n, c.R.F.h = false, false, false
}

func op08(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.bus.Write(addr, uint8(c.R.SP))
	c.bus.Write(addr+1, uint8(c.R.SP>>8))
}

func op09(c *SM83) {
	hl := c.R.HL.Pack()
	bc := c.R.BC.Pack()
	c.R.HL.Unpack(hl + bc)
	c.R.F.n, c.R.F.h, c.R.F.c = false, ((hl&0x0FFF)+(bc&0x0FFF) > 0x0FFF), (uint(hl)+uint(bc) > 0xFFFF)
}

func op0A(c *SM83) { c.R.A = c.bus.Read(c.R.BC.Pack()) }

func op0B(c *SM83) { c.R.BC.Unpack(c.R.BC.Pack() - 1) }

func op0C(c *SM83) {
	c.R.BC.Lo++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.BC.Lo == 0), false, (c.R.BC.Lo&0x0F == 0x00)
}

func op0D(c *SM83) {
	c.R.BC.Lo--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.BC.Lo == 0), true, (c.R.BC.Lo&0x0F == 0x0F)
}

func op0E(c *SM83) { c.R.BC.Lo = c.fetch() }

// rrca
func op0F(c *SM83) {
	lsb := c.R.A & 1
	c.R.F.c = lsb != 0
	c.R.A = (c.R.A >> 1) | (lsb << 7)
	c.R.F.z, c.R.F.n, c.R.F.h = false, false, false
}

func op10(c *SM83) {
	c.R.PC++ // NOTE: 遊戯王DM4はこれをしっかりしないと動かない
	c.stop()
}

func op11(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.R.DE.Unpack((hi << 8) | lo)
}

func op12(c *SM83) { c.bus.Write(c.R.DE.Pack(), c.R.A) }

func op13(c *SM83) { c.R.DE.Unpack(c.R.DE.Pack() + 1) }

func op14(c *SM83) {
	c.R.DE.Hi++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.DE.Hi == 0), false, (c.R.DE.Hi&0x0F == 0x00)
}

func op15(c *SM83) {
	c.R.DE.Hi--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.DE.Hi == 0), true, (c.R.DE.Hi&0x0F == 0x0F)
}

func op16(c *SM83) { c.R.DE.Hi = c.fetch() }

// rla
func op17(c *SM83) {
	carry := btou8(c.R.F.c)
	c.R.F.c = util.Bit(c.R.A, 7)
	c.R.A = (c.R.A << 1) | carry
	c.R.F.z, c.R.F.n, c.R.F.h = false, false, false
}

func op18(c *SM83) {
	rel := int8(c.fetch())
	c.branch(c.R.PC + uint16(rel))
}

func op19(c *SM83) {
	hl := c.R.HL.Pack()
	de := c.R.DE.Pack()
	c.R.HL.Unpack(hl + de)
	c.R.F.n = false
	c.R.F.h = (hl&0x0FFF)+(de&0x0FFF) > 0x0FFF
	c.R.F.c = uint(hl)+uint(de) > 0xFFFF
}

func op1A(c *SM83) { c.R.A = c.bus.Read(c.R.DE.Pack()) }

func op1B(c *SM83) { c.R.DE.Unpack(c.R.DE.Pack() - 1) }

func op1C(c *SM83) {
	c.R.DE.Lo++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.DE.Lo == 0), false, (c.R.DE.Lo&0x0F == 0x00)
}

func op1D(c *SM83) {
	c.R.DE.Lo--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.DE.Lo == 0), true, (c.R.DE.Lo&0x0F == 0x0F)
}

func op1E(c *SM83) { c.R.DE.Lo = c.fetch() }

// rra
func op1F(c *SM83) {
	carry := btou8(c.R.F.c)
	c.R.F.c = util.Bit(c.R.A, 0)
	c.R.A = (c.R.A >> 1) | (carry << 7)
	c.R.F.z, c.R.F.n, c.R.F.h = false, false, false
}

func op20(c *SM83) {
	rel := int8(c.fetch())
	if !c.R.F.z {
		c.branch(c.R.PC + uint16(rel))
	}
}

func op21(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.R.HL.Unpack((hi << 8) | lo)
}

func op22(c *SM83) {
	c.bus.Write(c.R.HL.Pack(), c.R.A)
	c.R.HL.Unpack(c.R.HL.Pack() + 1)
}

func op23(c *SM83) { c.R.HL.Unpack(c.R.HL.Pack() + 1) }

func op24(c *SM83) {
	c.R.HL.Hi++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.HL.Hi == 0), false, (c.R.HL.Hi&0x0F == 0x00)
}

func op25(c *SM83) {
	c.R.HL.Hi--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.HL.Hi == 0), true, (c.R.HL.Hi&0x0F == 0x0F)
}

func op26(c *SM83) { c.R.HL.Hi = c.fetch() }

// daa
func op27(c *SM83) {
	carry := c.R.F.c
	if !c.R.F.n {
		if carry || c.R.A > 0x99 {
			c.R.A += 0x60
			c.R.F.c = true
		}
		if c.R.F.h || (c.R.A&0xF) > 0x09 {
			c.R.A += 0x06
		}
	} else {
		if carry {
			c.R.A -= 0x60
		}
		if c.R.F.h {
			c.R.A -= 0x06
		}
	}
	c.R.F.z, c.R.F.h = (c.R.A == 0), false
}

func op28(c *SM83) {
	rel := int8(c.fetch())
	if c.R.F.z {
		c.branch(c.R.PC + uint16(rel))
	}
}

// add hl, hl
func op29(c *SM83) {
	hl := c.R.HL.Pack()
	result := uint32(hl) + uint32(hl)
	c.R.HL.Unpack(uint16(result))
	c.R.F.n, c.R.F.h, c.R.F.c = false, ((hl&0x0FFF)+(hl&0x0FFF) > 0x0FFF), (result > 0xFFFF)
}

func op2A(c *SM83) {
	c.R.A = c.bus.Read(c.R.HL.Pack())
	c.R.HL.Unpack(c.R.HL.Pack() + 1)
}

func op2B(c *SM83) { c.R.HL.Unpack(c.R.HL.Pack() - 1) }

func op2C(c *SM83) {
	c.R.HL.Lo++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.HL.Lo == 0), false, (c.R.HL.Lo&0x0F == 0x00)
}

// dec l
func op2D(c *SM83) {
	c.R.HL.Lo--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.HL.Lo == 0), true, (c.R.HL.Lo&0x0F == 0x0F)
}

func op2E(c *SM83) { c.R.HL.Lo = c.fetch() }

func op2F(c *SM83) {
	c.R.A = ^c.R.A
	c.R.F.n, c.R.F.h = true, true
}

func op30(c *SM83) {
	rel := int8(c.fetch())
	if !c.R.F.c {
		c.branch(c.R.PC + uint16(rel))
	}
}

func op31(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.R.SP = (hi << 8) | lo
}

func op32(c *SM83) {
	c.bus.Write(c.R.HL.Pack(), c.R.A)
	c.R.HL.Unpack(c.R.HL.Pack() - 1)
}

func op33(c *SM83) { c.R.SP++ }

func op34(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	val++
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h = (val == 0), false, (val&0x0F == 0x00)
}

func op35(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.bus.Read(hl)
	val--
	c.bus.Write(hl, val)
	c.R.F.z, c.R.F.n, c.R.F.h = (val == 0), true, (val&0x0F == 0x0F)
}

func op36(c *SM83) {
	hl := c.R.HL.Pack()
	val := c.fetch()
	c.bus.Write(hl, val)
}

func op37(c *SM83) { c.R.F.n, c.R.F.h, c.R.F.c = false, false, true }

func op38(c *SM83) {
	rel := int8(c.fetch())
	if c.R.F.c {
		c.branch(c.R.PC + uint16(rel))
	}
}

func op39(c *SM83) {
	sp := c.R.SP
	hl := c.R.HL.Pack()
	result := uint32(sp) + uint32(hl)
	c.R.HL.Unpack(uint16(result))
	c.R.F.n, c.R.F.h, c.R.F.c = false, ((sp&0x0FFF)+(hl&0x0FFF) > 0x0FFF), (result > 0xFFFF)
}

func op3A(c *SM83) {
	c.R.A = c.bus.Read(c.R.HL.Pack())
	c.R.HL.Unpack(c.R.HL.Pack() - 1)
}

func op3B(c *SM83) { c.R.SP-- }

func op3C(c *SM83) {
	c.R.A++
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.A == 0), false, (c.R.A&0x0F == 0x00)
}

func op3D(c *SM83) {
	c.R.A--
	c.R.F.z, c.R.F.n, c.R.F.h = (c.R.A == 0), true, (c.R.A&0x0F == 0x0F)
}

func op3E(c *SM83) { c.R.A = c.fetch() }

func op3F(c *SM83) { c.R.F.n, c.R.F.h, c.R.F.c = false, false, !c.R.F.c }

func op40(c *SM83) { /* ld b, b */ }

func op41(c *SM83) { c.R.BC.Hi = c.R.BC.Lo }

func op42(c *SM83) { c.R.BC.Hi = c.R.DE.Hi }

func op43(c *SM83) { c.R.BC.Hi = c.R.DE.Lo }

func op44(c *SM83) { c.R.BC.Hi = c.R.HL.Hi }

func op45(c *SM83) { c.R.BC.Hi = c.R.HL.Lo }

func op46(c *SM83) { c.R.BC.Hi = c.bus.Read(c.R.HL.Pack()) }

func op47(c *SM83) { c.R.BC.Hi = c.R.A }

func op48(c *SM83) { c.R.BC.Lo = c.R.BC.Hi }

func op49(c *SM83) { /* ld c, c */ }

func op4A(c *SM83) { c.R.BC.Lo = c.R.DE.Hi }

func op4B(c *SM83) { c.R.BC.Lo = c.R.DE.Lo }

func op4C(c *SM83) { c.R.BC.Lo = c.R.HL.Hi }

func op4D(c *SM83) { c.R.BC.Lo = c.R.HL.Lo }

func op4E(c *SM83) { c.R.BC.Lo = c.bus.Read(c.R.HL.Pack()) }

func op4F(c *SM83) { c.R.BC.Lo = c.R.A }

func op50(c *SM83) { c.R.DE.Hi = c.R.BC.Hi }

func op51(c *SM83) { c.R.DE.Hi = c.R.BC.Lo }

func op52(c *SM83) { /* ld d, d */ }

func op53(c *SM83) { c.R.DE.Hi = c.R.DE.Lo }

func op54(c *SM83) { c.R.DE.Hi = c.R.HL.Hi }

func op55(c *SM83) { c.R.DE.Hi = c.R.HL.Lo }

func op56(c *SM83) { c.R.DE.Hi = c.bus.Read(c.R.HL.Pack()) }

func op57(c *SM83) { c.R.DE.Hi = c.R.A }

func op58(c *SM83) { c.R.DE.Lo = c.R.BC.Hi }

func op59(c *SM83) { c.R.DE.Lo = c.R.BC.Lo }

// ld e, d
func op5A(c *SM83) { c.R.DE.Lo = c.R.DE.Hi }

func op5B(c *SM83) { /* ld e, e */ }

func op5C(c *SM83) { c.R.DE.Lo = c.R.HL.Hi }

func op5D(c *SM83) { c.R.DE.Lo = c.R.HL.Lo }

func op5E(c *SM83) { c.R.DE.Lo = c.bus.Read(c.R.HL.Pack()) }

func op5F(c *SM83) { c.R.DE.Lo = c.R.A }

func op60(c *SM83) { c.R.HL.Hi = c.R.BC.Hi }

func op61(c *SM83) { c.R.HL.Hi = c.R.BC.Lo }

func op62(c *SM83) { c.R.HL.Hi = c.R.DE.Hi }

func op63(c *SM83) { c.R.HL.Hi = c.R.DE.Lo }

func op64(c *SM83) { /* ld h, h */ }

func op65(c *SM83) { c.R.HL.Hi = c.R.HL.Lo }

func op66(c *SM83) { c.R.HL.Hi = c.bus.Read(c.R.HL.Pack()) }

func op67(c *SM83) { c.R.HL.Hi = c.R.A }

func op68(c *SM83) { c.R.HL.Lo = c.R.BC.Hi }

func op69(c *SM83) { c.R.HL.Lo = c.R.BC.Lo }

func op6A(c *SM83) { c.R.HL.Lo = c.R.DE.Hi }

func op6B(c *SM83) { c.R.HL.Lo = c.R.DE.Lo }

func op6C(c *SM83) { c.R.HL.Lo = c.R.HL.Hi }

func op6D(c *SM83) { /* ld l, l */ }

func op6E(c *SM83) { c.R.HL.Lo = c.bus.Read(c.R.HL.Pack()) }

func op6F(c *SM83) { c.R.HL.Lo = c.R.A }

func op70(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.BC.Hi) }

func op71(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.BC.Lo) }

func op72(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.DE.Hi) }

func op73(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.DE.Lo) }

func op74(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.HL.Hi) }

func op75(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.HL.Lo) }

func op76(c *SM83) { c.halt() }

func op77(c *SM83) { c.bus.Write(c.R.HL.Pack(), c.R.A) }

func op78(c *SM83) { c.R.A = c.R.BC.Hi }

func op79(c *SM83) { c.R.A = c.R.BC.Lo }

func op7A(c *SM83) { c.R.A = c.R.DE.Hi }

func op7B(c *SM83) { c.R.A = c.R.DE.Lo }

func op7C(c *SM83) { c.R.A = c.R.HL.Hi }

func op7D(c *SM83) { c.R.A = c.R.HL.Lo }

func op7E(c *SM83) { c.R.A = c.bus.Read(c.R.HL.Pack()) }

func op7F(c *SM83) { /* ld a, a */ }

func op80(c *SM83) { c.add(c.R.BC.Hi, false) }

func op81(c *SM83) { c.add(c.R.BC.Lo, false) }

func op82(c *SM83) { c.add(c.R.DE.Hi, false) }

func op83(c *SM83) { c.add(c.R.DE.Lo, false) }

func op84(c *SM83) { c.add(c.R.HL.Hi, false) }

func op85(c *SM83) { c.add(c.R.HL.Lo, false) }

func op86(c *SM83) { c.add(c.bus.Read(c.R.HL.Pack()), false) }

func op87(c *SM83) { c.add(c.R.A, false) }

func op88(c *SM83) { c.add(c.R.BC.Hi, c.R.F.c) }

func op89(c *SM83) { c.add(c.R.BC.Lo, c.R.F.c) }

func op8A(c *SM83) { c.add(c.R.DE.Hi, c.R.F.c) }

func op8B(c *SM83) { c.add(c.R.DE.Lo, c.R.F.c) }

func op8C(c *SM83) { c.add(c.R.HL.Hi, c.R.F.c) }

func op8D(c *SM83) { c.add(c.R.HL.Lo, c.R.F.c) }

func op8E(c *SM83) { c.add(c.bus.Read(c.R.HL.Pack()), c.R.F.c) }

func op8F(c *SM83) { c.add(c.R.A, c.R.F.c) }

func op90(c *SM83) { c.sub(c.R.BC.Hi, false) }

func op91(c *SM83) { c.sub(c.R.BC.Lo, false) }

func op92(c *SM83) { c.sub(c.R.DE.Hi, false) }

func op93(c *SM83) { c.sub(c.R.DE.Lo, false) } // sub a, e

func op94(c *SM83) { c.sub(c.R.HL.Hi, false) }

func op95(c *SM83) { c.sub(c.R.HL.Lo, false) }

func op96(c *SM83) { c.sub(c.bus.Read(c.R.HL.Pack()), false) }

func op97(c *SM83) { c.sub(c.R.A, false) }

func op98(c *SM83) { c.sub(c.R.BC.Hi, c.R.F.c) }

func op99(c *SM83) { c.sub(c.R.BC.Lo, c.R.F.c) }

func op9A(c *SM83) { c.sub(c.R.DE.Hi, c.R.F.c) } // sbc a, d

func op9B(c *SM83) { c.sub(c.R.DE.Lo, c.R.F.c) }

func op9C(c *SM83) { c.sub(c.R.HL.Hi, c.R.F.c) }

func op9D(c *SM83) { c.sub(c.R.HL.Lo, c.R.F.c) }

func op9E(c *SM83) { c.sub(c.bus.Read(c.R.HL.Pack()), c.R.F.c) }

func op9F(c *SM83) { c.sub(c.R.A, c.R.F.c) }

func opA0(c *SM83) {
	c.R.A &= c.R.BC.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA1(c *SM83) {
	c.R.A &= c.R.BC.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA2(c *SM83) {
	c.R.A &= c.R.DE.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA3(c *SM83) {
	c.R.A &= c.R.DE.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA4(c *SM83) {
	c.R.A &= c.R.HL.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA5(c *SM83) {
	c.R.A &= c.R.HL.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA6(c *SM83) {
	c.R.A &= c.bus.Read(c.R.HL.Pack())
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA7(c *SM83) {
	c.R.A &= c.R.A
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opA8(c *SM83) {
	c.R.A ^= c.R.BC.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opA9(c *SM83) {
	c.R.A ^= c.R.BC.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opAA(c *SM83) {
	c.R.A ^= c.R.DE.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opAB(c *SM83) {
	c.R.A ^= c.R.DE.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opAC(c *SM83) {
	c.R.A ^= c.R.HL.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opAD(c *SM83) {
	c.R.A ^= c.R.HL.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opAE(c *SM83) {
	c.R.A ^= c.bus.Read(c.R.HL.Pack())
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opAF(c *SM83) {
	c.R.A ^= c.R.A
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB0(c *SM83) {
	c.R.A |= c.R.BC.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB1(c *SM83) {
	c.R.A |= c.R.BC.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB2(c *SM83) {
	c.R.A |= c.R.DE.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB3(c *SM83) {
	c.R.A |= c.R.DE.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB4(c *SM83) {
	c.R.A |= c.R.HL.Hi
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB5(c *SM83) {
	c.R.A |= c.R.HL.Lo
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB6(c *SM83) {
	c.R.A |= c.bus.Read(c.R.HL.Pack())
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB7(c *SM83) {
	c.R.A |= c.R.A
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opB8(c *SM83) { c.cp(c.R.BC.Hi) }

func opB9(c *SM83) { c.cp(c.R.BC.Lo) }

func opBA(c *SM83) { c.cp(c.R.DE.Hi) }

func opBB(c *SM83) { c.cp(c.R.DE.Lo) }

func opBC(c *SM83) { c.cp(c.R.HL.Hi) }

func opBD(c *SM83) { c.cp(c.R.HL.Lo) }

func opBE(c *SM83) { c.cp(c.bus.Read(c.R.HL.Pack())) }

func opBF(c *SM83) { c.cp(c.R.A) }

func opC0(c *SM83) {
	if !c.R.F.z {
		c.ret()
	}
}

func opC1(c *SM83) { c.R.BC.Unpack(c.pop16()) }

func opC2(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.R.F.z {
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
	if !c.R.F.z {
		c.call((hi << 8) | lo)
	}
}

func opC5(c *SM83) { c.push16(c.R.BC.Pack()) }

func opC6(c *SM83) {
	a := c.R.A
	val := c.fetch()
	result := uint16(a) + uint16(val)
	c.R.A = a + val
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, ((a&0x0F)+(val&0x0F) > 0x0F), (result > 0xFF)
}

func opC7(c *SM83) { c.call(0x00) }

func opC8(c *SM83) {
	if c.R.F.z {
		c.ret()
	}
}

func opC9(c *SM83) { c.ret() }

func opCA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.R.F.z {
		c.branch((hi << 8) | lo)
	}
}

func opCB(c *SM83) {
	opcode := c.fetch()
	c.inst.Opcode = opcode
	c.inst.CB = true

	(cbTable[opcode])(c)
	c.tick(cbCycles[opcode])
}

func opCC(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.R.F.z {
		c.call((hi << 8) | lo)
	}
}

func opCD(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	c.call((hi << 8) | lo)
}

func opCE(c *SM83) { c.add(c.fetch(), c.R.F.c) } // adc a, u8

func opCF(c *SM83) { c.call(0x08) }

func opD0(c *SM83) {
	if !c.R.F.c {
		c.ret()
	}
}

func opD1(c *SM83) { c.R.DE.Unpack(c.pop16()) }

func opD2(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.R.F.c {
		c.branch((hi << 8) | lo)
	}
}

func opD4(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if !c.R.F.c {
		c.call((hi << 8) | lo)
	}
}

func opD5(c *SM83) { c.push16(c.R.DE.Pack()) }

// sub a, u8
func opD6(c *SM83) {
	a, val := c.R.A, c.fetch()
	diff := int(a) - int(val)
	c.R.A = uint8(diff)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), true, (int(a&0x0F)-int(val&0x0F) < 0), (diff < 0)
}

func opD7(c *SM83) { c.call(0x10) }

func opD8(c *SM83) {
	if c.R.F.c {
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
	if c.R.F.c {
		c.branch((hi << 8) | lo)
	}
}

// call c, u16
func opDC(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	if c.R.F.c {
		c.call((hi << 8) | lo)
	}
}

// sbc a, u8
func opDE(c *SM83) {
	carry := btou8(c.R.F.c)
	a, val := c.R.A, c.fetch()
	result := int(a) - int(val) - int(carry)
	c.R.A = uint8(result)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), true, (int(a&0x0F)-int(val&0x0F)-int(carry) < 0), (result < 0)
}

func opDF(c *SM83) { c.call(0x18) }

func opE0(c *SM83) {
	addr := 0xFF00 | uint16(c.fetch())
	c.bus.Write(addr, c.R.A)
}

func opE1(c *SM83) { c.R.HL.Unpack(c.pop16()) } // pop hl

func opE2(c *SM83) {
	addr := 0xFF00 | uint16(c.R.BC.Lo)
	c.bus.Write(addr, c.R.A)
}

func opE5(c *SM83) { c.push16(c.R.HL.Pack()) } // push hl

// and a, u8
func opE6(c *SM83) {
	c.R.A &= c.fetch()
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, true, false
}

func opE7(c *SM83) { c.call(0x20) }

func opE8(c *SM83) {
	sp := c.R.SP
	rel := int8(c.fetch())
	val := sp + uint16(rel)
	c.R.SP = val
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = false, false, ((sp&0x0F)+(uint16(rel)&0x0F) > 0x0F), ((val & 0xFF) < (sp & 0xFF))
}

func opE9(c *SM83) { c.branch(c.R.HL.Pack()) }

func opEA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.bus.Write(addr, c.R.A)
}

// xor a, u8
func opEE(c *SM83) {
	c.R.A ^= c.fetch()
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opEF(c *SM83) { c.call(0x28) }

func opF0(c *SM83) {
	addr := 0xFF00 | uint16(c.fetch())
	c.R.A = c.bus.Read(addr)
}

func opF1(c *SM83) {
	af := c.pop16()
	c.R.A = uint8(af >> 8)
	c.R.F.Unpack(uint8(af))
}

func opF2(c *SM83) {
	addr := 0xFF00 | uint16(c.R.BC.Lo)
	c.R.A = c.bus.Read(addr)
}

func opF3(c *SM83) { c.IME = false }

// push af
func opF5(c *SM83) {
	a := uint16(c.R.A)
	f := uint16(c.R.F.Pack())
	af := (a << 8) | f
	c.push16(af)
}

func opF6(c *SM83) {
	c.R.A |= c.fetch()
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (c.R.A == 0), false, false, false
}

func opF7(c *SM83) { c.call(0x30) }

func opF8(c *SM83) {
	rel := int8(c.fetch())
	val := c.R.SP + uint16(rel)
	c.R.HL.Unpack(val)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = false, false, ((c.R.SP&0x0F)+(uint16(rel)&0x0F) > 0x0F), ((int(c.R.SP)&0xFF)+int(rel)&0xFF) > 0xFF
}

func opF9(c *SM83) { c.R.SP = c.R.HL.Pack() }

func opFA(c *SM83) {
	lo := uint16(c.fetch())
	hi := uint16(c.fetch())
	addr := (hi << 8) | lo
	c.R.A = c.bus.Read(addr)
}

func opFB(c *SM83) { c.IME = true }

func opFE(c *SM83) {
	val := c.fetch()
	diff := int(c.R.A) - int(val)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (diff == 0), true, ((c.R.A & 0x0F) < (val & 0x0F)), (diff < 0)
}

func opFF(c *SM83) { c.call(0x38) }
