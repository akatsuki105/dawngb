package gb

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"runtime"

	"github.com/akatsuki105/dawngb/core/gb/apu"
	"github.com/akatsuki105/dawngb/core/gb/cartridge"
	"github.com/akatsuki105/dawngb/core/gb/cpu"
	"github.com/akatsuki105/dawngb/core/gb/ppu"
	"github.com/akatsuki105/dawngb/internal/debugger"
)

const KB, MB = 1024, 1024 * 1024

type Model uint8

// ハードウェアの種類
const (
	MODEL_DMG Model = iota
	MODEL_SGB
	MODEL_CGB
	MODEL_AGB
)

type LoadCmd = uint8

const (
	LOAD_ROM LoadCmd = iota
	LOAD_SAVE
	LOAD_BIOS
)

type DumpCmd = uint8

const (
	DUMP_SAVE DumpCmd = iota
)

var buttons = [8]string{"A", "B", "SELECT", "START", "RIGHT", "LEFT", "UP", "DOWN"}

type GB struct {
	Model  Model // ハードウェアの種類
	CPU    *cpu.CPU
	PPU    *ppu.PPU
	APU    *apu.APU
	Cart   *cartridge.Cartridge
	inputs uint8 // 押されている時にビットを立てる; bit0: A, bit1: B, bit2: SELECT, bit3: START, bit4: RIGHT, bit5: LEFT, bit6: UP, bit7: DOWN
	WRAM   WRAM
	Snap   Snapshot
	debugger.Debugger
}

type WRAM struct {
	Data [(4 * KB) * 8]uint8
	Bank uint8 // SVBK(0xFF70, 0..7, CGB only)
}

func New(model Model, audioBuffer io.Writer) *GB {
	g := &GB{
		Model: model,
		Snap:  *NewSnapshot(0),
	}
	g.CPU = cpu.New(g.IsColor(), g)
	g.PPU = ppu.New(g.CPU)
	g.APU = apu.New(audioBuffer)
	g.WRAM.Bank = 1
	return g
}

func (g *GB) Reset() {
	if g.Cart != nil {
		clear(g.WRAM.Data[:])
		g.WRAM.Bank = 1
		g.CPU.Reset()
		g.PPU.Reset()
		g.APU.Reset()
		g.inputs = 0
	}
}

func (g *GB) DirectBoot() error {
	g.skipBIOS()
	return nil
}

func (g *GB) Quit() {}

func (g *GB) onPanic() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "Frame: %d, y: %d, PC: 0x%04X\n", g.PPU.Frame, g.PPU.Ly, g.CPU.R.PC)
		for depth := 0; ; depth++ {
			_, file, line, ok := runtime.Caller(depth)
			if !ok {
				break
			}
			fmt.Fprintf(os.Stderr, "======> %d: %v:%d\n", depth, file, line)
		}
		panic(r)
	}
}

func (g *GB) skipBIOS() {
	g.CPU.SkipBIOS()
	g.PPU.SkipBIOS()
	g.APU.SkipBIOS()
	g.Write(0xFF02, 0x7F) // SC
	g.Write(0xFF0F, 0xE1) // IF
	cgbflag := g.Cart.ROM[0x143]
	if cgbflag&0x80 == 0 {
		g.Write(0xFF4C, 4) // KEY0
	}
	g.Write(0xFF4D, 0x7E) // KEY1
	g.Write(0xFF4F, 0xFE) // VBK

	if g.IsColor() && cgbflag&0x80 == 0 {
		g.PPU.ColorizeDMG()
	}
}

var errInvalidCmd = fmt.Errorf("invalid command")

func (g *GB) Load(cmd LoadCmd, args ...any) error {
	switch cmd {
	case LOAD_ROM:
		if len(args) != 1 {
			return fmt.Errorf("LOAD_ROM command requires []uint8")
		}
		data, ok := args[0].([]uint8)
		if !ok {
			return fmt.Errorf("LOAD_ROM command requires []uint8")
		}
		cartridge, err := cartridge.New(data)
		if err != nil {
			return err
		}
		g.Cart = cartridge

	case LOAD_SAVE:
		if len(args) != 1 {
			return fmt.Errorf("LOAD_SAVE command requires []uint8")
		}
		if g.Cart == nil {
			return fmt.Errorf("no cartridge loaded")
		}
		sram, ok := args[0].([]uint8)
		if !ok {
			return fmt.Errorf("LOAD_SAVE command requires []uint8")
		}
		err := g.Cart.LoadSRAM(sram)
		if err != nil {
			return err
		}

	case LOAD_BIOS:
		if len(args) != 1 {
			return fmt.Errorf("LOAD_BIOS command requires []uint8")
		}
		bios, ok := args[0].([]uint8)
		if !ok {
			return fmt.Errorf("LOAD_BIOS command requires []uint8")
		}
		err := g.CPU.LoadBIOS(bios)
		if err != nil {
			return err
		}

	default:
		return errInvalidCmd
	}

	return nil
}

func (g *GB) LoadROM(rom []uint8) error {
	err := g.Load(LOAD_ROM, rom)
	if err != nil {
		return err
	}
	g.Reset()
	g.DirectBoot()
	return nil
}

func (g *GB) LoadSave(savedata []uint8) error {
	return g.Load(LOAD_SAVE, savedata)
}

func (g *GB) Dump(cmd DumpCmd, args ...any) ([]uint8, error) {
	switch cmd {
	case DUMP_SAVE:
		if g.Cart == nil {
			return []uint8{}, fmt.Errorf("no cartridge loaded")
		}
		return g.Cart.SRAM(), nil
	default:
		return nil, errInvalidCmd
	}
}

func (g *GB) RunFrame() {
	if g.Cart != nil {
		// defer g.onPanic()

		g.CPU.SendInputs(g.inputs ^ 0xFF) // ボタンの状態をCPUに送る(ただし、押されてないボタンのビットを立てる)
		g.inputs = 0

		g.CPU.Usage = 0

		const FRAME = 70224 * ppu.CYCLE
		start := g.CPU.Cycles

		frame := g.PPU.Frame
		for frame == g.PPU.Frame && ((g.CPU.Cycles - start) < FRAME) {
			g.step()
		}
		g.APU.FlushSamples()
	}
}

func (g *GB) step() {
	delta := g.CPU.Step() // CPUで1命令実行して、その後に他のコンポーネントを同期させる
	g.PPU.Run(delta)
	g.APU.Run(delta)
}

func (g *GB) Resolution() (w int, h int) { return 160, 144 }
func (g *GB) Screen() []color.NRGBA      { return g.PPU.Screen() }

func (g *GB) SetKeyInput(key string, press bool) {
	if press {
		for i, b := range buttons {
			if b == key {
				g.inputs |= (1 << uint(i))
				break
			}
		}
	}
}

// IsCGBMode returns true if the hardware has CGB features (i.e., it's a CGB or AGB).
func (g *GB) IsColor() bool {
	return g.Model == MODEL_CGB || g.Model == MODEL_AGB
}
