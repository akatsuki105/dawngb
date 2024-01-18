package gb

import (
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

func (d *dmaController) Reset(hasBIOS bool) {
	if !hasBIOS {
		d.Write(0xFF51, 0xFF)
		d.Write(0xFF52, 0xFF)
		d.Write(0xFF53, 0xFF)
		d.Write(0xFF54, 0xFF)
		d.Write(0xFF55, 0xFF)
	}
}

func (d *dmaController) Read(addr uint16) uint8 {
	val := uint8(0xFF)
	if addr == 0xFF55 {
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
		d.src &= 0b1111_1111_1111_0000
	case 0xFF52: // lower src
		d.src = (d.src & 0xFF00) | uint16(val)
		d.src &= 0b1111_1111_1111_0000
	case 0xFF53: // upper dst
		d.dst = (d.dst & 0x00FF) | (uint16(val) << 8)
		d.dst &= 0b0001_1111_1111_0000
		d.dst |= 0x8000
	case 0xFF54: // lower dst
		d.dst = (d.dst & 0xFF00) | uint16(val)
		d.dst &= 0b0001_1111_1111_0000
		d.dst |= 0x8000
	case 0xFF55: // control
		wasCompleted := d.completed
		d.length = (uint16(val&0b111_1111) + 1) * 16 // 16~2048バイトまで指定可能
		d.mode = val >> 7
		d.completed = (d.mode == GDMA)
		d.g.runHDMA = nil
		if wasCompleted && d.mode == GDMA {
			// Trigger GDMA
			period := int64(d.length) * 4
			d.g.dma.Callback = func(cyclesLate int64) {
				for d.length > 0 {
					for i := uint16(0); i < 16; i++ {
						d.g.video.Write(d.dst+i, d.g.m.Read(d.src+i))
					}
					d.src += 16
					d.dst += 16
					d.length -= 16
				}
				d.g.blocked = false
			}
			d.g.blocked = true
			d.g.s.Schedule(&d.g.dma, period)
		} else {
			if d.mode == HDMA {
				// Trigger HDMA
				d.g.runHDMA = d.runHDMA
			}
		}
	}
}

func (d *dmaController) runHDMA() {
	d.g.dma.Callback = func(cyclesLate int64) {
		for i := uint16(0); i < 16; i++ {
			d.g.video.Write(d.dst, d.g.m.Read(d.src))
			d.src++
			d.dst++
			d.length--
		}
		if d.length == 0 {
			d.g.runHDMA = nil
			d.completed = true
		}
		d.g.blocked = false
	}
	d.g.blocked = true
	d.g.s.Schedule(&d.g.dma, 64)
}
