package expressions

import "github.com/awgraves/go-lox/tokens"

type Expression interface {
	Accept(v Visitor) (interface{}, error)
}

type Binary struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (e Binary) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinary(e)
}

type Grouping struct {
	Expression Expression
}

func (e Grouping) Accept(v Visitor) (interface{}, error) {
	return v.VisitGrouping(e)
}

type Literal struct {
	Value interface{}
}

func (e Literal) Accept(v Visitor) (interface{}, error) {
	return v.VisitLiteral(e)
}

type Unary struct {
	Operator tokens.Token
	Right    Expression
}

func (e Unary) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnary(e)
}
