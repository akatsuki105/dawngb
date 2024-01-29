package main

import (
	"fmt"
	"syscall/js"
)

func init() {
	js.Global().Set("press", js.FuncOf(press))
}

func press(this js.Value, args []js.Value) any {
	fmt.Println("press", args[0].String())
	inputMap[args[0].String()] = true
	return nil
}
