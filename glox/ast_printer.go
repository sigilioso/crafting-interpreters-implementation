package main

import (
	"fmt"
	"glox/expr"
	"glox/tokens"
)

type AstPrinter struct{}

func (p AstPrinter) Print(e expr.Expr[string]) string {
	return e.Accept(p)
}

func (p AstPrinter) VisitForBinary(e expr.Binary[string]) string {
	return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (p AstPrinter) VisitForGrouping(e expr.Grouping[string]) string {
	return p.parenthesize("group", e.Expression)
}

func (p AstPrinter) VisitForLiteral(e expr.Literal[string]) string {
	// TODO: check nil literals
	if e.Value == tokens.NilLiteral {
		return "nil"
	}
	return fmt.Sprintf("%v", e.Value)
}

func (p AstPrinter) VisitForUnary(e expr.Unary[string]) string {
	return p.parenthesize(e.Operator.Lexeme, e.Right)
}

func (p AstPrinter) parenthesize(name string, exprs ...expr.Expr[string]) string {
	s := "(" + name
	for _, e := range exprs {
		s += " " + e.Accept(p)
	}
	s += ")"
	return s
}
