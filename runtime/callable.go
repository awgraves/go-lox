package runtime

import "time"

type LoxCallable interface {
	Arity() int
	Call(interp *interpreter, args []interface{}) (interface{}, error)
}

type Clock struct{}

func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(interp *interpreter, args []interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}
