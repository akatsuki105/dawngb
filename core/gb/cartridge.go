package gb

import (
	"encoding/binary"

	. "github.com/akatsuki105/dugb/util/datasize"
)

type cartridge struct {
	entry [2]uint16
	title string
	rom   []uint8
	ram   [8 * KB]uint8 // SRAM
}

func (g *GB) loadCartridge(rom []uint8) {
	c := &cartridge{
		entry: [2]uint16{binary.LittleEndian.Uint16(rom[0x100:0x102]), binary.LittleEndian.Uint16(rom[0x102:0x104])},
		title: string(rom[0x134:0x144]),
	}

	c.rom = make([]uint8, (32*KB)<<rom[0x148])
	copy(c.rom, rom[0x100:])

	g.cartridge = c
}
