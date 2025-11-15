// Generated via tools/generate-ast
package stmt

import (
	"glox/expr"
	"glox/tokens"
)

type Stmt[T any] interface {
	Accept(Visitor[T]) (T, error)
}

type Block[T any] struct {
	Statements []Stmt[T]
}

func (e *Block[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForBlock(e)
}

type Expression[T any] struct {
	Expression expr.Expr[T]
}

func (e *Expression[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForExpression(e)
}

type Function[T any] struct {
	Name   tokens.Token
	Params []tokens.Token
	Body   []Stmt[T]
}

func (e *Function[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForFunction(e)
}

type If[T any] struct {
	Condition  expr.Expr[T]
	ThenBranch Stmt[T]
	ElseBranch Stmt[T]
}

func (e *If[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForIf(e)
}

type Print[T any] struct {
	Expression expr.Expr[T]
}

func (e *Print[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForPrint(e)
}

type Return[T any] struct {
	Keyword tokens.Token
	Value   expr.Expr[T]
}

func (e *Return[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForReturn(e)
}

type Var[T any] struct {
	Name        tokens.Token
	Initializer expr.Expr[T]
}

func (e *Var[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForVar(e)
}

type While[T any] struct {
	Condition expr.Expr[T]
	Body      Stmt[T]
}

func (e *While[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForWhile(e)
}

type Visitor[T any] interface {
	VisitForBlock(*Block[T]) (T, error)
	VisitForExpression(*Expression[T]) (T, error)
	VisitForFunction(*Function[T]) (T, error)
	VisitForIf(*If[T]) (T, error)
	VisitForPrint(*Print[T]) (T, error)
	VisitForReturn(*Return[T]) (T, error)
	VisitForVar(*Var[T]) (T, error)
	VisitForWhile(*While[T]) (T, error)
}
