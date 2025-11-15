// Generated via tools/generate-ast
package expr

import "glox/tokens"

type Expr[T any] interface {
	Accept(Visitor[T]) (T, error)
}

type Assign[T any] struct {
	Name  tokens.Token
	Value Expr[T]
}

func (e *Assign[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForAssign(e)
}

type Binary[T any] struct {
	Left     Expr[T]
	Operator tokens.Token
	Right    Expr[T]
}

func (e *Binary[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForBinary(e)
}

type Call[T any] struct {
	Callee    Expr[T]
	Paren     tokens.Token
	Arguments []Expr[T]
}

func (e *Call[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForCall(e)
}

type Grouping[T any] struct {
	Expression Expr[T]
}

func (e *Grouping[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForGrouping(e)
}

type Literal[T any] struct {
	Value any
}

func (e *Literal[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForLiteral(e)
}

type Unary[T any] struct {
	Operator tokens.Token
	Right    Expr[T]
}

func (e *Unary[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForUnary(e)
}

type Logical[T any] struct {
	Left     Expr[T]
	Operator tokens.Token
	Right    Expr[T]
}

func (e *Logical[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForLogical(e)
}

type Variable[T any] struct {
	Name tokens.Token
}

func (e *Variable[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForVariable(e)
}

type Visitor[T any] interface {
	VisitForAssign(*Assign[T]) (T, error)
	VisitForBinary(*Binary[T]) (T, error)
	VisitForCall(*Call[T]) (T, error)
	VisitForGrouping(*Grouping[T]) (T, error)
	VisitForLiteral(*Literal[T]) (T, error)
	VisitForUnary(*Unary[T]) (T, error)
	VisitForLogical(*Logical[T]) (T, error)
	VisitForVariable(*Variable[T]) (T, error)
}
