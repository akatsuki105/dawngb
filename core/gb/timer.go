package gb

import (
	"github.com/akatsuki105/dugb/util"
	"github.com/akatsuki105/dugb/util/sched"
)

var timaClock = [4]int64{64, 1, 4, 16}

type timer struct {
	g              *GB
	tima, tma, tac uint8

	updateEvent sched.Event
	counter     int64 // 524288Hz(一番細かいのが524288Hzなのであとはそれの倍数で数えれば良い)
}

func newTimer(g *GB) *timer {
	t := &timer{
		g: g,
	}
	t.updateEvent = *sched.NewEvent("GB_TIMER_UPDATE", t.update)
	return t
}

func (t *timer) Reset() {
	t.counter = 0
	t.tima = 0
	t.tma = 0
	t.tac = 0
	t.g.s.Schedule(&t.updateEvent, 16)
}

func (t *timer) Read(addr uint16) uint8 {
	x := t.g.cpu.Cycle / 4
	switch addr {
	case 0xFF04:
		div := t.counter / (16 * x)
		return uint8(div)
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	}
	return 0
}

func (t *timer) Write(addr uint16, val uint8) {
	switch addr {
	case 0xFF04:
		t.counter = 0
	case 0xFF05:
		t.tima = val
	case 0xFF06:
		t.tma = val
	case 0xFF07:
		t.tac = val & 0b111
	}
}

func (t *timer) update(cyclesLate int64) {
	x := t.g.cpu.Cycle / 4

	t.counter++
	if util.Bit(t.tac, 2) {
		if (t.counter % (timaClock[t.tac&0b11] * x)) == 0 {
			t.tima++
			if t.tima == 0 {
				t.tima = t.tma
				t.g.requestInterrupt(2)
			}
		}
	}

	t.g.s.Schedule(&t.updateEvent, 16-cyclesLate) // 524288Hz
}
