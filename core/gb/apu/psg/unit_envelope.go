package psg

// 音の三要素 のうち、 音の大きさ (振幅)
type Envelope struct {
	initialVolume uint8 // 初期音量(リスタート時にセット)
	volume        uint8 // 0..15
	direction     bool  // 音量変更の方向(trueで大きくなっていく)

	// 音量変更の速さ(0に近いほど速い、ただし0だと変化なし)
	// speed が n のとき、 音量変更は (n / 64) 秒 ごとに行われる
	speed uint8
	step  uint8 // 音量変更を行うタイミングをカウントするためのカウンタ
}

func newEnvelope() *Envelope {
	return &Envelope{
		step: 8,
	}
}

func (e *Envelope) reset() {
	e.initialVolume, e.volume = 0, 0
	e.direction = false
	e.speed, e.step = 0, 8
}

func (e *Envelope) TurnOff() {
	e.speed = 0
	e.direction = false
	e.initialVolume = 0
}

func (e *Envelope) reload() {
	e.volume = e.initialVolume
	e.step = e.speed
	if e.speed == 0 {
		e.step = 8
	}
}

func (e *Envelope) update() {
	if e.speed != 0 {
		e.step--
		if e.step == 0 {
			e.updateVolume()
			e.step = e.speed
		}
	}
}

func (e *Envelope) updateVolume() {
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

type EnvelopeSnapshot struct {
	Header                uint64
	InitialVolume, Volume uint8
	Direction             bool
	Speed, Step           uint8
	Reserved              [7]uint8
}

func (e *Envelope) CreateSnapshot() EnvelopeSnapshot {
	s := EnvelopeSnapshot{
		InitialVolume: e.initialVolume,
		Volume:        e.volume,
		Direction:     e.direction,
		Speed:         e.speed,
		Step:          e.step,
	}
	return s
}

func (e *Envelope) RestoreSnapshot(snap EnvelopeSnapshot) bool {
	e.initialVolume, e.volume = snap.InitialVolume, snap.Volume
	e.direction = snap.Direction
	e.speed = snap.Speed
	e.step = snap.Step
	return true
}
