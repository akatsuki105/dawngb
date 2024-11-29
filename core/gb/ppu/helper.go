package ppu

type rgb555 = uint16 // Litte Endian (0b0BBBBBGGGGGRRRRR); e.g. 0x6180 is RGB555(0, 12, 24)

var dmgPalette = [4]rgb555{
	0x7FFF, // (31, 31, 31)
	0x56B5, // (15, 15, 15)
	0x294A, // (10, 10, 10)
	0x0000, // (0, 0, 0)
}

// DMGのゲームをCGBで起動した場合のデフォルトのパレット
var cgbPalette = [8]rgb555{
	// BG (Green)
	0x7FFF, // (31, 31, 31)
	0x1BEF, // (15, 31, 6)
	0x6180, // (0, 12, 24)
	0x0000, // (0, 0, 0)

	// OBJ (Red)
	0x7FFF, // (31, 31, 31)
	0x421F, // (31, 15, 15)
	0x1CF2, // (18, 7, 7)
	0x0000, // (0, 0, 0)
}
