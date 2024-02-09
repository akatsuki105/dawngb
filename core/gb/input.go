package gb

import (
	"io"

	"github.com/akatsuki105/dawngb/util"
)

type Input struct {
	g        *GB
	p14, p15 bool
	joyp     uint8
}

func newInput(g *GB) *Input {
	return &Input{
		g: g,
	}
}

func (i *Input) Reset(hasBIOS bool) {
	i.p14, i.p15 = false, false
	i.joyp = 0x0F
	if !hasBIOS {
		i.Write(0xFF00, 0x30)
		i.Write(0xFF00, 0xCF)
	}
}

func (i *Input) poll() {
	dpad := uint8(0x0)
	dpad = util.SetBit(dpad, 0, !i.g.inputs[4+0])
	dpad = util.SetBit(dpad, 1, !i.g.inputs[4+1])
	dpad = util.SetBit(dpad, 2, !i.g.inputs[4+2])
	dpad = util.SetBit(dpad, 3, !i.g.inputs[4+3])

	button := uint8(0x0)
	button = util.SetBit(button, 0, !i.g.inputs[0])
	button = util.SetBit(button, 1, !i.g.inputs[1])
	button = util.SetBit(button, 2, !i.g.inputs[2])
	button = util.SetBit(button, 3, !i.g.inputs[3])

	i.joyp = 0x0F
	if !i.p14 {
		i.joyp &= dpad
	}
	if !i.p15 {
		i.joyp &= button
	}

	if i.joyp != 0x0F {
		i.g.requestInterrupt(4)
	}
}

func (i *Input) Read(addr uint16) uint8 {
	i.poll()
	val := i.joyp
	val = util.SetBit(val, 4, i.p14)
	val = util.SetBit(val, 5, i.p15)

	/*
		NOTE:
			ポケモンの赤・緑・青・ピカチュウ版では、上位4ビットが0xCになっている。
			そうしないと、ポケモンのゲームはハードウェアをスーパゲームボーイとして認識し、上入力が効かなくなる。
	*/
	val |= 0xC0
	return val
}

func (i *Input) Write(addr uint16, val uint8) {
	i.p14 = util.Bit(val, 4)
	i.p15 = util.Bit(val, 5)
}

func (i *Input) Serialize(s io.Writer) {
	s.Write([]byte{util.Btou8(i.p14), util.Btou8(i.p15), i.joyp})
}

func (i *Input) Deserialize(s io.Reader) {
	buf := make([]byte, 3)
	s.Read(buf)
	i.p14, i.p15 = buf[0] != 0, buf[1] != 0
	i.joyp = buf[2]
}
