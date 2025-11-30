package main

import (
	"fmt"
)

type functionType int

const (
	functionTypeNone functionType = iota
	functionTypeFunction
	functionTypeInitializer
	functionTypeMethod 
)

type classType int 

const (
	classTypeNone classType = iota
	classTypeClass
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
	currentFunctionType functionType
	currentClassType classType
}

func NewResolver(runtime LoxRuntime, interpreter *Interpreter) *Resolver {
	return &Resolver{
		runtime:         runtime,
		interpreter:     interpreter,
		scopes:          make([]map[string]*varDecl, 0),
		currentFunctionType: functionTypeNone,
		currentClassType: classTypeNone,
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

	enclosingClass := r.currentClassType
	r.currentClassType = classTypeClass 

	// Declare and define class itself 
	if err := r.declare(stmt.className); err != nil {
		return err
	}
	r.define(stmt.className)

	// Start a new scope into which 'this' keyword can be injected, so that 
	// methods can be bound to a class instance, and then inject 'this'
	r.beginScope()
	r.injectThis()

	// Declare and define class methods
	methodNames := make(map[string]bool)
	for _, method := range stmt.methods {
		// Prevent multiple declarations of methods with the same name 
		if _, ok := methodNames[method.functionName.lexeme]; ok {
			_ = r.endScope()
			r.runtime.parseError(method.functionName,"method with this name already exists")
			return fmt.Errorf("method with name %s already exists", method.functionName.lexeme)
		} else {
			methodNames[method.functionName.lexeme] = true 
		}

		fnType := functionTypeMethod 
		if method.functionName.lexeme == "init" {
			fnType = functionTypeInitializer 
		}
		if err := r.resolveFunction(method, fnType); err != nil {
			return err 
		}
	}

	if err := r.endScope(); err != nil {
		return err 
	}

	r.currentClassType = enclosingClass
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
	if r.currentFunctionType == functionTypeNone {
		r.runtime.parseError(stmt.keyword, "Can't return from top-level code.")
		return fmt.Errorf("resolver error ")
	}

	if stmt.returnValue != nil { // resolve return value, if there is one
		// initializers can't return values 
		if r.currentFunctionType == functionTypeInitializer {
			r.interpreter.lox.error(stmt.keyword.line, "can't return a value from an initializer")
			return fmt.Errorf("can't return a value from an initializer")
		}

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

func (r *Resolver) VisitThisExpr(t *ThisExpr) (any, error) {
	// Can only reference 'this' inside a class
	if r.currentClassType == classTypeNone {
		r.interpreter.lox.error(t.keyword.line, "can't use 'this' outside a class")
		return nil, fmt.Errorf("can't use 'this' outside a class")
	}

	// "this" is treated like a local variable that gets injected 
	// by the resolver when the class is defined  
	r.resolveLocal(t, t.keyword)
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

	enclosingFunction := r.currentFunctionType
	r.currentFunctionType = fnType

	// Function parameters and function body are in a new scope
	var err error
	r.beginScope()

	// If it's a getter function, need to be inside a class
	if function.isGetter && r.currentClassType == classTypeNone {
		_ = r.endScope() // might return error, but already in error case
		r.runtime.parseError(function.functionName,"getter function has to be inside a class")
		return fmt.Errorf("getter function has to be inside a class")
	}

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

	r.currentFunctionType = enclosingFunction
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

// injectThis defines 'this' as a local variable in the current scope
func (r *Resolver) injectThis() {
	currentScope := r.scopes[len(r.scopes) - 1]
	dummyThisToken := Token{THIS, "this", nil, 0}
	currentScope["this"] = &varDecl{dummyThisToken, isUsed}
}
