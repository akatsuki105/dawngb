package cartridge

import (
	"bytes"
	"encoding/binary"
)

var tmp = bytes.NewBuffer(make([]uint8, 0, 512)) // Needed ?

type Snapshot struct {
	Header   uint64 // バージョン番号とかなんか持たせたいとき用に確保
	Buffer   [512]uint8
	Reserved [32]uint8
}

func (c *Cartridge) CreateSnapshot() Snapshot {
	s := Snapshot{}

	switch mapper := c.MBC.(type) {
	case *MBC1:
		snap := mapper.CreateSnapshot()
		binary.Write(tmp, binary.LittleEndian, snap)
		copy(s.Buffer[:], tmp.Bytes())
		tmp.Reset()
	case *MBC2:
		snap := mapper.CreateSnapshot()
		binary.Write(tmp, binary.LittleEndian, snap)
		copy(s.Buffer[:], tmp.Bytes())
		tmp.Reset()
	}
	return s
}

func (c *Cartridge) RestoreSnapshot(snap Snapshot) bool {
	switch mapper := c.MBC.(type) {
	case *MBC1:
		tmp.Write(snap.Buffer[:])
		var s MBC1Snapshot
		binary.Read(tmp, binary.LittleEndian, &s)
		mapper.RestoreSnapshot(s)
		tmp.Reset()
	case *MBC2:
		tmp.Write(snap.Buffer[:])
		var s MBC2Snapshot
		binary.Read(tmp, binary.LittleEndian, &s)
		mapper.RestoreSnapshot(s)
		tmp.Reset()
	}
	return true
}
