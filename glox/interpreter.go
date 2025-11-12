package main

import (
	"fmt"
	"reflect"
)

// The Interpreter object interprets/evaluates the ASTs produced by the Parser
type Interpreter struct {
	lox     LoxRuntime
	globals *Environment
	env     *Environment
	locals  map[Expr]int
}

func NewInterpreter(lox LoxRuntime) *Interpreter {
	globals := NewEnvironment(nil)
	globals.defineVarValue("clock", clockFn{})
	return &Interpreter{
		lox:     lox,
		globals: globals,
		env:     globals,
		locals:  make(map[Expr]int),
	}
}

func (i *Interpreter) interpret(statements []Stmt) []any {
	results := make([]any, 0)
	for _, stmt := range statements {
		// Collect the results of evaluating any top-level statements that are
		// expressions, used for REPL mode
		if expr_stmt, ok := stmt.(*ExpressionStmt); ok {
			if value, err := i.evaluate(expr_stmt.expression); err == nil {
				results = append(results, value)
			} else {
				i.lox.runtimeError(err)
				return nil
			}
		} else {
			if err := i.execute(stmt); err != nil {
				i.lox.runtimeError(err)
				return nil
			}
		}
	}

	return results
}

func (i *Interpreter) execute(stmt Stmt) error {
	return stmt.Accept(i)
}

func (i *Interpreter) evaluate(e Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	//fmt.Printf("Storing %v @ %v at depth %d\n", expr, &expr, depth)
	i.locals[expr] = depth
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) error {
	var value any
	var err error

	// If variable has an initialization expression, need to evaluate the expression
	// so the value can be assigned to the variable
	if stmt.initializer != nil {
		if value, err = i.evaluate(stmt.initializer); err != nil {
			return err
		}
	}

	// Store the variable name and associated value
	i.env.defineVarValue(stmt.variable.lexeme, value)
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) error {
	_, err := i.evaluate(stmt.expression) // Expression statements don't produce a value
	return err
}

func (i *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) error {
	loxFn := &LoxFunction{stmt, i.env}
	i.env.defineVarValue(stmt.functionName.lexeme, loxFn)
	return nil
}

// Execute print statement
func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) error {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", value) // Print statement outputs result of evaluating expression
	return nil
}

// Execute return statement
func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) error {
	var value any
	var err error

	if stmt.returnValue != nil {
		// Evaluate the expression to be returned, and then wrap it in a
		// special sentinel type that conforms to the Error() interface, but
		// actually wraps the return value. This is used to short-cut execution
		// and unwind the call stack. Not great, but it's what we've got.
		if value, err = i.evaluate(stmt.returnValue); err != nil {
			return err
		} else {
			return &ReturnValue{value}
		}
	}

	// Return without value - still need to wrap in ReturnValue to stop execution
	return &ReturnValue{nil}
}

// Execute 'if' statement
func (i *Interpreter) VisitIfStmt(stmt *IfStmt) error {
	var err error
	var condition any
	if condition, err = i.evaluate(stmt.condition); err != nil {
		return err
	}
	if isTruthy(condition) {
		return i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		return i.execute(stmt.elseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) error {
	for {
		var condition any
		var err error
		if condition, err = i.evaluate(stmt.condition); err != nil {
			return err
		}
		if isTruthy(condition) {
			if err = i.execute(stmt.body); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// Execute statements within a block ie { ... }
func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) error {
	// When interpreting a block, create a new environment to handle
	// the lexical scope for that block, and use it to evaluate statements
	// inside the block
	blockEnv := NewEnvironment(i.env)
	return i.executeBlock(stmt.statements, blockEnv)
}

// Execute block of statements, within the supplied environment
func (i *Interpreter) executeBlock(statements []Stmt, blockEnv *Environment) error {

	prevEnv := i.env
	i.env = blockEnv // use new environment to evaluate statements in the block
	var err error = nil
	for _, stmt := range statements {
		if err = i.execute(stmt); err != nil {
			break
		}
	}
	// Done evaluating the block, restore previous environment
	i.env = prevEnv
	return err
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) (any, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	// If local variable, assign to the right scope
	if distance, ok := i.locals[expr]; ok {
		i.env.assignAt(distance, expr.variable, value)
		return value, nil
	}
	// Else, it's a variable in the global scope
	if err = i.globals.assignVarValue(expr.variable, value); err != nil {
		return nil, err
	}

	return value, nil // Assignment expressions return the value on the RHS
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.token_type {
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil

	case GREATER:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val > right_val), nil
		} else {
			return nil, err
		}

	case GREATER_EQUAL:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val >= right_val), nil
		} else {
			return nil, err
		}

	case LESS:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val < right_val), nil
		} else {
			return nil, err
		}

	case LESS_EQUAL:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val <= right_val), nil
		} else {
			return nil, err
		}

	case MINUS:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val - right_val), nil
		} else {
			return nil, err
		}

	case PLUS:
		switch left.(type) {
		case float64:
			if right_val, ok := right.(float64); ok {
				left_val, _ := left.(float64)
				return (left_val + right_val), nil
			}

		case string:
			if right_val, ok := right.(string); ok {
				left_val, _ := left.(string)
				return (left_val + right_val), nil
			}
		}

		return nil, RuntimeError{expr.Operator, "operands to operator + must be numbers/strings"}

	case SLASH:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			if right_val == 0 {
				return nil, RuntimeError{expr.Operator, "illegal operation: division by zero"}
			}
			return (left_val / right_val), nil
		} else {
			return nil, err
		}

	case STAR:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val * right_val), nil
		} else {
			return nil, err
		}

	default:
		//NOP
	}

	return nil, nil
}

