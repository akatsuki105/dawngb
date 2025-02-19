package sm83

import "github.com/akatsuki105/dawngb/util"

type pair struct {
	Lo, Hi uint8
}

func (p *pair) Pack() uint16 {
	return uint16(p.Hi)<<8 | uint16(p.Lo)
}

func (p *pair) Unpack(val uint16) {
	p.Lo = uint8(val)
	p.Hi = uint8(val >> 8)
}

type Registers struct {
	A          uint8
	F          psr
	BC, DE, HL pair
	SP, PC     uint16
}

func (r *Registers) reset() {
	r.A = 0x00
	r.F.Unpack(0x00)
	r.BC.Unpack(0x0000)
	r.DE.Unpack(0x0000)
	r.HL.Unpack(0x0000)
	r.SP, r.PC = 0x0000, 0x0000
}

// ZNHC----
type psr struct {
	z, n, h, c bool
}

func (p *psr) Pack() uint8 {
	packed := uint8(0)
	packed = util.SetBit(packed, 7, p.z)
	packed = util.SetBit(packed, 6, p.n)
	packed = util.SetBit(packed, 5, p.h)
	packed = util.SetBit(packed, 4, p.c)
	return packed
}

func (p *psr) Unpack(val uint8) {
	p.z = util.Bit(val, 7)
	p.n = util.Bit(val, 6)
	p.h = util.Bit(val, 5)
	p.c = util.Bit(val, 4)
}
