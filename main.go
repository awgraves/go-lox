package main

import (
	"fmt"
	"os"

	"github.com/awgraves/go-lox/runtime"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		runtime.RunPrompt()
		return
	case 1:
		runtime.RunFile(args[0])
		return
	default:
		fmt.Println("Usage: lox [path/to/script.lx]")
	}
}
