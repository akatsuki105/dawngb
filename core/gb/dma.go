package gb

import (
	"fmt"

	"github.com/akatsuki105/dugb/util"
)

const (
	GDMA = iota
	HDMA
)

// DMA Controller(CGB only)
type dmaController struct {
	g         *GB
	mode      uint8
	src, dst  uint16
	length    uint16
	completed bool
}

func newDMAController(g *GB) *dmaController {
	return &dmaController{
		g:         g,
		completed: true,
	}
}

func (d *dmaController) Reset(hasBIOS bool) {}

func (d *dmaController) Read(addr uint16) uint8 {
	val := uint8(0xFF)
	if addr == 0xFF55 {
		val = uint8(0)
		val = util.SetBit(val, 7, d.completed)
		length := uint8(((d.length / 16) - 1) & 0x7F)
		val |= length
	}
	return val
}

func (d *dmaController) Write(addr uint16, val uint8) {
	switch addr {
	case 0xFF51: // upper src
		d.src = (d.src & 0x00FF) | (uint16(val) << 8)
	case 0xFF52: // lower src
		d.src = (d.src & 0xFF00) | uint16(val)
	case 0xFF53: // upper dst
		d.dst = (d.dst & 0x00FF) | (uint16(val) << 8)
	case 0xFF54: // lower dst
		d.dst = (d.dst & 0xFF00) | uint16(val)
	case 0xFF55: // control
		d.length = (uint16(val&0b111_1111) + 1) * 16 // 16~2048バイトまで指定可能
		d.mode = val >> 7
		d.completed = false
		if d.mode == GDMA {
			// Trigger GDMA
			period := int64(d.length) * 4
			length := d.length
			d.g.dma.Callback = func(cyclesLate int64) {
				for i := uint16(0); i < length; i++ {
					for j := 0; j < 16; j++ {
						d.g.video.Write(d.dst+i, d.g.m.Read(d.src+i))
					}
					d.length -= 16
				}
				d.completed = true
				d.g.blocked = false
			}
			d.g.blocked = true
			d.g.s.Schedule(&d.g.dma, period)
		} else {
			// Trigger HDMA
			fmt.Println("Trigger HDMA")
			d.completed = true
		}
	}
}
