package gb

import (
	"io"

	"github.com/akatsuki105/dawngb/util"
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
	doHDMA    bool
}

func newDMAController(g *GB) *dmaController {
	return &dmaController{
		g:         g,
		completed: true,
	}
}

func (d *dmaController) Reset(hasBIOS bool) {
	d.mode = GDMA
	d.src, d.dst, d.length = 0, 0, 0
	d.completed = true
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
		if wasCompleted && d.mode == GDMA { // Trigger GDMA
			d.runGDMA()
		} else if d.mode == HDMA { // Trigger HDMA
			d.g.runHDMA = d.runHDMA
		}
	}
}

func (d *dmaController) runGDMA() {
	period := int64(d.length) * 4
	for d.length > 0 {
		for i := uint16(0); i < 16; i++ {
			d.g.video.Write(d.dst+i, d.g.m.Read(d.src+i))
		}
		d.src += 16
		d.dst += 16
		d.length -= 16
	}
	d.g.tick(period)
}

// HBlank になるたびに実行される
func (d *dmaController) runHDMA() {
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
	d.g.tick(64)
}

func (d *dmaController) Serialize(s io.Writer) {
	data := []uint8{uint8(d.mode), uint8(d.src >> 8), uint8(d.src), uint8(d.dst >> 8), uint8(d.dst), uint8(d.length >> 8), uint8(d.length), util.Btou8(d.completed)}
	s.Write(data)
	// TODO: schedule
}

func (d *dmaController) Deserialize(s io.Reader) {
	data := make([]uint8, 7)
	s.Read(data)
	d.mode, d.src, d.dst, d.length, d.completed = uint8(data[0]), uint16(data[1])<<8|uint16(data[2]), uint16(data[3])<<8|uint16(data[4]), uint16(data[5])<<8|uint16(data[6]), (data[7] != 0)
	// TODO: schedule
}
