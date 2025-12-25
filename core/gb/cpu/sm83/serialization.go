package sm83

import "errors"

type Snapshot struct {
	Header             uint64 // バージョン番号とかなんか持たせたいとき用に確保
	A, F               uint8
	BC, DE, HL, SP, PC uint16
	Inst               Context
	IME                bool
	Reserved           [8]uint8 // 拡張用
}

var errSnapshotNil = errors.New("SM83 snapshot is nil")

func (c *SM83) UpdateSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	snap.A, snap.F = c.R.A, c.R.F.Pack()
	snap.BC, snap.DE, snap.HL = c.R.BC.Pack(), c.R.DE.Pack(), c.R.HL.Pack()
	snap.SP, snap.PC = c.R.SP, c.R.PC
	snap.Inst = c.inst
	snap.IME = c.IME
	return nil
}

func (c *SM83) RestoreSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	c.R.A = snap.A
	c.R.F.Unpack(snap.F)
	c.R.BC.Unpack(snap.BC)
	c.R.DE.Unpack(snap.DE)
	c.R.HL.Unpack(snap.HL)
	c.R.SP = snap.SP
	c.R.PC = snap.PC
	c.inst = snap.Inst
	c.IME = snap.IME
	return nil
}
