package software

import "image/color"

type rgb555 struct {
	r, g, b uint8
}

func (c rgb555) RGBA() color.RGBA {
	r := (c.r << 3) | (c.r >> 2)
	g := (c.g << 3) | (c.g >> 2)
	b := (c.b << 3) | (c.b >> 2)
	return color.RGBA{r, g, b, 0xFF}
}
