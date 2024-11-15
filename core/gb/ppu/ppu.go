package ppu

import (
	"image/color"

	"github.com/akatsuki105/dawngb/core/gb/ppu/renderer"
	"github.com/akatsuki105/dawngb/util"
)

const KB = 1024

const CYCLE = 2

type VRAM struct {
	data [16 * KB]uint8
	bank uint
}

type Bus interface {
	Read(addr uint16) uint8
}

// OAM DMA
type DMA struct {
	active bool
	src    uint16
	until  int64
}

type PPU struct {
	bus             Bus
	cycles          int64 // 遅れているサイクル数(8.38MHzのマスターサイクル単位)
	screen          [160 * 144]color.NRGBA
	FrameCounter    uint64
	lx, ly          int
	r               renderer.Renderer
	ram             VRAM
	DMA             *DMA
	lcdc, stat, lyc uint8
	irq             func(id int)
	onHBlank        func()
	oam             [160]uint8
	ioreg           [0x30]uint8
	enableLatch     bool // LCDC.7をセットしてPPUを有効にすると、次のフレームから表示が開始される そうじゃないとゴミが表示される
	objCount        int
}

func New(bus Bus, irq func(id int), onHBlank func()) *PPU {
	p := &PPU{
		bus:      bus,
		DMA:      &DMA{},
		irq:      irq,
		onHBlank: onHBlank,
		stat:     0x80,
	}
	p.r = renderer.New("dummy", p.ram.data[:], p.oam[:], 0)
	return p
}

func (p *PPU) Reset(model int, hasBIOS bool) {
	p.r = renderer.New("software", p.ram.data[:], p.oam[:], model)
	p.lx, p.ly = 0, 0
	p.stat = 0x80
	p.ram.bank = 0
	p.objCount = 0
	p.DMA.active, p.DMA.src, p.DMA.until = false, 0, 0
	if !hasBIOS {
		p.skipBIOS()
	}
}

func (p *PPU) skipBIOS() {
	p.Write(0xFF40, 0x91) // LCDC
	p.Write(0xFF47, 0xFC) // BGP
}

func (p *PPU) Screen() []color.NRGBA {
	return p.screen[:]
}

func (p *PPU) Run(cycles8MHz int64) {
	if p.DMA.active {
		p.runDMA(cycles8MHz)
	}

	p.cycles += cycles8MHz
	for p.cycles >= 2 { // 1dot = 4MHz
		p.step()
		p.cycles -= 2
	}
}

func (p *PPU) step() {
	if util.Bit(p.lcdc, 7) {
		if p.ly < 144 {
			switch p.lx {
			case 0:
				p.scanOAM()
			case 80:
				p.drawing()
			case 252 + (p.objCount * 6):
				p.hblank()
			}
		}
		p.lx++
		if p.lx == 456 {
			p.lx = 0
			p.incrementLY()
		}
	}
}

func (p *PPU) incrementLY() {
	p.objCount = 0
	p.ly++
	switch p.ly {
	case 144:
		p.vblank()
	case 154:
		p.ly = 0
		p.enableLatch = false
		p.FrameCounter++
	}
	p.compareLYC()
}

func (p *PPU) compareLYC() {
	oldStat := p.stat
	p.stat = util.SetBit(p.stat, 2, p.ly == int(p.lyc))
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.irq(1)
	}
}

func statIRQAsserted(stat uint8) bool {
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

func (p *PPU) runDMA(cycles8MHz int64) {
	p.DMA.until -= cycles8MHz
	if p.DMA.until <= 0 {
		for i := uint16(0); i < 160; i++ {
			p.Write(0xFE00+i, p.bus.Read(p.DMA.src+i))
		}
		p.DMA.active = false
	}
}

func (p *PPU) TriggerDMA(src uint16, m int64) {
	if !p.DMA.active {
		p.DMA.active = true
		p.DMA.src = src
		p.DMA.until = 160 * m
	}
}
