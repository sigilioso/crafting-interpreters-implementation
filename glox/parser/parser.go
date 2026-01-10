package parser

import (
	"errors"
	"fmt"
	gloxErrors "glox/errors"
	"glox/expr"
	"glox/stmt"
	"glox/tokens"
)

const ARGUMENTS_LIMIT = 255

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

func (p *Parser[T]) Parse() ([]stmt.Stmt[T], error) {
	statements := []stmt.Stmt[T]{}
	for !p.isAtEnd() {
		statement := p.declaration()
		statements = append(statements, statement)
	}
	return statements, nil
}

func (p *Parser[T]) declaration() stmt.Stmt[T] {
	// Get a regular statement if no other declaration matches
	statementGetter := p.statement

	if p.match(tokens.Class) {
		stmt, err := p.classDeclaration()
		if err != nil {
			gloxErrors.AtToken(p.previous(), fmt.Sprintf("%s", err))
			return nil
		}
		return stmt
	}

	if p.match(tokens.Fun) {
		statementGetter = func() (stmt.Stmt[T], error) {
			return p.function("function")
		}
	} else if p.match(tokens.Var) {
		statementGetter = p.varDeclaration
	}

	statement, err := statementGetter()
	if err != nil {
		// Synchronize if we found any parsing error
		p.synchronize()
		return nil
	}
	return statement
}

func (p *Parser[T]) classDeclaration() (stmt.Stmt[T], error) {
	name, err := p.consume(tokens.Identifier, "Expect class name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.LeftBrace, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}
	methods := []*stmt.Function[T]{}
	for !p.check(tokens.RightBrace) && !p.isAtEnd() {
		f, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, f)
	}
	_, err = p.consume(tokens.RightBrace, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}
	return &stmt.Class[T]{Name: name, Methods: methods}, nil

}

