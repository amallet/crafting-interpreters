package main

// Expr represents an expression in the Lox language
type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

// ExprVisitor defines the visitor interface for expressions
type ExprVisitor interface {
	VisitAssignExpr(expr *AssignExpr) (any, error)
	VisitCallExpr(expr *CallExpr) (any, error)
	VisitPropGetExpr(expr *PropGetExpr) (any, error)
	VisitPropSetExpr(expr *PropSetExpr) (any, error)
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitLogicalExpr(expr *LogicalExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitVariableExpr(expr *VariableExpr) (any, error)
}

// AssignExpr represents an assignment expression
type AssignExpr struct {
	variable Token
	value    Expr
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

// CallExpr represents a call expression
type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (c *CallExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitCallExpr(c)
}

// PropGetExpr represents an expression retrieving a property on an object
type PropGetExpr struct {
	object   Expr
	propName Token
}

func (p *PropGetExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitPropGetExpr(p)
}

// PropSetExpr represents an expression setting a property on an object
type PropSetExpr struct {
	object    Expr
	propName  Token
	propValue Expr
}

func (p *PropSetExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitPropSetExpr(p)
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
	variable Token
}

func (v *VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}
