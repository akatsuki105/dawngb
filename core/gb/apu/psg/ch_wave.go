package psg

import (
	"encoding/binary"
	"io"
)

const waveBank = 2

// .volume が index
var volumeShift = [4]uint8{4, 0, 1, 2} // 波形は最大15なので4左シフトすれば0%

type wave struct {
	model   uint8
	enabled bool // NR52.2

	dacEnable bool   // NR30.7
	length    uint16 // NR31; 音の残り再生時間
	volume    uint8  // NR32.5-6; 0: 0%, 1: 100%, 2: 50%, 3: 25%
	stop      bool   // NR34.6; .length が 0 になったときに 音を止めるかどうか

	period      uint16 // NR33.0-7, NR34.0-2; GBでは周波数を指定するのではなく、周期の長さを指定する
	freqCounter uint16

	samples [16 * waveBank]uint8 // 4bitサンプル*32 で16バイト ; GBAの場合はバンクが2つある
	sample  uint8                // 0..15
	window  uint8                // 0 ~ 31

	output uint8 // 0..15

	// For GBA
	mode    uint8 // NR30.5; 0: 16バイト(32サンプル)を演奏に使い、裏のバンクでは読み書きを行う、 1: 32バイト(64サンプル)を全部演奏に使う
	bank    uint8 // NR30.6
	curBank uint8 // 現在演奏中のバンク、modeが1の場合は、 .bank の値と必ずしも一致しないので
}

func newWaveChannel(model uint8) *wave {
	return &wave{
		model: model,
	}
}

func (ch *wave) reset() {
	ch.enabled = false
	ch.dacEnable = false
	ch.volume, ch.stop, ch.length = 0, false, 0
	ch.period, ch.freqCounter = 0, 0
	clear(ch.samples[:])
	ch.window = 0
	ch.mode, ch.bank, ch.curBank = 0, 0, 0
	ch.output = 0
}

func (ch *wave) reload() {
	ch.enabled = ch.dacEnable
	ch.freqCounter = ch.windowStepCycle() + 2
	ch.window = 0
	ch.output = 0
	if ch.length == 0 {
		ch.length = 256
	}
}

func (ch *wave) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length == 0 {
			ch.enabled = false
		}
	}
}

func (ch *wave) clockTimer() {
	if ch.freqCounter > 0 {
		ch.freqCounter--
		if ch.freqCounter == 0 {
			ch.freqCounter = ch.windowStepCycle()
			ch.update()
			if ch.window == 0 {
				ch.curBank ^= ch.mode
			}
		}
	}
}

func (ch *wave) getOutput() uint8 {
	if ch.enabled {
		shift := volumeShift[ch.volume]
		return ch.output >> shift
	}
	return 0
}

func (ch *wave) update() {
	ch.window = (ch.window + 1) & 0x1F // 読み出す前にインクリメント(CH3のreload後に最初に読み出すのはsamples[0]の下位ニブル)

	upper := (ch.window & 0x1) == 0
	if upper {
		ch.output = ch.samples[ch.window>>1] >> 4
	} else {
		ch.output = ch.samples[ch.window>>1] & 0xF
	}
}

func (ch *wave) read(addr uint16) uint8 {
	if !ch.enabled {
		bank := uint16(0)
		if ch.model == MODEL_GBA {
			if ch.bank == 0 {
				bank = 16
			}
		}
		return ch.samples[bank|(addr&0xF)]
	}
	return 0xFF // AGB
}

func (ch *wave) write(addr uint16, val uint8) {
	if !ch.enabled {
		bank := uint16(0)
		if ch.model == MODEL_GBA {
			if ch.bank == 0 {
				bank = 16
			}
		}
		ch.samples[bank|(addr&0xF)] = val
	}
}

func (ch *wave) windowStepCycle() uint16 {
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
	binary.Write(s, binary.LittleEndian, ch.curBank)
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
	binary.Read(s, binary.LittleEndian, &ch.curBank)
	binary.Read(s, binary.LittleEndian, &ch.mode)
}
