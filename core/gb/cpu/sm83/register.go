package sm83

import "github.com/akatsuki105/dawngb/util"

type pair struct {
	lo, hi uint8
}

func (p *pair) pack() uint16 {
	return uint16(p.hi)<<8 | uint16(p.lo)
}

func (p *pair) unpack(val uint16) {
	p.lo = uint8(val)
	p.hi = uint8(val >> 8)
}

type Registers struct {
	a          uint8
	f          psr
	bc, de, hl pair
	sp, pc     uint16
}

// ZNHC----
type psr struct {
	z, n, h, c bool
}

func (p *psr) pack() uint8 {
	packed := uint8(0)
	packed = util.SetBit(packed, 7, p.z)
	packed = util.SetBit(packed, 6, p.n)
	packed = util.SetBit(packed, 5, p.h)
	packed = util.SetBit(packed, 4, p.c)
	return packed
}

func (p *psr) unpack(val uint8) {
	p.z = util.Bit(val, 7)
	p.n = util.Bit(val, 6)
	p.h = util.Bit(val, 5)
	p.c = util.Bit(val, 4)
}
