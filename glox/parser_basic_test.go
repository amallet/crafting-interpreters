package main

import (
	"testing"
)

func TestParserEmptyInput(t *testing.T) {
	lox := NewTestGLox()
	scanner := NewScanner(lox, "")
	tokens := scanner.scanTokens()

	parser := NewParser(lox, tokens)
	statements, err := parser.parse()

	if err != nil {
		t.Errorf("Expected no error for empty input, got: %v", err)
	}

	if len(statements) != 0 {
		t.Errorf("Expected 0 statements, got %d", len(statements))
	}
}

func TestParserLiteralExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"123;", 123.0},
		{"123.45;", 123.45},
		{`"hello";`, "hello"},
		{"true;", true},
		{"false;", false},
		{"nil;", nil},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			parser := NewParser(lox, tokens)
			statements, err := parser.parse()

			if err != nil {
				t.Errorf("Parse error for %s: %v", test.input, err)
				return
			}

			if len(statements) != 1 {
				t.Errorf("Expected 1 statement, got %d", len(statements))
				return
			}

			exprStmt, ok := statements[0].(*ExpressionStmt)
			if !ok {
				t.Errorf("Expected ExpressionStmt, got %T", statements[0])
				return
			}

			literalExpr, ok := exprStmt.expression.(*LiteralExpr)
			if !ok {
				t.Errorf("Expected LiteralExpr, got %T", exprStmt.expression)
				return
			}

			if literalExpr.Value != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, literalExpr.Value)
			}
		})
	}
}

func TestParserUnaryExpressions(t *testing.T) {
	tests := []struct {
		input       string
		operator    TokenType
		expectedRHS any // Expected right-hand side operand value
	}{
		{"!true;", BANG, true},
		{"-123;", MINUS, 123.0},
		{"!false;", BANG, false},
		{"-45.67;", MINUS, 45.67},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			parser := NewParser(lox, tokens)
			statements, err := parser.parse()

			if err != nil {
				t.Errorf("Parse error for %s: %v", test.input, err)
				return
			}

			if len(statements) != 1 {
				t.Errorf("Expected 1 statement, got %d", len(statements))
				return
			}

			exprStmt, ok := statements[0].(*ExpressionStmt)
			if !ok {
				t.Errorf("Expected ExpressionStmt, got %T", statements[0])
				return
			}

			unaryExpr, ok := exprStmt.expression.(*UnaryExpr)
			if !ok {
				t.Errorf("Expected UnaryExpr, got %T", exprStmt.expression)
				return
			}

			// Check operator
			if unaryExpr.Operator.token_type != test.operator {
				t.Errorf("Expected operator %v, got %v", test.operator, unaryExpr.Operator.token_type)
			}

			// Check right operand
			literalExpr, ok := unaryExpr.Right.(*LiteralExpr)
			if !ok {
				t.Errorf("Expected LiteralExpr as right operand, got %T", unaryExpr.Right)
				return
			}

			if literalExpr.Value != test.expectedRHS {
				t.Errorf("Expected right operand %v, got %v", test.expectedRHS, literalExpr.Value)
			}
		})
	}
}

func TestParserBinaryExpressions(t *testing.T) {
	tests := []struct {
		input       string
		operator    TokenType
		expectedLHS any // Expected left-hand side operand value
		expectedRHS any // Expected right-hand side operand value
	}{
		{"1 + 2;", PLUS, 1.0, 2.0},
		{"3 - 4;", MINUS, 3.0, 4.0},
		{"5 * 6;", STAR, 5.0, 6.0},
		{"7 / 8;", SLASH, 7.0, 8.0},
		{"9 > 10;", GREATER, 9.0, 10.0},
		{"11 >= 12;", GREATER_EQUAL, 11.0, 12.0},
		{"13 < 14;", LESS, 13.0, 14.0},
		{"15 <= 16;", LESS_EQUAL, 15.0, 16.0},
		{"17 == 18;", EQUAL_EQUAL, 17.0, 18.0},
		{"19 != 20;", BANG_EQUAL, 19.0, 20.0},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			parser := NewParser(lox, tokens)
			statements, err := parser.parse()

			if err != nil {
				t.Errorf("Parse error for %s: %v", test.input, err)
				return
			}

			if len(statements) != 1 {
				t.Errorf("Expected 1 statement, got %d", len(statements))
				return
			}

			exprStmt, ok := statements[0].(*ExpressionStmt)
			if !ok {
				t.Errorf("Expected ExpressionStmt, got %T", statements[0])
				return
			}

			binaryExpr, ok := exprStmt.expression.(*BinaryExpr)
			if !ok {
				t.Errorf("Expected BinaryExpr, got %T", exprStmt.expression)
				return
			}

			// Check operator
			if binaryExpr.Operator.token_type != test.operator {
				t.Errorf("Expected operator %v, got %v", test.operator, binaryExpr.Operator.token_type)
			}

			// Check left operand
			leftLiteral, ok := binaryExpr.Left.(*LiteralExpr)
			if !ok {
				t.Errorf("Expected LiteralExpr as left operand, got %T", binaryExpr.Left)
				return
			}

			if leftLiteral.Value != test.expectedLHS {
				t.Errorf("Expected left operand %v, got %v", test.expectedLHS, leftLiteral.Value)
			}

			// Check right operand
			rightLiteral, ok := binaryExpr.Right.(*LiteralExpr)
			if !ok {
				t.Errorf("Expected LiteralExpr as right operand, got %T", binaryExpr.Right)
				return
			}

			if rightLiteral.Value != test.expectedRHS {
				t.Errorf("Expected right operand %v, got %v", test.expectedRHS, rightLiteral.Value)
			}
		})
	}
}

func TestParserGroupingExpressions(t *testing.T) {
	input := "(1 + 2) * 3;"
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	parser := NewParser(lox, tokens)
	statements, err := parser.parse()

	if err != nil {
		t.Errorf("Parse error for %s: %v", input, err)
		return
	}

	if len(statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(statements))
		return
	}

	exprStmt, ok := statements[0].(*ExpressionStmt)
	if !ok {
		t.Errorf("Expected ExpressionStmt, got %T", statements[0])
		return
	}

	binaryExpr, ok := exprStmt.expression.(*BinaryExpr)
	if !ok {
		t.Errorf("Expected BinaryExpr, got %T", exprStmt.expression)
		return
	}

	// The outer expression should be multiplication
	if binaryExpr.Operator.token_type != STAR {
		t.Errorf("Expected STAR operator, got %v", binaryExpr.Operator.token_type)
	}

	// The left side should be a grouping expression
	groupingExpr, ok := binaryExpr.Left.(*GroupingExpr)
	if !ok {
		t.Errorf("Expected GroupingExpr on left, got %T", binaryExpr.Left)
		return
	}

	// The grouped expression should be addition
	innerBinaryExpr, ok := groupingExpr.Expression.(*BinaryExpr)
	if !ok {
		t.Errorf("Expected BinaryExpr inside grouping, got %T", groupingExpr.Expression)
		return
	}

	if innerBinaryExpr.Operator.token_type != PLUS {
		t.Errorf("Expected PLUS operator inside grouping, got %v", innerBinaryExpr.Operator.token_type)
	}
}
