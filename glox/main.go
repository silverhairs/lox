package main

import (
	"fmt"
	"glox/cmd"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		cmd.StartREPL(os.Stdin)

	} else if len(args) == 1 {
		cmd.RunFile(args[0])

	} else {
		fmt.Printf("Usage: glox [script]\n")
		os.Exit(64)
	}
}
