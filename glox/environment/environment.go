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
		return nil, errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
	}
	return value, nil
}

func (e *Environment) Assign(name tokens.Token, value any) error {
	if _, defined := e.values[name.Lexeme]; defined {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) AssignAt(distance int, name tokens.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	result := e
	for i := 0; i < distance; i++ {
		result = result.enclosing
	}
	return result
}

func (e *Environment) Print() {
	fmt.Printf("Values: %v\n", e.values)
	if e.enclosing != nil {
		fmt.Println("-->")
		e.enclosing.Print()
	}
}
