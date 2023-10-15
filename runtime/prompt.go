package runtime

import (
	"bufio"
	"fmt"
	"os"
)

func RunPrompt() {
	fmt.Print(GREEN)
	fmt.Println("Lox Shell v0.0")
	fmt.Print(BLUE)
	fmt.Println("Type 'exit' to quit.")
	fmt.Print(RESET_COLOR)
	fmt.Println()

	promptLoop()
}

func promptLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text()
		if line == "exit" {
			break
		}
		if line == "" {
			continue
		}
		run(line)
	}
}
