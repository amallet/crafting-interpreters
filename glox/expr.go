package main

// Expr represents an expression in the Lox language
type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

// ExprVisitor defines the visitor interface for expressions
type ExprVisitor interface {
	VisitAssignExpr(expr *AssignExpr) (any, error)
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitLogicalExpr(expr *LogicalExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitVariableExpr(expr *VariableExpr) (any, error)
}

// AssignExpr represents an assignment expression
type AssignExpr struct {
	name  Token
	value Expr
}

func (a *AssignExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssignExpr(a)
}

// BinaryExpr represents a binary expression: left operator right
type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(e)
}

// GroupingExpr represents a parenthesized expression: (expression)
type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(e)
}

// LiteralExpr represents a literal value expression
type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(e)
}

// LogicalExpr represents a logical expression (and, or)
type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l *LogicalExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

// UnaryExpr represents a unary expression: operator right
type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (e *UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(e)
}

// VariableExpr represents a variable expression: <variable name>
type VariableExpr struct {
	name Token
}

func (v *VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}
