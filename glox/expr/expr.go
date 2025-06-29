// Generated via tools/generate-ast
package expr

import "glox/tokens"

type Expr[T any] interface {
	Accept(Visitor[T]) (T, error)
}

type Binary[T any] struct {
	Left     Expr[T]
	Operator tokens.Token
	Right    Expr[T]
}

func (e Binary[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForBinary(e)
}

type Grouping[T any] struct {
	Expression Expr[T]
}

func (e Grouping[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForGrouping(e)
}

type Literal[T any] struct {
	Value any
}

func (e Literal[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForLiteral(e)
}

type Unary[T any] struct {
	Operator tokens.Token
	Right    Expr[T]
}

func (e Unary[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForUnary(e)
}

type Visitor[T any] interface {
	VisitForBinary(Binary[T]) (T, error)
	VisitForGrouping(Grouping[T]) (T, error)
	VisitForLiteral(Literal[T]) (T, error)
	VisitForUnary(Unary[T]) (T, error)
}
