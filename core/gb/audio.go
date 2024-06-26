package gb

import (
	"io"

	"github.com/akatsuki105/dawngb/core/gb/apu"
)

type audio struct {
	cycles       int64 // 8.3MHzのマスターサイクル単位
	apu          apu.APU
	sampleWriter io.Writer
	sampleTimer  int64
}

func newAudio(audioBuffer io.Writer) *audio {
	return &audio{
		apu:          apu.New(apu.MODEL_GB),
		sampleWriter: audioBuffer,
	}
}

func (a *audio) reset(hasBIOS bool) {
	a.apu.Reset(hasBIOS)
	a.cycles = 0
	a.sampleTimer = 0
}

// cycles は 8.3MHzのマスターサイクル単位
func (a *audio) tick(cycles int64) {
	if cycles > 0 {
		for i := int64(0); i < cycles; i++ {
			a.cycles++
			if a.cycles%2 == 0 {
				a.apu.Step()
			}

			a.sampleTimer--
			if a.sampleTimer <= 0 {
				a.sampleTimer += 256 // 32768Hzにダウンサンプリングしたい = 32768Hzごとにサンプルを生成したい = 256マスターサイクルごとにサンプルを生成する (8.3MHz / 32768Hz = 256)
				left, right := a.apu.Sample()
				a.sampleWriter.Write([]byte{0, left, 0, right})
			}
		}
	}
}
