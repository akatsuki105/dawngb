package software

import (
	"image/color"

	"github.com/akatsuki105/dawngb/util"
	"golang.org/x/exp/constraints"
)

const KB = 1024

const (
	Z_BD = iota
	Z_BG
	Z_SPR
	Z_MASTER_SPR = 100 // CGBモードでLCDC.0が0のときは必ずスプライトが前面に来る
)

type rgb555 = uint16 // 0b0_BBBBB_GGGGG_RRRRR

var dmgPalette = [4]rgb555{
	0b11111_11111_11111,
	0b10001_10001_10001,
	0b01010_01010_01010,
	0b00000_00000_00000,
}

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
	color   rgb555
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

func (s *Software) DrawScanline(y int, scanline []color.NRGBA) {
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
		r5, g5, b5 := uint8(s.scanline[i].color&0x1F), uint8((s.scanline[i].color>>5)&0x1F), uint8((s.scanline[i].color>>10)&0x1F)
		scanline[i].R = (r5 << 3) | (r5 >> 2)
		scanline[i].G = (g5 << 3) | (g5 >> 2)
		scanline[i].B = (b5 << 3) | (b5 >> 2)
		scanline[i].A = 0xFF
	}
}

func (s *Software) SetLCDC(val uint8) {
	s.bg.active = (val & (1 << 0)) != 0
	if s.model == 1 {
		s.bg.active = true
		if (val & (1 << 0)) != 0 {
			s.sprite.z = 0
		} else {
			s.sprite.z = Z_MASTER_SPR
		}
	}
	s.bg.tilemap = [2]uint16{0x1800, 0x1C00}[(val>>3)&1]
	s.bg.tiledata = [2]int{0x800, 0x0}[(val>>4)&1]

	s.win.active = s.bg.active && util.Bit(val, 5)
	s.win.tilemap = [2]uint16{0x1800, 0x1C00}[(val>>6)&1]

	s.sprite.active = util.Bit(val, 1)
	s.sprite.height = [2]int{8, 16}[(val>>2)&1]
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
	if s.model == 1 {
		palID := int((s.bg.bgpi & 0x3F) / 8)
		colorID := int(s.bg.bgpi&7) >> 1
		rgb := s.bg.palette[palID*4+colorID]
		isHi := util.Bit(s.bg.bgpi, 0)
		if isHi {
			return uint8(rgb >> 8)
		} else {
			return uint8(rgb)
		}
	}
	return 0xFF
}

func (s *Software) SetBGPD(val uint8) uint8 {
	if s.model == 1 {
		palID := int((s.bg.bgpi & 0x3F) / 8)
		colorID := int(s.bg.bgpi&7) >> 1
		idx := palID*4 + colorID
		rgb555 := s.bg.palette[idx]
		isHi := util.Bit(s.bg.bgpi, 0)
		if isHi {
			rgb555 = (rgb555 & 0x00FF) | (uint16(val) << 8)
		} else {
			rgb555 = (rgb555 & 0xFF00) | uint16(val)
		}
		s.bg.palette[idx] = rgb555

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
	if s.model == 1 {
		palID := int((s.sprite.obpi & 0x3F) / 8)
		colorID := int(s.sprite.obpi&7) >> 1
		rgb := s.sprite.palette[palID*4+colorID]
		isHi := util.Bit(s.sprite.obpi, 0)
		if isHi {
			return uint8(rgb >> 8)
		} else {
			return uint8(rgb)
		}
	}
	return 0xFF
}

func (s *Software) SetOBPD(val uint8) uint8 {
	if s.model == 1 {
		palID := int((s.sprite.obpi & 0x3F) / 8)
		colorID := int(s.sprite.obpi&7) >> 1
		idx := palID*4 + colorID
		rgb555 := s.sprite.palette[palID*4+colorID]
		isHi := util.Bit(s.sprite.obpi, 0)
		if isHi {
			rgb555 = (rgb555 & 0x00FF) | (uint16(val) << 8)
		} else {
			rgb555 = (rgb555 & 0xFF00) | uint16(val)
		}
		s.sprite.palette[idx] = rgb555

		if util.Bit(s.sprite.obpi, 7) {
			obpi := (s.sprite.obpi + 1) & 0x3F
			s.sprite.obpi &= 0xC0
			s.sprite.obpi |= obpi
		}
	}
	return s.sprite.obpi
}

// Helper functions
func flip[V constraints.Integer, W constraints.Integer](size V, b bool, i W) V {
	if b {
		return V(int(size-1) - int(i))
	}
	return V(i)
}
