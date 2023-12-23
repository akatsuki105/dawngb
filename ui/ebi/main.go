package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// ExitCode represents program's status code
type exitCode int

// exit code
const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func main() {
	os.Exit(int(Run()))
}

func Run() exitCode {
	flag.Parse()

	if flag.NArg() == 0 {
		msg := "rom path is not specified"
		fmt.Fprintln(os.Stderr, msg)
		return exitCodeError
	}

	romPath := flag.Arg(0)
	e := createEmu()
	e.LoadROM(romPath)

	w, h := e.Layout(0, 0)
	ebiten.SetWindowSize(w*2, h*2)
	ebiten.SetWindowTitle(fmt.Sprintf("DuGB - %s", e.Title()))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(e); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	return exitCodeOK
}
