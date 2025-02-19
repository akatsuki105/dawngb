package ppu

type Snapshot struct {
	Header          uint64
	Cycles          int64
	Frames          uint64
	Lx, Ly          int16
	VRAM            VRAM
	DMA             DMA
	LCDC, STAT, LYC uint8
	OAM             [160]uint8
	Palette         [(4 * 8) * 2]uint16
	IOReg           [0x30]uint8
	EnableLatch     bool
	ObjCount        uint8
	BGPI, OBPI      uint8
	Reserved        [64]uint8
}

func (p *PPU) CreateSnapshot() Snapshot {
	s := Snapshot{
		Cycles:      p.cycles,
		Frames:      p.frameCounter,
		Lx:          int16(p.lx),
		Ly:          int16(p.ly),
		VRAM:        p.RAM,
		DMA:         p.DMA,
		LCDC:        p.lcdc,
		STAT:        p.stat,
		LYC:         p.lyc,
		OAM:         p.OAM,
		Palette:     p.Palette,
		IOReg:       p.ioreg,
		EnableLatch: p.enableLatch,
		ObjCount:    p.objCount,
		BGPI:        p.bgpi,
		OBPI:        p.obpi,
	}
	return s
}

func (p *PPU) RestoreSnapshot(snap Snapshot) bool {
	p.cycles = snap.Cycles
	p.frameCounter = snap.Frames
	p.lx, p.ly = int(snap.Lx), int(snap.Ly)
	copy(p.RAM.Data[:], snap.VRAM.Data[:])
	p.RAM.Bank = snap.VRAM.Bank
	p.DMA.Active, p.DMA.Src, p.DMA.Until = snap.DMA.Active, snap.DMA.Src, snap.DMA.Until
	p.lcdc, p.stat, p.lyc = snap.LCDC, snap.STAT, snap.LYC
	p.r.SetLCDC(p.lcdc)
	copy(p.OAM[:], snap.OAM[:])
	copy(p.Palette[:], snap.Palette[:])
	copy(p.ioreg[:], snap.IOReg[:])
	p.r.SetSCX(p.ioreg[0x3])
	p.r.SetSCY(p.ioreg[0x2])
	p.r.SetWX(p.ioreg[0xB])
	p.r.SetWY(p.ioreg[0xA])
	p.r.SetBGP(p.ioreg[0x7])
	p.r.SetOBP(0, p.ioreg[0x8])
	p.r.SetOBP(1, p.ioreg[0x9])
	p.enableLatch = snap.EnableLatch
	p.objCount = snap.ObjCount
	p.bgpi, p.obpi = snap.BGPI, snap.OBPI
	return true
}
