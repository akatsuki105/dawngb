package psg

import (
	"encoding/binary"
	"io"
)

type noise struct {
	enabled bool

	length int32 // 音の残り再生時間
	stop   bool  // NR44.6

	envelope *envelope

	lfsr uint16 // Noiseの疑似乱数(lfsr: Linear Feedback Shift Register = 疑似乱数生成アルゴリズム)

	// この2つでノイズの周波数(疑似乱数の生成頻度)を決める
	divisor uint8 // ノイズ周波数1(カウント指定)
	octave  uint8 // ノイズ周波数2(オクターブ指定)
	period  int32

	narrow bool // NR43.3; 0: 15bit, 1: 7bit

	output uint8 // 0..15
}

func newNoiseChannel() *noise {
	return &noise{
		envelope: newEnvelope(),
		lfsr:     0,
	}
}

func (ch *noise) reset() {
	ch.enabled = false
	ch.length = 0
	ch.stop = false
	ch.envelope.reset()
	ch.lfsr = 0
	ch.divisor, ch.octave = 0, 0
	ch.period = 0
	ch.narrow = false
	ch.output = 0
}

func (ch *noise) clock64Hz() {
	if ch.enabled {
		ch.envelope.update()
	}
}

func (ch *noise) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length <= 0 {
			ch.enabled = false
		}
	}
}

func (ch *noise) clockTimer() {
	// ch.enabledに関わらず、乱数は生成される
	ch.period--
	if ch.period <= 0 {
		ch.period = ch.calcFreqency()
		if ch.octave < 14 {
			bit := ((ch.lfsr ^ (ch.lfsr >> 1)) & 1)
			if ch.narrow {
				ch.lfsr = (ch.lfsr >> 1) ^ (bit << 6)
			} else {
				ch.lfsr = (ch.lfsr >> 1) ^ (bit << 14)
			}
		}
	}

	result := uint8(0)
	if (ch.lfsr & 1) == 0 {
		result = ch.envelope.volume
	}

	if !ch.enabled {
		result = 0
	}

	ch.output = result
}

var noisePeriodTable = []uint8{4, 8, 16, 24, 32, 40, 48, 56}

func (ch *noise) calcFreqency() int32 {
	return int32(noisePeriodTable[ch.divisor]) << ch.octave
}

func (ch *noise) getOutput() uint8 {
	if ch.enabled {
		return ch.output
	}
	return 0
}

func (ch *noise) tryRestart() {
	ch.enabled = ch.dacEnable()
	ch.envelope.reload()
	if ch.length == 0 {
		ch.length = 64
	}
	ch.lfsr = 0x7FFF
}

func (ch *noise) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
}

func (ch *noise) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, ch.enabled)
	binary.Write(s, binary.LittleEndian, ch.length)
	binary.Write(s, binary.LittleEndian, ch.stop)
	ch.envelope.serialize(s)
	binary.Write(s, binary.LittleEndian, ch.lfsr)
	binary.Write(s, binary.LittleEndian, ch.octave)
	binary.Write(s, binary.LittleEndian, ch.divisor)
	binary.Write(s, binary.LittleEndian, ch.period)
	binary.Write(s, binary.LittleEndian, ch.narrow)
}

func (ch *noise) deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &ch.enabled)
	binary.Read(s, binary.LittleEndian, &ch.length)
	binary.Read(s, binary.LittleEndian, &ch.stop)
	ch.envelope.deserialize(s)
	binary.Read(s, binary.LittleEndian, &ch.lfsr)
	binary.Read(s, binary.LittleEndian, &ch.octave)
	binary.Read(s, binary.LittleEndian, &ch.divisor)
	binary.Read(s, binary.LittleEndian, &ch.period)
	binary.Read(s, binary.LittleEndian, &ch.narrow)
}
