package software

import (
	"github.com/akatsuki105/dawngb/util"
)

type spriteLayer struct {
	active  bool
	r       *Software
	height  int // 8 or 16
	palette [4 * 8]rgb555
	obpi    uint8
	z       int // スプライト全体に履かせる下駄となるz-index(CGBのLCDC.0で決定)
}

type sprite struct {
	x, y         int
	tileID       int
	xflip, yflip bool
	palID        int
	bank         uint
	z            int
}

func newSpriteLayer(r *Software) *spriteLayer {
	l := &spriteLayer{
		r:      r,
		height: 8,
	}
	return l
}

func (l *spriteLayer) drawScanline(y int, scanline []pixel) {
	if l.active {
		// 1行に描画されるスプライトの数は最大10個
		sprites := [10]int{}
		amount := 0
		for i := 0; i < 40; i++ {
			spriteIdx := i
			spriteY := int(l.r.oam[spriteIdx*4+0]) - 16
			if (spriteY <= y) && (y < spriteY+l.height) {
				if amount < 10 {
					sprites[amount] = spriteIdx
					amount++
				}
			}
		}

		for i := amount - 1; i >= 0; i-- {
			spriteIdx := sprites[i]
			switch l.height {
			case 8:
				l.drawObjScanline8(spriteIdx, scanline, y)
			case 16:
				l.drawObjScanline16(spriteIdx, scanline, y)
			}
		}
	}
}

func (l *spriteLayer) drawObjScanline8(spriteIdx int, scanline []pixel, y int) {
	s := l.getSprite(spriteIdx)
	z := (s.z + l.z)

	tiledata := l.r.vram[(s.bank * (8 * KB)) : (s.bank*(8*KB))+0x1000]
	tile := tiledata[s.tileID*16 : (s.tileID+1)*16] // 2bpp = 16byte

	row := util.Flip(8, s.yflip, y-s.y) // (スプライトの一番上を0行目として)上から何行目か

	planes := [2]uint8{tile[(row&0b111)*2], tile[(row&0b111)*2+1]}

	palette := l.palette[s.palID*4 : (s.palID+1)*4]

	for i := 0; i < 8; i++ {
		lo := (planes[0] >> (7 - uint(i))) & 0b1
		hi := (planes[1] >> (7 - uint(i))) & 0b1
		colorID := int((hi << 1) | lo)
		if colorID != 0 {
			idx := s.x + util.Flip(8, s.xflip, i)
			if (0 <= idx) && (idx < 160) {
				if z >= scanline[idx].z || (scanline[idx].colorID == 0) {
					scanline[idx].rgba = palette[colorID]
					scanline[idx].z = z
					scanline[idx].colorID = colorID
				}
			}
		}
	}
}

func (l *spriteLayer) drawObjScanline16(spriteIdx int, scanline []pixel, y int) {
	s := l.getSprite(spriteIdx)
	z := (s.z + l.z)

	tiledata := l.r.vram[(s.bank * (8 * KB)) : (s.bank*(8*KB))+0x1000]
	tileID := s.tileID & 0xFE
	tile := tiledata[tileID*16 : (tileID+2)*16] // 2bpp

	row := util.Flip(16, s.yflip, y-s.y) // (スプライトの一番上を0行目として)上から何行目か

	var planes [2]uint8
	if row < 8 {
		planes = [2]uint8{tile[(row&0b111)*2], tile[(row&0b111)*2+1]}
	} else {
		planes = [2]uint8{tile[(row&0b111)*2+16], tile[(row&0b111)*2+17]}
	}

	palette := l.palette[s.palID*4 : (s.palID+1)*4]

	for i := 0; i < 8; i++ {
		lo := (planes[0] >> (7 - uint(i))) & 0b1
		hi := (planes[1] >> (7 - uint(i))) & 0b1
		colorID := int((hi << 1) | lo)
		if colorID != 0 {
			idx := s.x + util.Flip(8, s.xflip, i)
			if (0 <= idx) && (idx < 160) {
				if z >= scanline[idx].z || (scanline[idx].colorID == 0) {
					scanline[idx].rgba = palette[colorID]
					scanline[idx].z = z
					scanline[idx].colorID = colorID
				}
			}
		}
	}
}

func (l *spriteLayer) getSprite(spriteIdx int) *sprite {
	byte0 := l.r.oam[spriteIdx*4+0]
	byte1 := l.r.oam[spriteIdx*4+1]
	byte2 := l.r.oam[spriteIdx*4+2]
	byte3 := l.r.oam[spriteIdx*4+3]

	z := Z_SPR
	if util.Bit(byte3, 7) {
		z = Z_BD
	}

	bank := uint(0)
	palID := (int(byte3>>4) & 0b1)
	if l.r.model == 1 {
		bank = uint(util.Btou8(util.Bit(byte3, 3)))
		palID = int(byte3) & 0b111
	}
	return &sprite{
		y:      int(byte0) - 16,
		x:      int(byte1) - 8,
		tileID: int(byte2),
		xflip:  util.Bit(byte3, 5),
		yflip:  util.Bit(byte3, 6),
		palID:  palID,
		bank:   bank,
		z:      z,
	}
}