// Evaluate function calls
func (i *Interpreter) VisitCallExpr(e *CallExpr) (any, error) {
	var callee any
	var arguments []any
	var callable LoxCallable
	var err error

	// Resolve function being called
	if callee, err = i.evaluate(e.Callee); err != nil {
		return nil, err
	}

	// Generate values for all the arguments
	arguments = make([]any, 0)
	for _, arg := range e.Arguments {
		if value, err := i.evaluate(arg); err != nil {
			return nil, err
		} else {
			arguments = append(arguments, value)
		}
	}

	// Make actual call to function, if it is callable
	var ok bool
	if callable, ok = callee.(LoxCallable); !ok {
		return nil, RuntimeError{e.Paren, "Can only call functions and classes."}
	}
	if callable.arity() != len(arguments) {
		return nil, RuntimeError{e.Paren,
			fmt.Sprintf("Expected %d arguments but got %d", callable.arity(), len(arguments))}
	}

	return callable.call(i, arguments)
}

// Evaluate expressions in parentheses
func (i *Interpreter) VisitGroupingExpr(e *GroupingExpr) (any, error) {
	return i.evaluate(e.Expression)
}

// Evaluate a literal expression
func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value, nil // Interpreting/evaluating a literal expression just returns the actual value
}

// Evaluate logical (AND, OR) expression
func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) (any, error) {

	if left, err := i.evaluate(expr.Left); err != nil {
		return nil, err
	} else {
		if isTruthy(left) {
			if expr.Operator.token_type == OR {
				return left, nil
			}
		} else if expr.Operator.token_type == AND {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

// Evaluate unary expr eg -5, !foo
func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.token_type {
	case BANG:
		return !isTruthy(right), nil

	case MINUS:
		if value, ok := right.(float64); ok {
			return (-value), nil
		} else {
			return nil, RuntimeError{expr.Operator, "operand to operator - must be a number"}
		}
	}

	return nil, nil
}

// Evaluate variable
func (i *Interpreter) VisitVariableExpr(expr *VariableExpr) (any, error) {
	return i.lookupVariable(expr.variable, expr)
}

func (i *Interpreter) lookupVariable(name Token, expr Expr) (any, error) {
	if dist, ok := i.locals[expr]; ok {
		return i.env.getAt(dist, name.lexeme), nil
	} else {
		value, err := i.globals.getVarValue(name)
		return value, err
	}
}

func isTruthy(a any) bool {
	if a == nil {
		return false
	}

	if val, ok := a.(bool); ok {
		return val
	}

	return true
}

func convertNumberOperands(operator Token, a, b any) (float64, float64, error) {
	va, a_ok := a.(float64)
	vb, b_ok := b.(float64)

	if !a_ok || !b_ok {
		return 0, 0, RuntimeError{operator,
			fmt.Sprintf("operands to operator %s must be numbers", operator.lexeme)}
	}

	return va, vb, nil
}

func isEqual(a any, b any) bool {
	return reflect.DeepEqual(a, b)
}
