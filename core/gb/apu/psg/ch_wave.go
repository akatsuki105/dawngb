package psg

import (
	"encoding/binary"
	"io"
)

type wave struct {
	enabled bool

	dacEnable bool  // NR30's bit7
	stop      bool  // .length が 0 になったときに 音を止めるかどうか(NR34's bit6)
	length    int32 // 音の残り再生時間
	volume    uint8 // NR32's bit6-5 (0: 0%, 1: 100%, 2: 50%, 3: 25%)

	period      int32 // GBでは周波数を指定するのではなく、周期の長さを指定する
	freqCounter int32

	samples [32]uint8 // 4bit sample
	window  int8      // 0 ~ 31

	// For GBA
	bank     uint8 // NR30.6
	usedBank uint8 // 現在演奏中のバンク、modeが1の場合は、 .bank の値と必ずしも一致しないので
	mode     uint8 //　 0: 16バイト(32サンプル)を演奏に使い、裏のバンクでは読み書きを行う、 1: 32バイト(64サンプル)を全部演奏に使う
}

func newWaveChannel() *wave {
	return &wave{}
}

func (ch *wave) reset() {
	ch.enabled = false
	ch.dacEnable = false
	ch.stop, ch.length = false, 0
	ch.volume = 0
	ch.period, ch.freqCounter = 0, 0
	clear(ch.samples[:])
	ch.window = 0
	ch.bank, ch.usedBank, ch.mode = 0, 0, 0
}

func (ch *wave) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length <= 0 {
			ch.enabled = false
		}
	}
}

func (ch *wave) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
	} else {
		ch.freqCounter = ch.windowStepCycle()
		ch.window = (ch.window + 1) & 0x1F
		if ch.window == 0 {
			ch.usedBank ^= ch.mode
		}
	}
}

func (ch *wave) getOutput() uint8 {
	if ch.enabled && ch.dacEnable {
		isHi := ch.window&1 == 0 // 上位4bit -> 下位4bit -> 上位4bit -> 下位4bit -> ...
		sample := uint8(0)
		if isHi {
			sample = ch.samples[ch.window>>1] >> 4
		} else {
			sample = ch.samples[ch.window>>1] & 0xF
		}
		return sample >> ch.volume
	}
	return 0
}

func (ch *wave) windowStepCycle() int32 {
	return 2 * (2048 - ch.period)
}

func (ch *wave) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, ch.enabled)
	binary.Write(s, binary.LittleEndian, ch.dacEnable)
	binary.Write(s, binary.LittleEndian, ch.stop)
	binary.Write(s, binary.LittleEndian, ch.length)
	binary.Write(s, binary.LittleEndian, ch.volume)
	binary.Write(s, binary.LittleEndian, ch.period)
	binary.Write(s, binary.LittleEndian, ch.freqCounter)
	binary.Write(s, binary.LittleEndian, ch.samples)
	binary.Write(s, binary.LittleEndian, ch.window)
	binary.Write(s, binary.LittleEndian, ch.bank)
	binary.Write(s, binary.LittleEndian, ch.usedBank)
	binary.Write(s, binary.LittleEndian, ch.mode)
}

func (ch *wave) deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &ch.enabled)
	binary.Read(s, binary.LittleEndian, &ch.dacEnable)
	binary.Read(s, binary.LittleEndian, &ch.stop)
	binary.Read(s, binary.LittleEndian, &ch.length)
	binary.Read(s, binary.LittleEndian, &ch.volume)
	binary.Read(s, binary.LittleEndian, &ch.period)
	binary.Read(s, binary.LittleEndian, &ch.freqCounter)
	binary.Read(s, binary.LittleEndian, &ch.samples)
	binary.Read(s, binary.LittleEndian, &ch.window)
	binary.Read(s, binary.LittleEndian, &ch.bank)
	binary.Read(s, binary.LittleEndian, &ch.usedBank)
	binary.Read(s, binary.LittleEndian, &ch.mode)
}
