package psg

import (
	"encoding/binary"
	"io"
)

// 音の三要素 のうち、 音の大きさ (振幅)
type envelope struct {
	initialVolume uint8 // 初期音量(リスタート時にセット)
	volume        uint8 // 0..15
	direction     bool  // 音量変更の方向(trueで大きくなっていく)

	// 音量変更の速さ(0に近いほど速い、ただし0だと変化なし)
	// speed が n のとき、 音量変更は (n / 64) 秒 ごとに行われる
	speed uint8
	step  uint8 // 音量変更を行うタイミングをカウントするためのカウンタ
}

func newEnvelope() *envelope {
	return &envelope{
		step: 8,
	}
}

func (e *envelope) reset() {
	e.initialVolume, e.volume = 0, 0
	e.direction = false
	e.speed, e.step = 0, 8
}

func (e *envelope) reload() {
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
		}
	}
}

func (e *envelope) updateVolume() {
	if e.direction {
		if e.volume < 15 {
			e.volume++
		}
	} else {
		if e.volume > 0 {
			e.volume--
		}
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
