package runtime

import "fmt"

func printError(message string) {
	fmt.Print(RED)
	fmt.Println(message)
	fmt.Print(RESET_COLOR)
}

type ErrorReporter interface {
	AddError(lineNum int, charIdx int, message string)
	HasError() bool
	Report()
}

type basicErrorReporter struct {
	errorMsgs []string
}

func newBasicErrorReporter() *basicErrorReporter {
	return &basicErrorReporter{
		errorMsgs: []string{},
	}
}

func (b *basicErrorReporter) AddError(lineNum int, charIdx int, message string) {
	m := fmt.Sprintf("[line %d pos %d] Error: %s", lineNum, charIdx, message)
	b.errorMsgs = append(b.errorMsgs, m)
}

func (b *basicErrorReporter) HasError() bool {
	return len(b.errorMsgs) > 0
}

func (b *basicErrorReporter) Report() {
	for _, m := range b.errorMsgs {
		printError(m)
	}
}
