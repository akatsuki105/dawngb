package software

import (
	"github.com/akatsuki105/dawngb/util"
)

type bgLayer struct {
	active   bool
	r        *Software
	tilemap  uint16 // 0x1800 or 0x1C00
	tiledata int    // 0x800 or 0x0
	palette  [4 * 8]rgb555
	scx, scy int
	bgpi     uint8
}

func newBG(r *Software) *bgLayer {
	l := &bgLayer{
		r:        r,
		tilemap:  0x1800,
		tiledata: 0x800,
	}
	return l
}

func (l *bgLayer) drawScanline(y int, scanline []pixel) {
	if l.active {
		y = (y + l.scy) % 256

		tilemap := l.r.vram[l.tilemap : l.tilemap+1024]
		attrmap := l.r.vram[(8*KB)+uint(l.tilemap) : (8*KB)+uint(l.tilemap)+1024]
		for i := 0; i < 160; i++ {
			x := (i + l.scx) % 256

			end := 8 - (x & 0b111)

			// 左端
			if i == 0 && (x&0b111 != 0) {
				x = (x - (x & 0b111)) % 256
			}

			// 8pxずつ描画
			if x&0b111 == 0 {
				z := Z_BG

				tileOffset := ((y/8)*32 + (x / 8)) % 1024

				var tileID int
				if l.tiledata == 0x0 {
					tileID = int(tilemap[tileOffset])
				} else {
					tileID = int(int8(tilemap[tileOffset])) + 256
				}
				attr := attrmap[tileOffset]
				palID := int(attr & 0b111)
				tileBank := uint(util.Btou8(util.Bit(attr, 3)))
				hflip := util.Bit(attr, 5)
				if util.Bit(attr, 7) {
					z += Z_SPR
				}

				tiledata := l.r.vram[(8*KB)*tileBank:]
				tile := tiledata[tileID*16 : (tileID+1)*16] // 2bpp = 16byte

				yy := util.Flip(8, util.Bit(attr, 6), (y & 0b111))
				planes := [2]uint8{tile[yy*2], tile[yy*2+1]}

				for j := 0; j < end; j++ {
					lo := (planes[0] >> ((end - 1) - j)) & 0b1
					hi := (planes[1] >> ((end - 1) - j)) & 0b1
					colorID := int((hi << 1) | lo)
					x := i + util.Flip(8, hflip, j)
					if x < len(scanline) {
						if z >= scanline[x].z {
							scanline[x].rgba = l.palette[(palID*4)+colorID]
							scanline[x].z = z
							scanline[x].colorID = colorID
							scanline[x].isBG = true
						}
					}
				}
			}
		}
	}
}
