package audio

type sweep struct {
	enabled bool
	square  *square

	speed int // NR10's bit6-4
	up    bool
	shift int

	step int // 0 ~ speed
}

func newSweep(ch *square) *sweep {
	return &sweep{
		square: ch,
		up:     true,
		speed:  8,
	}
}

func (s *sweep) reset() {
	s.step = s.speed
	s.up = true
	s.enabled = (s.speed != 8 || s.shift != 0)
}

func (s *sweep) update() bool {
	if s.enabled {
		s.step--
		if s.step <= 0 {
			if !s.updateFrequency(true) {
				return false
			}
			s.step = s.speed
		}
	}

	return true
}

func (s *sweep) updateFrequency(first bool) bool {
	if !first || s.speed != 8 {
		period := s.square.period

		if s.up {
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
