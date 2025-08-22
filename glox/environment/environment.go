package environment

import (
	"fmt"
	"glox/errors"
	"glox/tokens"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func New(enclosing *Environment) *Environment {
	return &Environment{values: map[string]any{}, enclosing: enclosing}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name tokens.Token) (any, error) {
	value, found := e.values[name.Lexeme]
	if !found {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable %q.", name.Lexeme))
	}
	return value, nil
}

func (e *Environment) Assign(name tokens.Token, value any) error {
	if _, defined := e.values[name.Lexeme]; defined {
		e.values[name.Lexeme] = value
		return nil
	}
	return errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable %q.", name.Lexeme))
}
