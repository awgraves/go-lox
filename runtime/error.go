package runtime

import "fmt"

func printError(message string) {
	fmt.Print(RED)
	fmt.Println(message)
	fmt.Print(RESET_COLOR)
}

func reportFileError(lineNum int, charIdx int, message string) {
	m := fmt.Sprintf("[line %d pos %d] Error: %s", lineNum, charIdx, message)
	printError(m)
}

func reportPromptError(charIdx int, message string) {
	m := fmt.Sprintf("[pos %d] Error: %s", charIdx, message)
	printError(m)
}
