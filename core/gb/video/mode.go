package video

// Mode 0
func (v *Video) hblank(cyclesLate int64) {
	v.mode = 0
	switch v.ly {
	case 143:
		v.s.Schedule(&v.events[1], ((204-v.renderingCycle)*CYCLE)-cyclesLate)
	default:
		v.s.Schedule(&v.events[2], ((204-v.renderingCycle)*CYCLE)-cyclesLate)
	}
}

// Mode 1
func (v *Video) vblank(cyclesLate int64) {
	v.mode = 1
	v.ly++
	switch v.ly {
	case 154:
		v.ly = 0
		v.FrameCounter++
		v.s.Schedule(&v.events[2], (456*CYCLE)-cyclesLate)
	default:
		v.s.ReSchedule(&v.events[1], (456*CYCLE)-cyclesLate)
	}
}

// Mode 2
func (v *Video) scanOAM(cyclesLate int64) {
	v.mode = 2
	v.ly++
	v.s.Schedule(&v.events[3], (80*CYCLE)-cyclesLate)
}

// Mode 3
func (v *Video) drawing(cyclesLate int64) {
	v.mode = 3
	v.renderingCycle = 0
	v.s.Schedule(&v.events[0], ((172+v.renderingCycle)*CYCLE)-cyclesLate)
}
