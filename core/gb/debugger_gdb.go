package gb

import (
	"fmt"
	"strconv"
	"strings"
)

// Reference: https://www.chciken.com/tlmboy/2022/04/03/gdb-z80.html

// data は $[data]#[checksum] の data部分
func (g *GB) HandleGDB(data string) string {
	reply := ""

	switch data {
	case "g": // 全部のレジスタを返す
		reply = "%04X%04X%04X%04X%04X%04Xxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // AF BC DE HL SP PC
		reply = fmt.Sprintf(reply, ((uint16(g.CPU.R.A) << 8) | uint16(g.CPU.R.F.Pack())), g.CPU.R.BC.Pack(), g.CPU.R.DE.Pack(), g.CPU.R.HL.Pack(), g.CPU.R.SP, g.CPU.R.PC)
	default:
		switch data[0] {
		case 'm': // e.g. 'm34,c': アドレス0x0034 から 12バイト読み込んで返す
			tmp := strings.Split(data[1:], ",")
			addr := parseHex(tmp[0])
			size := parseHex(tmp[1])
			for i := 0; i < int(size); i++ {
				addr := uint32(addr) + uint32(i)
				b := uint8(g.ViewMemory(0, addr, 1))
				reply += int2hex8(b)
			}
		}
	}

	return reply
}

func parseHex(s string) uint64 {
	val, _ := strconv.ParseUint(s, 16, 64)
	return val
}

func int2hex8(val uint8) string {
	const language = "0123456789abcdef"
	out := []uint8{0, 0}
	out[0] = language[val>>4]
	out[1] = language[val&0xf]
	return string(out)
}