func (p *Parser[T]) varDeclaration() (stmt.Stmt[T], error) {
	name, err := p.consume(tokens.Identifier, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer expr.Expr[T]
	if p.match(tokens.Equal) {
		initializer, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(tokens.Semicolon, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return &stmt.Var[T]{Name: name, Initializer: initializer}, nil
}

func (p *Parser[T]) statement() (stmt.Stmt[T], error) {
	if p.match(tokens.If) {
		return p.ifStatement()
	}
	if p.match(tokens.Return) {
		return p.returnStatement()
	}
	if p.match(tokens.Print) {
		return p.printStatement()
	}
	if p.match(tokens.While) {
		return p.whileStatemet()
	}
	if p.match(tokens.For) {
		return p.forStatement()
	}
	if p.match(tokens.LeftBrace) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &stmt.Block[T]{Statements: statements}, nil
	}
	return p.expressionStatement()
}

func (p *Parser[T]) function(functionType string) (f *stmt.Function[T], err error) {
	// function name
	name, err := p.consume(tokens.Identifier, fmt.Sprintf("Expect %s name.", functionType))
	if err != nil {
		return f, err
	}
	// (
	_, err = p.consume(tokens.LeftParen, fmt.Sprintf("Expect '(' after %s name.", functionType))
	if err != nil {
		return f, err
	}
	// parameters
	parameters := []tokens.Token{}
	if !p.check(tokens.RightParen) {
		param, err := p.consume(tokens.Identifier, "Expect parameter name.")
		if err != nil {
			return f, err
		}
		parameters = append(parameters, param)
		for p.match(tokens.Comma) {
			if len(parameters) >= 255 {
				return f, parseError(p.peek(), "Can't have more than 255 parameters.")
			}
			param, err := p.consume(tokens.Identifier, "Expect parameter name.")
			if err != nil {
				return f, err
			}
			parameters = append(parameters, param)
		}
	}
	// )
	_, err = p.consume(tokens.RightParen, "Expect ')' after parameters.")
	if err != nil {
		return f, err
	}
	// {
	_, err = p.consume(tokens.LeftBrace, fmt.Sprintf("Expect '{' before %s body.", functionType))
	if err != nil {
		return f, err
	}
	// function body
	body, err := p.block()
	if err != nil {
		return f, err
	}
	return &stmt.Function[T]{Name: name, Params: parameters, Body: body}, nil
}

func (p *Parser[T]) returnStatement() (stmt.Stmt[T], error) {
	keyword := p.previous()
	var value expr.Expr[T]
	if !p.check(tokens.Semicolon) {
		v, err := p.Expression()
		if err != nil {
			return nil, err
		}
		value = v
	}
	_, err := p.consume(tokens.Semicolon, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}
	return &stmt.Return[T]{Keyword: keyword, Value: value}, nil
}

func (p *Parser[T]) ifStatement() (stmt.Stmt[T], error) {
	_, err := p.consume(tokens.LeftParen, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.RightParen, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch stmt.Stmt[T]
	if p.match(tokens.Else) {
		elseB, err := p.statement()
		elseBranch = elseB
		if err != nil {
			return nil, err
		}
	}
	return &stmt.If[T]{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser[T]) whileStatemet() (stmt.Stmt[T], error) {
	_, err := p.consume(tokens.LeftParen, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.RightParen, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &stmt.While[T]{Condition: condition, Body: body}, nil
}

// forStatement syntactic sugar to support for syntax
func (p *Parser[T]) forStatement() (stmt.Stmt[T], error) {
	// for(var i = 0; i < 10; i++)
	_, err := p.consume(tokens.LeftParen, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer stmt.Stmt[T]
	if p.match(tokens.Semicolon) {
		initializer = nil
	} else if p.match(tokens.Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition expr.Expr[T]
	if !p.check(tokens.Semicolon) {
		condition, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var increment expr.Expr[T]
	if !p.check(tokens.RightParen) {
		increment, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(tokens.RightParen, "Expect ')' after loop condition")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// add the initializer statement if any
	if increment != nil {
		statements := []stmt.Stmt[T]{body, &stmt.Expression[T]{Expression: increment}}
		body = &stmt.Block[T]{Statements: statements}
	}
	// if no condition set true
	if condition == nil {
		condition = &expr.Literal[T]{Value: true}
	}
	body = &stmt.While[T]{Condition: condition, Body: body}

	// add the increment statement if any (before the while)
	if initializer != nil {
		statements := []stmt.Stmt[T]{initializer, body}
		body = &stmt.Block[T]{Statements: statements}
	}

	return body, nil
}

func (p *Parser[T]) printStatement() (stmt.Stmt[T], error) {
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(tokens.Semicolon, "Expect ';' after value."); err != nil {
		return nil, err
	}
	return &stmt.Print[T]{Expression: value}, nil
}

func (p *Parser[T]) expressionStatement() (stmt.Stmt[T], error) {
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(tokens.Semicolon, "Expect ';' after expression."); err != nil {
		return nil, err
	}
	return &stmt.Expression[T]{Expression: value}, nil
}

func (p *Parser[T]) Expression() (expr.Expr[T], error) {
	return p.assignment()
}

func (p *Parser[T]) assignment() (expr.Expr[T], error) {
	expression, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(tokens.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if expVar, isVariable := expression.(*expr.Variable[T]); isVariable {
			name := expVar.Name
			return &expr.Assign[T]{Name: name, Value: value}, nil
		} else if getExpr, isGet := expression.(*expr.Get[T]); isGet {
			return &expr.Set[T]{Name: getExpr.Name, Object: getExpr.Object, Value: value}, nil
		}
		return nil, parseError(equals, "Invalid assignment target.")
	}
	return expression, nil
}

func (p *Parser[T]) block() ([]stmt.Stmt[T], error) {
	statements := []stmt.Stmt[T]{}
	for (!p.check(tokens.RightBrace)) && !p.isAtEnd() {
		statement := p.declaration()
		statements = append(statements, statement)
	}
	_, err := p.consume(tokens.RightBrace, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
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
		expression = &expr.Binary[T]{Left: expression, Operator: operator, Right: right}
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
		expression = &expr.Binary[T]{Left: expression, Operator: operator, Right: right}
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
		expression = &expr.Binary[T]{Left: expression, Operator: operator, Right: right}
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
		expression = &expr.Binary[T]{Left: expression, Operator: operator, Right: right}
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
		return &expr.Unary[T]{Operator: operator, Right: right}, nil
	}
	return p.call()
}

func (p *Parser[T]) call() (expr.Expr[T], error) {
	expression, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(tokens.LeftParen) {
			expression, err = p.finishCall(expression)
			if err != nil {
				return nil, err
			}

		} else if p.match(tokens.Dot) {
			name, err := p.consume(tokens.Identifier, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expression = &expr.Get[T]{Name: name, Object: expression}
		} else {
			break
		}
	}
	return expression, nil
}

func (p *Parser[T]) finishCall(callee expr.Expr[T]) (expr.Expr[T], error) {
	arguments := []expr.Expr[T]{}
	if !p.check(tokens.RightParen) {
		arg, err := p.Expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
		for p.match(tokens.Comma) {
			if len(arguments) >= ARGUMENTS_LIMIT {
				return nil, parseError(p.peek(), fmt.Sprintf("Can't have more than %d arguments.", ARGUMENTS_LIMIT))
			}
			arg, err := p.Expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
		}
	}
	paren, err := p.consume(tokens.RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return &expr.Call[T]{Callee: callee, Paren: paren, Arguments: arguments}, nil
}

func (p *Parser[T]) primary() (expr.Expr[T], error) {
	switch {
	case p.match(tokens.False):
		return &expr.Literal[T]{Value: false}, nil
	case p.match(tokens.True):
		return &expr.Literal[T]{Value: true}, nil
	case p.match(tokens.Identifier):
		name := p.previous()
		return &expr.Variable[T]{Name: name}, nil
	case p.match(tokens.Nil):
		return &expr.Literal[T]{Value: tokens.NilLiteral}, nil
	case p.match(tokens.Number, tokens.String):
		return &expr.Literal[T]{Value: p.previous().Literal}, nil
	case p.match(tokens.This):
		return &expr.This[T]{Keyword: p.previous()}, nil
	case p.match(tokens.LeftParen):
		expression, err := p.Expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokens.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return &expr.Grouping[T]{Expression: expression}, nil
	}

	return nil, parseError(p.peek(), "Expect expression.")
}

func (p *Parser[T]) or() (expr.Expr[T], error) {
	expression, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expression = &expr.Logical[T]{Left: expression, Operator: operator, Right: right}
	}

	return expression, nil
}

func (p *Parser[T]) and() (expr.Expr[T], error) {
	expression, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expression = &expr.Logical[T]{Left: expression, Operator: operator, Right: right}
	}
	return expression, nil

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
