package apu

import (
	"encoding/binary"
	"io"
)

type noise struct {
	enabled bool
	ignored bool // Ignore sample output

	length int32 // 音の残り再生時間
	stop   bool

	envelope *envelope

	lfsr uint16 // Noiseの疑似乱数(lfsr: Linear Feedback Shift Register = 疑似乱数生成アルゴリズム)

	// この2つでノイズの周波数(疑似乱数の生成頻度)を決める
	octave  int32 // ノイズ周波数2(オクターブ指定)
	divisor int32 // ノイズ周波数1(カウント指定)
	period  int32

	width int32

	output int
}

func newNoiseChannel() *noise {
	return &noise{
		ignored:  true,
		envelope: newEnvelope(),
		lfsr:     1,
		width:    15,
	}
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
	result := 0
	ch.period--
	if ch.period <= 0 {
		ch.period = ch.calcFreqency()
		if ch.octave < 14 {
			mask := ((ch.lfsr ^ (ch.lfsr >> 1)) & 1)
			ch.lfsr = ((ch.lfsr >> 1) ^ (mask << (ch.width - 1))) & 0x7FFF
		}
	}

	if (ch.lfsr & 1) == 0 {
		result = int(ch.envelope.volume)
	}

	if !ch.enabled {
		result = 0
	}

	ch.output = result
}

func (ch *noise) calcFreqency() int32 {
	freq := int32(1)
	if ch.divisor != 0 {
		freq = 2 * ch.divisor
	}
	freq <<= ch.octave
	return freq * 8
}

func (ch *noise) getOutput() int {
	if !ch.ignored {
		if ch.enabled {
			return ch.output
		}
	}
	return 0
}

func (ch *noise) tryRestart() {
	ch.enabled = ch.dacEnable()
	ch.envelope.reset()
	if ch.length == 0 {
		ch.length = 64
	}
	ch.lfsr = 0x7FFF >> (15 - ch.width)
}

func (ch *noise) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
}

func (ch *noise) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, ch.enabled)
	binary.Write(s, binary.LittleEndian, ch.ignored)
	binary.Write(s, binary.LittleEndian, ch.length)
	binary.Write(s, binary.LittleEndian, ch.stop)
	ch.envelope.serialize(s)
	binary.Write(s, binary.LittleEndian, ch.lfsr)
	binary.Write(s, binary.LittleEndian, ch.octave)
	binary.Write(s, binary.LittleEndian, ch.divisor)
	binary.Write(s, binary.LittleEndian, ch.period)
	binary.Write(s, binary.LittleEndian, ch.width)
}

func (ch *noise) deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &ch.enabled)
	binary.Read(s, binary.LittleEndian, &ch.ignored)
	binary.Read(s, binary.LittleEndian, &ch.length)
	binary.Read(s, binary.LittleEndian, &ch.stop)
	ch.envelope.deserialize(s)
	binary.Read(s, binary.LittleEndian, &ch.lfsr)
	binary.Read(s, binary.LittleEndian, &ch.octave)
	binary.Read(s, binary.LittleEndian, &ch.divisor)
	binary.Read(s, binary.LittleEndian, &ch.period)
	binary.Read(s, binary.LittleEndian, &ch.width)
}
