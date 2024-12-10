package psg

import (
	"encoding/binary"
	"io"
)

type Noise struct {
	enabled bool // NR52.3

	length uint8 // 音の残り再生時間
	stop   bool  // NR44.6

	envelope *envelope

	LFSR uint16 // Noiseの疑似乱数(lfsr: Linear Feedback Shift Register = 疑似乱数生成アルゴリズム)

	// この2つでノイズの周波数(疑似乱数の生成頻度)を決める
	divisor uint8 // NR43.0-2; ノイズ周波数1(カウント指定)
	octave  uint8 // NR43.4-7; ノイズ周波数2(オクターブ指定)
	period  uint32

	narrow bool // NR43.3; 0: 15bit, 1: 7bit

	output uint8 // 0..15
}

func newNoiseChannel() *Noise {
	return &Noise{
		envelope: newEnvelope(),
		LFSR:     0,
	}
}

func (ch *Noise) reset() {
	ch.enabled = false
	ch.length, ch.stop = 0, false
	ch.envelope.reset()
	ch.LFSR = 0
	ch.divisor, ch.octave = 0, 0
	ch.period = 0
	ch.narrow = false
	ch.output = 0
}

func (ch *Noise) reload() {
	ch.enabled = ch.dacEnable()
	ch.envelope.reload()
	ch.LFSR = 0x7FFF
	if ch.length == 0 {
		ch.length = 64
	}
}

func (ch *Noise) clock64Hz() {
	if ch.enabled {
		ch.envelope.update()
	}
}

func (ch *Noise) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length == 0 {
			ch.enabled = false
		}
	}
}

func (ch *Noise) clockTimer() {
	// ch.enabledに関わらず、乱数は生成される
	ch.period--
	if ch.period == 0 {
		ch.period = ch.calcFreqency()
		ch.update()
	}

	result := uint8(0)
	if (ch.LFSR & 1) == 0 {
		result = ch.envelope.volume
	}

	if !ch.enabled {
		result = 0
	}

	ch.output = result
}

func (ch *Noise) update() {
	if ch.octave < 14 {
		bit := ((ch.LFSR ^ (ch.LFSR >> 1)) & 1)
		if ch.narrow {
			ch.LFSR = (ch.LFSR >> 1) ^ (bit << 6)
		} else {
			ch.LFSR = (ch.LFSR >> 1) ^ (bit << 14)
		}
	}
}

var noisePeriodTable = []uint8{4, 8, 16, 24, 32, 40, 48, 56}

func (ch *Noise) calcFreqency() uint32 {
	return uint32(noisePeriodTable[ch.divisor]) << ch.octave
}

// GetOutput gets 4bit sample (0..15)
func (ch *Noise) GetOutput() uint8 {
	if ch.enabled {
		return ch.output
	}
	return 0
}

func (ch *Noise) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
}

func (ch *Noise) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, ch.enabled)
	binary.Write(s, binary.LittleEndian, ch.length)
	binary.Write(s, binary.LittleEndian, ch.stop)
	ch.envelope.serialize(s)
	binary.Write(s, binary.LittleEndian, ch.LFSR)
	binary.Write(s, binary.LittleEndian, ch.octave)
	binary.Write(s, binary.LittleEndian, ch.divisor)
	binary.Write(s, binary.LittleEndian, ch.period)
	binary.Write(s, binary.LittleEndian, ch.narrow)
}

func (ch *Noise) deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &ch.enabled)
	binary.Read(s, binary.LittleEndian, &ch.length)
	binary.Read(s, binary.LittleEndian, &ch.stop)
	ch.envelope.deserialize(s)
	binary.Read(s, binary.LittleEndian, &ch.LFSR)
	binary.Read(s, binary.LittleEndian, &ch.octave)
	binary.Read(s, binary.LittleEndian, &ch.divisor)
	binary.Read(s, binary.LittleEndian, &ch.period)
	binary.Read(s, binary.LittleEndian, &ch.narrow)
}
