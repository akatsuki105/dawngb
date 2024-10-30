package gb

import (
	"fmt"
	"image/color"
	"io"

	"github.com/akatsuki105/dawngb/core/gb/apu"
	"github.com/akatsuki105/dawngb/core/gb/cartridge"
	"github.com/akatsuki105/dawngb/core/gb/cpu"
	"github.com/akatsuki105/dawngb/core/gb/video"
	"github.com/akatsuki105/dawngb/util"
)

const KB, MB = 1024, 1024 * 1024

var buttons = [8]string{"A", "B", "SELECT", "START", "RIGHT", "LEFT", "UP", "DOWN"}

type oamDmaController struct {
	active bool
	src    uint16
	until  int64
}

type GB struct {
	cycles    int64 // 8.3MHzのマスターサイクル単位
	cpu       *cpu.Cpu
	m         *Memory
	video     *video.Video
	cartridge *cartridge.Cartridge
	input     *input
	timer     *timer
	serial    *serial
	dmac      *dmaController
	ie        uint8
	interrupt [5]bool // IF
	halted    bool
	key1      bool    // FF4D's bit 0
	inputs    [8]bool // A, B, Select, Start, Right, Left, Up, Down
	runHDMA   func()
	oamDMA    oamDmaController
	apu       *apu.APU
}

func New(audioBuffer io.Writer) *GB {
	g := &GB{}
	g.m = newMemory(g)
	g.cpu = cpu.New(g.m, g.halt, g.stop, g.tick)
	g.video = video.New(g.requestInterrupt, g.triggerHDMA)
	g.timer = newTimer(g)
	g.apu = apu.New(audioBuffer)
	g.input = newInput(g)
	g.dmac = newDMAController(g)
	g.serial = newSerial(g)
	return g
}

func (g *GB) Reset(hasBIOS bool) {
	g.ie, g.interrupt = 0, [5]bool{}
	g.halted, g.key1, g.oamDMA.active = false, false, false

	model := 0
	if g.cartridge != nil && g.cartridge.IsCGB() {
		model = 1
	}

	g.m.Reset(hasBIOS)
	g.cpu.Reset(hasBIOS)
	g.video.Reset(model, hasBIOS)
	g.apu.Reset(hasBIOS)
	g.timer.Reset(hasBIOS)
	g.input.Reset(hasBIOS)
	g.dmac.Reset(hasBIOS)
	g.serial.Reset(hasBIOS)

	if !hasBIOS {
		g.m.Write(0xFF02, 0x7F)
		g.m.Write(0xFF0F, 0xE1)
		if model == 1 {
			g.m.Write(0xFF4D, 0x7E)
			g.m.Write(0xFF4F, 0xFE)
		}
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
	if g.cartridge != nil {
		const FRAME = 70224 * video.CYCLE
		start := g.cycles

		frame := g.video.FrameCounter
		for frame == g.video.FrameCounter && ((g.cycles - start) < FRAME) {
			g.run()
			g.video.CatchUp()
		}
		g.video.CatchUp()
		g.apu.FlushSamples()
	}
}

func (g *GB) run() {
	prev := g.cycles

	if g.dmac.doHDMA {
		g.dmac.doHDMA = false
		g.dmac.runHDMA()
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
			g.tick(1)
		} else {
			g.cpu.Step()
		}
	}

	g.catchUp(g.cycles - prev)
}

func (g *GB) tick(cycles int64) {
	g.cycles += cycles
}

func (g *GB) catchUp(cycles int64) {
	g.apu.Tick(cycles)
	g.video.Tick(cycles)
	g.timer.tick(cycles)
	g.serial.tick(cycles)

	if g.oamDMA.active {
		g.oamDMA.until -= cycles
		if g.oamDMA.until <= 0 {
			for i := uint16(0); i < 160; i++ {
				g.video.Write(0xFE00+i, g.m.Read(g.oamDMA.src+i))
			}
			g.oamDMA.active = false
		}
	}
}

func (g *GB) Resolution() (w int, h int) { return 160, 144 }

func (g *GB) Screen() []color.RGBA {
	if g.cartridge != nil {
		return g.video.Screen()
	}
	return []color.RGBA{}
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

func (g *GB) requestInterrupt(id int) { g.interrupt[id] = true }

func (g *GB) checkInterrupt() int {
	for i := 0; i < 5; i++ {
		if util.Bit(g.ie, i) && g.interrupt[i] {
			return i
		}
	}
	return -1
}

func (g *GB) triggerOAMDMA(src uint16) {
	if !g.oamDMA.active {
		g.oamDMA.active = true
		g.oamDMA.src = src
		g.oamDMA.until = 160 * g.cpu.Cycle
	}
}

func (g *GB) triggerHDMA() {
	if g.runHDMA != nil {
		g.dmac.doHDMA = true
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

func (g *GB) Serialize(state io.Writer) {
	// TODO: implement
	g.input.Serialize(state)
	g.timer.Serialize(state)
	g.dmac.Serialize(state)
}

func (g *GB) Deserialize(state io.Reader) {
	// TODO: implement
	g.input.Deserialize(state)
	g.timer.Deserialize(state)
	g.dmac.Deserialize(state)
}
