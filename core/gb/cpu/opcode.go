package cpu

import (
	"fmt"

	"github.com/akatsuki105/dugb/util"
)

func todo(c *Cpu) {
	if c.inst.cb {
		panic(fmt.Sprintf("todo opcode: 0xCB+0x%02X in 0x%04X", c.inst.opcode, c.inst.addr))
	} else {
		panic(fmt.Sprintf("todo opcode: 0x%02X in 0x%04X", c.inst.opcode, c.inst.addr))
	}
}

func (c *Cpu) branch(dst uint16) {
	c.r.pc = dst
	c.s.Add(c.Cycle)
}

func (c *Cpu) push8(val uint8) {
	c.r.sp--
	c.m.Write(c.r.sp, val)
	c.s.Add(c.Cycle)
}

func (c *Cpu) push16(val uint16) {
	c.push8(uint8(val >> 8))
	c.push8(uint8(val))
}

func (c *Cpu) pop8() uint8 {
	val := c.m.Read(c.r.sp)
	c.r.sp++
	c.s.Add(c.Cycle)
	return val
}

func (c *Cpu) pop16() uint16 {
	lo := uint16(c.pop8())
	hi := uint16(c.pop8())
	return (hi << 8) | lo
}

func (c *Cpu) Interrupt(id int) {
	c.s.Add(2 * c.Cycle)
	c.IME = false
	c.push16(c.r.pc)
	c.branch([5]uint16{0x40, 0x48, 0x50, 0x58, 0x60}[id])
}

func (c *Cpu) bit(val uint8, bit int) {
	c.r.f.z = !util.Bit(val, bit)
	c.r.f.n, c.r.f.h = false, true
}

func (c *Cpu) ret() {
	c.branch(c.pop16())
}

func (c *Cpu) call(dst uint16) {
	pc := c.r.pc
	c.branch(dst)
	c.push16(pc)
}

func (c *Cpu) cp(val uint8) {
	a := c.r.a
	diff := int(a) - int(val)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (diff == 0), true, (int(a&0x0F)-int(val&0x0F) < 0), (diff < 0)
}

func (c *Cpu) add(val uint8) {
	a := c.r.a
	sum := int(a) + int(val)
	c.r.a = uint8(sum)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, ((a&0x0F)+(val&0x0F) > 0x0F), (sum > 0xFF)
}

func (c *Cpu) sub(val uint8) {
	a := c.r.a
	diff := int(a) - int(val)
	c.r.a = uint8(diff)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), true, (int(a&0x0F)-int(val&0x0F) < 0), (diff < 0)
}

func (c *Cpu) rst(addr uint16) { c.call(addr) }

func (c *Cpu) adc(val uint8) {
	carry := util.Btou8(c.r.f.c)
	a := c.r.a
	result := int(a) + int(val) + int(carry)
	c.r.a = uint8(result)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), false, (int(a&0x0F)+int(val&0x0F)+int(carry) > 0x0F), (result > 0xFF)
}

func (c *Cpu) sbc(val uint8) {
	carry := util.Btou8(c.r.f.c)
	a := c.r.a
	result := int(a) - int(val) - int(carry)
	c.r.a = uint8(result)
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (c.r.a == 0), true, (int(a&0x0F)-int(val&0x0F)-int(carry) < 0), (result < 0)
}

func (c *Cpu) set_hl(bit int, b bool) {
	hl := c.r.hl.pack()
	val := c.m.Read(hl)
	val = util.SetBit(val, bit, b)
	c.m.Write(hl, val)
}

func (c *Cpu) rr(r *uint8) {
	carry := util.Btou8(c.r.f.c)
	c.r.f.c = util.Bit(*r, 0)
	*r = (*r >> 1) | (carry << 7)
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}

func (c *Cpu) rrc(r *uint8) {
	lsb := util.Bit(*r, 0)
	c.r.f.c = lsb
	*r = (*r >> 1) | (util.Btou8(lsb) << 7)
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}

func (c *Cpu) rl(r *uint8) {
	carry := c.r.f.c
	c.r.f.c = util.Bit(*r, 7)
	*r = (*r << 1) | util.Btou8(carry)
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}

func (c *Cpu) sla(r *uint8) {
	c.r.f.c = util.Bit(*r, 7)
	*r <<= 1
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}

func (c *Cpu) srl(r *uint8) {
	c.r.f.c = util.Bit(*r, 0)
	*r >>= 1
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}

func (c *Cpu) swap(r *uint8) {
	*r = (*r&0x0F)<<4 | (*r&0xF0)>>4
	c.r.f.z, c.r.f.n, c.r.f.h, c.r.f.c = (*r == 0), false, false, false
}

func (c *Cpu) sra(r *uint8) {
	msb := *r & 0b1000_0000
	c.r.f.c = util.Bit(*r, 0)
	*r = (*r >> 1) | msb
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}

func (c *Cpu) rlc(r *uint8) {
	msb := util.Bit(*r, 7)
	c.r.f.c = msb
	*r = (*r << 1) | util.Btou8(msb)
	c.r.f.z, c.r.f.n, c.r.f.h = (*r == 0), false, false
}
