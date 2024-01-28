package statements

import "github.com/awgraves/go-lox/expressions"

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
