// Generated via tools/generate-ast
package stmt

import "glox/expr"

type Stmt[T any] interface {
	Accept(Visitor[T]) (T, error)
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

type Visitor[T any] interface {
	VisitForExpression(Expression[T]) (T, error)
	VisitForPrint(Print[T]) (T, error)
}
