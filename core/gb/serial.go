package gb

import (
	"math"

	"github.com/akatsuki105/dawngb/core/gb/cpu"
	"github.com/akatsuki105/dawngb/util"
)

type serial struct {
	irq    func(int)
	until  int64
	sb, sc uint8
}

func newSerial(irq func(int)) *serial {
	return &serial{irq: irq}
}

func (s *serial) Reset(hasBIOS bool) {
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

func (s *serial) Read(addr uint16) uint8 {
	switch addr {
	case 0xFF01:
		return s.sb
	case 0xFF02:
		return s.sc
	}
	return 0
}

func (s *serial) Write(addr uint16, val uint8) {
	switch addr {
	case 0xFF01:
		s.sb = val
	case 0xFF02:
		s.sc = val
		// 一部のソフト(ポケモンクリスタル など)は起動時にシリアル通信を実装してないと動かないのでダミーで実装
		if util.Bit(val, 7) {
			s.until = math.MaxInt64
			if util.Bit(val, 0) {
				s.until = 512 * 8
			}
		}
	}
}

// ポケモンクリスタルの起動にシリアル通信機能が必要なので暫定措置
func (s *serial) dummyTransfer() {
	s.sc &= 0x7F
	s.sb = 0xFF
	s.irq(cpu.IRQ_SERIAL)
}
