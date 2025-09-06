package main

// Expr represents an expression in the Lox language
type Expr interface {
	Accept(visitor ExprVisitor) any
}

// ExprVisitor defines the visitor interface for expressions
type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitUnaryExpr(expr *Unary) any
}

// Binary represents a binary expression: left operator right
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(e)
}

// Grouping represents a parenthesized expression: (expression)
type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(e)
}

// Literal represents a literal value expression
type Literal struct {
	Value any
}

func (e *Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(e)
}

// Unary represents a unary expression: operator right
type Unary struct {
	Operator Token
	Right    Expr
}

func (e *Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}
