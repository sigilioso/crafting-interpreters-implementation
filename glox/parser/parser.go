package parser

import (
	"errors"
	gloxErrors "glox/errors"
	"glox/expr"
	"glox/tokens"
)

var ErrParse = errors.New("parse Error")

type Parser[T any] struct {
	tokens  []tokens.Token
	current int
}

func NewParser[T any](token_list []tokens.Token) Parser[T] {
	return Parser[T]{
		tokens:  token_list,
		current: 0,
	}
}

func (p *Parser[T]) Parse() expr.Expr[T] {
	expression, err := p.expression()
	if err != nil {
		return nil
	}
	return expression
}

func (p *Parser[T]) expression() (expr.Expr[T], error) {
	return p.equality()
}

func (p *Parser[T]) equality() (expr.Expr[T], error) {
	expression, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary[T]{Left: expression, Operator: operator, Right: right}
	}
	return expression, nil
}

func (p *Parser[T]) comparison() (expr.Expr[T], error) {
	expression, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary[T]{Left: expression, Operator: operator, Right: right}
	}
	return expression, nil
}

func (p *Parser[T]) term() (expr.Expr[T], error) {
	expression, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(tokens.Minus, tokens.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary[T]{Left: expression, Operator: operator, Right: right}
	}
	return expression, nil
}

func (p *Parser[T]) factor() (expr.Expr[T], error) {
	expression, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(tokens.Slash, tokens.Star) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expression = expr.Binary[T]{Left: expression, Operator: operator, Right: right}
	}
	return expression, nil
}

func (p *Parser[T]) unary() (expr.Expr[T], error) {
	if p.match(tokens.Bang, tokens.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expr.Unary[T]{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser[T]) primary() (expr.Expr[T], error) {
	switch {
	case p.match(tokens.False):
		return expr.Literal[T]{Value: false}, nil
	case p.match(tokens.True):
		return expr.Literal[T]{Value: true}, nil
	case p.match(tokens.Nil):

		return expr.Literal[T]{Value: tokens.NilLiteral}, nil
	case p.match(tokens.Number, tokens.String):
		return expr.Literal[T]{Value: p.previous().Literal}, nil
	case p.match(tokens.LeftParen):
		expression, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokens.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return expr.Grouping[T]{Expression: expression}, nil
	}

	return nil, parseError(p.peek(), "Expect expression.")
}

func (p *Parser[T]) match(tokenTypes ...tokens.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser[T]) check(tokenType tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser[T]) advance() tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser[T]) isAtEnd() bool {
	return p.peek().TokenType == tokens.Eof
}

func (p *Parser[T]) peek() tokens.Token {
	return p.tokens[p.current]
}

func (p *Parser[T]) previous() tokens.Token {
	return p.tokens[p.current-1]
}

func (p *Parser[T]) consume(tokenType tokens.TokenType, message string) (tokens.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	token := p.peek()
	return token, parseError(token, message)
}

func (p *Parser[T]) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == tokens.Semicolon {
			return
		}
		switch p.peek().TokenType {
		case tokens.Class, tokens.Fun, tokens.Var, tokens.For, tokens.If, tokens.While, tokens.Print, tokens.Return:
			return
		}
		p.advance()
	}
}

func parseError(token tokens.Token, message string) error {
	gloxErrors.AtToken(token, message)
	return ErrParse
}
