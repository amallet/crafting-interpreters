package main

import (
	"fmt"
)

// Parser implements a recursive descent parser for the following grammar:
//
// expression     → equality
// equality       → comparison ( ( "!=" | "==" ) comparison )*;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*;
// term           → factor ( ( "-" | "+" ) factor )*;
// factor         → unary ( ( "/" | "*" ) unary )*;
// unary          → ( "!" | "-" ) unary | primary;
// primary        → "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")";
//
// The grammar follows operator precedence with the following precedence levels
// (from lowest to highest):
// 1. Equality: ==, !=
// 2. Comparison: >, >=, <, <=
// 3. Term: +, -
// 4. Factor: *, /
// 5. Unary: !, -
// 6. Primary: literals, grouping

type Parser struct {
	lox     *GLox
	tokens  []Token
	current int
}

func NewParser(lox *GLox, tokens []Token) *Parser {
	return &Parser{
		lox:     lox,
		tokens:  append([]Token(nil), tokens...),
		current: 0,
	}
}


func (p *Parser) parse() (Expr, error) {
	return p.expression()
}

// expression → equality
func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )*;
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )*;
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

// term → factor ( ( "-" | "+" ) factor )*;
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

// factor → unary ( ( "/" | "*" ) unary )*;
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{expr, operator, right}
	}

	return expr, nil
}

// unary → ( "!" | "-" ) unary | primary;
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{operator, right}, nil
	}

	return p.primary()
}

// primary → "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")";
func (p *Parser) primary() (Expr, error) {
	if p.match(TRUE) {
		return &Literal{TRUE}, nil
	}

	if p.match(FALSE) {
		return &Literal{FALSE}, nil
	}

	if p.match(NIL) {
		return &Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression"); err != nil {
			return nil, err
		}
		return &Grouping{expr}, nil
	}

	return nil, p.constructError(p.peek(), "Expected expression")
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.constructError(p.peek(), message)
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().token_type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().token_type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) constructError(token Token, message string) error {
	p.lox.parseError(token, message)
	return fmt.Errorf("parse error")

}
