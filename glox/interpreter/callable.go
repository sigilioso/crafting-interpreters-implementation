package interpreter

type GloxCallable interface {
	// Arity determines the number of expected arguments
	Arity() int
	// Call performs a call using the interpreter
	Call(interpreter *Interpreter, arguments []any) (any, error)
}
