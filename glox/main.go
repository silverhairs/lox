package main

import (
	"craftinginterpreters/cmd"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Usage: glox [script]\n")
		os.Exit(64)

	} else if len(args) == 1 {
		cmd.RunFile(args[0])

	} else {
		cmd.StartREPL(os.Stdin)
	}
}
