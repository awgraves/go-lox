package expressions

import (
	"fmt"
	"strings"
)

type Visitor interface {
	VisitBinary(Binary) (interface{}, error)
	VisitGrouping(Grouping) (interface{}, error)
	VisitLiteral(Literal) (interface{}, error)
	VisitUnary(Unary) (interface{}, error)
}

// for testing purposes
// implements Visitor
type AstPrinter struct {
}

func (a *AstPrinter) Print(expr Expression) (interface{}, error) {
	return expr.Accept(a)
}

func (a *AstPrinter) VisitBinary(expr Binary) (interface{}, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a *AstPrinter) VisitGrouping(expr Grouping) (interface{}, error) {
	return a.parenthesize("group", expr.Expression), nil
}

func (a *AstPrinter) VisitLiteral(expr Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	switch expr.Value.(type) {
	case string:
		return expr.Value, nil
	case int, float64:
		return fmt.Sprintf("%v", expr.Value), nil
	default:
		panic("Unknown literal")
	}
}

func (a *AstPrinter) VisitUnary(expr Unary) (interface{}, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expression) string {
	builder := strings.Builder{}

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		val, _ := expr.Accept(a)
		builder.WriteString(val.(string))
	}
	builder.WriteString(")")
	return builder.String()
}
