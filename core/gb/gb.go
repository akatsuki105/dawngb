package gb

import (
	"fmt"
	"image/color"
	"io"

	"github.com/akatsuki105/dugb/core/gb/audio"
	"github.com/akatsuki105/dugb/core/gb/cartridge"
	"github.com/akatsuki105/dugb/core/gb/cpu"
	"github.com/akatsuki105/dugb/core/gb/video"
	"github.com/akatsuki105/dugb/util"
	"github.com/akatsuki105/dugb/util/sched"
)

var buttons = [8]string{"A", "B", "SELECT", "START", "RIGHT", "LEFT", "UP", "DOWN"}

type peripheral interface {
	Reset(hasBIOS bool)
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type GB struct {
	cpu       *cpu.Cpu
	m         *Memory
	video     *video.Video
	s         *sched.Sched
	cartridge *cartridge.Cartridge
	input     peripheral
	timer     peripheral
	audio     audio.Audio
	ie        uint8
	interrupt [5]bool // IF
	dma       sched.Event
	halted    bool
	blocked   bool // DMA
	key1      bool // FF4D's bit 0
	inOAMDMA  bool
	sb, sc    uint8
	inputs    [8]bool // A, B, Select, Start, Right, Left, Up, Down
	dmac      peripheral
	runHDMA   func()
	serial    sched.Event
}

func New(audioBuffer io.Writer) *GB {
	s := sched.New()
	g := &GB{
		s:   s,
		dma: *sched.NewEvent("GB_DMA", func(cycle int64) {}),
	}
	g.m = newMemory(g)
	g.cpu = cpu.New(s, g.m, g.halt, g.stop)
	g.video = video.New(g.requestInterrupt, g.triggerHDMA)
	g.timer = newTimer(g)
	g.audio = audio.New(audioBuffer)
	g.input = newInput(g)
	g.dmac = newDMAController(g)
	g.serial = *sched.NewEvent("GB_SERIAL", g.dummyTransfer)
	return g
}

func (g *GB) ID() string {
	return "GB"
}

func (g *GB) Reset(hasBIOS bool) {
	model := 0
	if g.cartridge != nil && g.cartridge.IsCGB() {
		model = 1
	}
	g.s.Reset()
	g.m.Reset(hasBIOS)
	g.cpu.Reset(hasBIOS)
	g.video.Reset(model, hasBIOS)
	g.audio.Reset(hasBIOS)
	g.timer.Reset(hasBIOS)
	g.input.Reset(hasBIOS)
	g.dmac.Reset(hasBIOS)

	if !hasBIOS {
		g.m.Write(0xFF02, 0x7F)
		g.m.Write(0xFF0F, 0xE1)
		g.m.Write(0xFF4D, 0xFF)
		g.m.Write(0xFF4F, 0xFF)
	}
}

func (g *GB) LoadROM(romData []byte) error {
	g.cartridge = cartridge.New(romData)
	g.Reset(false)

	return nil
}

func (g *GB) LoadSRAM(data []byte) error {
	if g.cartridge == nil {
		return fmt.Errorf("no cartridge loaded")
	}
	err := g.cartridge.LoadSRAM(data)
	if err != nil {
		return err
	}
	g.Reset(false)
	return nil
}

func (g *GB) SRAM() []byte {
	if g.cartridge == nil {
		return nil
	}
	return g.cartridge.SRAM()
}

func (g *GB) RunFrame() {
	const FRAME = 70224 * video.CYCLE
	start := g.s.Cycle()

	frame := g.video.FrameCounter
	for frame == g.video.FrameCounter && ((g.s.Cycle() - start) < FRAME) {
		g.run()
		g.video.CatchUp()
	}
	g.audio.CatchUp()
	g.video.CatchUp()
}

func (g *GB) run() {
	prev := g.s.Cycle()

	if g.blocked {
		g.s.Add(max(g.s.UntilNextEvent(), 1))
	} else {
		irqID := g.checkInterrupt()
		if irqID >= 0 {
			g.halted = false
			if g.cpu.IME {
				g.interrupt[irqID] = false
				g.cpu.Interrupt(irqID)
			} else {
				g.cpu.Step()
			}
		} else if g.halted {
			g.s.Add(max(g.s.UntilNextEvent(), 1))
		} else {
			g.cpu.Step()
		}
	}

	g.audio.Add(g.s.Cycle() - prev)
	g.video.Add(g.s.Cycle() - prev)

	g.s.Commit()
}

func (g *GB) Resolution() (w int, h int) {
	return 160, 144
}

func (g *GB) Screen() []color.RGBA {
	return g.video.Screen()
}

func (g *GB) SetKeyInput(key string, press bool) {
	for i, b := range buttons {
		if b == key {
			g.inputs[i] = press
		}
	}
}

func (g *GB) Title() string {
	if g.cartridge == nil {
		return ""
	}
	return g.cartridge.Title()
}

func (g *GB) requestInterrupt(id int) {
	g.interrupt[id] = true
}

func (g *GB) checkInterrupt() int {
	for i := 0; i < 5; i++ {
		if util.Bit(g.ie, i) && g.interrupt[i] {
			return i
		}
	}
	return -1
}

func (g *GB) triggerOAMDMA(src uint16) {
	g.dma.Callback = func(cyclesLate int64) {
		for i := uint16(0); i < 160; i++ {
			g.video.Write(0xFE00+i, g.m.Read(src+i))
		}
		g.inOAMDMA = false
	}
	g.inOAMDMA = true
	g.s.Schedule(&g.dma, 160*g.cpu.Cycle)
}

func (g *GB) triggerHDMA() {
	if g.runHDMA != nil {
		g.runHDMA()
	}
}

func (g *GB) halt() {
	if g.cpu.IME {
		g.halted = true
	} else {
		if g.checkInterrupt() < 0 {
			g.halted = true
		}
	}
}

func (g *GB) stop() {
	if g.key1 {
		if g.cpu.Cycle == 4 {
			g.cpu.Cycle = 8
		} else {
			g.cpu.Cycle = 4
		}
		g.key1 = false
	}
}

func (g *GB) dummyTransfer(cyclesLate int64) {
	// ポケモンクリスタルの起動にシリアル通信機能が必要なので暫定措置
	g.sc &= 0x7F
	g.sb = 0xFF
	g.requestInterrupt(3)
}
