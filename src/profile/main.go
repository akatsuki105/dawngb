package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/akatsuki105/dawngb/core/gb"
	"github.com/pkg/profile"
)

// ExitCode represents program's status code
type exitCode int

// exit code
const (
	ExitCodeOK exitCode = iota
	ExitCodeError
)

var (
	s = flag.Int("s", 30, "How many seconds to run the emulator.")
)

func main() {
	os.Exit(int(run()))
}

func run() exitCode {
	flag.Parse()
	if flag.NArg() > 0 {
		defer profile.Start(profile.ProfilePath("./build/profile")).Stop()

		c := gb.New(gb.MODEL_CGB, nil)

		rom, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitCodeError
		}
		err = c.Load(gb.LOAD_ROM, rom)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitCodeError
		}

		c.Reset(false)
		fmt.Printf("Run emulator for %d seconds\n", *s)

		for i := 0; i < (*s)*60; i++ {
			c.RunFrame()
			if i%60 == 0 {
				fmt.Printf("%d sec\n", i/60+1)
			}
		}
	}

	return ExitCodeOK
}
