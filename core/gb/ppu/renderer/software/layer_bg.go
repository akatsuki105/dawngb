package software

import (
	"unsafe"

	"github.com/akatsuki105/dawngb/util"
)

type bgLayer struct {
	enable   bool // LCDC.0 (CGBではattrmapのbit7を使用するかどうか(=つまりBGは常に有効))
	r        *Software
	tilemap  uint16 // 0x1800 or 0x1C00
	tiledata int    // 0x800 or 0x0
	bgp      uint8  // BGP(0xFF47)
	palette  []rgb555
	scx, scy uint8
	scanline [160]pixel
}

func newBG(r *Software, palette []rgb555) *bgLayer {
	l := &bgLayer{
		r:        r,
		tilemap:  0x1800,
		tiledata: 0x800,
		palette:  palette,
	}
	return l
}

func (l *bgLayer) drawScanline(y int) {
	enable := l.r.isCGB() || l.enable
	if enable {
		y = (y + int(l.scy)) % 256

		tilemap := l.r.vram[l.tilemap : l.tilemap+1024]
		attrmap := l.r.vram[(8*KB)+uint(l.tilemap) : (8*KB)+uint(l.tilemap)+1024] // CGBモードのみ
		for i := 0; i < 160; i++ {
			x := (i + int(l.scx)) % 256

			end := 8 - (x & 0b111)

			// 左端
			if i == 0 && (x&0b111 != 0) {
				x = (x - (x & 0b111)) % 256
			}

			// 8pxずつ描画
			if x&0b111 == 0 {
				tileOffset := ((y/8)*32 + (x / 8)) % 1024

				var tileID int
				if l.tiledata == 0x0 {
					tileID = int(tilemap[tileOffset])
				} else {
					tileID = int(int8(tilemap[tileOffset])) + 256
				}

				attr := attrmap[tileOffset]
				if !l.r.isCGB() {
					attr = 0 // DMGモードでは属性マップは常に0
				}
				palID := attr & 0b111
				tileBank := uint(util.Btou8(util.Bit(attr, 3)))
				hflip := util.Bit(attr, 5)

				tiledata := l.r.vram[(8*KB)*tileBank:]
				tile := tiledata[tileID*16 : (tileID+1)*16] // 2bpp = 16byte

				yy := flip(8, util.Bit(attr, 6), (y & 0b111))
				planes := [2]uint8{tile[yy*2], tile[yy*2+1]}

				for j := 0; j < end; j++ {
					lo := (planes[0] >> ((end - 1) - j)) & 0b1
					hi := (planes[1] >> ((end - 1) - j)) & 0b1
					colorID := uint8((hi << 1) | lo) // 0..3
					x := i + flip(8, hflip, j)
					if x < 160 {
						l.scanline[x].color = l.getColor(palID, colorID)
						l.scanline[x].colorID = colorID
						l.scanline[x].priority = attr&(1<<7) != 0
					}
				}
			}
		}
	}
}

func (l *bgLayer) drawTilemap(buffer unsafe.Pointer, n int) {
	var tilebase, mapbase uint16
	switch n {
	case 0: // Auto
		tilebase, mapbase = uint16(l.tiledata), l.tilemap
	case 1: // Tile: 0x0, Map: 0x1800
		tilebase, mapbase = 0, 0x1800
	case 2: // Tile: 0x800, Map: 0x1800
		tilebase, mapbase = 0x800, 0x1800
	case 3: // Tile: 0x0, Map: 0x1C00
		tilebase, mapbase = 0, 0x1C00
	case 4: // Tile: 0x800, Map: 0x1C00
		tilebase, mapbase = 0x800, 0x1C00
	}
	tilemap := l.r.vram[mapbase : mapbase+1024]
	attrmap := l.r.vram[(8*KB)+mapbase : (8*KB)+mapbase+1024] // CGBモードのみ

	for i := 0; i < 1024; i++ {
		tx, ty := (i%32)*8, (i/32)*8

		var tileID int
		if tilebase == 0x0 {
			tileID = int(tilemap[i])
		} else {
			tileID = int(int8(tilemap[i])) + 256
		}

		attr := attrmap[i]
		if !l.r.isCGB() {
			attr = 0 // DMGモードでは属性マップは常に0
		}
		palID := attr & 0b111
		tileBank := int((attr >> 3) & 0b1)
		xflip, yflip := (attr&(1<<5)) != 0, (attr&(1<<6)) != 0

		tile := l.r.vram[((8*KB)*tileBank)+(tileID*16) : ((8*KB)*tileBank)+((tileID+1)*16)]

		for py := 0; py < 8; py++ {
			y := ty + flip(8, yflip, py)

			lo := tile[py*2]
			hi := tile[py*2+1]
			for px := 0; px < 8; px++ {
				x := tx + flip(8, xflip, px)

				loBit := (lo >> (7 - px)) & 0b1
				hiBit := (hi >> (7 - px)) & 0b1
				colorID := uint8((hiBit << 1) | loBit)
				rgb555 := l.getColor(palID, colorID)

				r5, g5, b5 := uint8(rgb555&0x1F), uint8((rgb555>>5)&0x1F), uint8((rgb555>>10)&0x1F)
				pixel := (*[4]uint8)(unsafe.Pointer(uintptr(buffer) + uintptr((y*256+x)*4)))
				pixel[0] = (r5 << 3) | (r5 >> 2)
				pixel[1] = (g5 << 3) | (g5 >> 2)
				pixel[2] = (b5 << 3) | (b5 >> 2)
				pixel[3] = 0xFF
			}
		}
	}
}

func (l *bgLayer) getColor(palID, n uint8) rgb555 {
	if l.r.isCGB() {
		return l.palette[((palID&0b111)*4)+(n&0b11)]
	}

	return l.palette[(l.bgp>>((n&0b11)*2))&0b11]
}
