package gb

import (
	"encoding/binary"
	"io"

	"github.com/akatsuki105/dawngb/util"
	"github.com/akatsuki105/dawngb/util/sched"
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

func (t *timer) Reset(hasBIOS bool) {
	t.counter, t.tima, t.tma, t.tac = 0, 0, 0, 0
	if !hasBIOS {
		t.tac = 0xF8
	}
	t.g.s.Reschedule(&t.updateEvent, 16)
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

	t.g.s.Reschedule(&t.updateEvent, 16-cyclesLate) // 524288Hz
}

func (t *timer) Serialize(s io.Writer) {
	data := []uint8{t.tima, t.tma, t.tac}
	binary.LittleEndian.PutUint64(data, uint64(t.counter))
	binary.LittleEndian.PutUint64(data, uint64(t.g.s.Until(&t.updateEvent)))
	s.Write(data)
}

func (t *timer) Deserialize(s io.Reader) {
	data := make([]uint8, 8)
	s.Read(data)
	t.tima, t.tma, t.tac = data[0], data[1], data[2]
	t.counter = int64(binary.LittleEndian.Uint64(data))
	t.g.s.Reschedule(&t.updateEvent, int64(binary.LittleEndian.Uint64(data)))
}
