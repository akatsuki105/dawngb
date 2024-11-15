package psg

import (
	"encoding/binary"
	"io"
)

type sweep struct {
	enabled bool
	square  *square

	// スイープ開始時の周波数(スイープ中に0xFF13と0xFF14に書き込まれて.square.periodが変更されても影響を受けないようにするためのもの)
	periodLatch uint16

	shift    uint8 // NR10.0-2; スイープ量
	negate   bool  // NR10.3; スイープの方向(0: 加算, 1: 減算)
	interval uint8 // NR10.4-6; スイープ間隔

	step uint8 // スイープ間隔(.interval)をカウントするためのカウンタ
}

func newSweep(ch *square) *sweep {
	return &sweep{
		square: ch,
	}
}

func (s *sweep) reset() {
	s.enabled = false
	s.periodLatch, s.interval, s.negate, s.shift = 0, 0, false, 0
	s.step = 0
}

func (s *sweep) reload() {
	s.periodLatch = s.square.period
	s.step = s.interval
	if s.interval == 0 {
		s.step = 8
	}
	s.enabled = (s.interval != 0 || s.shift != 0)
	if s.shift != 0 {
		s.checkOverflow()
	}
}

func (s *sweep) update() {
	s.step--
	if s.step == 0 {
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
		delta := s.periodLatch >> s.shift
		freq := s.periodLatch
		if !s.negate {
			freq += delta
		} else {
			freq -= delta
		}

		if freq > 2047 {
			s.square.enabled = false
		} else if s.shift != 0 {
			s.periodLatch = freq
			s.square.period = freq
			s.square.freqCounter = s.square.dutyStepCycle()
		}
	}
}

func (s *sweep) checkOverflow() {
	if s.enabled {
		delta := s.periodLatch >> s.shift
		freq := s.periodLatch
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

func (s *sweep) serialize(w io.Writer) {
	binary.Write(w, binary.LittleEndian, s.enabled)
	binary.Write(w, binary.LittleEndian, s.periodLatch)
	binary.Write(w, binary.LittleEndian, s.interval)
	binary.Write(w, binary.LittleEndian, s.negate)
	binary.Write(w, binary.LittleEndian, s.shift)
	binary.Write(w, binary.LittleEndian, s.step)
}

func (s *sweep) deserialize(r io.Reader) {
	binary.Read(r, binary.LittleEndian, &s.enabled)
	binary.Read(r, binary.LittleEndian, &s.periodLatch)
	binary.Read(r, binary.LittleEndian, &s.interval)
	binary.Read(r, binary.LittleEndian, &s.negate)
	binary.Read(r, binary.LittleEndian, &s.shift)
	binary.Read(r, binary.LittleEndian, &s.step)
}
