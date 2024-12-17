package software

type spriteLayer struct {
	enable   bool // LCDC.1
	r        *Software
	height   uint8    // 8 or 16
	obp      [2]uint8 // OBP0: 0xFF48, OBP1: 0xFF49
	palette  []rgb555
	scanline [160]pixel
}

type sprite struct {
	x, y         int
	tileID       int
	xflip, yflip bool
	palID        uint8 // DMG: 0 or 1, CGB: 0-7
	bank         uint
	priority     bool // OAM Priority (3バイト目のbit7)
}

func newSpriteLayer(r *Software, palette []rgb555) *spriteLayer {
	l := &spriteLayer{
		r:       r,
		height:  8,
		palette: palette,
	}
	return l
}

func (l *spriteLayer) drawScanline(y int) {
	if l.enable {
		// 1行に描画されるスプライトの数は最大10個
		sprites := [10]int{}
		amount := 0
		for i := 0; i < 40; i++ {
			spriteIdx := i
			spriteY := int(l.r.oam[spriteIdx*4+0]) - 16
			if (spriteY <= y) && (y < spriteY+int(l.height)) {
				if amount < 10 {
					sprites[amount] = spriteIdx
					amount++
				}
			}
		}

		var spr sprite
		for i := amount - 1; i >= 0; i-- {
			l.getSprite(sprites[i], &spr)
			switch l.height {
			case 8:
				l.drawObjScanline8(&spr, y)
			case 16:
				l.drawObjScanline16(&spr, y)
			}
		}
	}
}

func (l *spriteLayer) drawObjScanline8(s *sprite, y int) {
	tiledata := l.r.vram[(s.bank * (8 * KB)) : (s.bank*(8*KB))+0x1000]
	tile := tiledata[s.tileID*16 : (s.tileID+1)*16] // 2bpp = 16byte

	row := flip(8, s.yflip, y-s.y) // (スプライトの一番上を0行目として)上から何行目か

	planes := [2]uint8{tile[(row&0b111)*2], tile[(row&0b111)*2+1]}

	for i := 0; i < 8; i++ {
		lo := (planes[0] >> (7 - uint(i))) & 0b1
		hi := (planes[1] >> (7 - uint(i))) & 0b1
		colorID := ((hi << 1) | lo)
		if colorID != 0 {
			idx := s.x + flip(8, s.xflip, i)
			if (0 <= idx) && (idx < 160) {
				l.scanline[idx].color = l.getColor(s.palID, colorID)
				l.scanline[idx].colorID = colorID
				l.scanline[idx].priority = s.priority
			}
		}
	}
}

func (l *spriteLayer) drawObjScanline16(s *sprite, y int) {
	tiledata := l.r.vram[(s.bank * (8 * KB)) : (s.bank*(8*KB))+0x1000]
	tileID := s.tileID & 0xFE
	tile := tiledata[tileID*16 : (tileID+2)*16] // 2bpp

	row := flip(16, s.yflip, y-s.y) // (スプライトの一番上を0行目として)上から何行目か

	var planes [2]uint8
	if row < 8 {
		planes = [2]uint8{tile[(row&0b111)*2], tile[(row&0b111)*2+1]}
	} else {
		planes = [2]uint8{tile[(row&0b111)*2+16], tile[(row&0b111)*2+17]}
	}

	for i := 0; i < 8; i++ {
		lo := (planes[0] >> (7 - i)) & 0b1
		hi := (planes[1] >> (7 - i)) & 0b1
		colorID := ((hi << 1) | lo)
		if colorID != 0 {
			idx := s.x + flip(8, s.xflip, i)
			if (0 <= idx) && (idx < 160) {
				l.scanline[idx].color = l.getColor(s.palID, colorID)
				l.scanline[idx].colorID = colorID
				l.scanline[idx].priority = s.priority
			}
		}
	}
}

func (l *spriteLayer) getSprite(spriteIdx int, s *sprite) {
	byte0 := l.r.oam[spriteIdx*4+0]
	byte1 := l.r.oam[spriteIdx*4+1]
	byte2 := l.r.oam[spriteIdx*4+2]
	byte3 := l.r.oam[spriteIdx*4+3]

	bank := uint(0)
	palID := ((byte3 >> 4) & 0b1) // 0: OBP0, 1: OBP1
	if l.r.isCGB() {
		bank = uint((byte3 >> 3) & 0b1)
		palID = byte3 & 0b111
	}

	s.x, s.y = int(byte1)-8, int(byte0)-16
	s.tileID = int(byte2)
	s.xflip, s.yflip = (byte3&(1<<5)) != 0, (byte3&(1<<6)) != 0
	s.palID, s.bank, s.priority = palID, bank, (byte3&(1<<7)) == 0
}

func (l *spriteLayer) getColor(palID, n uint8) rgb555 {
	if l.r.isCGB() {
		return l.palette[((palID&0b111)*4)+(n&0b11)]
	}

	obp := l.obp[palID&1]
	return l.palette[(obp>>((n&0b11)*2))&0b11]
}
