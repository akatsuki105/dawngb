package gb

import (
	"image/color"

	"github.com/akatsuki105/dugb/core/gb/cpu"
	"github.com/akatsuki105/dugb/core/gb/video"
	"github.com/akatsuki105/dugb/util/scheduler"
)

type GB struct {
	cpu       *cpu.Cpu
	m         Memory
	video     *video.Video
	s         *scheduler.Scheduler
	cartridge *cartridge
}

func New() *GB {
	s := scheduler.New()
	g := &GB{
		video: video.New(s),
		s:     s,
	}
	g.m = *newMemory(g)
	g.cpu = cpu.New(s, &g.m)
	return g
}

func (g *GB) ID() string {
	return "GB"
}

func (g *GB) Reset() {
	g.cpu.Reset()
	g.video.Reset()
}

func (g *GB) LoadROM(romData []byte) error {
	g.loadCartridge(romData)
	g.Reset()
	return nil
}

func (g *GB) RunFrame() {
	const FRAME = 70224 * video.CYCLE
	start := g.s.Cycle()

	frame := g.video.FrameCounter
	for frame == g.video.FrameCounter && ((g.s.Cycle() - start) < FRAME) {
		g.run()
	}
}

func (g *GB) run() {
	for g.cpu.Cycles < g.cpu.NextEvent {
		g.cpu.Step()
	}
	g.cpu.ProcessEvents()
}

func (g *GB) Resolution() (w int, h int) {
	return 160, 144
}

func (g *GB) FrameBuffer() []color.RGBA {
	return g.video.FrameBuffer()
}

func (g *GB) SetKeyInput(key string, press bool) {
}

func (g *GB) Title() string {
	if g.cartridge == nil {
		return ""
	}
	return g.cartridge.title
}
