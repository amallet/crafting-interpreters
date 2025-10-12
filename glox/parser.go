package main

import (
	"fmt"
)

// Parser implements a recursive descent parser for the following grammar. Note that
// this is in order of *increasing* precedence.
//
// program        → declaration* EOF;
// declaration    → varDecl | statement ;
// varDecl        → "var" IDENTIFIER ("=" expression)? ";" ;
// statement	  → exprStmt | ifStmt | printStmt | whileStmt | forStmt | block;
// exprStmt       → expression ";" ;
// ifStmt         → "if" "(" expression ")" statement ( else statement )? ;
// printStmt      → "print" expression ";" ;
// whileStmt      → "while" "(" expression ")" statement ";" ;
// forStmt        → "for" "(" ( varDecl | exprStmt | ";")
//                            expression? ";"
//                            expression? ";" ")" statement ;
// block          → "{" declaration* "}";
// expression     → assignment ";"
// assignment     → IDENTIFIER "=" assignment | logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
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
	lox     LoxRuntime
	tokens  []Token
	current int
}

func NewParser(lox LoxRuntime, tokens []Token) *Parser {
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

	if _, err = p.consume(SEMICOLON, "Expect ';' after variable declaration"); err != nil {
		return nil, err
	}

	return &VarStmt{name, init_expression}, nil
}

// statemenent → printStmt | ifStmt | exprStmt | block
func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		if statements, err := p.blockStatement(); err != nil {
			return nil, err
		} else {
			return &BlockStmt{statements}, nil
		}
	}

	return p.expressionStatement()
}

// ifStmt → "if" expression ("else" expression)?;
func (p *Parser) ifStatement() (Stmt, error) {
	var err error
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after 'if"); err != nil {
		return nil, err
	}

	var condition Expr
	if condition, err = p.expression(); err != nil {
		return nil, err
	}

	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after if condition"); err != nil {
		return nil, err
	}

	var thenBranch Stmt
	if thenBranch, err = p.statement(); err != nil {
		return nil, err
	}
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		if elseBranch, err = p.statement(); err != nil {
			return nil, err
		}
	}

	return &IfStmt{condition, thenBranch, elseBranch}, nil
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

	return &PrintStmt{expr}, nil
}

// whileStmt → "while" "(" expression ")" statement ";" ;
func (p *Parser) whileStatement() (Stmt, error) {
	var condition Expr
	var stmt Stmt
	var err error

	if _, err = p.consume(LEFT_PAREN, "Expect '(' after while condition"); err != nil {
		return nil, err
	}

	if condition, err = p.expression(); err != nil {
		return nil, err
	}

	if _, err = p.consume(RIGHT_PAREN, "Expect ')' after conditional in while statement"); err != nil {
		return nil, err
	}

	if stmt, err = p.statement(); err != nil {
		return nil, err
	}

	return &WhileStmt{condition, stmt}, nil
}

