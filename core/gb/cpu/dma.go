package cpu

import (
	"github.com/akatsuki105/dawngb/core/gb/cpu/sm83"
	"github.com/akatsuki105/dawngb/core/gb/internal"
)

const (
	GDMA = iota
	HDMA
)

// VRAM DMA (CGB only)
type DMA struct {
	bus         sm83.Bus
	Mode        uint8
	Src, Dst    uint16
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
	d.Mode = GDMA
	d.Src, d.Dst, d.length = 0, 0, 0
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
		val = internal.SetBit(val, 7, d.completed)
		length := uint8(((d.length / 16) - 1) & 0x7F)
		val |= length
	}
	return val
}

func (d *DMA) Write(addr uint16, val uint8) int64 {
	switch addr {
	case 0xFF51: // upper src
		d.Src = (d.Src & 0x00FF) | (uint16(val) << 8)
		d.Src &= 0b1111_1111_1111_0000
	case 0xFF52: // lower src
		d.Src = (d.Src & 0xFF00) | uint16(val)
		d.Src &= 0b1111_1111_1111_0000
	case 0xFF53: // upper dst
		d.Dst = (d.Dst & 0x00FF) | (uint16(val) << 8)
		d.Dst &= 0b0001_1111_1111_0000
		d.Dst |= 0x8000
	case 0xFF54: // lower dst
		d.Dst = (d.Dst & 0xFF00) | uint16(val)
		d.Dst &= 0b0001_1111_1111_0000
		d.Dst |= 0x8000
	case 0xFF55: // control
		wasCompleted := d.completed
		d.length = (uint16(val&0b111_1111) + 1) * 16 // 16~2048バイトまで指定可能
		d.Mode = val >> 7
		d.completed = (d.Mode == GDMA)
		d.pendingHDMA = false
		if wasCompleted && d.Mode == GDMA { // Trigger GDMA
			return d.runGDMA()
		} else if d.Mode == HDMA { // Trigger HDMA
			d.pendingHDMA = true
		}
	}
	return 0
}

func (d *DMA) runGDMA() int64 {
	period := int64(d.length) * 4
	for d.length > 0 {
		for i := uint16(0); i < 16; i++ {
			d.bus.Write(d.Dst+i, d.bus.Read(d.Src+i))
		}
		d.Src += 16
		d.Dst += 16
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
		d.bus.Write(d.Dst, d.bus.Read(d.Src))
		d.Src++
		d.Dst++
		d.length--
	}
	if d.length == 0 {
		d.pendingHDMA = false
		d.completed = true
	}
}

type DMASnapshot struct {
	Header      uint64
	Mode        uint8
	Src, Dst    uint16
	Length      uint16
	Completed   bool
	DoHDMA      bool
	PendingHDMA bool
	Reserved    [7]uint8
}

func (d *DMA) CreateSnapshot() DMASnapshot {
	return DMASnapshot{
		Mode:        d.Mode,
		Src:         d.Src,
		Dst:         d.Dst,
		Length:      d.length,
		Completed:   d.completed,
		DoHDMA:      d.doHDMA,
		PendingHDMA: d.pendingHDMA,
	}
}

func (d *DMA) RestoreSnapshot(snap DMASnapshot) bool {
	d.Mode = snap.Mode
	d.Src, d.Dst = snap.Src, snap.Dst
	d.length = snap.Length
	d.completed, d.doHDMA, d.pendingHDMA = snap.Completed, snap.DoHDMA, snap.PendingHDMA
	return true
}
