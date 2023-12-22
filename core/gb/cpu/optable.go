package cpu

type opcode = func(c *Cpu)

var opTable = [256]opcode{
	/* 0x00 */ op00,
}

var opCycles = [256]int64{}

func op00(c *Cpu) {}
