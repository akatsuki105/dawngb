package audio

import "io"

type Audio interface {
	Reset(hasBIOS bool)

	Tick(cycles int64)
	CatchUp()

	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type audio struct {
	enabled bool

	ch1, ch2 *square
	ch3      *wave
	ch4      *noise

	sampleBuffer io.Writer
	cycles       int64 // 遅れているサイクル数(8.3MHzのマスターサイクル単位)

	sequencerCounter int64 // (フレームシーケンサの)512Hzを生み出すためのカウンタ (ref: https://gbdev.io/pandocs/Audio_details.html#div-apu)
	sequencerStep    int64 // 512Hzから 64, 128, 256Hzなどの生み出すためのカウンタ

	sampleTimer int64 // 1サンプルを生み出すために44100Hzを生み出すためのカウンタ

	ioreg  [0x30]uint8
	volume [2]int // NR50(Left, Right)
}

func New(sampleBuffer io.Writer) Audio {
	return &audio{
		sampleBuffer: sampleBuffer,
	}
}

func (a *audio) Reset(hasBIOS bool) {
	a.enabled = false
	a.ch1 = newSquareChannel(true)
	a.ch2 = newSquareChannel(false)
	a.ch3 = newWaveChannel()
	a.ch4 = newNoiseChannel()
	a.cycles = 0
	a.sequencerCounter = 0
	a.sequencerStep = 0
	a.sampleTimer = 0
	a.volume = [2]int{7, 7}
	a.ioreg = [0x30]uint8{}
	if !hasBIOS {
		a.skipBIOS()
	}
}

func (a *audio) skipBIOS() {
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

func (a *audio) Tick(cycles int64) {
	a.cycles += cycles
}

func (a *audio) CatchUp() {
	apuCycles := a.cycles / 2 // APU　は 4.19MHz で動作する, マスターサイクルを 8.3MHz とすると 8.3MHz / 4.19MHz = 2

	for i := int64(0); i < apuCycles; i++ {
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

			// サンプルを生成
			if a.sampleTimer <= 0 {
				sample := (a.ch1.getOutput() + a.ch2.getOutput() + a.ch3.getOutput() + a.ch4.getOutput()) // 各チャンネルの出力(音量=波)を足し合わせたものがサンプル
				if a.sampleBuffer != nil {
					left := uint8((sample * a.volume[0]) / 7)
					right := uint8((sample * a.volume[1]) / 7)
					a.sampleBuffer.Write([]byte{0, left, 0, right})
				}

				a.sampleTimer = 95 // 44100Hzにダウンサンプリングしたい = 44100Hzごとにサンプルを生成したい = 95APUサイクルごとにサンプルを生成したい(4194304/44100 = 95)
			}
			a.sampleTimer--
		}
	}

	a.cycles -= apuCycles * 2
}
