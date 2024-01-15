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

var (
	turbo = flag.Int("t", 1, "Emulator speed xN")
)

func main() {
	os.Exit(int(Run()))
}

func Run() exitCode {
	flag.Parse()

	e := createEmu(false)

	if flag.NArg() > 0 {
		e.LoadROMFromPath(flag.Arg(0))
	}

	w, h := e.Layout(0, 0)
	ebiten.SetWindowSize(w*2, h*2)
	ebiten.SetWindowTitle(e.Title())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	defer func() {
		if e.music != nil {
			e.music.Close()
			e.context.Close()
		}
	}()

	if *turbo > 1 {
		e.Turbo(*turbo)
	}

	if err := ebiten.RunGame(e); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	return exitCodeOK
}
