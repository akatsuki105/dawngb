package psg

import (
	"encoding/binary"
	"io"
)

var squareDutyTable = [4][8]int{
	{0, 0, 0, 0, 0, 0, 0, 1}, // 12.5%
	{1, 0, 0, 0, 0, 0, 0, 1}, // 25%
	{1, 0, 0, 0, 0, 1, 1, 1}, // 50%
	{0, 1, 1, 1, 1, 1, 1, 0}, // 75%
}

type square struct {
	enabled bool

	length int32 // 音の残り再生時間
	stop   bool  // .length が 0 になったときに 音を止めるかどうか(NR14's bit6)

	envelope *envelope
	sweep    *sweep

	duty        uint8 // NR11's bit7-6, (squareDutyTable の index)
	dutyCounter uint8 // 0 ~ 7

	period      int32 // GBでは周波数を指定するのではなく、周期の長さを指定する, 実際の周波数は ((4194304/32)/(2048-period)) Hz (64~131072 Hz -> 65536~32 APUサイクル)
	freqCounter int32
}

func newSquareChannel(hasSweep bool) *square {
	ch := &square{
		envelope: newEnvelope(),
	}

	// スイープ機能があるのは ch1 のみなので区別する必要がある
	if hasSweep {
		ch.sweep = newSweep(ch)
	}
	return ch
}

func (ch *square) clock64Hz() {
	if ch.enabled {
		ch.envelope.update()
	}
}

func (ch *square) clock128Hz() {
	if ch.sweep != nil {
		ch.sweep.update()
	}
}

func (ch *square) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length <= 0 {
			ch.enabled = false
		}
	}
}

func (ch *square) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
	} else {
		ch.freqCounter = ch.dutyStepCycle()
		ch.dutyCounter = (ch.dutyCounter + 1) % 8
	}
}

func (ch *square) getOutput() uint8 {
	if ch.enabled {
		dutyTable := (squareDutyTable[ch.duty])[:]
		if dutyTable[ch.dutyCounter] != 0 {
			return ch.envelope.volume
		}
	}
	return 0
}

// デューティ比の1ステップの長さをAPUサイクル数で返す
func (ch *square) dutyStepCycle() int32 {
	// hz := (1048576 / (2048 - ch.period)) // freqency
	// return 4194304 / hz
	return 4 * (2048 - ch.period)
}

func (ch *square) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
}

func (ch *square) tryRestart() {
	ch.enabled = ch.dacEnable()
	ch.freqCounter = ch.dutyStepCycle()
	ch.envelope.reset()
	if ch.sweep != nil {
		ch.sweep.reset()
	}
	if ch.length == 0 {
		ch.length = 64
	}
}

func (ch *square) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, ch.enabled)
	binary.Write(s, binary.LittleEndian, ch.length)
	binary.Write(s, binary.LittleEndian, ch.stop)
	ch.envelope.serialize(s)
	if ch.sweep != nil {
		ch.sweep.serialize(s)
	}
	binary.Write(s, binary.LittleEndian, ch.duty)
	binary.Write(s, binary.LittleEndian, ch.dutyCounter)
	binary.Write(s, binary.LittleEndian, ch.period)
	binary.Write(s, binary.LittleEndian, ch.freqCounter)
}

func (ch *square) deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &ch.enabled)
	binary.Read(s, binary.LittleEndian, &ch.length)
	binary.Read(s, binary.LittleEndian, &ch.stop)
	ch.envelope.deserialize(s)
	if ch.sweep != nil {
		ch.sweep.deserialize(s)
	}
	binary.Read(s, binary.LittleEndian, &ch.duty)
	binary.Read(s, binary.LittleEndian, &ch.dutyCounter)
	binary.Read(s, binary.LittleEndian, &ch.period)
	binary.Read(s, binary.LittleEndian, &ch.freqCounter)
}
