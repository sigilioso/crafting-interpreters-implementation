package interpreter

import (
	"errors"
	"glox/expr"
	"glox/tokens"
)

// TODO: type checks and handle errors. We need to refactor Visitor to also
// return an error =(

var ErrRuntime = errors.New("runtime error")

type Expr = expr.Expr[any]
type Binary = expr.Binary[any]
type Literal = expr.Literal[any]
type Unary = expr.Unary[any]
type Grouping = expr.Grouping[any]
type Visitor = expr.Visitor[any]

type Interpreter struct{}

func (i Interpreter) VisitForGrouping(grouping Grouping) any {
	return i.evaluate(grouping.Expression)
}

func (i Interpreter) VisitForLiteral(literal Literal) any {
	return literal.Value
}

func (i Interpreter) VisitForBinary(binary Binary) any {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)

	switch binary.Operator.TokenType {
	case tokens.Plus:
		return i.sum(left, right)
	case tokens.Minus:
		return left.(float64) - right.(float64)
	case tokens.Slash:
		return left.(float64) / right.(float64)
	case tokens.Star:
		return left.(float64) * right.(float64)

	case tokens.Greater:
		return left.(float64) > right.(float64)
	case tokens.GreaterEqual:
		return left.(float64) >= right.(float64)
	case tokens.Less:
		return left.(float64) < right.(float64)
	case tokens.LessEqual:
		return left.(float64) <= right.(float64)

	case tokens.Equal:
		return left == right
	case tokens.BangEqual:
		return left != right

	}
	return nil // unreachable
}
func (i Interpreter) VisitForUnary(unary Unary) any {
	right := i.evaluate(unary.Right)

	switch unary.Operator.TokenType {
	case tokens.Minus:
		return -(right.(float64))
	case tokens.Bang:
		return !i.isTruthy(right)
	}
	return nil // unreachable
}

func (i Interpreter) evaluate(expression Expr) any {
	return expression.Accept(i)
}

// isTruthy consider anything but nil or false value as true
func (i Interpreter) isTruthy(v any) bool {
	if v == nil {
		return false
	}
	if value, isBool := v.(bool); isBool {
		return value
	}
	return true
}

func (i Interpreter) sum(left any, right any) any {
	lNum, lIsNumber := left.(float64)
	rNum, rIsNumber := right.(float64)
	if lIsNumber && rIsNumber {
		return lNum + rNum
	}
	lStr, lIsStr := left.(string)
	rStr, rIsStr := right.(string)
	if lIsStr && rIsStr {
		return lStr + rStr
	}

	return nil
}
