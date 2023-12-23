package video

import (
	"image/color"

	"github.com/akatsuki105/dugb/core/gb/video/renderer"
	"github.com/akatsuki105/dugb/util/scheduler"
)

const CYCLE = 1

type Video struct {
	screen         [160 * 144]color.RGBA
	s              *scheduler.Scheduler
	FrameCounter   uint64
	ly             int
	r              renderer.Renderer
	mode           int
	renderingCycle int64
	events         [4]scheduler.Event
}

func New(s *scheduler.Scheduler) *Video {
	v := &Video{
		s: s,
		r: renderer.New("software"),
	}
	v.events = [4]scheduler.Event{
		*scheduler.NewEvent("GB_HBLANK", v.hblank, 0x10),
		*scheduler.NewEvent("GB_VBLANK", v.vblank, 0x11),
		*scheduler.NewEvent("GB_SCAN_OAM", v.scanOAM, 0x12),
		*scheduler.NewEvent("GB_DRAWING", v.drawing, 0x13),
	}
	return v
}

func (v *Video) Reset() {
	v.ly = -1
	v.scanOAM(0)
}

func (v *Video) FrameBuffer() []color.RGBA {
	return v.screen[:]
}

func (v *Video) Scanline() uint8 {
	return uint8(v.ly)
}
