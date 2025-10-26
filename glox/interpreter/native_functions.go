package interpreter

import "time"

type clock struct{}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli()), nil
}

func (c *clock) String() string {
	return "<native fn>"
}
