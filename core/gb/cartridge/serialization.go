package cartridge

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var errSnapshotNil = errors.New("gb.Cartridge snapshot is nil")

var tmp = bytes.NewBuffer(make([]uint8, 0, 512)) // Needed ?

type Snapshot struct {
	Header   uint64 // バージョン番号とかなんか持たせたいとき用に確保
	Buffer   [512]uint8
	Reserved [32]uint8
}

func (c *Cartridge) CreateSnapshot() Snapshot {
	snap := Snapshot{}
	c.UpdateSnapshot(&snap)
	return snap
}

func (c *Cartridge) UpdateSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
	switch mapper := c.MBC.(type) {
	case *MBC1:
		mbc1 := mapper.CreateSnapshot()
		binary.Write(tmp, binary.LittleEndian, mbc1)
		copy(snap.Buffer[:], tmp.Bytes())
		tmp.Reset()
	case *MBC2:
		mbc2 := mapper.CreateSnapshot()
		binary.Write(tmp, binary.LittleEndian, mbc2)
		copy(snap.Buffer[:], tmp.Bytes())
		tmp.Reset()
	case *MBC3:
		mbc3 := mapper.CreateSnapshot()
		binary.Write(tmp, binary.LittleEndian, mbc3)
		copy(snap.Buffer[:], tmp.Bytes())
		tmp.Reset()
	case *MBC5:
		mbc5 := mapper.CreateSnapshot()
		binary.Write(tmp, binary.LittleEndian, mbc5)
		copy(snap.Buffer[:], tmp.Bytes())
		tmp.Reset()
	}
	return nil
}

func (c *Cartridge) RestoreSnapshot(snap *Snapshot) error {
	if snap == nil {
		return errSnapshotNil
	}
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
	case *MBC3:
		tmp.Write(snap.Buffer[:])
		var s MBC3Snapshot
		binary.Read(tmp, binary.LittleEndian, &s)
		mapper.RestoreSnapshot(&s)
		tmp.Reset()
	case *MBC5:
		tmp.Write(snap.Buffer[:])
		var s MBC5Snapshot
		binary.Read(tmp, binary.LittleEndian, &s)
		mapper.RestoreSnapshot(&s)
		tmp.Reset()
	}
	return nil
}
