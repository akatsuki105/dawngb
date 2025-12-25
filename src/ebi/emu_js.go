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
	App.Emu.Reset = true
	return nil
}

func setPaused(this js.Value, args []js.Value) any {
	App.Emu.Paused = args[0].Bool()
	return nil
}

func enableSound(this js.Value, args []js.Value) any {
	return nil
}

func press(this js.Value, args []js.Value) any {
	for key := range Inputs {
		if key == args[0].String() {
			Inputs[key] = args[1].Bool()
			break
		}
	}
	return nil
}

func loadROM(this js.Value, args []js.Value) any {
	raw := args[0]
	rom := make([]uint8, raw.Get("length").Int())
	js.CopyBytesToGo(rom, raw)
	App.Emu.LoadROM(rom)
	return nil
}

func loadSave(this js.Value, args []js.Value) any {
	sram, err := App.Emu.Core.Dump(gb.DUMP_SAVE)
	if err == nil {
		size := len(sram)
		newSram := make([]uint8, size)
		js.CopyBytesToGo(newSram, args[0])
		App.Emu.LoadSave(newSram)
	}
	return nil
}

func dumpSave(this js.Value, args []js.Value) any {
	sram, err := App.Emu.Core.Dump(gb.DUMP_SAVE)
	if err != nil {
		return nil
	}
	dst := js.Global().Get("Uint8Array").New(len(sram))
	js.CopyBytesToJS(dst, sram)
	return dst
}
