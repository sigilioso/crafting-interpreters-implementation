package main

import (
	"fmt"
	"glox/expr"
	"glox/tokens"
	"os"
)

// chap05Hack executes the example for chapter 5
func chap05Hack(chap05 bool) {
	if chap05 {
		expression := expr.Binary[string]{
			Left: &expr.Unary[string]{
				Operator: tokens.NewToken(tokens.Minus, "-", tokens.NilLiteral, 1),
				Right:    &expr.Literal[string]{Value: 123},
			},
			Operator: tokens.NewToken(tokens.Star, "*", tokens.NilLiteral, 1),
			Right:    &expr.Grouping[string]{Expression: &expr.Literal[string]{Value: 45.67}},
		}
		p := AstPrinter{}
		fmt.Println(p.Print(&expression))
		os.Exit(0)
	}
}
