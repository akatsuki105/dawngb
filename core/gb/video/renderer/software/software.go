package software

import (
	"image"
	"image/color"

	"github.com/akatsuki105/dugb/util"
)

const (
	Z_BD = iota
	Z_BG
	Z_SPR
)

var dmgPalette = [4]color.RGBA{
	{0xE0, 0xF8, 0xCF, 0xFF},
	{0x86, 0xC0, 0x6C, 0xFF},
	{0x30, 0x68, 0x50, 0xFF},
	{0x07, 0x18, 0x21, 0xFF},
}

var screen = image.NewRGBA(image.Rect(0, 0, 256, 256))

type Software struct {
	vram     []uint8
	oam      []uint8
	scanline [160]pixel
	bg       *bgLayer
	win      *windowLayer
	sprite   *spriteLayer
	model    int // 0: DMG, 1: CGB
}

type pixel struct {
	rgba    color.RGBA
	z       int // z-index
	colorID int
}

func New(vram, oam []uint8, model int) *Software {
	r := &Software{
		vram:  vram,
		oam:   oam,
		model: model,
	}
	r.bg = newBG(r)
	r.win = newWindow(r)
	r.sprite = newSpriteLayer(r)
	return r
}

func (s *Software) DrawScanline(y int, scanline []color.RGBA) {
	for i := 0; i < 160; i++ {
		s.scanline[i].z = -1
		s.scanline[i].colorID = 0
	}

	s.bg.drawScanline(y, s.scanline[:])
	s.win.drawScanline(y, s.scanline[:])
	s.sprite.drawScanline(y, s.scanline[:])

	for i := 0; i < 160; i++ {
		scanline[i] = s.scanline[i].rgba
	}
}

func (s *Software) SetLCDC(val uint8) {
	s.bg.active = util.Bit(val, 0)
	s.bg.tilemap = [2]uint16{0x1800, 0x1C00}[util.Btoi(util.Bit(val, 3))]
	s.bg.tiledata = [2]int{0x800, 0x0}[util.Btoi(util.Bit(val, 4))]

	s.win.active = s.bg.active && util.Bit(val, 5)
	s.win.tilemap = [2]uint16{0x1800, 0x1C00}[util.Btoi(util.Bit(val, 6))]

	s.sprite.active = util.Bit(val, 1)
	s.sprite.height = [2]int{8, 16}[util.Btoi(util.Bit(val, 2))]
}

func (s *Software) SetBGP(val uint8) {
	s.bg.palette[0] = dmgPalette[val&0b11]
	s.bg.palette[1] = dmgPalette[(val>>2)&0b11]
	s.bg.palette[2] = dmgPalette[(val>>4)&0b11]
	s.bg.palette[3] = dmgPalette[(val>>6)&0b11]
}

func (s *Software) SetOBP0(val uint8) {
	s.sprite.palette[0] = dmgPalette[val&0b11]
	s.sprite.palette[1] = dmgPalette[(val>>2)&0b11]
	s.sprite.palette[2] = dmgPalette[(val>>4)&0b11]
	s.sprite.palette[3] = dmgPalette[(val>>6)&0b11]
}

func (s *Software) SetOBP1(val uint8) {
	s.sprite.palette[4] = dmgPalette[val&0b11]
	s.sprite.palette[5] = dmgPalette[(val>>2)&0b11]
	s.sprite.palette[6] = dmgPalette[(val>>4)&0b11]
	s.sprite.palette[7] = dmgPalette[(val>>6)&0b11]
}

func (s *Software) SetSCX(val uint8) { s.bg.scx = int(val) }
func (s *Software) SetSCY(val uint8) { s.bg.scy = int(val) }

func (s *Software) SetWX(val uint8) { s.win.wx = int(val) - 7 }
func (s *Software) SetWY(val uint8) { s.win.wy = int(val) }

func (s *Software) SetBGPI(val uint8) { s.bg.bgpi = val }
func (s *Software) SetBGPD(val uint8) {
	if s.model == 1 {
		palID := int((s.bg.bgpi & 0x3F) / 8)
		colorID := int(s.bg.bgpi&7) >> 1
		rgba := &s.bg.palette[palID*4+colorID]
		isHi := util.Bit(s.bg.bgpi, 0)
		// val is rgb555 format
		if isHi {
			// 0b0BBBBBGG
			rgba.G = (((val & 0b11) << 3) | ((rgba.G >> 3) & 0b111)) << 3
			rgba.B = ((val >> 2) & 0b11111) << 3
		} else {
			// 0bGGGRRRRR
			rgba.R = (val & 0b11111) << 3
			rgba.G = (((rgba.G >> 3) & 0b11000) | ((val >> 5) & 0b111)) << 3
		}
		rgba.A = 0xFF

		if util.Bit(s.bg.bgpi, 7) {
			s.bg.bgpi++
		}
	}
}

func (s *Software) SetOBPI(val uint8) { s.sprite.obpi = val }
func (s *Software) SetOBPD(val uint8) {
	if s.model == 1 {
		palID := int((s.sprite.obpi & 0x3F) / 8)
		colorID := int(s.sprite.obpi&7) >> 1
		rgba := &s.sprite.palette[palID*4+colorID]
		isHi := util.Bit(s.sprite.obpi, 0)
		// val is rgb555 format
		if isHi {
			// 0b0BBBBBGG
			rgba.G = (((val & 0b11) << 3) | ((rgba.G >> 3) & 0b111)) << 3
			rgba.B = ((val >> 2) & 0b11111) << 3
		} else {
			// 0bGGGRRRRR
			rgba.R = (val & 0b11111) << 3
			rgba.G = (((rgba.G >> 3) & 0b11000) | ((val >> 5) & 0b111)) << 3
		}
		rgba.A = 0xFF

		if util.Bit(s.sprite.obpi, 7) {
			s.sprite.obpi++
		}
	}
}

func (s *Software) Debug() image.Image {
	for row := 0; row < 12; row++ {
		for col := 0; col < 32; col++ {
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					tileID := row*32 + col
					tile := s.vram[tileID*16 : (tileID+1)*16]
					for y := 0; y < 8; y++ {
						for x := 0; x < 8; x++ {
							planes := [2]uint8{tile[(y&0b111)*2], tile[(y&0b111)*2+1]}
							lo := (planes[0] >> (7 - uint(x))) & 0b1
							hi := (planes[1] >> (7 - uint(x))) & 0b1
							colorID := int((hi << 1) | lo)
							screen.Set(col*8+x, row*8+y, s.bg.palette[colorID])
						}
					}
				}
			}
		}
	}
	return screen
}
