package main

import "syscall/js"

func init() {
	js.Global().Set("increment", js.FuncOf(increment))
}

func increment(this js.Value, args []js.Value) any {
	return map[string]any{"message": 1}
}
