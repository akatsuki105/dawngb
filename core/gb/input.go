package gb

import "github.com/akatsuki105/dugb/util"

type Input struct {
	g   *GB
	val uint8
}

func newInput(g *GB) *Input {
	return &Input{
		g: g,
	}
}

func (i *Input) Reset(hasBIOS bool) {
	i.val = 0x0F
	if !hasBIOS {
		i.val = 0xCF
	}
}

func (i *Input) Read(addr uint16) uint8 {
	pressed := false

	val := i.val
	if !util.Bit(val, 5) {
		for key := 0; key < 4; key++ {
			if i.g.inputs[key] {
				val = util.SetBit(val, key, false)
				pressed = true
			}
		}
	}
	if !util.Bit(val, 4) {
		for key := 4; key < 8; key++ {
			if i.g.inputs[key] {
				val = util.SetBit(val, key-4, false)
				pressed = true
			}
		}
	}

	if pressed {
		i.g.requestInterrupt(4)
	}

	i.val = val
	return val
}

func (i *Input) Write(addr uint16, val uint8) {
	i.val = val | 0x0F
}
