package apu

import (
	"encoding/binary"
	"io"
)

const (
	MODEL_GB = iota
	MODEL_GBA
)

// GB/GBA の PSG (Programmable Sound Generator) ユニット
// 4.19MHz で動作する (GB/GBA両方)
type APU interface {
	Reset(hasBIOS bool)
	Step() // step 1 APU cycle (4.19MHz)

	Read(addr uint16) uint8
	Write(addr uint16, val uint8)

	Sample() (lsample, rsample uint8)

	Serialize(s io.Writer)
	Deserialize(s io.Reader)
}

type apu struct {
	enabled bool
	model   uint8

	ch1, ch2 *square
	ch3      *wave
	ch4      *noise

	sequencerCounter int64 // (フレームシーケンサの)512Hzを生み出すためのカウンタ (ref: https://gbdev.io/pandocs/Audio_details.html#div-apu)
	sequencerStep    int64 // 512Hzから 64, 128, 256Hzなどの生み出すためのカウンタ

	ioreg  [0x30]uint8
	volume [2]uint8 // NR50(Left, Right)
}

func New(model int) APU {
	return &apu{
		model: uint8(model),
	}
}

func (a *apu) Reset(hasBIOS bool) {
	a.enabled = false
	a.ch1 = newSquareChannel(true)
	a.ch2 = newSquareChannel(false)
	a.ch3 = newWaveChannel()
	a.ch4 = newNoiseChannel()
	a.sequencerCounter = 0
	a.sequencerStep = 0
	a.volume = [2]uint8{7, 7}
	a.ioreg = [0x30]uint8{}
	if !hasBIOS {
		a.skipBIOS()
	}
}

func (a *apu) skipBIOS() {
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
func (a *apu) Step() {
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

			a.sequencerStep = (a.sequencerStep + 1) % 8
			a.sequencerCounter = 8192 // 512Hz = 4194304/8192
		}

		a.ch1.clockTimer()
		a.ch2.clockTimer()
		a.ch3.clockTimer()
		a.ch4.clockTimer()
	}
}

func (a *apu) Sample() (lsample, rsample uint8) {
	sample := (a.ch1.getOutput() + a.ch2.getOutput() + a.ch3.getOutput() + a.ch4.getOutput()) // 各チャンネルの出力(音量=波)を足し合わせたものがサンプル
	left := uint8((sample * int(a.volume[0])) / 7)
	right := uint8((sample * int(a.volume[1])) / 7)
	return left, right
}

func (a *apu) Serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, a.enabled)
	binary.Write(s, binary.LittleEndian, a.model)
	a.ch1.serialize(s)
	a.ch2.serialize(s)
	a.ch3.serialize(s)
	a.ch4.serialize(s)
	binary.Write(s, binary.LittleEndian, a.sequencerCounter)
	binary.Write(s, binary.LittleEndian, a.sequencerStep)
	binary.Write(s, binary.LittleEndian, a.ioreg)
	binary.Write(s, binary.LittleEndian, a.volume)
}

func (a *apu) Deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &a.enabled)
	binary.Read(s, binary.LittleEndian, &a.model)
	a.ch1.deserialize(s)
	a.ch2.deserialize(s)
	a.ch3.deserialize(s)
	a.ch4.deserialize(s)
	binary.Read(s, binary.LittleEndian, &a.sequencerCounter)
	binary.Read(s, binary.LittleEndian, &a.sequencerStep)
	binary.Read(s, binary.LittleEndian, &a.ioreg)
	binary.Read(s, binary.LittleEndian, &a.volume)
}
