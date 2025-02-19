package cpu

import (
	"math"
)

type serial struct {
	irq    func(int)
	until  int64
	sb, sc uint8
}

func newSerial(irq func(int)) *serial {
	return &serial{irq: irq}
}

func (s *serial) reset() {
	s.until = math.MaxInt64
	s.sb, s.sc = 0, 0
}

func (s *serial) run(cycles8MHz int64) {
	for i := int64(0); i < cycles8MHz; i++ {
		s.until--
		if s.until <= 0 {
			s.until = math.MaxInt64
			s.dummyTransfer()
		}
	}
}

func (s *serial) setSC(val uint8) {
	s.sc = val
	// 一部のソフト(ポケモンクリスタル など)は起動時にシリアル通信を実装してないと動かないのでダミーで実装
	if (val & (1 << 7)) != 0 {
		s.until = math.MaxInt64
		if (s.sc & (1 << 0)) != 0 {
			s.until = 512 * 8
		}
	}
}

// ポケモンクリスタルの起動にシリアル通信機能が必要なので暫定措置
func (s *serial) dummyTransfer() {
	s.sc &= 0x7F
	s.sb = 0xFF
	s.irq(IRQ_SERIAL)
}

type SerialSnapshot struct {
	Header   uint64
	Until    int64
	SB, SC   uint8
	Reserved [14]uint8
}

func (s *serial) CreateSnapshot() SerialSnapshot {
	return SerialSnapshot{
		Until: s.until,
		SB:    s.sb,
		SC:    s.sc,
	}
}

func (s *serial) RestoreSnapshot(snap SerialSnapshot) bool {
	s.until = snap.Until
	s.sb, s.sc = snap.SB, snap.SC
	return true
}
