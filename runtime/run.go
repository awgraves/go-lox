package runtime

import (
	"bufio"
	"fmt"
	"strings"
)

func run(input string) {
	reader := strings.NewReader(input)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	// TODO: skip comment lines and parse ; separately
	for {
		if hasNext := scanner.Scan(); !hasNext {
			break
		}
		token := scanner.Text()
		fmt.Println(token)
	}
}
