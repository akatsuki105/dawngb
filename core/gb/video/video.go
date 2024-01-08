package video

import (
	"image/color"

	"github.com/akatsuki105/dugb/core/gb/video/renderer"
	"github.com/akatsuki105/dugb/util"
	. "github.com/akatsuki105/dugb/util/datasize"
	"github.com/akatsuki105/dugb/util/sched"
)

const CYCLE = 2

type VRAM struct {
	data [16 * KB]uint8
	bank int
}

type Video struct {
	screen                                             [160 * 144]color.RGBA
	s                                                  *sched.Sched
	FrameCounter                                       uint64
	ly                                                 int
	r                                                  renderer.Renderer
	renderingCycle                                     int64
	events                                             [4]sched.Event
	ram                                                VRAM
	lcdc, stat, lyc, scx, scy, wx, wy, bgp, obp0, obp1 uint8
	onInterrupt                                        func(id int)
	OAM                                                [160]uint8
	ioreg                                              [0x40]uint8
}

func New(s *sched.Sched, onInterrupt func(id int)) *Video {
	v := &Video{
		s:           s,
		onInterrupt: onInterrupt,
	}
	v.r = renderer.New("dummy", v.ram.data[:], v.OAM[:], 0)
	v.events = [4]sched.Event{
		*sched.NewEvent("GB_HBLANK", v.hblank),
		*sched.NewEvent("GB_VBLANK", v.vblank),
		*sched.NewEvent("GB_SCAN_OAM", v.scanOAM),
		*sched.NewEvent("GB_DRAWING", v.drawing),
	}
	return v
}

func (v *Video) Reset(model int, hasBIOS bool) {
	v.r = renderer.New("software", v.ram.data[:], v.OAM[:], model)
	v.ly = -1
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
