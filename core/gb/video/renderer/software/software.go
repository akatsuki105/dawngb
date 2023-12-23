package software

import "image/color"

type Software struct{}

func New() *Software {
	return &Software{}
}

func (s *Software) DrawScanline(y int, scanline []color.RGBA) {
	for i := 0; i < 160; i++ {
		scanline[i] = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}
}
