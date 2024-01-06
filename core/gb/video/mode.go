package video

import (
	"github.com/akatsuki105/dugb/util"
)

// Mode 0
func (v *Video) hblank(cyclesLate int64) {
	v.stat = (v.stat & 0xFC)
	if util.Bit(v.lcdc, 7) {
		v.r.DrawScanline(v.ly, v.screen[v.ly*160:(v.ly+1)*160])
	}
	if util.Bit(v.stat, 3) {
		v.onInterrupt(1)
	}

	switch v.ly {
	case 143:
		v.s.Schedule(&v.events[1], ((204-v.renderingCycle)*CYCLE)-cyclesLate)
	default:
		v.s.Schedule(&v.events[2], ((204-v.renderingCycle)*CYCLE)-cyclesLate)
	}
}

// Mode 1
func (v *Video) vblank(cyclesLate int64) {
	v.stat = (v.stat & 0xFC) | 1
	if v.ly == 144 {
		v.onInterrupt(0)
		if util.Bit(v.stat, 4) {
			v.onInterrupt(1)
		}
	}

	v.setLy(v.ly + 1)
	switch v.ly {
	case 153:
		v.s.Schedule(&v.events[2], (456*CYCLE)-cyclesLate)
	default:
		v.s.Schedule(&v.events[1], (456*CYCLE)-cyclesLate)
	}
}

// Mode 2
func (v *Video) scanOAM(cyclesLate int64) {
	v.stat = (v.stat & 0xFC) | 2
	if util.Bit(v.stat, 5) {
		v.onInterrupt(1)
	}
	v.setLy(v.ly + 1)
	v.s.Schedule(&v.events[3], (80*CYCLE)-cyclesLate)
}

// Mode 3
func (v *Video) drawing(cyclesLate int64) {
	v.stat = (v.stat & 0xFC) | 3
	v.renderingCycle = 0
	v.s.Schedule(&v.events[0], ((172+v.renderingCycle)*CYCLE)-cyclesLate)
}
