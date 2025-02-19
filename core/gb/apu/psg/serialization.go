package psg

type Snapshot struct {
	Header                    uint64
	Model                     uint8
	Enabled                   bool
	CH1, CH2                  SquareSnapshot
	CH3                       WaveSnapshot
	CH4                       NoiseSnapshot
	SequencerCounter          int16
	SequencerStep             uint8
	IOReg                     [0x30]uint8
	LeftVolume, RightVolume   uint8
	LeftEnables, RightEnables [4]bool
	Reserved                  [16]uint8
}

func (p *PSG) CreateSnapshot() Snapshot {
	s := Snapshot{
		Model:            p.model,
		Enabled:          p.enabled,
		CH1:              p.CH1.CreateSnapshot(),
		CH2:              p.CH2.CreateSnapshot(),
		CH3:              p.CH3.CreateSnapshot(),
		CH4:              p.CH4.CreateSnapshot(),
		SequencerCounter: p.sequencerCounter,
		SequencerStep:    p.sequencerStep,
		IOReg:            p.ioreg,
		LeftVolume:       p.leftVolume,
		RightVolume:      p.rightVolume,
		LeftEnables:      p.leftEnables,
		RightEnables:     p.rightEnables,
	}
	return s
}

func (p *PSG) RestoreSnapshot(snap Snapshot) bool {
	p.model = snap.Model
	p.enabled = snap.Enabled
	p.CH1.RestoreSnapshot(snap.CH1)
	p.CH2.RestoreSnapshot(snap.CH2)
	p.CH3.RestoreSnapshot(snap.CH3)
	p.CH4.RestoreSnapshot(snap.CH4)
	p.sequencerCounter = snap.SequencerCounter
	p.sequencerStep = snap.SequencerStep
	copy(p.ioreg[:], snap.IOReg[:])
	p.leftVolume, p.rightVolume = snap.LeftVolume, snap.RightVolume
	p.leftEnables = snap.LeftEnables
	p.rightEnables = snap.RightEnables
	return true
}
