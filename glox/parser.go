package main

import (
	"fmt"
)

// Parser implements a recursive descent parser for the following grammar. Note that
// this is (roughly) in order of *increasing* precedence.
//
// program        → declaration* EOF;
// declaration    → classDecl | funDecl | varDecl | statement ;
// classDecl      → "class" IDENTIFIER "(" function* ")" ;
// funDecl        → "fun" function;
// function       → IDENTIFIER ("(" parameters? ")")? block ;
// parameters     → IDENTIFIER ("," IDENTIFIER)* ;
// varDecl        → "var" IDENTIFIER ("=" expression)? ";" ;
// statement	  → exprStmt | ifStmt | printStmt | whileStmt | forStmt | returnStmt | block;
// exprStmt       → expression ";" ;
// ifStmt         → "if" "(" expression ")" statement ( else statement )? ;
// printStmt      → "print" expression ";" ;
// whileStmt      → "while" "(" expression ")" statement ";" ;
// forStmt        → "for" "(" ( varDecl | exprStmt | ";")
//                            expression? ";"
//                            expression? ";" ")" statement ;
// returnStmt     → "return" expression? ";" ;
// block          → "{" declaration* "}";
// expression     → assignmentOrValue ";"
// assignmentOrValue     → ( call ".")? IDENTIFIER "=" assignment | logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
// equality       → comparison ( ( "!=" | "==" ) comparison )*;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*;
// term           → factor ( ( "-" | "+" ) factor )*;
// factor         → unary ( ( "/" | "*" ) unary )*;
// unary          → ( "!" | "-" ) unary | | call
// call           → primary ( "(" arguments? ")" | "." IDENTIFIER )*;
// arguments      → expression ( "," expression )* ;
// primary        → "true" | "false" | "nil" | "this" | NUMBER | STRING | "(" expression ")" | IDENTIFIER;
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

// declaration -> classDecl | funDecl | varDecl | statement
func (p *Parser) declaration() (Stmt, error) {
	var stmt Stmt
	var err error

	if p.matches(CLASS) {
		stmt, err = p.classDeclaration()
	} else if p.matches(FUN) {
		stmt, err = p.function("function")
	} else if p.matches(VAR) {
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

// class → "class" IDENTIFIER "{" function* "}";
func (p *Parser) classDeclaration() (Stmt, error) {
	var err error
	var className Token
	methods := make([]*FunctionStmt, 0)

	if className, err = p.consume(IDENTIFIER, "Expect class name"); err != nil {
		return nil, err
	}

	// Parse class methods
	if _, err := p.consume(LEFT_BRACE, "Expect '{' after class name"); err != nil {
		return nil, err
	}

	for !p.nextTokenTypeIs(RIGHT_BRACE) && !p.isAtEnd() {
		var method *FunctionStmt
		if method, err = p.function("method"); err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}

	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after class body"); err != nil {
		return nil, err
	}

	return &ClassStmt{className: className, methods: methods}, nil
}

// function       → IDENTIFIER ("(" parameters? ")")? block ;
// parameters     → IDENTIFIER ("," IDENTIFIER)* ;
func (p *Parser) function(kind string) (*FunctionStmt, error) {
	var err error
	var fnName Token
	var fnParams []Token

	// Parse function name
	if fnName, err = p.consume(IDENTIFIER, "Expect "+kind+" name."); err != nil {
		return nil, err
	}

	// A 'getter' is modeled as a function with no parameter list, so parsing needs to handle
	// functions with and without parameter lists 
	isGetter := true 
	if p.nextTokenTypeIs(LEFT_PAREN) { 
		isGetter = false // regular function, with (possibly empty) parameter list

		// Parse function params
		if _, err = p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name."); err != nil {
			return nil, err
		}

		fnParams = make([]Token, 0)
		if !p.nextTokenTypeIs(RIGHT_PAREN) {
			var parameter Token
			if parameter, err = p.consume(IDENTIFIER, "Expect parameter name."); err != nil {
				return nil, err
			}
			fnParams = append(fnParams, parameter)

			for p.matches(COMMA) {
				if len(fnParams) >= 255 {
					return nil, p.constructError(p.peek(), "Can't have more than 255 parameters.")
				}
				if parameter, err = p.consume(IDENTIFIER, "Expect parameter name."); err != nil {
					return nil, err
				}
				fnParams = append(fnParams, parameter)
			}

		}
		if _, err = p.consume(RIGHT_PAREN, "Expect ')' after parameters. "); err != nil {
			return nil, err
		}
	}

	// Special 'init' method that serves as constructor must always have parameter list, even if it's empty 
	if fnName.lexeme == "init" && isGetter {
		return nil, p.constructError(p.previous(),"init function must have parameter list")
	}

	// Parse function body
	if _, err = p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body."); err != nil {
		return nil, err
	}

	var fnBody []Stmt
	if fnBody, err = p.blockStatement(); err != nil {
		return nil, err
	}

	return &FunctionStmt{fnName, isGetter, fnParams, fnBody}, nil
}

