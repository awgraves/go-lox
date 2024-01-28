package runtime

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awgraves/go-lox/expressions"
	"github.com/awgraves/go-lox/statements"
	"github.com/awgraves/go-lox/tokens"
)

type interpreter struct {
	errReporter ErrorReporter
}

func newIntepreter(errReporter ErrorReporter) *interpreter {
	return &interpreter{
		errReporter: errReporter,
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
		return strings.TrimRight(fmt.Sprintf("%v", num), ".0")
	}
	return fmt.Sprintf("%v", v)
}

func (i *interpreter) VisitExpressionStmt(stmt statements.ExpStmt) error {
	_, err := i.evaluate(stmt.Expression)
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

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
