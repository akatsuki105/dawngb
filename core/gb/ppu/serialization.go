package ppu

import "errors"

var errSnapshotNil = errors.New("PPU snapshot is nil")

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

func (p *PPU) UpdateSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	snap.Cycles, snap.Frames = p.cycles, p.Frame
	snap.Lx, snap.Ly = int16(p.Lx), int16(p.Ly)

	copy(snap.VRAM.Data[:], p.RAM.Data[:])
	snap.VRAM.Bank = p.RAM.Bank

	snap.DMA.Active, snap.DMA.Src, snap.DMA.Until = p.DMA.Active, p.DMA.Src, p.DMA.Until
	snap.LCDC, snap.STAT, snap.LYC = p.LCDC, p.STAT, p.LYC
	copy(snap.OAM[:], p.OAM[:])
	copy(snap.Palette[:], p.Palette[:])
	copy(snap.IOReg[:], p.ioreg[:])
	snap.EnableLatch = p.enableLatch
	snap.ObjCount = p.objCount
	snap.BGPI, snap.OBPI = p.BGPI, p.OBPI
	return nil
}

func (p *PPU) RestoreSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}

	p.cycles = snap.Cycles
	p.Frame = snap.Frames
	p.Lx, p.Ly = int(snap.Lx), int(snap.Ly)
	copy(p.RAM.Data[:], snap.VRAM.Data[:])
	p.RAM.Bank = snap.VRAM.Bank
	p.DMA.Active, p.DMA.Src, p.DMA.Until = snap.DMA.Active, snap.DMA.Src, snap.DMA.Until
	p.LCDC, p.STAT, p.LYC = snap.LCDC, snap.STAT, snap.LYC
	p.r.SetLCDC(p.LCDC)
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
	p.BGPI, p.OBPI = snap.BGPI, snap.OBPI
	return nil
}
