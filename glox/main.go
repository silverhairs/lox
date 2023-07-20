package main

import (
	"fmt"
	"glox/lox"
	"os"
)

func main() {
	args := os.Args[1:]
	runner := lox.NewRunner(os.Stderr, os.Stdout)

	if len(args) < 1 {
		runner.StartREPL(os.Stdin)

	} else if len(args) == 1 {
		runner.RunFile(args[0])

	} else {
		fmt.Printf("Usage: glox [script]\n")
		os.Exit(64)
	}
}
