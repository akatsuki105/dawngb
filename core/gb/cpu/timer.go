package cpu

var timaClock = [4]int64{64, 1, 4, 16}

type Timer struct {
	irq            func(n int)
	clock          *int64
	cycles         int64 // CPUから見て遅れているマスターサイクル数
	TIMA, TMA, TAC uint8
	counter        int64 // 524288Hz(一番細かいのが524288Hzなのであとはそれの倍数で数えれば良い)
}

func newTimer(irq func(n int), clock *int64) *Timer {
	return &Timer{
		irq:   irq,
		clock: clock,
	}
}

func (t *Timer) reset() {
	t.cycles = 0
	t.counter, t.TIMA, t.TMA, t.TAC = 0, 0, 0, 0
}

func (t *Timer) run(cycles8MHz int64) {
	t.cycles += cycles8MHz
	for t.cycles >= 16 {
		t.update()
		t.cycles -= 16
	}
}

// 524288Hz
func (t *Timer) update() {
	x := (*t.clock) / 4

	t.counter++
	if (t.TAC & (1 << 2)) != 0 {
		if (t.counter % (timaClock[t.TAC&0b11] * x)) == 0 {
			t.TIMA++
			if t.TIMA == 0 {
				t.TIMA = t.TMA
				t.irq(IRQ_TIMER)
			}
		}
	}
}

func (t *Timer) Read(addr uint16) uint8 {
	x := (*t.clock) / 4
	switch addr {
	case 0xFF04:
		div := t.counter / (16 * x)
		return uint8(div)
	case 0xFF05:
		return t.TIMA
	case 0xFF06:
		return t.TMA
	case 0xFF07:
		return t.TAC
	}
	return 0
}

func (t *Timer) Write(addr uint16, val uint8) {
	switch addr {
	case 0xFF04:
		t.counter = 0
	case 0xFF05:
		t.TIMA = val
	case 0xFF06:
		t.TMA = val
	case 0xFF07:
		t.TAC = val & 0b111
	}
}

type TimerSnapshot struct {
	Header         uint64
	Cycles         int64
	Tima, Tma, Tac uint8
	Counter        int64
	Reserved       [7]uint8
}

func (t *Timer) CreateSnapshot() TimerSnapshot {
	return TimerSnapshot{
		Cycles:  t.cycles,
		Tima:    t.TIMA,
		Tma:     t.TMA,
		Tac:     t.TAC,
		Counter: t.counter,
	}
}

func (t *Timer) RestoreSnapshot(snap TimerSnapshot) bool {
	t.cycles = snap.Cycles
	t.TIMA, t.TMA, t.TAC = snap.Tima, snap.Tma, snap.Tac
	t.counter = snap.Counter
	return true
}
