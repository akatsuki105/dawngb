package cpu

import (
	"github.com/akatsuki105/dawngb/core/gb/cpu/sm83"
	"github.com/akatsuki105/dawngb/util"
)

const (
	GDMA = iota
	HDMA
)

// VRAM DMA (CGB only)
type DMA struct {
	bus         sm83.Bus
	mode        uint8
	src, dst    uint16
	length      uint16
	completed   bool
	doHDMA      bool
	pendingHDMA bool
}

func newDMA(bus sm83.Bus) *DMA {
	return &DMA{
		bus:       bus,
		completed: true,
	}
}

func (d *DMA) reset() {
	d.mode = GDMA
	d.src, d.dst, d.length = 0, 0, 0
	d.completed = true
}

func (d *DMA) skipBIOS() {
	d.Write(0xFF51, 0xFF)
	d.Write(0xFF52, 0xFF)
	d.Write(0xFF53, 0xFF)
	d.Write(0xFF54, 0xFF)
	d.Write(0xFF55, 0xFF)
}

func (d *DMA) Read(addr uint16) uint8 {
	val := uint8(0xFF)
	if addr == 0xFF55 {
		val = util.SetBit(val, 7, d.completed)
		length := uint8(((d.length / 16) - 1) & 0x7F)
		val |= length
	}
	return val
}

func (d *DMA) Write(addr uint16, val uint8) int64 {
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
		d.pendingHDMA = false
		if wasCompleted && d.mode == GDMA { // Trigger GDMA
			return d.runGDMA()
		} else if d.mode == HDMA { // Trigger HDMA
			d.pendingHDMA = true
		}
	}
	return 0
}

func (d *DMA) runGDMA() int64 {
	period := int64(d.length) * 4
	for d.length > 0 {
		for i := uint16(0); i < 16; i++ {
			d.bus.Write(d.dst+i, d.bus.Read(d.src+i))
		}
		d.src += 16
		d.dst += 16
		d.length -= 16
	}
	return period
}

func (d *DMA) startHDMA() {
	if d.pendingHDMA {
		d.doHDMA = true
	}
}

// HBlank になるたびに実行される
func (d *DMA) runHDMA() {
	for i := uint16(0); i < 16; i++ {
		d.bus.Write(d.dst, d.bus.Read(d.src))
		d.src++
		d.dst++
		d.length--
	}
	if d.length == 0 {
		d.pendingHDMA = false
		d.completed = true
	}
}
