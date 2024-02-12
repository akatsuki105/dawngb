package gb

import (
	"encoding/binary"
	"io"

	"github.com/akatsuki105/dawngb/util"
)

var timaClock = [4]int64{64, 1, 4, 16}

type timer struct {
	g              *GB
	cycles         int64 // CPUから見て遅れているマスターサイクル数
	tima, tma, tac uint8
	counter        int64 // 524288Hz(一番細かいのが524288Hzなのであとはそれの倍数で数えれば良い)
}

func newTimer(g *GB) *timer { return &timer{g: g} }

func (t *timer) Reset(hasBIOS bool) {
	t.cycles = 0
	t.counter, t.tima, t.tma, t.tac = 0, 0, 0, 0
	if !hasBIOS {
		t.tac = 0xF8
	}
}

func (t *timer) tick(cycles int64) {
	t.cycles += cycles
	for t.cycles >= 16 {
		t.update()
		t.cycles -= 16
	}
}

func (t *timer) Read(addr uint16) uint8 {
	x := t.g.cpu.Cycle / 4
	switch addr {
	case 0xFF04:
		div := t.counter / (16 * x)
		return uint8(div)
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	}
	return 0
}

func (t *timer) Write(addr uint16, val uint8) {
	switch addr {
	case 0xFF04:
		t.counter = 0
	case 0xFF05:
		t.tima = val
	case 0xFF06:
		t.tma = val
	case 0xFF07:
		t.tac = val & 0b111
	}
}

// 524288Hz
func (t *timer) update() {
	x := t.g.cpu.Cycle / 4

	t.counter++
	if util.Bit(t.tac, 2) {
		if (t.counter % (timaClock[t.tac&0b11] * x)) == 0 {
			t.tima++
			if t.tima == 0 {
				t.tima = t.tma
				t.g.requestInterrupt(2)
			}
		}
	}
}

func (t *timer) Serialize(s io.Writer) {
	data := []uint8{}
	binary.LittleEndian.PutUint64(data, uint64(t.cycles))  // 8
	data = append(data, t.tima, t.tma, t.tac)              // 3
	binary.LittleEndian.PutUint64(data, uint64(t.counter)) // 8
	s.Write(data)
}

func (t *timer) Deserialize(s io.Reader) {
	data := make([]uint8, 19)
	s.Read(data)
	t.cycles = int64(binary.LittleEndian.Uint64(data[0:8]))
	t.tima, t.tma, t.tac = data[8], data[9], data[10]
	t.counter = int64(binary.LittleEndian.Uint64(data[11:19]))
}
