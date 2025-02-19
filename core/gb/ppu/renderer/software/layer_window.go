package software

type windowLayer struct {
	enable  bool // LCDC.5
	r       *Software
	tilemap uint16 // LCDC.6; 0x1800 or 0x1C00
	wx, wy  int

	/*
		ウィンドウは、機能的にはLYに似た内部ラインカウンタを保持している。
		これはウィンドウが表示されているときのみインクリメントされる。 このラインカウンタによって、現在の走査線上にレンダリングされるウィンドウの行が決定される。
	*/
	ly uint8
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
		if (l.wx < 160) && (l.wy >= 0 && l.wy < 144) {
			if y >= l.wy {
				y = int(l.ly)

				tilemap := l.r.vram[l.tilemap : l.tilemap+1024]
				attrmap := l.r.vram[(8*KB)+uint(l.tilemap) : (8*KB)+uint(l.tilemap)+1024]

				for i := 0; i < 160; i++ {
					if i == l.wx {
						rendered = true // wx が 0..159 の範囲にある場合、ウィンドウが表示されているとみなす (wx=-6とかでもウィンドウは表示される?が、表示されたと見なされない; SaGa1などでタイトル画面の描画に必要)
					}

					if i >= l.wx {
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
		l.ly++
	}
}
