package cpu

import (
	"errors"

	"github.com/akatsuki105/dawngb/core/gb/cpu/sm83"
)

type Snapshot struct {
	Header           uint64 // バージョン番号とかなんか持たせたいとき用に確保
	IsCGB            bool
	Cycles           int64
	SM83             sm83.Snapshot
	Clock            int64
	Timer            TimerSnapshot
	DMA              DMASnapshot
	P14, P15         bool
	JoyP, Inputs     uint8
	Serial           SerialSnapshot
	FF50             bool
	HRAM             [0x7F]uint8
	Halted           bool
	IE, IF           uint8
	Key0, Key1       uint8
	FF72, FF73, FF74 uint8
}

var errSnapshotNil = errors.New("CPU snapshot is nil")

func (c *CPU) CreateSnapshot() Snapshot {
	s := Snapshot{
		IsCGB:  c.isCGB,
		Cycles: c.Cycles,
		Clock:  c.Clock,
		Timer:  c.Timer.CreateSnapshot(),
		DMA:    c.DMA.CreateSnapshot(),
		P14:    c.Joypad.P14,
		P15:    c.Joypad.P15,
		JoyP:   c.Joypad.JOYP,
		Inputs: c.Joypad.inputs,
		Serial: c.Serial.CreateSnapshot(),
		FF50:   c.BIOS.FF50,
		HRAM:   c.HRAM,
		Halted: c.Halted,
		IE:     c.IE,
		IF:     c.IF,
		Key0:   c.Key0,
		Key1:   c.Key1,
		FF72:   c.FF72,
		FF73:   c.FF73,
		FF74:   c.FF74,
	}
	c.SM83.UpdateSnapshot(&s.SM83)
	return s
}

func (c *CPU) RestoreSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	c.isCGB = snap.IsCGB
	c.Cycles = snap.Cycles
	if err := c.SM83.RestoreSnapshot(&snap.SM83); err != nil {
		return err
	}
	c.Clock = snap.Clock
	c.Timer.RestoreSnapshot(snap.Timer)
	c.DMA.RestoreSnapshot(snap.DMA)
	c.Joypad.P14, c.Joypad.P15, c.Joypad.JOYP, c.Joypad.inputs = snap.P14, snap.P15, snap.JoyP, snap.Inputs
	c.Serial.RestoreSnapshot(snap.Serial)
	c.BIOS.FF50 = snap.FF50
	copy(c.HRAM[:], snap.HRAM[:])
	c.Halted = snap.Halted
	c.IE, c.IF, c.Key0, c.Key1 = snap.IE, snap.IF, snap.Key0, snap.Key1
	c.FF72, c.FF73, c.FF74 = snap.FF72, snap.FF73, snap.FF74
	return nil
}
