package runtime

import (
	"fmt"
)

func run(input string) {
	errReporter := newBasicErrorReporter()
	scanner := newScanner(input, errReporter)
	scanner.ScanTokens()

	if errReporter.HasError() {
		printError("Errors found - Would not attempt to execute this code.")
		errReporter.Report()
	}

	for _, t := range scanner.Tokens {
		fmt.Println(t)
	}
}
