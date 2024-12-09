package sm83

type Snapshot struct {
	A, F               uint8
	BC, DE, HL, SP, PC uint16
	IME                bool
}

func (c *SM83) CreateSnapshot() Snapshot {
	return Snapshot{
		A:   c.R.A,
		F:   c.R.F.Pack(),
		BC:  c.R.BC.Pack(),
		DE:  c.R.DE.Pack(),
		HL:  c.R.HL.Pack(),
		SP:  c.R.SP,
		PC:  c.R.PC,
		IME: c.IME,
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
	c.IME = snap.IME
	return true
}
