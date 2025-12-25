package sm83

import (
	"github.com/akatsuki105/dawngb/core/gb/internal"
)

type Pair struct {
	Lo, Hi uint8
}

func (p *Pair) Pack() uint16 {
	return uint16(p.Hi)<<8 | uint16(p.Lo)
}

func (p *Pair) Unpack(val uint16) {
	p.Lo = uint8(val)
	p.Hi = uint8(val >> 8)
}

type Registers struct {
	A          uint8
	F          PSR
	BC, DE, HL Pair
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
type PSR struct {
	z, n, h, c bool
}

func (p *PSR) Pack() uint8 {
	packed := uint8(0)
	packed = internal.SetBit(packed, 7, p.z)
	packed = internal.SetBit(packed, 6, p.n)
	packed = internal.SetBit(packed, 5, p.h)
	packed = internal.SetBit(packed, 4, p.c)
	return packed
}

func (p *PSR) Unpack(val uint8) {
	p.z = internal.Bit(val, 7)
	p.n = internal.Bit(val, 6)
	p.h = internal.Bit(val, 5)
	p.c = internal.Bit(val, 4)
}
