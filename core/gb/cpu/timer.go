package cpu

var timaClock = [4]int64{64, 1, 4, 16}

type timer struct {
	irq            func(n int)
	clock          *int64
	cycles         int64 // CPUから見て遅れているマスターサイクル数
	tima, tma, tac uint8
	counter        int64 // 524288Hz(一番細かいのが524288Hzなのであとはそれの倍数で数えれば良い)
}

func newTimer(irq func(n int), clock *int64) *timer {
	return &timer{
		irq:   irq,
		clock: clock,
	}
}

func (t *timer) reset() {
	t.cycles = 0
	t.counter, t.tima, t.tma, t.tac = 0, 0, 0, 0
}

func (t *timer) run(cycles8MHz int64) {
	t.cycles += cycles8MHz
	for t.cycles >= 16 {
		t.update()
		t.cycles -= 16
	}
}

// 524288Hz
func (t *timer) update() {
	x := (*t.clock) / 4

	t.counter++
	if (t.tac & (1 << 2)) != 0 {
		if (t.counter % (timaClock[t.tac&0b11] * x)) == 0 {
			t.tima++
			if t.tima == 0 {
				t.tima = t.tma
				t.irq(IRQ_TIMER)
			}
		}
	}
}

func (t *timer) Read(addr uint16) uint8 {
	x := (*t.clock) / 4
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

type timerSnapshot struct {
	Cycles         int64
	Counter        int64
	Tima, Tma, Tac uint8
}

func (t *timer) CreateSnapshot() timerSnapshot {
	return timerSnapshot{
		Cycles:  t.cycles,
		Counter: t.counter,
		Tima:    t.tima,
		Tma:     t.tma,
		Tac:     t.tac,
	}
}

func (t *timer) RestoreSnapshot(snap timerSnapshot) bool {
	t.cycles = snap.Cycles
	t.counter = snap.Counter
	t.tima = snap.Tima
	t.tma = snap.Tma
	t.tac = snap.Tac
	return true
}
