package video

import (
	"github.com/akatsuki105/dugb/util"
)

// Mode 0
func (v *Video) hblank(cyclesLate int64) {
	oldStat := v.stat
	v.stat = (v.stat & 0xFC)
	if util.Bit(v.lcdc, 7) {
		v.r.DrawScanline(v.ly, v.screen[v.ly*160:(v.ly+1)*160])
	}
	if !statIRQAsserted(oldStat) && statIRQAsserted(v.stat) {
		v.onInterrupt(1)
	}
}

// Mode 1
func (v *Video) vblank(cyclesLate int64) {
	oldStat := v.stat
	v.stat = (v.stat & 0xFC) | 1
	v.onInterrupt(0)

	if !statIRQAsserted(oldStat) && statIRQAsserted(v.stat) {
		v.onInterrupt(1)
	}
}

// Mode 2
func (v *Video) scanOAM(cyclesLate int64) {
	oldStat := v.stat
	v.stat = (v.stat & 0xFC) | 2
	if !statIRQAsserted(oldStat) && statIRQAsserted(v.stat) {
		v.onInterrupt(1)
	}
}

// Mode 3
func (v *Video) drawing(cyclesLate int64) {
	v.stat = (v.stat & 0xFC) | 3
	v.renderingCycle = 0
}
