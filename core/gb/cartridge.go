package gb

import (
	. "github.com/akatsuki105/dugb/util/datasize"
)

type cartridge struct {
	entry [4]uint8
	title string
	rom   []uint8
	ram   [8 * KB]uint8 // SRAM
}

func (g *GB) loadCartridge(rom []uint8) {
	c := &cartridge{
		entry: [4]uint8{rom[0x100], rom[0x101], rom[0x102], rom[0x103]},
		title: string(rom[0x134:0x144]),
	}

	c.rom = make([]uint8, (32*KB)<<rom[0x148])
	copy(c.rom, rom)

	g.cartridge = c
}
