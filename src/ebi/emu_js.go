package main

import "syscall/js"

func init() {
	js.Global().Set("press", js.FuncOf(press))
}

func press(this js.Value, args []js.Value) any {
	inputMap[args[0].String()] = true
	return nil
}
