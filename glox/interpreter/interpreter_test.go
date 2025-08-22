package interpreter

import (
	"fmt"
	"glox/parser"
	"glox/scanner"
	"glox/tokens"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpreterExpressions(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "(5 - (3 - 1)) + -1",
			expected: "2",
		},
		{
			input:    "\"a\" + \"b\"",
			expected: "ab",
		},
		{
			input:    "4 > 3",
			expected: "true",
		},
		{
			input:    "1 + 10 * 2",
			expected: "21",
		},
		{
			input:    "(1 + 10) * 2",
			expected: "22",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			expression := buildExpression(t, tc.input)

			i := New()
			result, err := i.interpret(expression)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestInterpreterExpressionErrors(t *testing.T) {
	cases := []struct {
		input string
	}{
		{
			input: "\"a\" + 4",
		},
		{
			input: "4 > true",
		},
		{
			input: "1 + 10 * false",
		},
		{
			input: "(1 + \"3\") * 2",
		},
		{
			input: "0 / 0", // Chap07 latest challenge
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			expression := buildExpression(t, tc.input)
			i := New()
			result, err := i.interpret(expression)
			t.Logf("Result: %v", result)
			require.Error(t, err)
		})
	}
}

func buildExpression(t *testing.T, input string) Expr {
	scanner := scanner.NewScanner(input)
	scanner.ScanTokens()
	parser := parser.NewParser[any](scanner.Tokens())
	expression, err := parser.Expression()
	require.NoError(t, err)
	return expression

}

func TestNumOperation(t *testing.T) {
	op := tokens.Token{}
	result, err := numOperation(op, 1.0, 2.0, func(l, r float64) any { return l + r })
	require.NoError(t, err)
	require.Equal(t, 3.0, result)
	result, err = numOperation(op, 1.0, 2.0, func(l, r float64) any { return l > r })
	require.NoError(t, err)
	require.Equal(t, false, result)

	_, err = numOperation(op, "1", "2", func(l, r float64) any { return l + r })
	require.Error(t, err)
	_, err = numOperation(op, 1.0, "2", func(l, r float64) any { return l + r })
	require.Error(t, err)
	_, err = numOperation(op, "1", 2.0, func(l, r float64) any { return l + r })
	require.Error(t, err)
}

func TestIsTruthy(t *testing.T) {
	cases := []struct {
		v        any
		expected bool
	}{
		{
			v:        nil,
			expected: false,
		},
		{
			v:        false,
			expected: false,
		},
		{
			v:        true,
			expected: true,
		},
		{
			v:        "",
			expected: true,
		},
		{
			v:        "something",
			expected: true,
		},
		{
			v:        0.0,
			expected: true,
		},
		{
			v:        42.0,
			expected: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		name := fmt.Sprintf("%v", tc.v)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, isTruthy(tc.v))
		})
	}

}

func TestSumOk(t *testing.T) {
	cases := []struct {
		a, b     any
		expected any
	}{
		{
			a:        3.0,
			b:        5.0,
			expected: 8.0,
		},
		{
			a:        "a",
			b:        "b",
			expected: "ab",
		},
	}

	for _, tc := range cases {
		tc := tc
		name := fmt.Sprintf("%v + %v", tc.a, tc.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			op := tokens.Token{}
			result, err := sum(op, tc.a, tc.b)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}

}

func TestSumError(t *testing.T) {
	cases := []struct {
		a, b any
	}{
		{
			a: 3.0,
			b: "5.0",
		},
		{
			a: true,
			b: false,
		},
		{
			a: "a",
			b: 3.0,
		},
		{
			a: true,
			b: 3.0,
		},
		{
			a: "b",
			b: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		name := fmt.Sprintf("%v + %v", tc.a, tc.b)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			op := tokens.Token{}
			_, err := sum(op, tc.a, tc.b)
			require.Error(t, err)
		})
	}

}

func TestStringify(t *testing.T) {
	cases := []struct {
		v        any
		expected string
	}{
		{
			v:        nil,
			expected: "nil",
		},
		{
			v:        .25,
			expected: "0.25",
		},
		{
			v:        42,
			expected: "42",
		},
		{
			v:        "towel",
			expected: "towel",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.expected, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, stringify(tc.v))
		})
	}

}
