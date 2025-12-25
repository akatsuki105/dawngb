package apu

import (
	"encoding/binary"
	"io"

	"github.com/akatsuki105/dawngb/core/gb/apu/psg"
)

// SoCに組み込まれているため、`/cpu`にある方が正確ではある
type APU struct {
	cycles int64 // 8MHzのマスターサイクル単位
	*psg.PSG
	sampleWriter io.Writer

	samples     [560 * 2]int16 // [[left, right]...], 549 = 32768 / 59.7275
	sampleCount uint16
	Mask        uint8
}

func New(audioBuffer io.Writer) *APU {
	if audioBuffer == nil {
		audioBuffer = io.Discard
	}

	return &APU{
		PSG:          psg.New(psg.MODEL_GB),
		sampleWriter: audioBuffer,
		Mask:         0b1111, // (CH4, CH3, CH2, CH1)
	}
}

func (a *APU) Reset() {
	a.PSG.Reset()
	a.cycles = 0
	clear(a.samples[:])
	a.sampleCount = 0
}

func (a *APU) Run(cycles8MHz int64) {
	for i := int64(0); i < cycles8MHz; i++ {
		a.cycles++
		if a.cycles&0b11 == 0 { // 2MHz
			a.PSG.Step()
		}
		if (a.cycles & 0xFF) == 0 { // 32768Hzにダウンサンプリングしたい = 32768Hzごとにサンプルを生成したい = 256マスターサイクルごとにサンプルを生成する (8MHz / 32768Hz = 256)
			if int(a.sampleCount) < len(a.samples)/2 {
				left, right := a.PSG.Sample(a.Mask)
				lvolume, rvolume := a.PSG.Volume()
				lsample, rsample := (int(left)*512)-16384, (int(right)*512)-16384
				lsample, rsample = (lsample*int(lvolume+1))/8, (rsample*int(rvolume+1))/8
				a.samples[a.sampleCount*2] = int16(lsample) / 2
				a.samples[a.sampleCount*2+1] = int16(rsample) / 2
				a.sampleCount++
			}
		}
	}
}

func (a *APU) FlushSamples() {
	binary.Write(a.sampleWriter, binary.LittleEndian, a.samples[:a.sampleCount*2])
	a.sampleCount = 0
}

type Snapshot struct {
	Header   uint64
	Cycles   int64
	PSG      psg.Snapshot
	Reserved [16]uint8
}

func (a *APU) CreateSnapshot() Snapshot {
	return Snapshot{
		Cycles: a.cycles,
		PSG:    a.PSG.CreateSnapshot(),
	}
}

func (a *APU) RestoreSnapshot(snap Snapshot) bool {
	a.cycles = snap.Cycles
	ok := a.PSG.RestoreSnapshot(snap.PSG)
	return ok
}
