package cpu

import (
	"math"
)

type Serial struct {
	irq    func(int)
	until  int64
	SB, SC uint8
}

func newSerial(irq func(int)) *Serial {
	return &Serial{irq: irq}
}

func (s *Serial) reset() {
	s.until = math.MaxInt64
	s.SB, s.SC = 0, 0
}

func (s *Serial) run(cycles8MHz int64) {
	for i := int64(0); i < cycles8MHz; i++ {
		s.until--
		if s.until <= 0 {
			s.until = math.MaxInt64
			s.dummyTransfer()
		}
	}
}

func (s *Serial) setSC(val uint8) {
	s.SC = val
	// 一部のソフト(ポケモンクリスタル など)は起動時にシリアル通信を実装してないと動かないのでダミーで実装
	if (val & (1 << 7)) != 0 {
		s.until = math.MaxInt64
		if (s.SC & (1 << 0)) != 0 {
			s.until = 512 * 8
		}
	}
}

// ポケモンクリスタルの起動にシリアル通信機能が必要なので暫定措置
func (s *Serial) dummyTransfer() {
	s.SC &= 0x7F
	s.SB = 0xFF
	s.irq(IRQ_SERIAL)
}

type SerialSnapshot struct {
	Header   uint64
	Until    int64
	SB, SC   uint8
	Reserved [14]uint8
}

func (s *Serial) CreateSnapshot() SerialSnapshot {
	return SerialSnapshot{
		Until: s.until,
		SB:    s.SB,
		SC:    s.SC,
	}
}

func (s *Serial) RestoreSnapshot(snap SerialSnapshot) bool {
	s.until = snap.Until
	s.SB, s.SC = snap.SB, snap.SC
	return true
}
