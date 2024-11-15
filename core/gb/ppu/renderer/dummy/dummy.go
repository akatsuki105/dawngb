package dummy

import (
	"image/color"
)

type Renderer struct{}

func New() *Renderer {
	return &Renderer{}
}

func (r *Renderer) DrawScanline(y int, scanline []color.NRGBA) {}
func (r *Renderer) SetLCDC(val uint8)                          {}
func (r *Renderer) SetBGP(val uint8)                           {}
func (r *Renderer) SetOBP0(val uint8)                          {}
func (r *Renderer) SetOBP1(val uint8)                          {}
func (r *Renderer) SetSCX(val uint8)                           {}
func (r *Renderer) SetSCY(val uint8)                           {}
func (r *Renderer) SetWX(val uint8)                            {}
func (r *Renderer) SetWY(val uint8)                            {}

func (r *Renderer) SetBGPI(val uint8)       {}
func (r *Renderer) GetBGPD() uint8          { return 0 }
func (r *Renderer) SetBGPD(val uint8) uint8 { return 0 }

func (r *Renderer) SetOBPI(val uint8)       {}
func (r *Renderer) GetOBPD() uint8          { return 0 }
func (r *Renderer) SetOBPD(val uint8) uint8 { return 0 }
