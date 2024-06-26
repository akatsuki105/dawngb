package apu

type sweep struct {
	enabled bool
	square  *square

	// スイープ開始時の周波数(スイープ中に0xFF13と0xFF14に書き込まれて.square.periodが変更されても影響を受けないようにするためのもの)
	periodShadow int

	interval int // NR10's bit6-4(スイープ間隔)
	negate   bool
	shift    int

	step int // スイープ間隔(.interval)をカウントするためのカウンタ
}

func newSweep(ch *square) *sweep {
	return &sweep{
		square:   ch,
		interval: 0,
		step:     8,
	}
}

func (s *sweep) reset() {
	s.periodShadow = s.square.period
	s.step = s.interval
	if s.interval == 0 {
		s.step = 8
	}
	s.enabled = (s.interval != 0 || s.shift != 0)
	if s.shift != 0 {
		s.updateFrequency()
		s.checkOverflow()
	}
}

func (s *sweep) update() {
	s.step--
	if s.step <= 0 {
		s.step = s.interval
		if s.interval == 0 {
			s.step = 8
		}
		if s.enabled && s.interval != 0 {
			s.updateFrequency()
			s.checkOverflow()
		}
	}
}

func (s *sweep) updateFrequency() {
	if s.enabled {
		delta := s.periodShadow >> s.shift
		freq := s.periodShadow
		if !s.negate {
			freq += delta
		} else {
			freq -= delta
		}

		if freq > 2047 {
			s.square.enabled = false
		} else if s.shift != 0 && freq >= 0 {
			s.periodShadow = freq
			s.square.period = freq
			s.square.freqCounter = s.square.dutyStepCycle()
		}
	}
}

func (s *sweep) checkOverflow() {
	if s.enabled {
		delta := s.periodShadow >> s.shift
		freq := s.periodShadow
		if !s.negate {
			freq += delta
		} else {
			freq -= delta
		}

		if freq > 2047 {
			s.square.enabled = false
		}
	}
}
