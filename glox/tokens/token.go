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
	if t.kind == Eof {
		return "EOF  null"
	}
	if t.kind == Number {
		return fmt.Sprintf("%v %s %s", t.kind, t.lexeme, formatFloat(t.literal))
	}
	return fmt.Sprintf("%v %s %v", t.kind, t.lexeme, t.literal)
}

func formatFloat(v any) string {
	s := strconv.FormatFloat(v.(float64), 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s = s + ".0"
	}
	return s
}
