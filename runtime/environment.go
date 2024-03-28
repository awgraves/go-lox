package runtime

import (
	"errors"
	"fmt"

	"github.com/awgraves/go-lox/tokens"
)

type Environment interface {
	define(name string, value interface{})
	get(name tokens.Token) (interface{}, error)
	getAt(distance int, name tokens.Token) (interface{}, error)
	assign(name tokens.Token, value interface{}) error
	assignAt(distance int, name tokens.Token, value interface{}) error
	getEnclosing() Environment
}

type environment struct {
	enclosing   Environment
	values      map[string]interface{}
	errReporter ErrorReporter
}

func newEnvironment(enclosing Environment) *environment {
	return &environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *environment) getEnclosing() Environment {
	return e.enclosing
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) get(name tokens.Token) (interface{}, error) {
	val, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		err := errors.New(fmt.Sprintf("Undefined variable '%s' when getting.", name))
		return nil, err
	}
	return val, nil
}

func (e *environment) getAt(distance int, name tokens.Token) (interface{}, error) {
	env := e.ancestor(distance)
	return env.get(name)
}

func (e *environment) ancestor(distance int) Environment {
	var env Environment = e
	for i := 0; i < distance; i++ {
		env = env.getEnclosing()
	}

	return env
}

func (e *environment) assign(name tokens.Token, value interface{}) error {
	_, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			e.enclosing.assign(name, value)
		}
		err := errors.New(fmt.Sprintf("Undefined variable '%s' when assigning.", name))
		return err
	}

	e.values[name.Lexeme] = value
	return nil
}

func (e *environment) assignAt(distance int, name tokens.Token, value interface{}) error {
	env := e.ancestor(distance)
	return env.assign(name, value)
}
