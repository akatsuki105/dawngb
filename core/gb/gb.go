package gb

import (
	"fmt"
	"image/color"
	"io"

	"github.com/akatsuki105/dawngb/core/gb/apu"
	"github.com/akatsuki105/dawngb/core/gb/cartridge"
	"github.com/akatsuki105/dawngb/core/gb/cpu"
	"github.com/akatsuki105/dawngb/core/gb/ppu"
)

const KB, MB = 1024, 1024 * 1024

const (
	MODEL_DMG = iota
	MODEL_CGB
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
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	apu       *apu.APU
	cartridge *cartridge.Cartridge
	timer     *timer
	serial    *serial
	inputs    uint8 // 押されている時にビットを立てる; bit0: A, bit1: B, bit2: SELECT, bit3: START, bit4: RIGHT, bit5: LEFT, bit6: UP, bit7: DOWN
	wram      [(4 * KB) * 8]uint8
	wramBank  uint // SVBK(0xFF70, CGB only)
}

func New(audioBuffer io.Writer) *GB {
	g := &GB{}
	g.cpu = cpu.New(g)
	g.ppu = ppu.New(g, g.cpu.IRQ, g.cpu.StartHDMA)
	g.apu = apu.New(audioBuffer)
	g.timer = newTimer(g)
	g.serial = newSerial(g.cpu.IRQ)
	g.wramBank = 1
	return g
}

func (g *GB) Reset(hasBIOS bool) {
	model := MODEL_DMG
	if g.cartridge != nil && g.cartridge.IsCGB() {
		model = MODEL_CGB
	}

	clear(g.wram[:])
	g.wramBank = 1
	g.cpu.Reset(hasBIOS)
	g.ppu.Reset(model, hasBIOS)
	g.apu.Reset(hasBIOS)
	g.timer.Reset(hasBIOS)
	g.serial.Reset(hasBIOS)
	g.inputs = 0

	if !hasBIOS {
		g.Write(0xFF02, 0x7F) // SC
		g.Write(0xFF0F, 0xE1) // IF
		if model == MODEL_CGB {
			g.Write(0xFF4D, 0x7E) // KEY1
			g.Write(0xFF4F, 0xFE) // VBK
		}
	}
}

var errInvalidCmd = fmt.Errorf("invalid command")

func (g *GB) Load(cmd LoadCmd, args ...any) error {
	switch cmd {
	case LOAD_ROM:
		if len(args) != 1 {
			return fmt.Errorf("LOAD_ROM command requires []uint8")
		}
		rom, ok := args[0].([]uint8)
		if !ok {
			return fmt.Errorf("LOAD_ROM command requires []uint8")
		}
		cartridge, err := cartridge.New(rom)
		if err != nil {
			return err
		}
		g.cartridge = cartridge

	case LOAD_SAVE:
		if len(args) != 1 {
			return fmt.Errorf("LOAD_SAVE command requires []uint8")
		}
		if g.cartridge == nil {
			return fmt.Errorf("no cartridge loaded")
		}
		sram, ok := args[0].([]uint8)
		if !ok {
			return fmt.Errorf("LOAD_SAVE command requires []uint8")
		}
		err := g.cartridge.LoadSRAM(sram)
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
		err := g.cpu.LoadBIOS(bios)
		if err != nil {
			return err
		}

	default:
		return errInvalidCmd
	}

	return nil
}

func (g *GB) Dump(cmd DumpCmd, args ...any) ([]uint8, error) {
	switch cmd {
	case DUMP_SAVE:
		if g.cartridge == nil {
			return []uint8{}, fmt.Errorf("no cartridge loaded")
		}
		return g.cartridge.SRAM(), nil
	default:
		return nil, errInvalidCmd
	}
}

func (g *GB) RunFrame() {
	if g.cartridge != nil {
		g.cpu.SendInputs(g.inputs ^ 0xFF) // ボタンの状態をCPUに送る(ただし、押されてないボタンのビットを立てる)
		g.inputs = 0

		const FRAME = 70224 * ppu.CYCLE
		start := g.cpu.Cycles

		frame := g.ppu.FrameCounter
		for frame == g.ppu.FrameCounter && ((g.cpu.Cycles - start) < FRAME) {
			g.step()
		}
		g.apu.FlushSamples()
	}
}

func (g *GB) step() {
	delta := g.cpu.Step() // CPUで1命令実行して、その後に他のコンポーネントを同期させる
	g.ppu.Run(delta)
	g.apu.Run(delta)
	g.timer.run(delta)
	g.serial.run(delta)
}

func (g *GB) Resolution() (w int, h int) { return 160, 144 }

func (g *GB) Screen() []color.NRGBA {
	return g.ppu.Screen()
}

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

func (g *GB) Title() string {
	if g.cartridge == nil {
		return ""
	}
	return g.cartridge.Title()
}

func (g *GB) Serialize(state io.Writer) {
	// TODO: implement
	g.timer.Serialize(state)
}

func (g *GB) Deserialize(state io.Reader) {
	// TODO: implement
	g.timer.Deserialize(state)
}
