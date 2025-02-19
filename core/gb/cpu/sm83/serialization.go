package sm83

type Snapshot struct {
	Header             uint64 // バージョン番号とかなんか持たせたいとき用に確保
	A, F               uint8
	BC, DE, HL, SP, PC uint16
	Inst               Context
	IME                bool
	Reserved           [8]uint8 // 拡張用
}

func (c *SM83) CreateSnapshot() Snapshot {
	return Snapshot{
		Header: 0,
		A:      c.R.A,
		F:      c.R.F.Pack(),
		BC:     c.R.BC.Pack(),
		DE:     c.R.DE.Pack(),
		HL:     c.R.HL.Pack(),
		SP:     c.R.SP,
		PC:     c.R.PC,
		Inst:   c.inst,
		IME:    c.IME,
	}
}

func (c *SM83) RestoreSnapshot(snap Snapshot) bool {
	c.R.A = snap.A
	c.R.F.Unpack(snap.F)
	c.R.BC.Unpack(snap.BC)
	c.R.DE.Unpack(snap.DE)
	c.R.HL.Unpack(snap.HL)
	c.R.SP = snap.SP
	c.R.PC = snap.PC
	c.inst = snap.Inst
	c.IME = snap.IME
	return true
}
