package main

import "fmt"

var NilLiteral = NilLiteralType{}

type NilLiteralType struct{}

func (n NilLiteralType) String() string { return "nil" }

type Token struct {
	kind    TokenType
	lexeme  string
	literal any
	line    int
}

func NewToken(kind TokenType, lexeme string, literal any, line int) Token {
	return Token{
		kind:    kind,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s %v", t.kind, t.lexeme, t.literal)
}
