package resolver

import (
	"glox/expr"
	"glox/tokens"
	"testing"
)

func TestLocals(t *testing.T) {
	locals := map[expr.Expr[any]]int{}

	v1 := expr.Variable[any]{Name: tokens.NewToken(tokens.Identifier, "a", 0, 1)}
	v2 := expr.Variable[any]{Name: tokens.NewToken(tokens.Identifier, "a", 0, 1)}

	store(locals, &v1, 1)
	store(locals, &v2, 2)

	t.Logf("%+v", locals)
	t.Fail()
}

func store(m map[expr.Expr[any]]int, v expr.Expr[any], n int) {
	m[v] = n
}
