package software

import (
	"image"
	"image/color"

	"github.com/akatsuki105/dawngb/util"
)

const KB = 1024

const (
	Z_BD = iota
	Z_BG
	Z_SPR
	Z_MASTER_SPR = 100 // CGBモードでLCDC.0が0のときは必ずスプライトが前面に来る
)

var dmgPalette = [4]rgb555{
	{0b11111, 0b11111, 0b11111},
	{0b10001, 0b10001, 0b10001},
	{0b01010, 0b01010, 0b01010},
	{0b00000, 0b00000, 0b00000},
}

var screen = image.NewRGBA(image.Rect(0, 0, 256, 256))

type Software struct {
	vram     []uint8
	oam      []uint8
	scanline [160]pixel
	bg       *bgLayer
	win      *windowLayer
	sprite   *spriteLayer
	model    int // 0: DMG, 1: CGB, 2: SGB
}

type pixel struct {
	rgba    rgb555
	z       int // z-index
	colorID int
	isBG    bool
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
	if y == 0 {
		s.win.y = 0
	}
	for i := 0; i < 160; i++ {
		s.scanline[i].z = -1
		s.scanline[i].colorID = 0
		s.scanline[i].isBG = false
	}

	s.bg.drawScanline(y, s.scanline[:])
	s.win.drawScanline(y, s.scanline[:])
	s.sprite.drawScanline(y, s.scanline[:])

	for i := 0; i < 160; i++ {
		scanline[i] = s.scanline[i].rgba.RGBA()
	}
}

func (s *Software) SetLCDC(val uint8) {
	s.bg.active = util.Bit(val, 0)
	if s.model == 1 {
		s.bg.active = true
		if util.Bit(val, 0) {
			s.sprite.z = 0
		} else {
			s.sprite.z = Z_MASTER_SPR
		}
	}
	s.bg.tilemap = [2]uint16{0x1800, 0x1C00}[util.Btoi(util.Bit(val, 3))]
	s.bg.tiledata = [2]int{0x800, 0x0}[util.Btoi(util.Bit(val, 4))]

	s.win.active = s.bg.active && util.Bit(val, 5)
	s.win.tilemap = [2]uint16{0x1800, 0x1C00}[util.Btoi(util.Bit(val, 6))]

	s.sprite.active = util.Bit(val, 1)
	s.sprite.height = [2]int{8, 16}[util.Btoi(util.Bit(val, 2))]
}

func (s *Software) SetBGP(val uint8) {
	if s.model != 1 {
		s.bg.palette[0] = dmgPalette[val&0b11]
		s.bg.palette[1] = dmgPalette[(val>>2)&0b11]
		s.bg.palette[2] = dmgPalette[(val>>4)&0b11]
		s.bg.palette[3] = dmgPalette[(val>>6)&0b11]
	}
}

func (s *Software) SetOBP0(val uint8) {
	if s.model != 1 {
		s.sprite.palette[0] = dmgPalette[val&0b11]
		s.sprite.palette[1] = dmgPalette[(val>>2)&0b11]
		s.sprite.palette[2] = dmgPalette[(val>>4)&0b11]
		s.sprite.palette[3] = dmgPalette[(val>>6)&0b11]
	}
}

func (s *Software) SetOBP1(val uint8) {
	if s.model != 1 {
		s.sprite.palette[4] = dmgPalette[val&0b11]
		s.sprite.palette[5] = dmgPalette[(val>>2)&0b11]
		s.sprite.palette[6] = dmgPalette[(val>>4)&0b11]
		s.sprite.palette[7] = dmgPalette[(val>>6)&0b11]
	}
}

func (s *Software) SetSCX(val uint8) { s.bg.scx = int(val) }
func (s *Software) SetSCY(val uint8) { s.bg.scy = int(val) }

func (s *Software) SetWX(val uint8) { s.win.wx = int(val) - 7 }
func (s *Software) SetWY(val uint8) { s.win.wy = int(val) }

func (s *Software) SetBGPI(val uint8) { s.bg.bgpi = val }

func (s *Software) GetBGPD() uint8 {
	val := uint8(0xFF)
	if s.model == 1 {
		palID := int((s.bg.bgpi & 0x3F) / 8)
		colorID := int(s.bg.bgpi&7) >> 1
		rgb := &s.bg.palette[palID*4+colorID]
		isHi := util.Bit(s.bg.bgpi, 0)
		if isHi {
			// 0b0BBBBBGG
			val = (rgb.b << 2)
			val |= ((rgb.g >> 3) & 0b11)
		} else {
			// 0bGGGRRRRR
			val = rgb.r
			val |= (rgb.g << 5)
		}
	}
	return val
}

func (s *Software) SetBGPD(val uint8) uint8 {
	if s.model == 1 {
		palID := int((s.bg.bgpi & 0x3F) / 8)
		colorID := int(s.bg.bgpi&7) >> 1
		rgb := &s.bg.palette[palID*4+colorID]
		isHi := util.Bit(s.bg.bgpi, 0)
		if isHi {
			// 0b0BBBBBGG
			rgb.g &= 0b111
			rgb.g |= ((val & 0b11) << 3)
			rgb.b = ((val >> 2) & 0b11111)
		} else {
			// 0bGGGRRRRR
			rgb.r = (val & 0b11111)
			rgb.g &= 0b11000
			rgb.g |= ((val >> 5) & 0b111)
		}

		if util.Bit(s.bg.bgpi, 7) {
			bgpi := (s.bg.bgpi + 1) & 0x3F
			s.bg.bgpi &= 0xC0
			s.bg.bgpi |= bgpi
		}
	}
	return s.bg.bgpi
}

func (s *Software) SetOBPI(val uint8) { s.sprite.obpi = val }

func (s *Software) GetOBPD() uint8 {
	val := uint8(0xFF)
	if s.model == 1 {
		palID := int((s.sprite.obpi & 0x3F) / 8)
		colorID := int(s.sprite.obpi&7) >> 1
		rgb := &s.sprite.palette[palID*4+colorID]
		isHi := util.Bit(s.sprite.obpi, 0)
		if isHi {
			// 0b0BBBBBGG
			val = (rgb.b << 2)
			val |= ((rgb.g >> 3) & 0b11)
		} else {
			// 0bGGGRRRRR
			val = rgb.r
			val |= (rgb.g << 5)
		}
	}
	return val
}

func (s *Software) SetOBPD(val uint8) uint8 {
	if s.model == 1 {
		palID := int((s.sprite.obpi & 0x3F) / 8)
		colorID := int(s.sprite.obpi&7) >> 1
		rgb := &s.sprite.palette[palID*4+colorID]
		isHi := util.Bit(s.sprite.obpi, 0)
		if isHi {
			// 0b0BBBBBGG
			rgb.g &= 0b111
			rgb.g |= ((val & 0b11) << 3)
			rgb.b = ((val >> 2) & 0b11111)
		} else {
			// 0bGGGRRRRR
			rgb.r = (val & 0b11111)
			rgb.g &= 0b11000
			rgb.g |= ((val >> 5) & 0b111)
		}

		if util.Bit(s.sprite.obpi, 7) {
			obpi := (s.sprite.obpi + 1) & 0x3F
			s.sprite.obpi &= 0xC0
			s.sprite.obpi |= obpi
		}
	}
	return s.sprite.obpi
}
