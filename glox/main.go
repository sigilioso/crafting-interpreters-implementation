package main

import (
	"bufio"
	"fmt"
	"glox/errors"
	"glox/interpreter"
	"glox/parser"
	"glox/scanner"
	"os"
)

func main() {

	chap05Hack(false) // switch to true to run chap05 hack only

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
	if errors.RuntimeErrorFound() {
		os.Exit(70)
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
	token_list := scanner.Tokens()
	parser := parser.NewParser[any](token_list)
	statements, _ := parser.Parse() // TODO: check returned error for parser
	if errors.ErrorFound() {
		return
	}
	loxInterpreter := interpreter.New()
	loxInterpreter.Interpret(statements)
}
