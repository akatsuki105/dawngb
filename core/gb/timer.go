package gb

import (
	"github.com/akatsuki105/dugb/util"
	"github.com/akatsuki105/dugb/util/sched"
)

var timaClock = [4]int64{256, 4, 16, 64}

type timer struct {
	g                   *GB
	div, tima, tma, tac uint8
	divEvent            sched.Event
	timaEvent           sched.Event
	overflowEvent       sched.Event
}

func newTimer(g *GB) *timer {
	t := &timer{
		g: g,
	}
	t.divEvent = *sched.NewEvent("GB_DIV", t.incrementDiv, 0x40)
	t.timaEvent = *sched.NewEvent("GB_TIMA", t.incrementTima, 0x41)
	t.overflowEvent = *sched.NewEvent("GB_OVERFLOW", t.overflowTima, 0x42)
	return t
}

func (t *timer) Reset() {
	t.div = 0x0
	t.tima = 0
	t.tma = 0
	t.tac = 0
	t.g.s.Schedule(&t.divEvent, 64*t.g.cpu.Cycle)
}

func (t *timer) ReadIO(addr uint16) uint8 {
	switch addr {
	case 0xFF04:
		return t.div
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	}
	return 0
}

func (t *timer) WriteIO(addr uint16, val uint8) {
	switch addr {
	case 0xFF04:
		t.div = 0
	case 0xFF05:
		t.tima = val
	case 0xFF06:
		t.tma = val
	case 0xFF07:
		enabled := util.Bit(t.tac, 2)
		t.tac = val & 0b111
		if enabled && !util.Bit(t.tac, 2) {
			t.g.s.Cancel(&t.timaEvent)
		} else if !enabled && util.Bit(t.tac, 2) {
			t.triggerTima()
		}
	}
}

func (t *timer) incrementDiv(cyclesLate int64) {
	t.div++
	t.g.s.Schedule(&t.divEvent, (64*t.g.cpu.Cycle)-cyclesLate)
}

func (t *timer) incrementTima(cyclesLate int64) {
	t.tima++
	if t.tima == 0 {
		t.g.s.Schedule(&t.overflowEvent, t.g.cpu.Cycle-cyclesLate)
	} else {
		clock := timaClock[t.tac&0b11] * t.g.cpu.Cycle
		t.g.s.Schedule(&t.timaEvent, clock-cyclesLate)
	}
}

func (t *timer) triggerTima() {
	clock := timaClock[t.tac&0b11] * t.g.cpu.Cycle
	t.g.s.Schedule(&t.timaEvent, clock)
}

func (t *timer) overflowTima(cyclesLate int64) {
	t.tima = t.tma
	t.g.requestInterrupt(2)

	clock := timaClock[t.tac&0b11] * t.g.cpu.Cycle
	t.g.s.Schedule(&t.timaEvent, clock-cyclesLate)
}
