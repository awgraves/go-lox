package runtime

import (
	"errors"
	"fmt"

	"github.com/awgraves/go-lox/expressions"
	"github.com/awgraves/go-lox/statements"
	"github.com/awgraves/go-lox/tokens"
)

type resolver struct {
	interpreter interpreter
	scopes      []map[string]bool
	errReporter ErrorReporter
}

func newResolver(i interpreter) *resolver {
	return &resolver{
		interpreter: i,
		errReporter: i.errReporter,
		scopes:      []map[string]bool{},
	}
}

func (r *resolver) resolveStmts(stmts []statements.Stmt) error {
	var err error
	for _, s := range stmts {
		err = r.resolveStmt(s)
		if err != nil {
			r.errReporter.AddError(0, 0, err.Error())
			return err
		}
	}
	return err
}

func (r *resolver) resolveStmt(stmt statements.Stmt) error {
	err := stmt.Accept(r)
	return err
}

func (r *resolver) resolveExprs(exprs []expressions.Expression) error {
	var err error
	for _, e := range exprs {
		err = r.resolveExpr(e)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *resolver) resolveExpr(expr expressions.Expression) error {
	_, err := expr.Accept(r)
	return err
}

func (r *resolver) beginScope() {
	newScope := make(map[string]bool)
	r.scopes = append([]map[string]bool{newScope}, r.scopes...)
	fmt.Println("beginning scope")
	fmt.Printf("total scopes: %d\n", len(r.scopes))
}

func (r *resolver) endScope() {
	r.scopes = r.scopes[1:]
	fmt.Println("ending scope")
	fmt.Printf("total scopes: %d\n", len(r.scopes))
}

func (r *resolver) declare(name tokens.Token) {
	if len(r.scopes) == 0 {
		return
	}
	fmt.Printf("declaring %s\n", name.Lexeme)
	scope := r.scopes[0]
	scope[name.Lexeme] = false
}

func (r *resolver) define(name tokens.Token) {
	if len(r.scopes) == 0 {
		return
	}
	fmt.Printf("defining %s\n", name.Lexeme)
	scope := r.scopes[0]
	scope[name.Lexeme] = true
}

func (r *resolver) resolveLocal(expr expressions.Expression, name tokens.Token) {
	for i := 0; i < len(r.scopes); i++ {
		scope := r.scopes[i]
		if _, ok := scope[name.Lexeme]; ok {
			fmt.Printf("%s found at level: %d\n", name.Lexeme, i)
			r.interpreter.resolve(expr, i)
			return
		}
	}
	fmt.Println("NOT FOUND")
}

func (r *resolver) resolveFunction(fun statements.FunctionStmt) {
	r.beginScope()
	for _, p := range fun.Params {
		r.declare(p)
		r.define(p)
	}
	r.resolveStmts(fun.Body)
	r.endScope()
}

func (r *resolver) VisitExpressionStmt(stmt statements.ExpStmt) error {
	err := r.resolveExpr(stmt.Expression)
	return err
}

func (r *resolver) VisitFunctionStmt(stmt statements.FunctionStmt) error {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt)
	return nil
}

func (r *resolver) VisitPrintStmt(stmt statements.PrintStmt) error {
	err := r.resolveExpr(stmt.Expression)
	return err
}

func (r *resolver) VisitReturnStmt(stmt statements.ReturnStmt) error {
	var err error
	if stmt.Value != nil {
		err = r.resolveExpr(stmt.Value)
	}
	return err
}

func (r *resolver) VisitVarStmt(stmt statements.VarStmt) error {
	fmt.Println("VAR STATEMENT")
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		err := r.resolveExpr(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	r.define(stmt.Name)
	return nil
}

func (r *resolver) VisitBlock(stmt statements.Block) error {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *resolver) VisitIfStmt(stmt statements.IfStmt) error {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		err := r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) VisitWhileStmt(stmt statements.WhileStmt) error {
	fmt.Println("while statement")
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *resolver) VisitBinary(expr expressions.Binary) (interface{}, error) {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(expr.Right)
	return nil, err
}

func (r *resolver) VisitGrouping(expr expressions.Grouping) (interface{}, error) {
	err := r.resolveExpr(expr.Expression)
	return nil, err
}

func (r *resolver) VisitLiteral(expr expressions.Literal) (interface{}, error) {
	return nil, nil
}

func (r *resolver) VisitUnary(expr expressions.Unary) (interface{}, error) {
	err := r.resolveExpr(expr.Right)
	return nil, err
}

func (r *resolver) VisitVariable(expr expressions.Variable) (interface{}, error) {
	if len(r.scopes) > 0 {
		scope := r.scopes[0]
		v, ok := scope[expr.Name.Lexeme]
		if ok && !v {
			// report the error
			err := errors.New("Can't read local variable in its own initializer.")
			return nil, err
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *resolver) VisitAssign(expr expressions.Assign) (interface{}, error) {
	fmt.Printf("assigning %s\n", expr.Name.Lexeme)
	err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *resolver) VisitLogical(expr expressions.Logical) (interface{}, error) {
	fmt.Println("LOGICAL")
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) VisitCall(expr expressions.Call) (interface{}, error) {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}

	for _, arg := range expr.Arguments {
		err = r.resolveExpr(arg)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
