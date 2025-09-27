//nolint:unused 
package main

import (
	"fmt"
	"strings"
)

type rpnPrinter struct {
	//NOP
}

func (a *rpnPrinter) print(expr Expr) string {
	result, _ := expr.Accept(a)
	str, _ := result.(string) // Accept returns any
	return str
}

func (a *rpnPrinter) VisitBinaryExpr(expr *Binary) (any, error) {
	return a.rpn_order(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (a *rpnPrinter) VisitGroupingExpr(expr *Grouping) (any,error) {
	return a.rpn_order("", expr.Expression)
}

func (a *rpnPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil 
	}
	return fmt.Sprintf("%v", expr.Value), nil 
}

func (a *rpnPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	res, _ := expr.Right.Accept(a)
	str := res.(string)
	if expr.Operator.token_type == MINUS {
        return fmt.Sprintf("0 %s -", str), nil
    } else {
		return fmt.Sprintf("%s %s", str, expr.Operator.lexeme), nil 
	}
}

func (a *rpnPrinter) VisitVariableExpr(expr *Variable) (any, error) {
	return expr.name.lexeme, nil 
 }
 

func (a *rpnPrinter) rpn_order(operation string, expressions ...Expr) (string, error) {
	var builder strings.Builder

	for _, expression := range expressions {
		builder.WriteString(" ")
		res, _ := expression.Accept(a)
		str, _ := res.(string) // Accept() returns any
		builder.WriteString(str)
	}
	if operation != "" {
		builder.WriteString(" ")
		builder.WriteString(operation)	
	}

	return builder.String(), nil 
}
