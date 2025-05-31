package main

import (
	"fmt"
	"glox/errors"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

type Scanner struct {
	tokens []Token
	source string

	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, line: 1}
}

func (s *Scanner) scanTokens() {

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(Eof, "", nil, s.line))
}

func (s *Scanner) scanToken() {
	r := s.advance()
	switch r {
	case '(':
		s.addNilToken(LeftParen)
	case ')':
		s.addNilToken(RightParen)
	case '{':
		s.addNilToken(LeftBrace)
	case '}':
		s.addNilToken(RightBrace)
	case ',':
		s.addNilToken(Comma)
	case '.':
		s.addNilToken(Dot)
	case ';':
		s.addNilToken(Semicolon)
	case '*':
		s.addNilToken(Star)

	case '!':
		if s.advanceIfMatches('=') {
			s.addNilToken(BangEqual)
		} else {
			s.addNilToken(Bang)
		}
	case '=':
		if s.advanceIfMatches('=') {
			s.addNilToken(EqualEqual)
		} else {
			s.addNilToken(Equal)
		}
	case '<':
		if s.advanceIfMatches('=') {
			s.addNilToken(LessEqual)
		} else {
			s.addNilToken(Less)
		}
	case '>':
		if s.advanceIfMatches('=') {
			s.addNilToken(GreaterEqual)
		} else {
			s.addNilToken(Greater)
		}
	case '/':
		if s.advanceIfMatches('/') { // a comment "// .." goes until the end of line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addNilToken(Slash)
		}
	case ' ', '\r', '\t': // Ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.handleString()

	default:
		if unicode.IsDigit(r) {
			s.handleNumber()
		} else if unicode.IsLetter(r) || r == '_' {
			s.handleIdentifier()
		} else {
			errors.Error(s.line, "Unexpected character.")
		}
	}

}

func (s *Scanner) PrintTokens() {
	for _, token := range s.tokens {
		fmt.Printf("%s\n", token)
	}
}

// peek gets current character without advancing
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	currentRune, _ := s.currentRune()
	return currentRune
}

// peekNext gets the char after
func (s *Scanner) peekNext() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	_, width := s.currentRune()
	if s.current+width >= len(s.source) {
		return rune(0)
	}
	nextRune, _ := utf8.DecodeRuneInString(s.source[s.current+width:])
	return nextRune
}

func (s *Scanner) currentRune() (rune, int) {
	return utf8.DecodeRuneInString(s.source[s.current:])
}

func (s *Scanner) advance() rune {
	currentRune, width := s.currentRune()
	s.current += width
	return currentRune
}

func (s *Scanner) addNilToken(tokenType TokenType) {
	s.addToken(tokenType, NilLiteral)
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) advanceIfMatches(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	currentRune, width := s.currentRune()
	if currentRune != expected {
		return false
	}

	s.current += width
	return true
}

func (s *Scanner) handleString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		errors.Error(s.line, "Unterminated string.")
		return
	}

	s.advance() // The closing "

	value := s.source[s.start+1 : s.current-1]
	s.addToken(String, value)
}

func (s *Scanner) handleNumber() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}
	// check for fractional part
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()                     // the "."
		for unicode.IsDigit(s.peek()) { // The digits after
			s.advance()
		}
	}
	n, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		panic("unreachable, we have just checked it is number")
	}
	s.addToken(Number, n)
}

func (s *Scanner) handleIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, fould := keywords[text]
	if !fould {
		tokenType = Identifier
	}
	s.addNilToken(tokenType)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || unicode.IsDigit(r)
}
