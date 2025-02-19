package gb

import (
	"encoding/binary"
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

func (g *GB) CreateSnapshot() Snapshot {
	return Snapshot{
		Header: Header{
			Magic:   [4]uint8{'D', 'A', 'W', 'N'},
			Version: 0,
		},
		Model:    uint8(g.Model),
		CPU:      g.CPU.CreateSnapshot(),
		PPU:      g.PPU.CreateSnapshot(),
		APU:      g.APU.CreateSnapshot(),
		Cart:     g.Cart.CreateSnapshot(),
		WRAM:     g.wram,
		WRAMBank: g.wramBank,
	}
}

func (g *GB) RestoreSnapshot(snap Snapshot) bool {
	if snap.Magic != [4]uint8{'D', 'A', 'W', 'N'} {
		return false
	}

	g.Model = Model(snap.Model)
	g.CPU.RestoreSnapshot(snap.CPU)
	g.PPU.RestoreSnapshot(snap.PPU)
	g.APU.RestoreSnapshot(snap.APU)
	g.Cart.RestoreSnapshot(snap.Cart)
	copy(g.wram[:], snap.WRAM[:])
	g.wramBank = snap.WRAMBank
	return true
}

func (g *GB) Serialize(state io.Writer) bool {
	snap := g.CreateSnapshot()
	binary.Write(state, binary.LittleEndian, snap)
	return true
}

func (g *GB) Deserialize(state io.Reader) bool {
	var snap Snapshot
	binary.Read(state, binary.LittleEndian, &snap)
	ok := g.RestoreSnapshot(snap)
	if ok {
		g.inputs = 0
	}

	return ok
}
