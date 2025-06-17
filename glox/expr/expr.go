// Generated via tools/generate-ast
package expr

import "glox/tokens"

type Expr[T any] interface {
	Accept(Visitor[T]) T
}

type Binary[T any] struct {
	Left     Expr[T]
	Operator tokens.Token
	Right    Expr[T]
}

func (e Binary[T]) Accept(v Visitor[T]) T {
	return v.VisitForBinary(e)
}

type Grouping[T any] struct {
	Expression Expr[T]
}

func (e Grouping[T]) Accept(v Visitor[T]) T {
	return v.VisitForGrouping(e)
}

type Literal[T any] struct {
	Value any
}

func (e Literal[T]) Accept(v Visitor[T]) T {
	return v.VisitForLiteral(e)
}

type Unary[T any] struct {
	Operator tokens.Token
	Right    Expr[T]
}

func (e Unary[T]) Accept(v Visitor[T]) T {
	return v.VisitForUnary(e)
}

type Visitor[T any] interface {
	VisitForBinary(Binary[T]) T
	VisitForGrouping(Grouping[T]) T
	VisitForLiteral(Literal[T]) T
	VisitForUnary(Unary[T]) T
}
