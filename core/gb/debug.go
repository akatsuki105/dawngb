package gb

import (
	"image"
)

func (g *GB) DebugVRAM() image.Image {
	return g.video.Debug()
}
