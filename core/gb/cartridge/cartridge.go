package cartridge

import (
	"fmt"

	. "github.com/akatsuki105/dugb/util/datasize"
)

var ramSizes = map[uint8]uint{
	2: 8 * KB,
	3: 32 * KB,
	4: 128 * KB,
	5: 64 * KB,
}

type mbc interface {
	read(addr uint16) uint8
	write(addr uint16, val uint8)
}

type Cartridge struct {
	title string
	rom   []uint8
	ram   []uint8 // SRAM
	mbc
}

func New(rom []uint8) *Cartridge {
	isCGB := rom[0x143] == 0x80 || rom[0x143] == 0xC0
	title := string(rom[0x134:0x13F])
	if !isCGB {
		title = string(rom[0x134:0x144])
	}

	c := &Cartridge{
		title: title,
		ram:   make([]uint8, 0),
	}

	c.rom = make([]uint8, (32*KB)<<rom[0x148])
	copy(c.rom, rom)

	ramSize, ok := ramSizes[rom[0x149]]
	if ok {
		c.ram = make([]uint8, ramSize)
	}

	c.mbc = createMBC(c)
	return c
}

func createMBC(c *Cartridge) mbc {
	mbcType := c.rom[0x147]
	switch mbcType {
	case 0:
		return newMBC0(c)
	case 1, 3:
		return newMBC1(c)
	case 16, 19:
		return newMBC3(c)
	case 27:
		return newMBC5(c)
	default:
		panic(fmt.Sprintf("unsupported mbc type: 0x%02X", mbcType))
	}
}

func (c *Cartridge) Title() string {
	return c.title
}

func (c *Cartridge) Read(addr uint16) uint8 {
	return c.mbc.read(addr)
}

func (c *Cartridge) Write(addr uint16, val uint8) {
	c.mbc.write(addr, val)
}

func (c *Cartridge) IsCGB() bool {
	return c.rom[0x143] == 0x80 || c.rom[0x143] == 0xC0
}