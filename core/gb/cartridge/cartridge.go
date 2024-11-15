package cartridge

import (
	"fmt"

	"github.com/akatsuki105/dawngb/util"
)

const KB, MB = 1024, 1024 * 1024

var RAM_SIZES = map[uint8]uint{
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

func New(rom []uint8) (*Cartridge, error) {
	isCGB := util.Bit(rom[0x143], 7)

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

	ramSize, ok := RAM_SIZES[rom[0x149]]
	if ok {
		c.ram = make([]uint8, ramSize)
	}

	mbc, err := createMBC(c)
	if err != nil {
		return nil, err
	}
	c.mbc = mbc
	fmt.Println("MapperID:", c.rom[0x147])
	return c, nil
}

func createMBC(c *Cartridge) (mbc, error) {
	mbcType := c.rom[0x147]
	switch mbcType {
	case 0:
		return newMBC0(c), nil
	case 1, 2, 3:
		return newMBC1(c), nil
	case 16, 19:
		return newMBC3(c), nil
	case 25, 26, 27:
		return newMBC5(c), nil
	default:
		return nil, fmt.Errorf("unsupported mbc type: 0x%02X", mbcType)
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
	return (c.rom[0x143] & (1 << 7)) != 0
}

func (c *Cartridge) LoadSRAM(data []uint8) error {
	copy(c.ram, data)
	return nil
}

func (c *Cartridge) SRAM() []uint8 {
	return c.ram
}
