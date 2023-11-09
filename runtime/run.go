package runtime

import (
	"fmt"
)

func run(input string) {
	scanner := newScanner(input)
	scanner.ScanTokens()

	for _, t := range scanner.Tokens {
		fmt.Println(t)
	}
}
