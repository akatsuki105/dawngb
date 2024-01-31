package main

import (
	"syscall/js"
)

func init() {
	js.Global().Set("press", js.FuncOf(press))
	js.Global().Set("save", js.FuncOf(save))
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

func save(this js.Value, args []js.Value) any {
	if emu != nil {
		return emu.c.SRAM()
	}
	return nil
}
