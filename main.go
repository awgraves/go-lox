package main

import (
	"fmt"
	"os"

	"github.com/awgraves/go-lox/input"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		input.RunPrompt()
		return
	case 1:
		input.RunFile(args[0])
		return
	default:
		fmt.Println("Usage: lox [path/to/script.lx]")
	}
}
