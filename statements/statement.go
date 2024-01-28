package statements

import (
	"github.com/awgraves/go-lox/expressions"
	"github.com/awgraves/go-lox/tokens"
)

type Stmt interface {
	Accept(v Visitor) error
}

type ExpStmt struct {
	Expression expressions.Expression
}

func (s ExpStmt) Accept(v Visitor) error {
	return v.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression expressions.Expression
}

func (s PrintStmt) Accept(v Visitor) error {
	return v.VisitPrintStmt(s)
}

type VarStmt struct {
	Name        tokens.Token
	Initializer expressions.Expression
}

func (s VarStmt) Accept(v Visitor) error {
	return v.VisitVarStmt(s)
}

type Block struct {
	Statements []Stmt
}

func (s Block) Accept(v Visitor) error {
	return v.VisitBlock(s)
}
