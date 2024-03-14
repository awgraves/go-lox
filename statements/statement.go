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

type FunctionStmt struct {
	Name   tokens.Token
	Params []tokens.Token
	Body   []Stmt
}

func (s FunctionStmt) Accept(v Visitor) error {
	return v.VisitFunctionStmt(s)
}

type PrintStmt struct {
	Expression expressions.Expression
}

func (s PrintStmt) Accept(v Visitor) error {
	return v.VisitPrintStmt(s)
}

type ReturnStmt struct {
	Keyword tokens.Token
	Value   expressions.Expression // might be nil!
}

func (s ReturnStmt) Accept(v Visitor) error {
	return v.VisitReturnStmt(s)
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

type IfStmt struct {
	Condition  expressions.Expression
	ThenBranch Stmt
	ElseBranch Stmt // possibly nil
}

func (s IfStmt) Accept(v Visitor) error {
	return v.VisitIfStmt(s)
}

type WhileStmt struct {
	Condition expressions.Expression
	Body      Stmt
}

func (s WhileStmt) Accept(v Visitor) error {
	return v.VisitWhileStmt(s)
}