// varDecl → "var" IDENTIFIER ("=" expression)? ";" ;
func (p *Parser) varDeclaration() (Stmt, error) {
	// 'var' keyword has already been consumed, so start by trying to parse
	// the identifier ie the variable name
	varName, err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}

	// Parse initialization expression, if there is one
	var initExpression Expr = nil
	if p.matches(EQUAL) {
		if initExpression, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(SEMICOLON, "Expect ';' after variable declaration"); err != nil {
		return nil, err
	}

	return &VarStmt{varName, initExpression}, nil
}

// statemenent → printStmt | ifStmt | exprStmt | block
func (p *Parser) statement() (Stmt, error) {
	if p.matches(IF) {
		return p.ifStatement()
	}

	if p.matches(FOR) {
		return p.forStatement()
	}

	if p.matches(WHILE) {
		return p.whileStatement()
	}

	if p.matches(PRINT) {
		return p.printStatement()
	}

	if p.matches(RETURN) {
		return p.returnStatement()
	}

	if p.matches(LEFT_BRACE) {
		if statements, err := p.blockStatement(); err != nil {
			return nil, err
		} else {
			return &BlockStmt{statements}, nil
		}
	}

	return p.expressionStatement()
}

// ifStmt → "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() (Stmt, error) {

	// 'if' keyword has already been consumed, so start parsing what's supposed to come
	// next
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
	if p.matches(ELSE) {
		if elseBranch, err = p.statement(); err != nil {
			return nil, err
		}
	}

	return &IfStmt{condition, thenBranch, elseBranch}, nil
}

