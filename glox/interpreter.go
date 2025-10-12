package main

import (
	"fmt"
	"reflect"
)

// The Interpreter object interprets/evaluates the ASTs produced by the Parser
type Interpreter struct {
	lox LoxRuntime
	env *Environment
}

func NewInterpreter(lox LoxRuntime) *Interpreter {
	return &Interpreter{
		lox: lox,
		env: NewEnvironment(nil),
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
	i.env.defineVarValue(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) error {
	_, err := i.evaluate(stmt.expression) // Expression statements don't produce a value
	return err
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) error {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", value) // Print statement outputs result of evaluating expression
	return nil
}

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
	for true {
		var condition any
		var err error
		if condition, err = i.evaluate(stmt.condition); err != nil {
			return err
		}
		if isTruthy(condition) {
			if err = i.execute(stmt.statement); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) error {
	// When interpreting a block, create a new environment to handle
	// the lexical scope for that block, and use it to evaluate statements
	// inside the block
	prevEnv := i.env
	blockEnv := NewEnvironment(i.env)
	i.env = blockEnv // use new environment to evaluate statements in the block
	var err error = nil
	for _, stmt := range stmt.statements {
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
	if err = i.env.assignVarValue(expr.name, value); err != nil {
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

// Evaluate expressions in parentheses
func (i *Interpreter) VisitGroupingExpr(e *GroupingExpr) (any, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value, nil // Interpreting/evaluating a literal expression just returns the actual value
}

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

func (i *Interpreter) VisitVariableExpr(expr *VariableExpr) (any, error) {
	return i.env.getVarValue(expr.name) // Evaluating a variable expression just returns the associated value

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
