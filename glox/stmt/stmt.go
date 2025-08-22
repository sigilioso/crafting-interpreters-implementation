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

func (e Block[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForBlock(e)
}

type Expression[T any] struct {
	Expression expr.Expr[T]
}

func (e Expression[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForExpression(e)
}

type Print[T any] struct {
	Expression expr.Expr[T]
}

func (e Print[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForPrint(e)
}

type Var[T any] struct {
	Name        tokens.Token
	Initializer expr.Expr[T]
}

func (e Var[T]) Accept(v Visitor[T]) (T, error) {
	return v.VisitForVar(e)
}

type Visitor[T any] interface {
	VisitForBlock(Block[T]) (T, error)
	VisitForExpression(Expression[T]) (T, error)
	VisitForPrint(Print[T]) (T, error)
	VisitForVar(Var[T]) (T, error)
}
