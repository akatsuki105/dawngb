package software

import (
	"github.com/akatsuki105/dugb/util"
	. "github.com/akatsuki105/dugb/util/datasize"
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

			// 8pxずつ描画
			if x&0b111 == 0 {
				z := Z_BG

				var tileID int
				if l.tiledata == 0x0 {
					tileID = int(tilemap[(y/8)*32+(x/8)])
				} else {
					tileID = int(int8(tilemap[(y/8)*32+(x/8)])) + 256
				}
				attr := attrmap[(y/8)*32+(x/8)]
				palID := int(attr & 0b111)
				tileBank := uint(util.Btou8(util.Bit(attr, 3)))

				tiledata := l.r.vram[(8*KB)*tileBank:]
				tile := tiledata[tileID*16 : (tileID+1)*16] // 2bpp = 16byte

				planes := [2]uint8{tile[(y&0b111)*2], tile[(y&0b111)*2+1]}

				for j := 0; j < 8; j++ {
					lo := (planes[0] >> (7 - uint(j))) & 0b1
					hi := (planes[1] >> (7 - uint(j))) & 0b1
					colorID := int((hi << 1) | lo)
					if (i + j) < len(scanline) {
						if z >= scanline[i+j].z {
							scanline[i+j].rgba = l.palette[(palID*4)+colorID].RGBA()
							scanline[i+j].z = z
							scanline[i+j].colorID = colorID
						}
					}
				}
			}
		}
	}
}
