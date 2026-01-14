package interpreter

import (
	"fmt"
	"glox/errors"
	"glox/tokens"
)

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]*LoxFunction
}

func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) Arity() int {
	if initializer := c.FindMethod("init"); initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewInstance(c)
	if initializer := c.FindMethod("init"); initializer != nil {
		if _, err := initializer.Bind(instance).Call(interpreter, arguments); err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (c *LoxClass) FindMethod(key string) *LoxFunction {
	method := c.Methods[key]
	if method != nil {
		return method
	}
	if c.Superclass != nil {
		return c.Superclass.FindMethod(key)
	}
	return nil
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]any
}

func NewInstance(c *LoxClass) *LoxInstance {
	return &LoxInstance{class: c, fields: map[string]any{}}
}

func (i *LoxInstance) String() string {
	return i.class.Name + " instance"
}

func (i *LoxInstance) Get(name tokens.Token) (any, error) {
	if field, exists := i.fields[name.Lexeme]; exists {
		return field, nil
	}
	if method := i.class.FindMethod(name.Lexeme); method != nil {
		return method.Bind(i), nil
	}
	return nil, errors.NewRuntimeError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme))
}

func (i *LoxInstance) Set(name tokens.Token, value any) {
	i.fields[name.Lexeme] = value
}
