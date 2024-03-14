package runtime

import (
	"fmt"
	"time"

	"github.com/awgraves/go-lox/statements"
)

type LoxCallable interface {
	Arity() int
	Call(interp *interpreter, args []interface{}) (interface{}, error)
	String() string
}

type Clock struct{}

func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(interp *interpreter, args []interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}

func (c Clock) String() string {
	return "<native fn>"
}

type LoxFunction struct {
	Closure     Environment
	Declaration statements.FunctionStmt
}

func (l LoxFunction) Arity() int {
	return len(l.Declaration.Params)
}

func (l LoxFunction) Call(interp *interpreter, args []interface{}) (interface{}, error) {
	env := newEnvironment(l.Closure)

	for i := 0; i < len(l.Declaration.Params); i++ {
		env.define(l.Declaration.Params[i].Lexeme, args[i])
	}

	err := interp.executeBlock(l.Declaration.Body, env)

	val, ok := err.(*ReturnValue)
	if ok {
		return val.Value, nil
	}

	return nil, err
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", l.Declaration.Name.Lexeme)
}
