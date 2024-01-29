package main

import (
	"syscall/js"
)

func init() {
	js.Global().Set("press", js.FuncOf(press))
}

func press(this js.Value, args []js.Value) any {
	for key := range inputMapWeb {
		if key == args[0].String() {
			inputMapWeb[key] = true
			break
		}
	}
	return nil
}
