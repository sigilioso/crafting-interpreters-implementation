package interpreter

import (
	"fmt"
	"glox/environment"
)

type LoxFunction struct {
	Declaration FunctionStmt
	Closure     *environment.Environment
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	env := environment.New(f.Closure)
	for i, arg := range arguments {
		env.Define(f.Declaration.Params[i].Lexeme, arg)
	}
	err := interpreter.executeBlock(f.Declaration.Body, env)
	if err != nil {
		if returnHolder, isReturn := err.(*Return); isReturn {
			return returnHolder.Value, nil
		}
		return nil, err
	}
	return nil, nil
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

// Return type represents a Golang error that holds the return value.
// This is needed because returns are handled as error although it is merely for control-flow.
type Return struct {
	Value any
}

func (e *Return) Error() string {
	return "not-really-an-error"
}
