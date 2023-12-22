package video

import (
	"image/color"

	"github.com/akatsuki105/dugb/util/scheduler"
)

type Video struct {
	screen [160 * 144]color.RGBA
	s      *scheduler.Scheduler
}

func New(s *scheduler.Scheduler) *Video {
	return &Video{
		s: s,
	}
}

func (v *Video) FrameBuffer() []color.RGBA {
	return v.screen[:]
}
