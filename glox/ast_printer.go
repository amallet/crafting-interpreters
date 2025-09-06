package main

import (
	"fmt"
	"strings"
)

type astPrinter struct {
	//NOP
}

func (a *astPrinter) print(expr Expr) string {
	result := expr.Accept(a)
	str, _ := result.(string) // Accept returns any
	return str
}

func (a *astPrinter) VisitBinaryExpr(expr *Binary) any {
	return a.parenthesize(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (a *astPrinter) VisitGroupingExpr(expr *Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *astPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (a *astPrinter) VisitUnaryExpr(expr *Unary) any {
	return a.parenthesize(expr.Operator.lexeme, expr.Right)
}

func (a *astPrinter) parenthesize(name string, expressions ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)

	for _, expression := range expressions {
		builder.WriteString(" ")
		res := expression.Accept(a)
		str, _ := res.(string) // Accept() returns any
		builder.WriteString(str)
	}
	builder.WriteString(")")

	return builder.String()
}
