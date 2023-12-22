package software

import "image/color"

type Software struct{}

func New() *Software {
	return &Software{}
}

func (s *Software) DrawScanline(y int, scanline []color.RGBA) {}