// forStmt → "for" "(" ( varDecl | exprStmt | ";")
//
//	expression? ";"
//	expression? ";" ")" statement ;
func (p *Parser) forStatement() (Stmt, error) {
	var loopVarInit Stmt = nil
	var loopCondition Expr
	var loopVarUpdate Expr
	var body Stmt
	var err error

	// 'for' statements are parsed into a combination of an initialization
	// expression and a 'while' statement that checks the loop condition,
	// with the loop variable update put into the body of the 'while' statement
	if _, err = p.consume(LEFT_PAREN, "Expect '(' after 'for'"); err != nil {
		return nil, err
	}

	// Parse initializer
	if p.match(SEMICOLON) {
		loopVarInit = nil
	} else if p.match(VAR) {
		if loopVarInit, err = p.varDeclaration(); err != nil {
			return nil, err
		}
	} else {
		if loopVarInit, err = p.expressionStatement(); err != nil {
			return nil, err
		}
	}

	// Parse the loop condition
	if !p.check(SEMICOLON) {
		if loopCondition, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(SEMICOLON, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	// Parse loop variable update
	if !p.check(RIGHT_PAREN) {
		if loopVarUpdate, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(RIGHT_PAREN, "Expect ')' after 'for' clause."); err != nil {
		return nil, err
	}

	// Parse loop body
	if body, err = p.statement(); err != nil {
		return nil, err
	}

	// Insert loop variable update, if there is one, as last statement in loop body
	if loopVarUpdate != nil {
		body = &BlockStmt{[]Stmt{
			body,
			&ExpressionStmt{loopVarUpdate},
		},
		}
	}

	// Make a while statement with the loop condition
	if loopCondition == nil {
		loopCondition = &LiteralExpr{true}
	}
	body = &WhileStmt{loopCondition, body}

	// Insert loop variable initialization before while loop
	if loopVarInit != nil {
		body = &BlockStmt{[]Stmt{
			loopVarInit,
			body,
		}}
	}

	return body, nil
}

// block → "{" declaration* "}";
func (p *Parser) blockStatement() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		if statement, err := p.declaration(); err == nil {
			statements = append(statements, statement)
		} else {
			return nil, err
		}
	}
	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}

	return statements, nil
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

	return &ExpressionStmt{expr}, nil
}

// expression → assignment;
func (p *Parser) expression() (Expr, error) {
	return p.assignmentExpr()
}

// assignment → IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignmentExpr() (Expr, error) {
	lhs, err := p.logicalOr()
	if err != nil {
		return nil, err
	}

	// Handle assignment case
	if p.match(EQUAL) {
		equals := p.previous()
		rvalue, err := p.assignmentExpr()
		if err != nil {
			return nil, err
		}

		if lvalue, ok := lhs.(*VariableExpr); ok { // LHS of assignment must be a variable
			name := lvalue.name
			return &AssignExpr{name, rvalue}, nil
		} else {
			return nil, p.constructError(equals, "Invalid assignment target")
		}
	}

	return lhs, nil
}

// logic_or → logic_and ("or" logic_and )* ;
func (p *Parser) logicalOr() (Expr, error) {
	var expr Expr
	var err error
	if expr, err = p.logicalAnd(); err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		if right, err := p.logicalAnd(); err != nil {
			return nil, err
		} else {
			expr = &LogicalExpr{expr, operator, right}
		}
	}

	return expr, nil
}

// logic_and → equality ("and" equality)* ;
func (p *Parser) logicalAnd() (Expr, error) {
	var expr Expr
	var err error
	if expr, err = p.equalityExpr(); err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		if right, err := p.equalityExpr(); err != nil {
			return nil, err
		} else {
			expr = &LogicalExpr{expr, operator, right}
		}
	}

	return expr, nil
}

// equality→ comparison ( ( "!=" | "==" ) comparison )*;
func (p *Parser) equalityExpr() (Expr, error) {
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
		expr = &BinaryExpr{expr, operator, right}
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
		expr = &BinaryExpr{expr, operator, right}
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
		expr = &BinaryExpr{expr, operator, right}
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
		expr = &BinaryExpr{expr, operator, right}
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
		return &UnaryExpr{operator, right}, nil
	}

	return p.primary()
}

// primary → "true" | "false" | "nil" | NUMBER | STRING |"(" expression ")" | IDENTIFIER;
func (p *Parser) primary() (Expr, error) {
	if p.match(TRUE) {
		return &LiteralExpr{true}, nil
	}

	if p.match(FALSE) {
		return &LiteralExpr{false}, nil
	}

	if p.match(NIL) {
		return &LiteralExpr{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{p.previous().literal}, nil
	}

	if p.match(IDENTIFIER) {
		return &VariableExpr{p.previous()}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression"); err != nil {
			return nil, err
		}
		return &GroupingExpr{expr}, nil
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

	for !p.isAtEnd() {
		if p.previous().token_type == SEMICOLON {
			return
		}

		switch p.peek().token_type {
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
