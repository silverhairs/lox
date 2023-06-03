package main

import (
	"craftinginterpreters/glox"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Printf("Number or args: %d\n The argument is %q", len(args), args)
	if len(args) < 1 {
		fmt.Printf("Usage: glox [script]\n")
		os.Exit(64)

	} else if len(args) == 1 {
		glox.RunFile(args[0])

	} else {
		glox.StartREPL(os.Stdin)
	}
}
