package cpu

import "github.com/akatsuki105/dugb/util"

type Registers struct {
	a  uint8
	f  psr
	bc uint16
	de uint16
	hl uint16
	sp uint16
	pc uint16
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
