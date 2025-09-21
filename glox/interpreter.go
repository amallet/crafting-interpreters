package main

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	lox *GLox
}

func NewInterpreter(lox *GLox) *Interpreter {
	return &Interpreter{
		lox:     lox,
	}
}

func (i *Interpreter) interpret(expr Expr) {
	if value, err := i.evaluate(expr); err == nil {
		fmt.Printf("%v\n", value)
	} else {
		i.lox.runtimeError(err)
	}
}

func (i *Interpreter) evaluate(e Expr) (any, error) {
	return e.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err 
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err 
	}

	switch(expr.Operator.token_type) {
	case BANG_EQUAL: return !isEqual(left, right), nil
	case EQUAL_EQUAL: return isEqual(left, right), nil 

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
		
		return nil, RuntimeError{ expr.Operator, "operands to operator + must be numbers/strings" }

	case SLASH:
		if left_val, right_val, err := convertNumberOperands(expr.Operator, left, right); err == nil {
			return (left_val/right_val), nil 
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

func (i *Interpreter) VisitGroupingExpr(e *Grouping) (any, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) (any, error) {
	return expr.Value, nil 
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err 
	}

	switch (expr.Operator.token_type) {
	case BANG:
		return !isTruthy(right), nil 

	case MINUS:
		if value, ok := right.(float64); ok {
			return (-value), nil 
		} else {
			return nil, RuntimeError{ expr.Operator, "operand to operator - must be a number"}
		}
	}

	return nil, nil 
}

func isTruthy(a any) bool {
	if (a == nil) {
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


