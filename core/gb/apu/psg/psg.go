package psg

// .model
const (
	MODEL_GB = iota
	MODEL_GBA
)

// GB/GBA の PSG (Programmable Sound Generator) ユニット, 4.19MHz で動作する (GB/GBA両方)
type PSG struct {
	model uint8

	Enabled bool // NR52.7

	CH1, CH2 *Square
	CH3      *Wave
	CH4      *Noise

	sequencerCounter int16 // (フレームシーケンサの)512Hzを生み出すためのカウンタ (ref: https://gbdev.io/pandocs/Audio_details.html#div-apu)
	sequencerStep    uint8 // 512Hzから 64, 128, 256Hzなどの生み出すためのカウンタ

	ioreg                     [0x30]uint8 // APUが勝手に状態を変えないレジスタ　はread時にここの値を返す
	leftVolume, rightVolume   uint8       // NR50 (n: 0..7)
	leftEnables, rightEnables [4]bool     // NR51, ch1, ch2, ch3, ch4 の左右出力を有効にするかどうか
}

func New(model uint8) *PSG {
	return &PSG{
		model: model,
		CH1:   newSquareChannel(true),
		CH2:   newSquareChannel(false),
		CH3:   newWaveChannel(model),
		CH4:   newNoiseChannel(),
	}
}

func (a *PSG) Reset() {
	a.Enabled = false
	a.CH1.Reset()
	a.CH2.Reset()
	a.CH3.Reset()
	a.CH4.Reset()
	a.sequencerCounter, a.sequencerStep = 0, 0
	clear(a.ioreg[:])
	a.leftVolume, a.rightVolume = 7, 7
	a.leftEnables, a.rightEnables = [4]bool{}, [4]bool{}
}

func (a *PSG) SkipBIOS() {
	a.Write(NR10, 0x80)
	a.Write(NR11, 0xBF)
	a.Write(NR12, 0xF3)
	a.Write(NR13, 0xFF)
	a.Write(NR14, 0xBF)
	a.Write(NR21, 0x3F)
	a.Write(NR22, 0x00)
	a.Write(NR23, 0xFF)
	a.Write(NR24, 0xBF)
	a.Write(NR30, 0x7F)
	a.Write(NR31, 0xFF)
	a.Write(NR32, 0x9F)
	a.Write(NR33, 0xFF)
	a.Write(NR34, 0xBF)
	a.Write(NR41, 0xFF)
	a.Write(NR42, 0x00)
	a.Write(NR43, 0x00)
	a.Write(NR44, 0xBF)
	a.Write(NR50, 0x77)
	a.Write(NR51, 0xF3)
	a.Write(NR52, 0xF1)
}

// 4MHz で1サイクル進める
func (a *PSG) Step() {
	if a.Enabled {
		if a.sequencerCounter > 0 {
			a.sequencerCounter--
		} else {
			is64Hz := a.sequencerStep == 7                                                                          // Envelope sweep
			is128Hz := a.sequencerStep == 2 || a.sequencerStep == 6                                                 // CH1 freq sweep
			is256Hz := a.sequencerStep == 0 || a.sequencerStep == 2 || a.sequencerStep == 4 || a.sequencerStep == 6 // Sound length

			if is256Hz {
				a.CH1.clock256Hz()
				a.CH2.clock256Hz()
				a.CH3.clock256Hz()
				a.CH4.clock256Hz()
			}
			if is128Hz {
				a.CH1.clock128Hz()
			}
			if is64Hz {
				a.CH1.clock64Hz()
				a.CH2.clock64Hz()
				a.CH4.clock64Hz()
			}
			a.sequencerStep = (a.sequencerStep + 1) & 7
			a.sequencerCounter = 8192 // 512Hz = 4194304/8192
		}

		a.CH1.clockTimer()
		a.CH2.clockTimer()
		a.CH3.clockTimer()
		a.CH4.clockTimer()
	}
}

// 0..63 の値を返す
func (a *PSG) Sample(mask uint8) (lsample, rsample uint8) {
	left, right := uint8(0), uint8(0)

	if a.Enabled {
		ch1, ch2, ch3, ch4 := a.CH1.GetOutput(), a.CH2.GetOutput(), a.CH3.GetOutput(), a.CH4.GetOutput()
		mask1, mask2, mask3, mask4 := (mask&(1<<0)) != 0, (mask&(1<<1)) != 0, (mask&(1<<2)) != 0, (mask&(1<<3)) != 0
		if mask1 && a.leftEnables[0] {
			left += ch1
		}
		if mask2 && a.leftEnables[1] {
			left += ch2
		}
		if mask3 && a.leftEnables[2] {
			left += ch3
		}
		if mask4 && a.leftEnables[3] {
			left += ch4
		}
		if mask1 && a.rightEnables[0] {
			right += ch1
		}
		if mask2 && a.rightEnables[1] {
			right += ch2
		}
		if mask3 && a.rightEnables[2] {
			right += ch3
		}
		if mask4 && a.rightEnables[3] {
			right += ch4
		}
	}
	return left, right
}

// Volume returns the volume of the NR50 (n: 0..7)
func (a *PSG) Volume() (left, right uint8) {
	return a.leftVolume, a.rightVolume
}
