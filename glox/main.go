package main

import (
	"bufio"
	"fmt"
	"glox/errors"
	"glox/expr"
	"glox/scanner"
	"glox/tokens"
	"os"
)

func main() {

	// chap05 hack
	expression := expr.Binary[string]{
		Left: expr.Unary[string]{
			Operator: tokens.NewToken(tokens.Minus, "-", tokens.NilLiteral, 1),
			Right:    expr.Literal[string]{Value: 123},
		},
		Operator: tokens.NewToken(tokens.Star, "*", tokens.NilLiteral, 1),
		Right:    expr.Grouping[string]{Expression: expr.Literal[string]{Value: 45.67}},
	}
	p := AstPrinter{}
	fmt.Println(p.Print(expression))
	os.Exit(0)

	switch len(os.Args) {
	case 1:
		runPrompt()
	case 2:
		runFile(os.Args[1])
	default:
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
		os.Exit(64)
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file in %q: %s\n", path, err)
		os.Exit(64)
	}
	run(string(bytes))
	if errors.ErrorFound() {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s", err)
			os.Exit(64)
		}
		run(line)
		errors.ResetError()
	}

}

func run(source string) {
	scanner := scanner.NewScanner(source)
	scanner.ScanTokens()
	scanner.PrintTokens()
}
