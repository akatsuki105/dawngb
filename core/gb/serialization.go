package gb

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/akatsuki105/dawngb/core/gb/apu"
	"github.com/akatsuki105/dawngb/core/gb/cartridge"
	"github.com/akatsuki105/dawngb/core/gb/cpu"
	"github.com/akatsuki105/dawngb/core/gb/ppu"
)

type Header struct {
	Magic    [4]uint8  // "DAWN"
	Version  uint64    // 今は0
	Reserved [64]uint8 // 拡張用
}

type Snapshot struct {
	Header
	Model    uint8
	CPU      cpu.Snapshot
	PPU      ppu.Snapshot
	APU      apu.Snapshot
	Cart     cartridge.Snapshot
	WRAM     [4 * KB * 8]uint8
	WRAMBank uint8
	Reserved [2 * KB]uint8 // 拡張用
}

func NewSnapshot(version uint64) *Snapshot {
	return &Snapshot{
		Header: Header{
			Magic:   [4]uint8{'D', 'A', 'W', 'N'},
			Version: version,
		},
	}
}

var errSnapshotNil = errors.New("GB snapshot is nil")

func (g *GB) UpdateSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	if snap.Magic != [4]uint8{'D', 'A', 'W', 'N'} {
		return errors.New("invalid snapshot ID")
	}
	snap.Model = uint8(g.Model)
	snap.CPU = g.CPU.CreateSnapshot()
	if err := g.PPU.UpdateSnapshot(&snap.PPU); err != nil {
		return err
	}
	snap.APU = g.APU.CreateSnapshot()
	if err := g.Cart.UpdateSnapshot(&snap.Cart); err != nil {
		return err
	}
	copy(snap.WRAM[:], g.WRAM.Data[:])
	snap.WRAMBank = g.WRAM.Bank
	return nil
}

func (g *GB) RestoreSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	if snap.Magic != [4]uint8{'D', 'A', 'W', 'N'} {
		return errors.New("invalid snapshot ID")
	}

	g.Model = Model(snap.Model)
	if err := g.CPU.RestoreSnapshot(&snap.CPU); err != nil {
		return err
	}
	g.PPU.RestoreSnapshot(&snap.PPU)
	g.APU.RestoreSnapshot(snap.APU)
	g.Cart.RestoreSnapshot(&snap.Cart)
	copy(g.WRAM.Data[:], snap.WRAM[:])
	g.WRAM.Bank = snap.WRAMBank
	return nil
}

func (g *GB) Serialize(w io.Writer) bool {
	if w != nil {
		err := g.UpdateSnapshot(&g.Snap)
		if err != nil {
			return false
		}
		binary.Write(w, binary.LittleEndian, g.Snap)
		return true
	}
	return false
}

func (g *GB) Deserialize(r io.Reader) bool {
	if r != nil {
		binary.Read(r, binary.LittleEndian, &g.Snap)
		err := g.RestoreSnapshot(&g.Snap)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
