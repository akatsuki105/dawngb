package cpu

import (
	"fmt"
	"math"

	"github.com/akatsuki105/dugb/util/scheduler"
)

const CLOCK = 1

type Cpu struct {
	r         Registers
	Cycles    int64
	NextEvent int64
	Halted    bool
	Blocked   bool
	EarlyExit bool // 描画やサウンドの準備が整った時に、残ったイベントの処理(.Cyclesのコミット)をせずに ProcessEvents を終了して、フロントエンド側に処理を移すためのフラグ
	inst      struct {
		opcode uint8
		addr   uint16
	}
	s *scheduler.Scheduler
}

func New(s *scheduler.Scheduler) *Cpu {
	return &Cpu{
		s: s,
	}
}

func (c *Cpu) Reset() {
	c.Cycles = 0
	c.NextEvent = 0
	c.s.Reset(&c.Cycles, &c.NextEvent)
}

func (c *Cpu) Step() {
	pc := c.r.pc
	c.inst.addr = pc
	opcode := c.fetch()
	c.inst.opcode = opcode

	fn := opTable[opcode]
	if fn != nil {
		fn(c)
	} else {
		panic(fmt.Sprintf("illegal opcode: 0x%02X in 0x%04X", opcode, pc))
	}

	c.Cycles += int64(opCycles[opcode] * CLOCK)
}

// This is called when `c.Cycles >= c.NextEvent` is satisfied.
func (c *Cpu) ProcessEvents() {
	nextEvent := c.NextEvent
	for c.Cycles >= nextEvent {
		c.NextEvent = math.MaxInt64
		nextEvent = 0

		// Blockedを考慮しているせいで複雑になっているが、基本的にc.Cyclesをコミットしているだけ
		first := true
		for first || (c.Blocked && !c.EarlyExit) {
			first = false

			cycles := c.Cycles
			c.Cycles = 0
			if cycles < 0 {
				panic(fmt.Sprintf("negative cycles passed: %d", cycles))
			}

			if cycles < nextEvent {
				// こっちにくるのは first == false　のとき、つまりCPUがDMAにBlockedされて停止しているとき
				nextEvent = c.s.Add(nextEvent) // CPUが停止しているので一気に進めてOK
			} else {
				// first == true のときは、(nextEvent が 0なので)、 こっち
				nextEvent = c.s.Add(cycles)
			}
		}

		c.NextEvent = nextEvent
		if c.Halted {
			c.Cycles = nextEvent
			break
		} else if nextEvent < 0 {
			panic(fmt.Sprintf("negative cycles passed: %d", nextEvent))
		}

		if c.EarlyExit {
			break
		}
	}

	c.EarlyExit = false
	if c.Blocked {
		c.Cycles = c.NextEvent
	}
}

func (c *Cpu) fetch() uint8 {
	c.r.pc++
	return 0
}
