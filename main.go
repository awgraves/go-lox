package main

import (
	"fmt"
	"os"

	"github.com/awgraves/go-lox/expressions"
	"github.com/awgraves/go-lox/runtime"
	"github.com/awgraves/go-lox/tokens"
)

func main() {
	args := os.Args[1:]

	// TMP testing purposes
	astPrinter := expressions.AstPrinter{}

	expr := expressions.Binary{
		Left: expressions.Unary{
			Operator: *tokens.NewToken(tokens.MINUS, "-", 0),
			Right:    expressions.Literal{Value: 123},
		},
		Operator: *tokens.NewToken(tokens.STAR, "*", 0),
		Right: expressions.Grouping{
			Expression: expressions.Literal{Value: 45.67},
		},
	}
	fmt.Println(astPrinter.Print(expr))

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
