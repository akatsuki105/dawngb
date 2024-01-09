package gb

import "github.com/akatsuki105/dugb/util"

type Input struct {
	inputs      [8]bool // A, B, Select, Start, Right, Left, Up, Down
	onInterrupt func(id int)
	val         uint8
}

func newInput(onInterrupt func(id int)) *Input {
	return &Input{
		onInterrupt: onInterrupt,
	}
}

func (i *Input) Reset() {
	i.val = 0x0F
}

func (i *Input) Read(addr uint16) uint8 {
	pressed := false

	val := i.val
	if !util.Bit(val, 5) {
		for key := 0; key < 4; key++ {
			if i.inputs[key] {
				val = util.SetBit(val, key, false)
				pressed = true
			}
		}
	}
	if !util.Bit(val, 4) {
		for key := 4; key < 8; key++ {
			if i.inputs[key] {
				val = util.SetBit(val, key-4, false)
				pressed = true
			}
		}
	}

	if pressed {
		i.onInterrupt(4)
	}

	i.val = val
	return val
}

func (i *Input) Write(addr uint16, val uint8) {
	i.val = val | 0x0F
}
