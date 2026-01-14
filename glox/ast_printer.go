package main

import (
	"fmt"
	"glox/expr"
	"glox/tokens"
	"strconv"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(e expr.Expr[string]) (string, error) {
	return e.Accept(p)
}

func (p AstPrinter) VisitForThis(t *expr.This[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForSuper(t *expr.Super[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForGet(e *expr.Get[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForSet(e *expr.Set[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForBinary(e *expr.Binary[string]) (string, error) {
	return p.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (p AstPrinter) VisitForGrouping(e *expr.Grouping[string]) (string, error) {
	return p.parenthesize("group", e.Expression)
}

func (p AstPrinter) VisitForLiteral(e *expr.Literal[string]) (string, error) {
	if e.Value == tokens.NilLiteral {
		return "nil", nil
	}
	return format(e.Value), nil
}

func (p AstPrinter) VisitForUnary(e *expr.Unary[string]) (string, error) {
	return p.parenthesize(e.Operator.Lexeme, e.Right)
}

func (p AstPrinter) VisitForVariable(e *expr.Variable[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForAssign(e *expr.Assign[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForLogical(e *expr.Logical[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) VisitForCall(e *expr.Call[string]) (string, error) {
	panic("Not implemented")
}

func (p AstPrinter) parenthesize(name string, exprs ...expr.Expr[string]) (string, error) {
	s := "(" + name
	for _, e := range exprs {
		v, err := e.Accept(p)
		if err != nil {
			return "", err
		}
		s += " " + v
	}
	s += ")"
	return s, nil
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
