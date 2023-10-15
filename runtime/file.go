package runtime

import (
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
