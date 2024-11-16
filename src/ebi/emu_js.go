package main

import (
	"syscall/js"

	"github.com/akatsuki105/dawngb/core/gb"
)

func init() {
	js.Global().Set("reset", js.FuncOf(reset))
	js.Global().Set("setPaused", js.FuncOf(setPaused))
	js.Global().Set("sound", js.FuncOf(enableSound))
	js.Global().Set("press", js.FuncOf(press))
	js.Global().Set("loadROM", js.FuncOf(loadROM))
	js.Global().Set("loadSave", js.FuncOf(loadSave))
	js.Global().Set("dumpSave", js.FuncOf(dumpSave))
}

func reset(this js.Value, args []js.Value) any {
	if emu != nil {
		emu.c.Reset(false)
	}
	return nil
}

func setPaused(this js.Value, args []js.Value) any {
	if emu != nil {
		emu.setPaused(args[0].Bool())
	}
	return nil
}

func enableSound(this js.Value, args []js.Value) any {
	if emu != nil {
		emu.enableSound(args[0].Bool())
	}
	return nil
}

func press(this js.Value, args []js.Value) any {
	for key := range inputMapWeb {
		if key == args[0].String() {
			inputMapWeb[key] = args[1].Bool()
			break
		}
	}
	return nil
}

func loadROM(this js.Value, args []js.Value) any {
	if emu != nil {
		raw := args[0]
		rom := make([]uint8, raw.Get("length").Int())
		js.CopyBytesToGo(rom, raw)
		emu.soundEnabled = true
		emu.LoadROM(rom)
		emu.c.Reset(false)
	}
	return nil
}

func loadSave(this js.Value, args []js.Value) any {
	if emu != nil {
		sram, err := emu.c.Dump(gb.DUMP_SAVE)
		if err == nil {
			size := len(sram)
			newSram := make([]uint8, size)
			js.CopyBytesToGo(newSram, args[0])
			emu.c.Load(gb.LOAD_SAVE, newSram)
			emu.c.Reset(false)
		}
	}
	return nil
}

func dumpSave(this js.Value, args []js.Value) any {
	if emu != nil {
		sram, err := emu.c.Dump(gb.DUMP_SAVE)
		if err == nil {
			dst := js.Global().Get("Uint8Array").New(len(sram))
			js.CopyBytesToJS(dst, sram)
			return dst
		}
	}
	return nil
}
