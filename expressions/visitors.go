package expressions

import (
	"fmt"
	"strings"
)

type Visitor interface {
	visitBinary(Binary) interface{}
	visitGrouping(Grouping) interface{}
	visitLiteral(Literal) interface{}
	visitUnary(Unary) interface{}
}

// for testing purposes
// implements Visitor
type AstPrinter struct {
}

func (a *AstPrinter) Print(expr Expression) interface{} {
	return expr.accept(a)
}

func (a *AstPrinter) visitBinary(expr Binary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) visitGrouping(expr Grouping) interface{} {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) visitLiteral(expr Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	switch expr.Value.(type) {
	case string:
		return expr.Value
	case int, float64:
		return fmt.Sprintf("%v", expr.Value)
	default:
		panic("Unknown literal")
	}
}

func (a *AstPrinter) visitUnary(expr Unary) interface{} {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expression) string {
	builder := strings.Builder{}

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.accept(a).(string))
	}
	builder.WriteString(")")
	return builder.String()
}
