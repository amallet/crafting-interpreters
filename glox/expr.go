package main

// Expr represents an expression in the Lox language
type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

// ExprVisitor defines the visitor interface for expressions
type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) (any, error)
	VisitGroupingExpr(expr *Grouping) (any, error)
	VisitLiteralExpr(expr *Literal) (any, error)
	VisitUnaryExpr(expr *Unary) (any, error)
	VisitVariableExpr(expr *Variable) (any, error)
}

// Binary represents a binary expression: left operator right
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *Binary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(e)
}

// Grouping represents a parenthesized expression: (expression)
type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(e)
}

// Literal represents a literal value expression
type Literal struct {
	Value any
}

func (e *Literal) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(e)
}

// Unary represents a unary expression: operator right
type Unary struct {
	Operator Token
	Right    Expr
}

func (e *Unary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(e)
}

// Variable represents a variable expression: <variable name>
type Variable struct {
	name Token
}

func (v *Variable) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}
