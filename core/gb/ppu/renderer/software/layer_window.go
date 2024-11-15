package software

import (
	"github.com/akatsuki105/dawngb/util"
)

type windowLayer struct {
	active  bool
	r       *Software
	tilemap uint16 // 0x1800 or 0x1C00
	wx, wy  int

	/*
		ウィンドウは、機能的にはLYに似た内部ラインカウンタを保持している。
		これはウィンドウが表示されているときのみインクリメントされる。 この行カウンタによって、現在の走査線上にレンダリングされるウィンドウの行が決定される。
	*/
	y int
}

func newWindow(r *Software) *windowLayer {
	return &windowLayer{
		r:       r,
		tilemap: 0x1800,
	}
}

func (l *windowLayer) drawScanline(y int, scanline []pixel) {
	rendered := false
	if l.active {
		if (l.wx >= 0 && l.wx < 160) && (l.wy >= 0 && l.wy < 144) {
			if y >= l.wy {
				y = l.y
				tilemap := l.r.vram[l.tilemap : l.tilemap+1024]
				attrmap := l.r.vram[(8*KB)+uint(l.tilemap) : (8*KB)+uint(l.tilemap)+1024]
				for i := 0; i < 160; i++ {
					if i >= l.wx {
						rendered = true

						x := i - l.wx

						// 8pxずつ描画
						if x&0b111 == 0 {
							z := Z_BG

							tileOffset := ((y/8)*32 + (x / 8)) % 1024

							var tileID int
							if l.r.bg.tiledata == 0x0 {
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

							yy := flip(8, util.Bit(attr, 6), (y & 0b111))
							planes := [2]uint8{tile[yy*2], tile[yy*2+1]}

							for j := 0; j < 8; j++ {
								lo := (planes[0] >> (7 - uint(j))) & 0b1
								hi := (planes[1] >> (7 - uint(j))) & 0b1
								colorID := int((hi << 1) | lo)
								x := i + flip(8, hflip, j)
								if x < len(scanline) {
									if scanline[x].z <= z || scanline[x].isBG {
										scanline[x].color = l.r.bg.palette[(palID*4)+colorID]
										scanline[x].z = z
										scanline[x].colorID = colorID
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// ウィンドウが表示されているときのみラインカウンタをインクリメント
	if rendered {
		l.y++
	}
}
