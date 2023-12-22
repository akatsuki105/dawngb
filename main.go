package main

import "os"

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
	return exitCodeOK
}
