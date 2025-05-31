package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

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
	}

}

func run(source string) {
	// TODO: errors
	scanner := NewScanner(source)
	scanner.scanTokens()
	fmt.Println("Tokens:")
	scanner.PrintTokens()
}
