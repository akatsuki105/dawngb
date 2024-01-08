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
	cycles                                             int64 // 遅れているサイクル数(8.3MHzのマスターサイクル単位)
	dot                                                int
	screen                                             [160 * 144]color.RGBA
	FrameCounter                                       uint64
	ly                                                 int
	r                                                  renderer.Renderer
	renderingCycle                                     int64
	ram                                                VRAM
	lcdc, stat, lyc, scx, scy, wx, wy, bgp, obp0, obp1 uint8
	onInterrupt                                        func(id int)
	OAM                                                [160]uint8
	ioreg                                              [0x40]uint8
}

func New(onInterrupt func(id int)) *Video {
	v := &Video{
		onInterrupt: onInterrupt,
	}
	v.r = renderer.New("dummy", v.ram.data[:], v.OAM[:], 0)
	return v
}

func (v *Video) Reset(model int, hasBIOS bool) {
	v.r = renderer.New("software", v.ram.data[:], v.OAM[:], model)
	v.ly = 0
	v.ram.bank = 0
	v.scanOAM(0)
	if !hasBIOS {
		v.skipBIOS()
	}
}

func (v *Video) skipBIOS() {
	v.lcdc = 0x91
	v.r.SetLCDC(v.lcdc)
}

func (v *Video) Screen() []color.RGBA {
	return v.screen[:]
}

func (v *Video) VRAM() []uint8 {
	bank := uint(v.ram.bank) * (8 * KB)
	return v.ram.data[bank : bank+(8*KB)]
}

func (v *Video) Add(cycles int64) {
	v.cycles += cycles
}

func (v *Video) CatchUp() {
	dotCycles := v.cycles / 2 // 1dot = 4.19MHz, 1マスターサイクル = 8.3MHz

	for i := 0; i < int(dotCycles); i++ {
		if v.ly < 144 {
			switch v.dot {
			case 0:
				v.scanOAM(0)
			case 80:
				v.drawing(0)
			case 252:
				v.hblank(0)
			}
		} else if v.ly == 144 {
			if v.dot == 0 {
				v.vblank(0)
			}
		}
		v.dot++
		if v.dot == 456 {
			v.dot = 0
			v.setLy(v.ly + 1)
		}
	}

	v.cycles -= dotCycles * 2
}

func (v *Video) setLy(ly int) {
	if ly == 154 {
		ly = 0
		v.FrameCounter++
	}
	v.ly = ly
	v.stat = util.SetBit(v.stat, 2, ly == int(v.lyc))
	if ly == int(v.lyc) && util.Bit(v.stat, 6) {
		v.onInterrupt(1)
	}
}

func (v *Video) Debug() image.Image {
	return v.r.Debug()
}
