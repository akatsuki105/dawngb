package cpu

import "github.com/akatsuki105/dawngb/core/gb/cpu/sm83"

type Snapshot struct {
	sm83.Snapshot
	timerSnapshot timerSnapshot
}

func (c *CPU) CreateSnapshot() Snapshot {
	return Snapshot{
		Snapshot:      c.SM83.CreateSnapshot(),
		timerSnapshot: c.timer.CreateSnapshot(),
	}
}

func (c *CPU) RestoreSnapshot(snap Snapshot) bool {
	c.SM83.RestoreSnapshot(snap.Snapshot)
	c.timer.RestoreSnapshot(snap.timerSnapshot)
	return true
}
