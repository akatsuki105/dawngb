package psg

var squareDutyTable = [4][8]uint8{
	{0, 0, 0, 0, 0, 0, 0, 1}, // 12.5%
	{1, 0, 0, 0, 0, 0, 0, 1}, // 25%
	{1, 0, 0, 0, 0, 1, 1, 1}, // 50%
	{0, 1, 1, 1, 1, 1, 1, 0}, // 75%
}

type Square struct {
	enabled bool // NR52.0(ch1), NR52.1(ch2)

	length uint8 // NR11.0-5; 音の残り再生時間
	stop   bool  // NR14.6; .length が 0 になったときに音を止めるかどうか

	envelope *Envelope
	sweep    *sweep

	duty        uint8 // NR11.6-7, (squareDutyTable の index)
	dutyCounter uint8 // 0 ~ 7

	period      uint16 // NR13.0-7, NR14.0-2; GBでは周波数を指定するのではなく、周期の長さを指定する
	freqCounter uint16

	output bool // 0 or 1; 矩形波の出力が1かどうか
}

func newSquareChannel(hasSweep bool) *Square {
	ch := &Square{
		envelope: newEnvelope(),
	}

	if hasSweep { // スイープ機能があるのは ch1 のみなので区別する必要がある
		ch.sweep = newSweep(ch)
	}
	return ch
}

func (ch *Square) Reset() {
	ch.TurnOff()
	ch.envelope.reset()
	if ch.sweep != nil {
		ch.sweep.reset()
	}
	ch.period, ch.freqCounter = 0, 0
}

func (ch *Square) TurnOff() {
	ch.enabled = false
	ch.length, ch.stop = 0, false
	ch.duty, ch.dutyCounter = 0, 0
	ch.envelope.TurnOff()
	if ch.sweep != nil {
		ch.sweep.TurnOff()
	}
	ch.output = false
}

func (ch *Square) reload() {
	ch.enabled = ch.dacEnable()
	ch.freqCounter = ch.dutyStepCycle()
	ch.envelope.reload()
	if ch.sweep != nil {
		ch.sweep.reload()
	}
	if ch.length == 0 {
		ch.length = 64
	}
}

func (ch *Square) clock64Hz() {
	if ch.enabled {
		ch.envelope.update()
	}
}

func (ch *Square) clock128Hz() {
	if ch.sweep != nil {
		ch.sweep.update()
	}
}

func (ch *Square) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length == 0 {
			ch.enabled = false
		}
	}
}

func (ch *Square) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
		if ch.freqCounter == 0 {
			ch.freqCounter = ch.dutyStepCycle()
			ch.update()
		}
	}
}

func (ch *Square) update() {
	ch.dutyCounter = (ch.dutyCounter + 1) & 7
	dutyTable := (squareDutyTable[ch.duty])[:]
	ch.output = dutyTable[ch.dutyCounter] != 0
}

// GetOutput gets 4bit sample (0..15)
func (ch *Square) GetOutput() uint8 {
	if ch.enabled && ch.output {
		return ch.envelope.volume
	}
	return 0
}

// デューティ比の1ステップの長さをAPUサイクル数で返す
func (ch *Square) dutyStepCycle() uint16 {
	// hz := (1048576 / (2048 - ch.period)) // freqency
	// return 4194304 / hz
	return 4 * (2048 - ch.period)
}

func (ch *Square) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
}

type SquareSnapshot struct {
	Header              uint64
	Enabled             bool
	Length              uint8
	Stop                bool
	Envelope            EnvelopeSnapshot
	Duty, DutyCounter   uint8
	Period, FreqCounter uint16
	Output              bool
	Sweep               SweepSnapshot
	Reserved            [16]uint8
}

func (ch *Square) CreateSnapshot() SquareSnapshot {
	s := SquareSnapshot{
		Enabled:     ch.enabled,
		Length:      ch.length,
		Stop:        ch.stop,
		Envelope:    ch.envelope.CreateSnapshot(),
		Duty:        ch.duty,
		DutyCounter: ch.dutyCounter,
		Period:      ch.period,
		FreqCounter: ch.freqCounter,
		Output:      ch.output,
	}
	if ch.sweep != nil {
		s.Sweep = ch.sweep.CreateSnapshot()
	}
	return s
}

func (ch *Square) RestoreSnapshot(snap SquareSnapshot) bool {
	ch.enabled = snap.Enabled
	ch.length, ch.stop = snap.Length, snap.Stop
	ch.duty, ch.dutyCounter = snap.Duty, snap.DutyCounter
	ch.period, ch.freqCounter = snap.Period, snap.FreqCounter
	ch.output = snap.Output
	if ch.sweep != nil {
		ch.sweep.RestoreSnapshot(snap.Sweep)
	}
	ch.envelope.RestoreSnapshot(snap.Envelope)
	return true
}
