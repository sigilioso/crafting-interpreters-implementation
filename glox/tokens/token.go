package tokens

import (
	"fmt"
	"strconv"
	"strings"
)

var NilLiteral = NilLiteralType{}

type NilLiteralType struct{}

func (n NilLiteralType) String() string { return "null" }

type Token struct {
	Kind    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(kind TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Kind:    kind,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	if t.Kind == Eof {
		return "EOF  null"
	}
	if t.Kind == Number {
		return fmt.Sprintf("%v %s %s", t.Kind, t.Lexeme, formatFloat(t.Literal))
	}
	return fmt.Sprintf("%v %s %v", t.Kind, t.Lexeme, t.Literal)
}

func formatFloat(v any) string {
	s := strconv.FormatFloat(v.(float64), 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s = s + ".0"
	}
	return s
}
