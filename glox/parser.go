package main

import (
	"fmt"
)

// Parser implements a recursive descent parser for the following grammar
// 
// program        → declaration* EOF;
// declaration    → varDecl | statement ;
// varDecl        → "var" IDENTIFIER ("=" expression)? ";" ;
// statement	  → exprStmt | printStmt ;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;
// expression     → equality
// equality       → comparison ( ( "!=" | "==" ) comparison )*;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*;
// term           → factor ( ( "-" | "+" ) factor )*;
// factor         → unary ( ( "/" | "*" ) unary )*;
// unary          → ( "!" | "-" ) unary | primary;
// primary        → "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")" | IDENTIFIER;
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


// Generate list of ASTs representing the expressions being parsed 
func (p *Parser) parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		if stmt, err := p.declaration(); err != nil {
			return nil, err 
		} else {
			statements = append(statements, stmt)
		}
	}
	return statements, nil 
}

// declaration -> varDecl | statement
func (p *Parser) declaration() (Stmt, error) {
	var stmt Stmt
	var err error 

	if p.match(VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}

	// If parsing encountered an error, update parser state to a place where parsing can contine 
	// and report the error 
	if err != nil {
		p.synchronize()
		return nil, err 
	}

	return stmt, nil
}

// varDecl → "var" IDENTIFIER ("=" expression)? ";" ;
func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err 
	}

	// Parse initialization expression, if there is one 
	var init_expression Expr = nil 
	if p.match(EQUAL) {
		if init_expression, err = p.expression(); err != nil {
			return nil, err 
		}
	}
	
	if _, err = p.consume(SEMICOLON,"Expect ';' after variable declaration"); err != nil {
		return nil, err 
	}

	return &VarStmt { name, init_expression}, nil 
}

// statemenent → printStmt | exprStmt 
func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

// printStmt → "print" expression ";"
func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err 
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err 
	}

	return &PrintStmt{ expr }, nil 
}

// exprStmt → expression ";"
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err 
	}

	return &ExpressionStmt{ expr }, nil 
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

// primary → "true" | "false" | "nil" | NUMBER | STRING |"(" expression ")" | IDENTIFIER;
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

	if p.match(IDENTIFIER) {
		return &Variable{p.previous()}, nil 
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

// Called in case of parse error, advances through tokens until start of next statement
func (p *Parser) synchronize() {
	_ = p.advance()

	for (!p.isAtEnd()) {
		if (p.previous().token_type == SEMICOLON) {
			return 
		}

		switch (p.peek().token_type) {
		case CLASS:
			fallthrough
		case FUN:
			fallthrough
		case VAR:
			fallthrough
		case FOR:
			fallthrough
		case IF:
			fallthrough
		case WHILE:
			fallthrough
		case PRINT:
			fallthrough
		case RETURN:
			return 
		}

		_ = p.advance()
	}
}

func (p *Parser) constructError(token Token, message string) error {
	p.lox.parseError(token, message)
	return fmt.Errorf("parse error")

}
