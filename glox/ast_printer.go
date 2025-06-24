package main

import (
	"fmt"
	"glox/expr"
	"glox/tokens"
	"strconv"
	"strings"
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
	return format(e.Value)
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

// format performs ugly formatting to match Java implementation used in book's tests
func format(v any) string {
	if i, ok := v.(int64); ok {
		return fmt.Sprintf("%d.0", i)
	}
	if f, ok := v.(float64); ok {
		s := strconv.FormatFloat(f, 'f', -1, 64)
		if !strings.Contains(s, ".") {
			s = s + ".0"
		}
		return s
	}
	return fmt.Sprintf("%v", v)
}
