package gb

import (
	"image/color"

	"github.com/akatsuki105/dugb/core/gb/cpu"
	"github.com/akatsuki105/dugb/core/gb/video"
	"github.com/akatsuki105/dugb/util/scheduler"
)

type GB struct {
	cpu   *cpu.Cpu
	m     Memory
	video *video.Video
	s     *scheduler.Scheduler
}

func New() *GB {
	s := scheduler.New()
	g := &GB{
		cpu:   cpu.New(s),
		video: video.New(s),
		s:     s,
	}
	g.m = *newMemory(g)
	return g
}

func (g *GB) ID() string {
	return "GB"
}

func (g *GB) LoadROM(romData []byte) error {
	return nil
}

func (g *GB) RunFrame() {
}

func (g *GB) Resolution() (w int, h int) {
	return 160, 144
}

func (g *GB) FrameBuffer() []color.RGBA {
	return g.video.FrameBuffer()
}

func (g *GB) SetKeyInput(key string, press bool) {
}
