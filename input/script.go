package input

import (
	"fmt"
	"os"
)

func RunFile(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Print(RED)
		fmt.Printf("Invalid file path: %s\n", filePath)
		os.Exit(1)
	}

	fmt.Print(string(bytes))
}
