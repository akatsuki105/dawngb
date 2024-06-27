package apu

import (
	"encoding/binary"
	"io"
)

// 音の三要素 のうち、 音の大きさ (振幅)
type envelope struct {
	initialVolume int32 // 初期音量(リスタート時にセット)
	volume        int32
	direction     bool // 音量変更の方向(trueで大きくなっていく)

	// 音量変更の速さ(0に近いほど速い、ただし0だと変化なし)
	// speed が n のとき、 音量変更は (n / 64) 秒 ごとに行われる
	speed int32
	step  int32 // 音量変更を行うタイミングをカウントするためのカウンタ
}

func newEnvelope() *envelope {
	return &envelope{
		step: 8,
	}
}

func (e *envelope) reset() {
	e.volume = e.initialVolume
	e.step = e.speed
	if e.speed == 0 {
		e.step = 8
	}
}

func (e *envelope) update() {
	if e.speed != 0 {
		e.step--
		if e.step == 0 {
			e.updateVolume()
			e.step = e.speed
			if e.speed == 0 {
				e.step = 8
			}
		}
	}
}

func (e *envelope) updateVolume() {
	if e.direction {
		e.volume = min(e.volume+1, 15)
	} else {
		e.volume = max(e.volume-1, 0)
	}
}

func (e *envelope) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, e.initialVolume)
	binary.Write(s, binary.LittleEndian, e.volume)
	binary.Write(s, binary.LittleEndian, e.direction)
	binary.Write(s, binary.LittleEndian, e.speed)
	binary.Write(s, binary.LittleEndian, e.step)
}

func (e *envelope) deserialize(d io.Reader) {
	binary.Read(d, binary.LittleEndian, &e.initialVolume)
	binary.Read(d, binary.LittleEndian, &e.volume)
	binary.Read(d, binary.LittleEndian, &e.direction)
	binary.Read(d, binary.LittleEndian, &e.speed)
	binary.Read(d, binary.LittleEndian, &e.step)
}
