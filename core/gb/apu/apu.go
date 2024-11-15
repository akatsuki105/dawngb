package apu

import (
	"encoding/binary"
	"io"

	"github.com/akatsuki105/dawngb/core/gb/apu/psg"
)

type APU struct {
	cycles int64 // 8MHzのマスターサイクル単位
	*psg.PSG
	sampleWriter io.Writer

	samples     [547 * 2]int16 // [[left, right]...], 547 = 32768 / 60
	sampleCount uint16
}

func New(audioBuffer io.Writer) *APU {
	if audioBuffer == nil {
		audioBuffer = io.Discard
	}

	return &APU{
		PSG:          psg.New(psg.MODEL_GB),
		sampleWriter: audioBuffer,
	}
}

func (a *APU) Reset(hasBIOS bool) {
	a.PSG.Reset(hasBIOS)
	a.cycles = 0
	clear(a.samples[:])
	a.sampleCount = 0
}

func (a *APU) Run(cycles8MHz int64) {
	for i := int64(0); i < cycles8MHz; i++ {
		a.cycles++
		if a.cycles%2 == 0 {
			a.PSG.Step()
		}
		if a.cycles%256 == 0 { // 32768Hzにダウンサンプリングしたい = 32768Hzごとにサンプルを生成したい = 256マスターサイクルごとにサンプルを生成する (8MHz / 32768Hz = 256)
			if int(a.sampleCount) < len(a.samples)/2 {
				left, right := a.PSG.Sample()
				lvolume, rvolume := a.PSG.Volume()
				lsample, rsample := (int(left)*512)-16384, (int(right)*512)-16384
				lsample, rsample = (lsample*int(lvolume+1))/8, (rsample*int(rvolume+1))/8
				a.samples[a.sampleCount*2] = int16(lsample)
				a.samples[a.sampleCount*2+1] = int16(rsample)
				a.sampleCount++
			}
		}
	}
}

func (a *APU) FlushSamples() {
	binary.Write(a.sampleWriter, binary.LittleEndian, a.samples[:a.sampleCount*2])
	a.sampleCount = 0
}
