package psg

import (
	"encoding/binary"
	"io"
)

// .model
const (
	MODEL_GB = iota
	MODEL_GBA
)

// GB/GBA の PSG (Programmable Sound Generator) ユニット, 4.19MHz で動作する (GB/GBA両方)
type PSG struct {
	enabled bool
	model   uint8

	ch1, ch2 *square
	ch3      *wave
	ch4      *noise

	sequencerCounter int64 // (フレームシーケンサの)512Hzを生み出すためのカウンタ (ref: https://gbdev.io/pandocs/Audio_details.html#div-apu)
	sequencerStep    int64 // 512Hzから 64, 128, 256Hzなどの生み出すためのカウンタ

	ioreg                     [0x30]uint8
	leftVolume, rightVolume   uint8   // NR50 (n: 0..7)
	leftEnables, rightEnables [4]bool // NR51, ch1, ch2, ch3, ch4 の左右出力を有効にするかどうか
}

func New(model uint8) *PSG {
	return &PSG{
		model: model,
	}
}

func (a *PSG) Reset(hasBIOS bool) {
	a.enabled = false
	a.ch1 = newSquareChannel(true)
	a.ch2 = newSquareChannel(false)
	a.ch3 = newWaveChannel()
	a.ch4 = newNoiseChannel()
	a.sequencerCounter = 0
	a.sequencerStep = 0
	a.leftVolume, a.rightVolume = 7, 7
	a.leftEnables, a.rightEnables = [4]bool{}, [4]bool{}
	a.ioreg = [0x30]uint8{}
	if !hasBIOS {
		a.skipBIOS()
	}
}

func (a *PSG) skipBIOS() {
	a.Write(0xFF10, 0x80)
	a.Write(0xFF11, 0xBF)
	a.Write(0xFF12, 0xF3)
	a.Write(0xFF13, 0xFF)
	a.Write(0xFF14, 0xBF)
	a.Write(0xFF16, 0x3F)
	a.Write(0xFF17, 0x00)
	a.Write(0xFF18, 0xFF)
	a.Write(0xFF19, 0xBF)
	a.Write(0xFF1A, 0x7F)
	a.Write(0xFF1B, 0xFF)
	a.Write(0xFF1C, 0x9F)
	a.Write(0xFF1D, 0xFF)
	a.Write(0xFF1E, 0xBF)
	a.Write(0xFF20, 0xFF)
	a.Write(0xFF21, 0x00)
	a.Write(0xFF22, 0x00)
	a.Write(0xFF23, 0xBF)
	a.Write(0xFF24, 0x77)
	a.Write(0xFF25, 0xF3)
	a.Write(0xFF26, 0xF1)
}

// 4.19 MHz で1サイクル進める
func (a *PSG) Step() {
	if a.enabled {
		if a.sequencerCounter > 0 {
			a.sequencerCounter--
		} else {
			is64Hz := a.sequencerStep == 7                                                                          // Envelope sweep
			is128Hz := a.sequencerStep == 2 || a.sequencerStep == 6                                                 // CH1 freq sweep
			is256Hz := a.sequencerStep == 0 || a.sequencerStep == 2 || a.sequencerStep == 4 || a.sequencerStep == 6 // Sound length

			if is256Hz {
				a.ch1.clock256Hz()
				a.ch2.clock256Hz()
				a.ch3.clock256Hz()
				a.ch4.clock256Hz()
			}
			if is128Hz {
				a.ch1.clock128Hz()
			}
			if is64Hz {
				a.ch1.clock64Hz()
				a.ch2.clock64Hz()
				a.ch4.clock64Hz()
			}
			a.sequencerStep = (a.sequencerStep + 1) & 7
			a.sequencerCounter = 8192 // 512Hz = 4194304/8192
		}

		a.ch1.clockTimer()
		a.ch2.clockTimer()
		a.ch3.clockTimer()
		a.ch4.clockTimer()
	}
}

// 0..63 の値を返す
func (a *PSG) Sample() (lsample, rsample uint8) {
	ch1 := a.ch1.getOutput()
	ch2 := a.ch2.getOutput()
	ch3 := a.ch3.getOutput()
	ch4 := a.ch4.getOutput()
	left, right := uint8(0), uint8(0)
	if a.leftEnables[0] {
		left += ch1
	}
	if a.leftEnables[1] {
		left += ch2
	}
	if a.leftEnables[2] {
		left += ch3
	}
	if a.leftEnables[3] {
		left += ch4
	}
	if a.rightEnables[0] {
		right += ch1
	}
	if a.rightEnables[1] {
		right += ch2
	}
	if a.rightEnables[2] {
		right += ch3
	}
	if a.rightEnables[3] {
		right += ch4
	}
	return left, right
}

// Volume returns the volume of the NR50 (n: 0..7)
func (a *PSG) Volume() (left, right uint8) {
	return a.leftVolume, a.rightVolume
}

func (a *PSG) Serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, a.enabled)
	binary.Write(s, binary.LittleEndian, a.model)
	a.ch1.serialize(s)
	a.ch2.serialize(s)
	a.ch3.serialize(s)
	a.ch4.serialize(s)
	binary.Write(s, binary.LittleEndian, a.sequencerCounter)
	binary.Write(s, binary.LittleEndian, a.sequencerStep)
	binary.Write(s, binary.LittleEndian, a.ioreg)
}

func (a *PSG) Deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &a.enabled)
	binary.Read(s, binary.LittleEndian, &a.model)
	a.ch1.deserialize(s)
	a.ch2.deserialize(s)
	a.ch3.deserialize(s)
	a.ch4.deserialize(s)
	binary.Read(s, binary.LittleEndian, &a.sequencerCounter)
	binary.Read(s, binary.LittleEndian, &a.sequencerStep)
	binary.Read(s, binary.LittleEndian, &a.ioreg)
}
