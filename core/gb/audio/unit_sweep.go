package audio

type sweep struct {
	enabled bool
	square  *square

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
	s.step = s.interval
	if s.interval == 0 {
		s.step = 8
	}
	s.negate = false
	s.enabled = (s.shift != 0)
}

func (s *sweep) update() bool {
	if s.enabled {
		s.step--
		if s.step <= 0 {
			if !s.updateFrequency(true) {
				return false
			}
			s.step = s.interval
			if s.interval == 0 {
				s.step = 8
			}
		}
	}

	return true
}

func (s *sweep) updateFrequency(first bool) bool {
	if !first || s.interval != 0 {
		period := s.square.period

		if !s.negate {
			period += (period >> s.shift)
			if period < 2048 {
				if first && s.shift != 0 {
					s.square.period = period
					s.square.freqCounter = s.square.dutyStepCycle()
					if !s.updateFrequency(false) {
						return false
					}
				}
			} else {
				return false
			}
		} else {
			period -= (period >> s.shift)
			if first && period >= 0 {
				s.square.period = period
				s.square.freqCounter = s.square.dutyStepCycle()
			}
		}
	}

	return true
}
