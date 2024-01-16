package video

import (
	"image"
	"image/color"

	"github.com/akatsuki105/dugb/core/gb/video/renderer"
	"github.com/akatsuki105/dugb/util"
	. "github.com/akatsuki105/dugb/util/datasize"
)

const CYCLE = 2

type VRAM struct {
	data [16 * KB]uint8
	bank int
}

type Video struct {
	cycles          int64 // 遅れているサイクル数(8.38MHzのマスターサイクル単位)
	dot             int
	screen          [160 * 144]color.RGBA
	FrameCounter    uint64
	ly              int
	r               renderer.Renderer
	renderingCycle  int64
	ram             VRAM
	lcdc, stat, lyc uint8
	onInterrupt     func(id int)
	oam             [160]uint8
	ioreg           [0x30]uint8
}

func New(onInterrupt func(id int)) *Video {
	v := &Video{
		onInterrupt: onInterrupt,
		stat:        0x80,
	}
	v.r = renderer.New("dummy", v.ram.data[:], v.oam[:], 0)
	return v
}

func (v *Video) Reset(model int, hasBIOS bool) {
	v.r = renderer.New("software", v.ram.data[:], v.oam[:], model)
	v.ly, v.dot = 0, 0
	v.stat = 0x80
	v.ram.bank = 0
	if !hasBIOS {
		v.skipBIOS()
	}
}

func (v *Video) skipBIOS() {
	v.Write(0xFF40, 0x91) // LCDC
	v.Write(0xFF47, 0xFC) // BGP
}

func (v *Video) Screen() []color.RGBA {
	return v.screen[:]
}

func (v *Video) Add(cycles int64) {
	v.cycles += cycles
}

func (v *Video) CatchUp() {
	dotCycles := v.cycles / 2 // 1dot = 4.19MHz, 1マスターサイクル = 8.38MHz

	for i := 0; i < int(dotCycles); i++ {
		if util.Bit(v.lcdc, 7) {
			if v.ly < 144 {
				switch v.dot {
				case 0:
					v.scanOAM()
				case 80:
					v.drawing()
				case 252:
					v.hblank()
				}
			}
			v.dot++
			if v.dot == 456 {
				v.dot = 0
				v.incrementLY()
			}
		}
	}

	v.cycles -= dotCycles * 2
}

func (v *Video) incrementLY() {
	v.ly++
	switch v.ly {
	case 144:
		v.vblank()
	case 154:
		v.ly = 0
		v.FrameCounter++
	}
	v.compareLYC()
}

func (v *Video) compareLYC() {
	oldStat := v.stat
	v.stat = util.SetBit(v.stat, 2, v.ly == int(v.lyc))
	if !statIRQAsserted(oldStat) && statIRQAsserted(v.stat) {
		v.onInterrupt(1)
	}
}

func (v *Video) Debug() image.Image {
	return v.r.Debug()
}

func statIRQAsserted(stat byte) bool {
	if util.Bit(stat, 6) && util.Bit(stat, 2) {
		return true
	}
	switch stat & 0b11 {
	case 0:
		return util.Bit(stat, 3)
	case 1:
		return util.Bit(stat, 4)
	case 2:
		return util.Bit(stat, 5)
	}
	return false
}
