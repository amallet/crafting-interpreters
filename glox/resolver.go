package main

import (
	"fmt"
)

type functionType int

const (
	functionTypeNone functionType = iota
	functionTypeFunction
)

type variableStatus int

const (
	isDeclared variableStatus = iota
	isDefined
	isUsed
)

type varDecl struct {
	token  Token
	status variableStatus
}

type Resolver struct {
	runtime         LoxRuntime
	interpreter     *Interpreter
	scopes          []map[string]*varDecl
	currentFunction functionType
}

func NewResolver(runtime LoxRuntime, interpreter *Interpreter) *Resolver {
	return &Resolver{
		runtime:         runtime,
		interpreter:     interpreter,
		scopes:          make([]map[string]*varDecl, 0),
		currentFunction: functionTypeNone,
	}
}

func (r *Resolver) resolveStmts(statements []Stmt) error {
	for _, stmt := range statements {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt Stmt) error {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) error {
	// A block statement gets its own scope
	r.beginScope()
	err := r.resolveStmts(stmt.statements)
	r.endScope()
	return err
}

func (r *Resolver) VisitClassStmt(stmt *ClassStmt) error {
	if err := r.declare(stmt.className); err != nil {
		return err
	}
	r.define(stmt.className)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExpressionStmt) error {
	return r.resolveExpr(stmt.expression)
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) error {
	if err := r.resolveExpr(stmt.condition); err != nil {
		return err
	}

	if err := r.resolveStmt(stmt.thenBranch); err != nil {
		return err
	}

	if stmt.elseBranch != nil {
		if err := r.resolveStmt(stmt.elseBranch); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *PrintStmt) error {
	return r.resolveExpr(stmt.expression)
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) error {

	// Can only have return statements inside a function
	if r.currentFunction == functionTypeNone {
		r.runtime.parseError(stmt.keyword, "Can't return from top-level code.")
		return fmt.Errorf("resolver error ")
	}

	if stmt.returnValue != nil { // resolve return value, if there is one
		if err := r.resolveExpr(stmt.returnValue); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitVarStmt(stmt *VarStmt) error {
	if err := r.declare(stmt.variable); err != nil {
		return err
	}

	if stmt.initializer != nil {
		if err := r.resolveExpr(stmt.initializer); err != nil {
			return err
		}
	}

	r.define(stmt.variable)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) error {
	if err := r.resolveExpr(stmt.condition); err != nil {
		return err
	}

	if err := r.resolveStmt(stmt.body); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) (any, error) {
	if err := r.resolveExpr(expr.value); err != nil {
		return nil, err
	}
	r.resolveLocal(expr, expr.variable)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) (any, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}

	for _, arg := range expr.Arguments {
		if err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitPropGetExpr(p *PropGetExpr) (any, error) {
	if err := r.resolveExpr(p.object); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitPropSetExpr(p *PropSetExpr) (any, error) {
	if err := r.resolveExpr(p.propValue); err != nil {
		return nil, err
	}

	if err := r.resolveExpr(p.object); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) (any, error) {
	return nil, r.resolveExpr(expr.Expression)
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.Right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) (any, error) {
	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitVariableExpr(expr *VariableExpr) (any, error) {

	// Check that variable isn't being referenced while still in its initializer ie
	// while it's been declared, but not yet defined
	if len(r.scopes) != 0 {
		top_scope := r.scopes[len(r.scopes)-1]
		if variable, ok := top_scope[expr.variable.lexeme]; ok {
			if variable.status == isDeclared {
				r.runtime.parseError(expr.variable, "Can't read local variable in its own initializer")
				return nil, fmt.Errorf("resolver error")
			}
		}
	}

	r.resolveLocal(expr, expr.variable)

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt *FunctionStmt) error {
	// Function name is declared and defined in the current scope
	if err := r.declare(stmt.functionName); err != nil {
		return err
	}
	r.define(stmt.functionName)
	return r.resolveFunction(stmt, functionTypeFunction)
}

func (r *Resolver) resolveFunction(function *FunctionStmt, fnType functionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = fnType

	// Function parameters and function body are in a new scope
	var err error
	r.beginScope()

	for _, param := range function.params {
		if err = r.declare(param); err != nil {
			_ = r.endScope() // might return error, but already in error case
			return err
		}
		r.define(param)
	}
	if err = r.resolveStmts(function.body); err != nil {
		_ = r.endScope() // might return error, but already in error case
		return err
	}
	err = r.endScope()

	r.currentFunction = enclosingFunction
	return err
}

func (r *Resolver) resolveLocal(expr Expr, token Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if variable, ok := r.scopes[i][token.lexeme]; ok {
			variable.status = isUsed // to keep track of used/unused variables
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) declare(token Token) error {

	if len(r.scopes) == 0 { // currently in global scope, don't need to declare it
		return nil
	}

	// Token always gets declared in the current (ie deepest) scope
	current_scope := r.scopes[len(r.scopes)-1]

	// Can't redeclare a variable if it's already been declared in this scope
	if _, ok := current_scope[token.lexeme]; ok {
		r.runtime.parseError(token, "Already a variable with this name in this scope.")
		return fmt.Errorf("resolver error")
	}

	current_scope[token.lexeme] = &varDecl{token: token, status: isDeclared}

	return nil
}

func (r *Resolver) define(token Token) {
	if len(r.scopes) == 0 { // in global scope, don't need to define token
		return
	}

	// Token is defined in current scope ie scope that's at the top of the stack
	r.scopes[len(r.scopes)-1][token.lexeme].status = isDefined
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]*varDecl))
}

func (r *Resolver) endScope() error {

	if len(r.scopes) > 0 {
		// Check that all variables defined in this scope were actually used
		for _, v := range r.scopes[len(r.scopes)-1] {
			if v.status != isUsed {
				r.runtime.parseError(v.token, "Unused variable")
				return fmt.Errorf("unused variable")
			}
		}

		r.scopes = r.scopes[:len(r.scopes)-1] // Pop top scope off the stack
	}

	return nil
}
