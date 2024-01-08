package gb

import "github.com/akatsuki105/dugb/util"

type Input struct {
	btn, dpad bool
	inputs    [8]bool // A, B, Select, Start, Right, Left, Up, Down
}

func newInput() *Input {
	return &Input{}
}

func (i *Input) ReadIO(addr uint16) uint8 {
	val := uint8(0b11_1111)
	if i.btn {
		val = util.SetBit(val, 5, false)
		val = util.SetBit(val, 0, !i.inputs[0])
		val = util.SetBit(val, 1, !i.inputs[1])
		val = util.SetBit(val, 2, !i.inputs[2])
		val = util.SetBit(val, 3, !i.inputs[3])
		return val
	}
	if i.dpad {
		val = util.SetBit(val, 4, false)
		val = util.SetBit(val, 0, !i.inputs[4])
		val = util.SetBit(val, 1, !i.inputs[5])
		val = util.SetBit(val, 2, !i.inputs[6])
		val = util.SetBit(val, 3, !i.inputs[7])
		return val
	}
	return val
}

func (i *Input) WriteIO(addr uint16, val uint8) {
	i.dpad = !util.Bit(val, 4)
	i.btn = !util.Bit(val, 5)
}
