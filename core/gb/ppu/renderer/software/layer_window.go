package software

import "unsafe"

type windowLayer struct {
	enable  bool // LCDC.5
	r       *Software
	tilemap uint16 // LCDC.6; 0x1800 or 0x1C00
	wx, wy  int

	/*
		ウィンドウは、機能的にはLYに似た内部ラインカウンタを保持している。
		これはウィンドウが表示されているときのみインクリメントされる。 この行カウンタによって、現在の走査線上にレンダリングされるウィンドウの行が決定される。
	*/
	y uint8
}

func newWindow(r *Software) *windowLayer {
	return &windowLayer{
		r:       r,
		tilemap: 0x1800,
	}
}

func (l *windowLayer) drawScanline(y int) {
	rendered := false
	if l.enable {
		if (l.wx >= 0 && l.wx < 160) && (l.wy >= 0 && l.wy < 144) {
			if y >= l.wy {
				y = int(l.y)
				tilemap := l.r.vram[l.tilemap : l.tilemap+1024]
				attrmap := l.r.vram[(8*KB)+uint(l.tilemap) : (8*KB)+uint(l.tilemap)+1024]
				for i := 0; i < 160; i++ {
					if i >= l.wx {
						rendered = true

						x := i - l.wx

						// 8pxずつ描画
						if x&0b111 == 0 {
							tileOffset := ((y/8)*32 + (x / 8)) % 1024

							var tileID int
							if l.r.bg.tiledata == 0x0 {
								tileID = int(tilemap[tileOffset])
							} else {
								tileID = int(int8(tilemap[tileOffset])) + 256
							}

							attr := attrmap[tileOffset]
							if !l.r.isCGB() {
								attr = 0 // DMGモードでは属性マップは常に0
							}
							palID := attr & 0b111
							tileBank := uint((attr >> 3) & 0b1)
							xflip, yflip := (attr&(1<<5)) != 0, (attr&(1<<6)) != 0

							tiledata := l.r.vram[(8*KB)*tileBank:]
							tile := tiledata[tileID*16 : (tileID+1)*16] // 2bpp = 16byte

							yy := flip(8, yflip, (y & 0b111))
							planes := [2]uint8{tile[yy*2], tile[yy*2+1]}

							for j := 0; j < 8; j++ {
								lo := (planes[0] >> (7 - j)) & 0b1
								hi := (planes[1] >> (7 - j)) & 0b1
								colorID := ((hi << 1) | lo)
								x := i + flip(8, xflip, j)
								if x < 160 {
									l.r.bg.scanline[x].color = l.r.bg.getColor(palID, colorID)
									l.r.bg.scanline[x].colorID = colorID
									l.r.bg.scanline[x].priority = attr&(1<<7) != 0
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

func (l *windowLayer) drawTilemap(buffer unsafe.Pointer, n int) {
	var tilebase, mapbase uint16
	switch n {
	case 0: // Auto
		tilebase, mapbase = uint16(l.r.bg.tiledata), l.tilemap
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
		tx, ty := i%32, i/32

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

		tile := l.r.vram[((8*KB)*tileBank)+(tileID*16) : ((8*KB)*tileBank)+((tileID+1)*16)]

		for py := 0; py < 8; py++ {
			y := int(ty*8 + py)
			lo := tile[py*2]
			hi := tile[py*2+1]

			for px := 0; px < 8; px++ {
				x := int(tx*8 + px)
				loBit := (lo >> (7 - px)) & 0b1
				hiBit := (hi >> (7 - px)) & 0b1
				colorID := uint8((hiBit << 1) | loBit)
				rgb555 := l.r.bg.getColor(palID, colorID)

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
