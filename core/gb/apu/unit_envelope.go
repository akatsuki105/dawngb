package apu

// 音の三要素 のうち、 音の大きさ (振幅)
type envelope struct {
	initialVolume int // 初期音量(リスタート時にセット)
	volume        int
	direction     bool // 音量変更の方向(trueで大きくなっていく)

	// 音量変更の速さ(0に近いほど速い、ただし0だと変化なし)
	// speed が n のとき、 音量変更は (n / 64) 秒 ごとに行われる
	speed int
	step  int // 音量変更を行うタイミングをカウントするためのカウンタ
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
