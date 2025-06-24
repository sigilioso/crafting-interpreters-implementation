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
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
	if t.TokenType == Eof {
		return "EOF  null"
	}
	if t.TokenType == Number {
		return fmt.Sprintf("%v %s %s", t.TokenType, t.Lexeme, FormatFloat(t.Literal))
	}
	return fmt.Sprintf("%v %s %v", t.TokenType, t.Lexeme, t.Literal)
}

func FormatFloat(v any) string {
	s := strconv.FormatFloat(v.(float64), 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s = s + ".0"
	}
	return s
}
