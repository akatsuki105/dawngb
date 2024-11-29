package sm83

import (
	"fmt"

	"github.com/akatsuki105/dawngb/util"
)

func todo(c *SM83) {
	if c.inst.cb {
		panic(fmt.Sprintf("todo opcode: 0xCB+0x%02X in 0x%04X", c.inst.opcode, c.inst.addr))
	} else {
		panic(fmt.Sprintf("todo opcode: 0x%02X in 0x%04X", c.inst.opcode, c.inst.addr))
	}
}

func (c *SM83) branch(dst uint16) {
	c.R.PC = dst
	c.tick(1)
}

func (c *SM83) push8(val uint8) {
	c.R.SP--
	c.bus.Write(c.R.SP, val)
	c.tick(1)
}

func (c *SM83) push16(val uint16) {
	c.push8(uint8(val >> 8))
	c.push8(uint8(val))
}

func (c *SM83) pop8() uint8 {
	val := c.bus.Read(c.R.SP)
	c.R.SP++
	c.tick(1)
	return val
}

func (c *SM83) pop16() uint16 {
	lo := uint16(c.pop8())
	hi := uint16(c.pop8())
	return (hi << 8) | lo
}

func (c *SM83) Interrupt(id int) {
	c.tick(2)
	c.IME = false
	c.push16(c.R.PC)
	c.branch([5]uint16{0x40, 0x48, 0x50, 0x58, 0x60}[id])
}

func (c *SM83) bit(val uint8, bit int) {
	c.R.F.z = !util.Bit(val, bit)
	c.R.F.n, c.R.F.h = false, true
}

func (c *SM83) ret() {
	c.branch(c.pop16())
}

func (c *SM83) call(dst uint16) {
	c.push16(c.R.PC)
	c.branch(dst)
}

func (c *SM83) cp(val uint8) {
	a := c.R.A
	x := uint16(a) - uint16(val)
	y := (a & 0xF) - (val & 0xF)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (uint8(x) == 0), true, (y > 0x0F), (x > 0xFF)
}

func (c *SM83) add(val uint8, carry bool) {
	x := uint16(c.R.A) + uint16(val) + uint16(util.Btou8(carry))
	y := uint16(c.R.A&0xF) + uint16(val&0xF) + uint16(util.Btou8(carry))
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (uint8(x) == 0), false, (y > 0x0F), (x > 0xFF)
	c.R.A = uint8(x)
}

func (c *SM83) sub(val uint8, carry bool) {
	cf := util.Btou8(carry)
	x := uint16(c.R.A) - uint16(val) - uint16(cf)
	y := (c.R.A & 0xF) - (val & 0xF) - cf
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (uint8(x) == 0), true, (y > 0x0F), (x > 0xFF)
	c.R.A = uint8(x)
}

func (c *SM83) set_hl(bit int, b bool) {
	hl := c.R.HL.pack()
	val := c.bus.Read(hl)
	val = util.SetBit(val, bit, b)
	c.bus.Write(hl, val)
}

func (c *SM83) rr(r *uint8) {
	carry := util.Btou8(c.R.F.c)
	c.R.F.c = util.Bit(*r, 0)
	*r = (*r >> 1) | (carry << 7)
	c.R.F.z, c.R.F.n, c.R.F.h = (*r == 0), false, false
}

func (c *SM83) rrc(r *uint8) {
	*r = (*r << 7) | (*r >> 1)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (*r == 0), false, false, util.Bit(*r, 7)
}

func (c *SM83) rl(r *uint8) {
	carry := util.Bit(*r, 7)
	*r = (*r << 1) | util.Btou8(c.R.F.c)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (*r == 0), false, false, carry
}

func (c *SM83) sla(r *uint8) {
	c.R.F.c = util.Bit(*r, 7)
	*r <<= 1
	c.R.F.z, c.R.F.n, c.R.F.h = (*r == 0), false, false
}

func (c *SM83) srl(r *uint8) {
	carry := util.Bit(*r, 0)
	*r >>= 1
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (*r == 0), false, false, carry
}

func (c *SM83) swap(r *uint8) {
	*r = (*r&0x0F)<<4 | (*r&0xF0)>>4
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (*r == 0), false, false, false
}

func (c *SM83) sra(r *uint8) {
	carry := util.Bit(*r, 0)
	*r = uint8(int8(*r) >> 1)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (*r == 0), false, false, carry
}

func (c *SM83) rlc(r *uint8) {
	*r = (*r << 1) | (*r >> 7)
	c.R.F.z, c.R.F.n, c.R.F.h, c.R.F.c = (*r == 0), false, false, util.Bit(*r, 0)
}
