package runtime

import (
	"errors"
	"fmt"

	"github.com/awgraves/go-lox/expressions"
	"github.com/awgraves/go-lox/statements"
	"github.com/awgraves/go-lox/tokens"
)

type interpreter struct {
	errReporter ErrorReporter
	environment Environment
}

func newIntepreter(errReporter ErrorReporter) *interpreter {
	globals := newEnvironment(nil)
	globals.define(tokens.Token{Lexeme: "clock"}, Clock{})

	return &interpreter{
		errReporter: errReporter,
		environment: globals,
	}
}

func (i *interpreter) interpret(statements []statements.Stmt) {
	for _, s := range statements {
		err := i.execute(s)
		if err != nil {
			//TODO: make this accurate
			i.errReporter.AddError(0, 0, err.Error())
			return
		}
	}
}

func (i *interpreter) execute(stmt statements.Stmt) error {
	return stmt.Accept(i)
}

func stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}

	num, ok := v.(float64)
	if ok {
		return fmt.Sprintf("%v", num)
	}
	return fmt.Sprintf("%v", v)
}

func (i *interpreter) VisitVarStmt(stmt statements.VarStmt) error {
	if stmt.Initializer != nil {
		value, err := i.evaluate(stmt.Initializer)
		if err != nil {
			return err
		}
		i.environment.define(stmt.Name, value)
		return nil
	}
	i.environment.define(stmt.Name, nil)
	return nil
}

func (i *interpreter) VisitWhileStmt(stmt statements.WhileStmt) error {
	res, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}
	for i.isTruthy(res) {
		err = i.execute(stmt.Body)
		if err != nil {
			return err
		}
		res, err = i.evaluate(stmt.Condition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) VisitAssign(expr expressions.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	i.environment.assign(expr.Name, value)
	return value, nil
}

func (i *interpreter) VisitVariable(expr expressions.Variable) (interface{}, error) {
	return i.environment.get(expr.Name)
}

func (i *interpreter) VisitBlock(stmt statements.Block) error {
	return i.executeBlock(stmt.Statements, newEnvironment(i.environment))
}

func (i *interpreter) executeBlock(statements []statements.Stmt, environment Environment) error {
	previous := i.environment
	i.environment = environment
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	i.environment = previous
	return nil
}

func (i *interpreter) VisitExpressionStmt(stmt statements.ExpStmt) error {
	_, err := i.evaluate(stmt.Expression)
	return err
}

func (i *interpreter) VisitIfStmt(stmt statements.IfStmt) error {
	res, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}
	if i.isTruthy(res) {
		err = i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		err = i.execute(stmt.ElseBranch)
	}
	return err
}

func (i *interpreter) VisitPrintStmt(stmt statements.PrintStmt) error {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return err
	}
	fmt.Println(stringify(value))
	return nil
}

func (i *interpreter) VisitLiteral(exp expressions.Literal) (interface{}, error) {
	return exp.Value, nil
}

func (i *interpreter) VisitLogical(exp expressions.Logical) (interface{}, error) {
	left, err := i.evaluate(exp.Left)
	if err != nil {
		return nil, err
	}

	if exp.Operator.TokenType == tokens.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(exp.Right)
}

func (i *interpreter) VisitGrouping(exp expressions.Grouping) (interface{}, error) {
	return i.evaluate(exp.Expression)
}

func (i *interpreter) evaluate(exp expressions.Expression) (interface{}, error) {
	return exp.Accept(i)
}

func castToFloat(i interface{}) (float64, error) {
	num, ok := i.(float64)
	if !ok {
		return 0, errors.New("Operand must be a number")
	}
	return num, nil
}

func (i *interpreter) VisitUnary(exp expressions.Unary) (interface{}, error) {
	right, err := i.evaluate(exp.Right)
	if err != nil {
		return nil, err
	}

	switch exp.Operator.TokenType {
	case tokens.BANG:
		return !i.isTruthy(right), nil
	case tokens.MINUS:
		num, err := castToFloat(right)
		if err != nil {
			return nil, err
		}

		return -num, nil
	}

	return nil, errors.New("TODO")
}

func (i *interpreter) isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}

	if b, ok := val.(bool); ok {
		return b
	}

	return true
}

func castToFloats(a, b interface{}) (float64, float64, error) {

	aFloat, aok := a.(float64)
	bFloat, bok := b.(float64)

	if aok && bok {
		return aFloat, bFloat, nil
	}

	return 0, 0, errors.New("not a number")
}

func (i *interpreter) VisitBinary(exp expressions.Binary) (interface{}, error) {
	left, err := i.evaluate(exp.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(exp.Right)
	if err != nil {
		return nil, err
	}

	switch exp.Operator.TokenType {
	case tokens.MINUS:
		left, right, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return left - right, nil
	case tokens.SLASH:
		left, right, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return left / right, nil
	case tokens.STAR:
		left, right, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return left * right, nil
	case tokens.PLUS:
		numLeft, numRight, err := castToFloats(left, right)
		if err == nil {
			return numLeft + numRight, nil
		}

		strLeft, lok := left.(string)
		strRight, rok := right.(string)

		if lok && rok {
			return strLeft + strRight, nil
		}
		// TODO: maybe define the types in msg?
		return nil, errors.New("operands must be two numbers or two strings")

	case tokens.GREATER:
		numLeft, numRight, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return numLeft > numRight, nil
	case tokens.GREATER_EQUAL:
		numLeft, numRight, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return numLeft >= numRight, nil
	case tokens.LESS:
		numLeft, numRight, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return numLeft < numRight, nil
	case tokens.LESS_EQUAL:
		numLeft, numRight, err := castToFloats(left, right)
		if err != nil {
			return nil, err
		}
		return numLeft <= numRight, nil
	case tokens.BANG_EQUAL:
		return !isEqual(left, right), nil
	case tokens.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}
	return nil, errors.New("TODO")
}

func (i *interpreter) VisitCall(expr expressions.Call) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := []interface{}{}
	for _, arg := range expr.Arguments {
		eval, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, eval)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, errors.New("Can only call functions and classes.")
	}

	arity := function.Arity()
	got := len(arguments)
	if got != arity {
		return nil, fmt.Errorf("Expected %d arguments but got %d", arity, got)
	}

	return function.Call(i, arguments)
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
