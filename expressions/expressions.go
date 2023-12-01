package expressions

import "github.com/awgraves/go-lox/tokens"

type Expression interface {
	accept(v Visitor) interface{}
}

type Binary struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (e Binary) accept(v Visitor) interface{} {
	return v.visitBinary(e)
}

type Grouping struct {
	Expression Expression
}

func (e Grouping) accept(v Visitor) interface{} {
	return v.visitGrouping(e)
}

type Literal struct {
	Value interface{} // TODO: not sure about this yet
}

func (e Literal) accept(v Visitor) interface{} {
	return v.visitLiteral(e)
}

type Unary struct {
	Operator tokens.Token
	Right    Expression
}

func (e Unary) accept(v Visitor) interface{} {
	return v.visitUnary(e)
}
