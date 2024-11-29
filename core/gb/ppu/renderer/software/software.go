package software

import (
	"image/color"

	"golang.org/x/exp/constraints"
)

const KB = 1024

type rgb555 = uint16 // 0b0_BBBBB_GGGGG_RRRRR

type Software struct {
	isCGB  func() bool // CGBモードかどうか (ハードがCGBでもDMGのゲームをする場合はfalse)
	vram   []uint8
	oam    []uint8
	bg     *bgLayer
	win    *windowLayer
	sprite *spriteLayer
}

type pixel struct {
	color    rgb555
	colorID  uint8
	priority bool
}

func New(vram []uint8, palette []rgb555, oam []uint8, isCGB func() bool) *Software {
	r := &Software{
		isCGB: isCGB,
		vram:  vram,
		oam:   oam,
	}
	r.bg = newBG(r, palette[:32])
	r.win = newWindow(r)
	r.sprite = newSpriteLayer(r, palette[32:])
	return r
}

func (s *Software) DrawScanline(y int, scanline []color.NRGBA) {
	if y == 0 {
		s.win.y = 0
	}
	for i := 0; i < 160; i++ {
		s.bg.scanline[i].colorID = 0
		s.bg.scanline[i].priority = false
		s.sprite.scanline[i].colorID = 0
		s.sprite.scanline[i].priority = false
	}

	s.bg.drawScanline(y)
	s.win.drawScanline(y)
	s.sprite.drawScanline(y)

	for i := 0; i < 160; i++ {
		rgb555 := s.mergeLayers(i)
		r5, g5, b5 := uint8(rgb555&0x1F), uint8((rgb555>>5)&0x1F), uint8((rgb555>>10)&0x1F)
		scanline[i].R = (r5 << 3) | (r5 >> 2)
		scanline[i].G = (g5 << 3) | (g5 >> 2)
		scanline[i].B = (b5 << 3) | (b5 >> 2)
		scanline[i].A = 0xFF
	}
}

// Merge BG and Object layers
func (s *Software) mergeLayers(x int) rgb555 {
	c := uint16(0x7FFF)
	bg, obj := &s.bg.scanline[x], &s.sprite.scanline[x]
	if obj.colorID == 0 {
		c = bg.color
	} else if bg.colorID == 0 {
		c = obj.color
	} else {
		if s.isCGB() {
			if !s.bg.enable {
				c = obj.color
			} else if bg.priority {
				c = bg.color
			} else if obj.priority {
				c = obj.color
			} else {
				c = bg.color
			}
		} else {
			if obj.priority {
				c = obj.color
			} else {
				c = bg.color
			}
		}
	}
	return c
}

func (s *Software) SetLCDC(val uint8) {
	s.bg.enable = (val & (1 << 0)) != 0
	s.bg.tilemap = [2]uint16{0x1800, 0x1C00}[(val>>3)&1]
	s.bg.tiledata = [2]int{0x800, 0x0}[(val>>4)&1]

	s.win.enable = ((val & (1 << 5)) != 0)
	s.win.tilemap = [2]uint16{0x1800, 0x1C00}[(val>>6)&1]

	s.sprite.enable = (val & (1 << 1)) != 0
	s.sprite.height = [2]uint8{8, 16}[(val>>2)&1]
}

func (s *Software) SetBGP(val uint8)       { s.bg.bgp = val }
func (s *Software) SetOBP(bank, val uint8) { s.sprite.obp[bank] = val }
func (s *Software) SetSCX(val uint8)       { s.bg.scx = val }
func (s *Software) SetSCY(val uint8)       { s.bg.scy = val }
func (s *Software) SetWX(val uint8)        { s.win.wx = int(val) - 7 }
func (s *Software) SetWY(val uint8)        { s.win.wy = int(val) }

// Helper functions
func flip[V constraints.Integer, W constraints.Integer](size V, b bool, i W) V {
	if b {
		return V(int(size-1) - int(i))
	}
	return V(i)
}
