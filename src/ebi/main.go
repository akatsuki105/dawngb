package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const EXPAND = 2.

// ExitCode represents program's status code
type exitCode int

// exit code
const (
	exitCodeOK exitCode = iota
	exitCodeError
)

var (
	turbo = flag.Int("t", 1, "Emulator speed xN")
	sound = flag.Bool("s", false, "Enable sound")
)

func main() {
	os.Exit(int(Run()))
}

func Run() exitCode {
	flag.Parse()

	e := createEmu()

	if flag.NArg() > 0 {
		e.LoadROMFromPath(flag.Arg(0))
	}

	w, h := e.Layout(0, 0)
	ebiten.SetWindowSize(int(float64(w)*EXPAND), int(float64(h)*EXPAND))
	ebiten.SetWindowTitle("DawnGB")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(false)

	defer func() {
		if e.music != nil {
			e.music.Close()
		}
	}()

	if *turbo > 1 {
		e.setTurbo(*turbo)
	}
	e.enableSound(*sound)

	if err := ebiten.RunGame(e); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}

	return exitCodeOK
}
