package psg

import (
	"encoding/binary"
	"io"
)

var squareDutyTable = [4][8]uint8{
	{0, 0, 0, 0, 0, 0, 0, 1}, // 12.5%
	{1, 0, 0, 0, 0, 0, 0, 1}, // 25%
	{1, 0, 0, 0, 0, 1, 1, 1}, // 50%
	{0, 1, 1, 1, 1, 1, 1, 0}, // 75%
}

type square struct {
	enabled bool // NR52.0(ch1), NR52.1(ch2)

	length uint8 // NR11.0-5; 音の残り再生時間
	stop   bool  // NR14.6; .length が 0 になったときに音を止めるかどうか

	envelope *envelope
	sweep    *sweep

	duty        uint8 // NR11.6-7, (squareDutyTable の index)
	dutyCounter uint8 // 0 ~ 7

	period      uint16 // NR13.0-7, NR14.0-2; GBでは周波数を指定するのではなく、周期の長さを指定する
	freqCounter uint16

	output bool // 0 or 1; 矩形波の出力が1かどうか
}

func newSquareChannel(hasSweep bool) *square {
	ch := &square{
		envelope: newEnvelope(),
	}

	if hasSweep { // スイープ機能があるのは ch1 のみなので区別する必要がある
		ch.sweep = newSweep(ch)
	}
	return ch
}

func (ch *square) reset() {
	ch.enabled = false
	ch.length, ch.stop = 0, false
	ch.envelope.reset()
	if ch.sweep != nil {
		ch.sweep.reset()
	}
	ch.duty, ch.dutyCounter = 0, 0
	ch.period, ch.freqCounter = 0, 0
	ch.output = false
}

func (ch *square) reload() {
	ch.enabled = ch.dacEnable()
	ch.freqCounter = ch.dutyStepCycle()
	ch.envelope.reload()
	if ch.sweep != nil {
		ch.sweep.reload()
	}
	if ch.length == 0 {
		ch.length = 64
	}
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
		if ch.length == 0 {
			ch.enabled = false
		}
	}
}

func (ch *square) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
		if ch.freqCounter == 0 {
			ch.freqCounter = ch.dutyStepCycle()
			ch.update()
		}
	}
}

func (ch *square) update() {
	ch.dutyCounter = (ch.dutyCounter + 1) & 7
	dutyTable := (squareDutyTable[ch.duty])[:]
	ch.output = dutyTable[ch.dutyCounter] != 0
}

func (ch *square) getOutput() uint8 {
	if ch.enabled && ch.output {
		return ch.envelope.volume
	}
	return 0
}

// デューティ比の1ステップの長さをAPUサイクル数で返す
func (ch *square) dutyStepCycle() uint16 {
	// hz := (1048576 / (2048 - ch.period)) // freqency
	// return 4194304 / hz
	return 4 * (2048 - ch.period)
}

func (ch *square) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
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
