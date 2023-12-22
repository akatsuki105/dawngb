package cpu

var cbTable = [256]opcode{
	/* 0x00 */ cb00,
}

func cb00(c *Cpu) {}