// printStmt → "print" expression ";"
func (p *Parser) printStatement() (Stmt, error) {

	// 'print' keyword has already been consumed, so start parsing what's supposed to come
	// next
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

	// 'while' keyword has already been consumed, so start parsing what's supposed to come
	// next
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

// forStmt → "for" "(" ( varDecl | exprStmt | ";") expression? ";" expression? ";" ")" statement ;
func (p *Parser) forStatement() (Stmt, error) {

	// 'for' statements are desugared by being parsed into a combination of an initialization
	// expression and a 'while' statement that checks the loop condition, with the loop
	// variable update put into the body of the 'while' statement
	var loopVarInit Stmt = nil
	var loopCondition Expr
	var loopVarUpdate Expr
	var body Stmt
	var err error

	// 'for' keyword has already been consumed, so start parsing what's supposed to come
	// next
	if _, err = p.consume(LEFT_PAREN, "Expect '(' after 'for'"); err != nil {
		return nil, err
	}

	// Parse initializer
	if p.matches(SEMICOLON) {
		loopVarInit = nil
	} else if p.matches(VAR) {
		if loopVarInit, err = p.varDeclaration(); err != nil {
			return nil, err
		}
	} else {
		if loopVarInit, err = p.expressionStatement(); err != nil {
			return nil, err
		}
	}

	// Parse the loop condition
	if !p.nextTokenTypeIs(SEMICOLON) {
		if loopCondition, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(SEMICOLON, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	// Parse loop variable update
	if !p.nextTokenTypeIs(RIGHT_PAREN) {
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

// return → "return" expression? ";" ;
func (p *Parser) returnStatement() (Stmt, error) {

	// 'return' keyword has already been consumed, so start parsing what's supposed to come
	// next
	keyword := p.previous()

	var value Expr = nil
	var err error
	if !p.nextTokenTypeIs(SEMICOLON) { // returning a value
		if value, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after return value."); err != nil {
		return nil, err
	}

	return &ReturnStmt{keyword, value}, nil
}

// block → "{" declaration* "}";
func (p *Parser) blockStatement() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.nextTokenTypeIs(RIGHT_BRACE) && !p.isAtEnd() {
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

// expression → assignmentOrValue;
func (p *Parser) expression() (Expr, error) {
	return p.assignmentOrValueExpr()
}

// assignmentOrValue → (call ".")? IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignmentOrValueExpr() (Expr, error) {
	// Have to handle expressions that are either assignments or 'just' expressions
	// that return a value. We don't know whether it's an assignment expression untl
	// parser sees an '=' sign, so start off by assumng it's not an assignment and
	// greedily parse as much of it as possible
	lhs, err := p.logicalOr()
	if err != nil {
		return nil, err
	}

	// IF there is an equal sign, it's an assignment
	if p.matches(EQUAL) {
		// What we have so far is the  l-value of the assignment, so need to
		// now parse the r-value of the assignment
		equals := p.previous()
		rvalue, err := p.assignmentOrValueExpr()
		if err != nil {
			return nil, err
		}

		// Only variables or instance properties can be assigned to
		switch lvalue := lhs.(type) {
		case *VariableExpr:
			name := lvalue.variable
			return &AssignExpr{name, rvalue}, nil
		case *PropGetExpr:
			return &PropSetExpr{lvalue.object, lvalue.propName, rvalue}, nil
		default:
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

	for p.matches(OR) {
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

	for p.matches(AND) {
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

	for p.matches(BANG_EQUAL, EQUAL_EQUAL) {
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

	for p.matches(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
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

	for p.matches(MINUS, PLUS) {
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

	for p.matches(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{expr, operator, right}
	}

	return expr, nil
}

// unary → ( "!" | "-" ) unary || call
func (p *Parser) unary() (Expr, error) {
	if p.matches(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{operator, right}, nil
	} else {
		return p.call()
	}
}

// call → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
func (p *Parser) call() (Expr, error) {
	var expr Expr
	var err error

	if expr, err = p.primary(); err != nil {
		return nil, err
	}

	for {
		if p.matches(LEFT_PAREN) {
			if expr, err = p.callArguments(expr); err != nil {
				return nil, err
			}
		} else if p.matches(DOT) {
			var propName Token
			if propName, err = p.consume(IDENTIFIER, "Expect property name after '.'"); err != nil {
				return nil, err
			}

			expr = &PropGetExpr{object: expr, propName: propName}
		} else {
			break
		}
	}

	return expr, nil
}

// arguments → expression ( "," expression )* ;
func (p *Parser) callArguments(callee Expr) (Expr, error) {
	var err error
	arguments := make([]Expr, 0)

	if !p.nextTokenTypeIs(RIGHT_PAREN) { // check for arguments
		// Parse first argument
		var expr Expr
		if expr, err = p.expression(); err != nil {
			return nil, err
		}
		arguments = append(arguments, expr)

		// Parse all the other arguments, if any
		for p.matches(COMMA) {
			if len(arguments) >= 255 {
				return nil, p.constructError(p.peek(), "Can't have more than 255 arguments.")
			}
			if expr, err = p.expression(); err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)
		}
	}

	var paren Token
	if paren, err = p.consume(RIGHT_PAREN, "Expect ')' after arguments."); err != nil {
		return nil, err
	}

	return &CallExpr{callee, paren, arguments}, nil
}

// primary → "true" | "false" | "nil" | "this" | NUMBER | STRING |"(" expression ")" | IDENTIFIER;
func (p *Parser) primary() (Expr, error) {
	if p.matches(TRUE) {
		return &LiteralExpr{true}, nil
	}

	if p.matches(FALSE) {
		return &LiteralExpr{false}, nil
	}

	if p.matches(NIL) {
		return &LiteralExpr{nil}, nil
	}

	if p.matches(NUMBER, STRING) {
		return &LiteralExpr{p.previous().literal}, nil
	}

	if p.matches(IDENTIFIER) {
		return &VariableExpr{p.previous()}, nil
	}

	if p.matches(THIS) {
		return &ThisExpr{p.previous()}, nil
	}

	if p.matches(LEFT_PAREN) {
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

func (p *Parser) matches(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.nextTokenTypeIs(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.nextTokenTypeIs(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.constructError(p.peek(), message)
}

func (p *Parser) nextTokenTypeIs(tokenType TokenType) bool {
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
