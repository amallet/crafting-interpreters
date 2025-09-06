package main

import (
	"fmt"
	"strings"
)

type rpnPrinter struct {
	//NOP
}

func (a *rpnPrinter) print(expr Expr) string {
	result := expr.Accept(a)
	str, _ := result.(string) // Accept returns any
	return str
}

func (a *rpnPrinter) VisitBinaryExpr(expr *Binary) any {
	return a.rpn_order(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (a *rpnPrinter) VisitGroupingExpr(expr *Grouping) any {
	return a.rpn_order("", expr.Expression)
}

func (a *rpnPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *rpnPrinter) VisitUnaryExpr(expr *Unary) any {
	res := expr.Right.Accept(a)
	str := res.(string)
	if expr.Operator.token_type == MINUS {
        return fmt.Sprintf("0 %s -", str)
    } else {
		return fmt.Sprintf("%s %s", str, expr.Operator.lexeme)
	}
}

func (a *rpnPrinter) rpn_order(operation string, expressions ...Expr) string {
	var builder strings.Builder

	for _, expression := range expressions {
		builder.WriteString(" ")
		res := expression.Accept(a)
		str, _ := res.(string) // Accept() returns any
		builder.WriteString(str)
	}
	if operation != "" {
		builder.WriteString(" ")
		builder.WriteString(operation)	
	}

	return builder.String()
}
