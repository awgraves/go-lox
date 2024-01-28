package runtime

import (
	"bufio"
	"fmt"
	"os"
)

func RunFile(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		printError(fmt.Sprintf("Invalid file path: %s\n", filePath))
		os.Exit(1)
	}

	run(string(bytes))
}

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

func run(input string) {
	errReporter := newBasicErrorReporter()
	scanner := newScanner(input, errReporter)
	scanner.ScanTokens()

	if errReporter.HasError() {
		printError("Errors found - runtime would not attempt to execute this code.")
		errReporter.Report()
		fmt.Println()
		return
	}

	parser := newParser(scanner.Tokens, errReporter)
	statements := parser.parse()

	if errReporter.HasError() {
		printError("Errors found - runtime would not attempt to execute this code.")
		errReporter.Report()
		fmt.Println()
		return
	}

	interpreter := newIntepreter(errReporter)
	interpreter.interpret(statements)

	if errReporter.HasError() {
		printError("Runtime error")
		errReporter.Report()
		fmt.Println()
		return
	}

	fmt.Println()
}
