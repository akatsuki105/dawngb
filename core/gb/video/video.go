package video

import (
	"image/color"

	"github.com/akatsuki105/dawngb/core/gb/video/renderer"
	"github.com/akatsuki105/dawngb/util"
)

const KB = 1024

const CYCLE = 2

type VRAM struct {
	data [16 * KB]uint8
	bank uint
}

type Video struct {
	cycles          int64 // 遅れているサイクル数(8.38MHzのマスターサイクル単位)
	screen          [160 * 144]color.RGBA
	FrameCounter    uint64
	lx, ly          int
	r               renderer.Renderer
	ram             VRAM
	lcdc, stat, lyc uint8
	onInterrupt     func(id int)
	onHBlank        func()
	oam             [160]uint8
	ioreg           [0x30]uint8
	enableLatch     bool // LCDC.7をセットしてPPUを有効にすると、次のフレームから表示が開始される そうじゃないとゴミが表示される
	objCount        int
}

func New(onInterrupt func(id int), onHBlank func()) *Video {
	v := &Video{
		onInterrupt: onInterrupt,
		onHBlank:    onHBlank,
		stat:        0x80,
	}
	v.r = renderer.New("dummy", v.ram.data[:], v.oam[:], 0)
	return v
}

func (v *Video) Reset(model int, hasBIOS bool) {
	v.r = renderer.New("software", v.ram.data[:], v.oam[:], model)
	v.lx, v.ly = 0, 0
	v.stat = 0x80
	v.ram.bank = 0
	v.objCount = 0
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

func (v *Video) Tick(cycles int64) {
	v.cycles += cycles
}

func (v *Video) CatchUp() {
	dotCycles := v.cycles / 2 // 1dot = 4.19MHz, 1マスターサイクル = 8.38MHz

	for i := 0; i < int(dotCycles); i++ {
		if util.Bit(v.lcdc, 7) {
			if v.ly < 144 {
				switch v.lx {
				case 0:
					v.scanOAM()
				case 80:
					v.drawing()
				case 252 + (v.objCount * 6):
					v.hblank()
				}
			}
			v.lx++
			if v.lx == 456 {
				v.lx = 0
				v.incrementLY()
			}
		}
	}

	v.cycles -= dotCycles * 2
}

func (v *Video) incrementLY() {
	v.objCount = 0
	v.ly++
	switch v.ly {
	case 144:
		v.vblank()
	case 154:
		v.ly = 0
		v.enableLatch = false
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
