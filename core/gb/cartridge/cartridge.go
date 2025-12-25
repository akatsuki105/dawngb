package cartridge

import (
	"fmt"
)

const KB, MB = 1024, 1024 * 1024

type MBC interface {
	read(addr uint16) uint8
	write(addr uint16, val uint8)
}

type Cartridge struct {
	ROM []uint8
	RAM []uint8 // SRAM
	MBC         // mapper
}

func New(rom []uint8) (*Cartridge, error) {
	c := &Cartridge{
		RAM: make([]uint8, 0),
	}

	c.ROM = make([]uint8, calcROMSize(rom[0x148]))
	copy(c.ROM, rom)

	c.RAM = make([]uint8, calcSRAMSize(rom[0x149])) // これはSRAMチップのサイズであって、MBC2のようなMBCチップにRAMが内蔵されている場合は0になるっぽい

	mbc, err := createMBC(c)
	if err != nil {
		return nil, err
	}
	c.MBC = mbc

	return c, nil
}

func createMBC(c *Cartridge) (MBC, error) {
	mbcType := c.ROM[0x147]
	switch mbcType {
	case 0:
		return newMBC0(c), nil
	case 1, 2, 3:
		return newMBC1(c), nil
	case 5, 6:
		c.RAM = make([]uint8, 512)
		return newMBC2(c), nil
	case 16, 19:
		return newMBC3(c), nil
	case 25, 26, 27:
		return newMBC5(c), nil
	default:
		return nil, fmt.Errorf("unsupported mbc type: 0x%02X", mbcType)
	}
}

func (c *Cartridge) Read(addr uint16) uint8 {
	return c.MBC.read(addr)
}

func (c *Cartridge) Write(addr uint16, val uint8) {
	c.MBC.write(addr, val)
}

func (c *Cartridge) CGBFlag() uint8 {
	return c.ROM[0x143]
}

func (c *Cartridge) LoadSRAM(data []uint8) error {
	copy(c.RAM, data)
	return nil
}

func (c *Cartridge) SRAM() []uint8 {
	return c.RAM
}

func (c *Cartridge) RAMSize() uint {
	switch c.MBC.(type) {
	case *MBC2:
		return 512 // 正確には512バイトの下位4ビットしか使わないので、Packすれば256バイト
	default: // SRAM
		return uint(len(c.RAM))
	}
}

func (c *Cartridge) ROMBankNumber() uint16 {
	switch mbc := c.MBC.(type) {
	case *MBC1:
		return uint16(mbc.ROMBank)
	case *MBC2:
		return uint16(mbc.ROMBank)
	case *MBC3:
		return uint16(mbc.ROMBank)
	case *MBC5:
		return mbc.ROMBank
	}
	return 1
}

func calcROMSize(n uint8) int {
	switch n {
	case 0x52:
		return 72 * (16 * KB)
	case 0x53:
		return 80 * (16 * KB)
	case 0x54:
		return 96 * (16 * KB)
	}
	return (32 * KB) << n
}

func calcSRAMSize(n uint8) int {
	switch n {
	case 1:
		return 2 * KB // 実際に商用ゲームで使われたことがなくHomebrewによる慣用的なもの
	case 2:
		return 8 * KB
	case 3:
		return 32 * KB
	case 4:
		return 128 * KB
	case 5:
		return 64 * KB
	}
	return 0
}
