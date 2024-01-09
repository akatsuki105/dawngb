package audio

import "io"

type Audio interface {
	Add(cycles int64)
	CatchUp()

	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type audio struct {
	enabled      bool
	ch1, ch2     *square
	ch3          *wave
	sampleBuffer io.Writer
	cycles       int64 // 遅れているサイクル数(8.3MHzのマスターサイクル単位)

	sequencerCounter int64 // (フレームシーケンサの)512Hzを生み出すためのカウンタ (ref: https://gbdev.io/pandocs/Audio_details.html#div-apu)
	sequencerStep    int64 // 512Hzから 64, 128, 256Hzなどの生み出すためのカウンタ

	sampleTimer int64 // 1サンプルを生み出すために44100Hzを生み出すためのカウンタ

	ioreg [0x30]uint8
}

func New(sampleBuffer io.Writer) Audio {
	return &audio{
		ch1:          newSquareChannel(true),
		ch2:          newSquareChannel(false),
		ch3:          newWaveChannel(),
		sampleBuffer: sampleBuffer,
	}
}

func (a *audio) Add(cycles int64) {
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
				}
				if is128Hz {
					a.ch1.clock128Hz()
				}
				if is64Hz {
					a.ch1.clock64Hz()
					a.ch2.clock64Hz()
				}

				a.sequencerStep = (a.sequencerStep + 1) % 8
				a.sequencerCounter = 8192 // 512Hz = 4194304/8192
			}

			a.ch1.clockTimer()
			a.ch2.clockTimer()
			a.ch3.clockTimer()

			// サンプルを生成
			sample := uint8(a.ch1.getOutput() + a.ch2.getOutput() + a.ch3.getOutput()) // 各チャンネルの出力(音量=波)を足し合わせたものがサンプル
			if a.sampleTimer <= 0 {
				if a.sampleBuffer != nil {
					a.sampleBuffer.Write([]byte{0, sample, 0, sample})
				}

				a.sampleTimer = 95 // 44100Hzにダウンサンプリングしたい = 44100Hzごとにサンプルを生成したい = 95APUサイクルごとにサンプルを生成したい(4194304/44100 = 95)
			}
			a.sampleTimer--
		}
	}

	a.cycles -= apuCycles * 2
}
