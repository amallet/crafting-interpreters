package main

import (
	"testing"
)

// Helper function to validate that an expression is a binary expression with specific operator
func validateBinaryExpr(t *testing.T, expr Expr, expectedOp TokenType) *BinaryExpr {
	binaryExpr, ok := expr.(*BinaryExpr)
	if !ok {
		t.Errorf("Expected BinaryExpr, got %T", expr)
		return nil
	}
	if binaryExpr.Operator.token_type != expectedOp {
		t.Errorf("Expected operator %v, got %v", expectedOp, binaryExpr.Operator.token_type)
		return nil
	}
	return binaryExpr
}

// Helper function to validate that an expression is a literal with specific value
func validateLiteralExpr(t *testing.T, expr Expr, expectedValue any) *LiteralExpr {
	literalExpr, ok := expr.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr, got %T", expr)
		return nil
	}
	if literalExpr.Value != expectedValue {
		t.Errorf("Expected literal %v, got %v", expectedValue, literalExpr.Value)
		return nil
	}
	return literalExpr
}

// Helper function to validate that an expression is a logical expression with specific operator
func validateLogicalExpr(t *testing.T, expr Expr, expectedOp TokenType) *LogicalExpr {
	logicalExpr, ok := expr.(*LogicalExpr)
	if !ok {
		t.Errorf("Expected LogicalExpr, got %T", expr)
		return nil
	}
	if logicalExpr.Operator.token_type != expectedOp {
		t.Errorf("Expected operator %v, got %v", expectedOp, logicalExpr.Operator.token_type)
		return nil
	}
	return logicalExpr
}

func TestParserOperatorPrecedence(t *testing.T) {
	t.Run("Multiplication has higher precedence than addition", func(t *testing.T) {
		// "1 + 2 * 3" should parse as "1 + (2 * 3)"
		input := "1 + 2 * 3;"
		lox := NewTestGLox()
		scanner := NewScanner(lox, input)
		tokens := scanner.scanTokens()

		parser := NewParser(lox, tokens)
		statements, err := parser.parse()

		if err != nil {
			t.Errorf("Parse error for %s: %v", input, err)
			return
		}

		exprStmt, ok := statements[0].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt, got %T", statements[0])
			return
		}

		// Root should be addition
		rootBinary := validateBinaryExpr(t, exprStmt.expression, PLUS)
		if rootBinary == nil {
			return
		}

		// Left side should be literal 1
		validateLiteralExpr(t, rootBinary.Left, 1.0)

		// Right side should be multiplication
		rightBinary := validateBinaryExpr(t, rootBinary.Right, STAR)
		if rightBinary == nil {
			return
		}

		// Multiplication operands should be 2 and 3
		validateLiteralExpr(t, rightBinary.Left, 2.0)
		validateLiteralExpr(t, rightBinary.Right, 3.0)
	})

	t.Run("Addition has higher precedence than equality", func(t *testing.T) {
		// "1 + 2 == 3" should parse as "(1 + 2) == 3"
		input := "1 + 2 == 3;"
		lox := NewTestGLox()
		scanner := NewScanner(lox, input)
		tokens := scanner.scanTokens()

		parser := NewParser(lox, tokens)
		statements, err := parser.parse()

		if err != nil {
			t.Errorf("Parse error for %s: %v", input, err)
			return
		}

		exprStmt, ok := statements[0].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt, got %T", statements[0])
			return
		}

		// Root should be equality
		rootBinary := validateBinaryExpr(t, exprStmt.expression, EQUAL_EQUAL)
		if rootBinary == nil {
			return
		}

		// Left side should be addition
		leftBinary := validateBinaryExpr(t, rootBinary.Left, PLUS)
		if leftBinary == nil {
			return
		}

		// Addition operands should be 1 and 2
		validateLiteralExpr(t, leftBinary.Left, 1.0)
		validateLiteralExpr(t, leftBinary.Right, 2.0)

		// Right side should be literal 3
		validateLiteralExpr(t, rootBinary.Right, 3.0)
	})

	t.Run("Comparison has higher precedence than logical AND", func(t *testing.T) {
		// "1 < 2 and 3 > 4" should parse as "(1 < 2) and (3 > 4)"
		input := "1 < 2 and 3 > 4;"
		lox := NewTestGLox()
		scanner := NewScanner(lox, input)
		tokens := scanner.scanTokens()

		parser := NewParser(lox, tokens)
		statements, err := parser.parse()

		if err != nil {
			t.Errorf("Parse error for %s: %v", input, err)
			return
		}

		exprStmt, ok := statements[0].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt, got %T", statements[0])
			return
		}

		// Root should be logical AND
		rootLogical := validateLogicalExpr(t, exprStmt.expression, AND)
		if rootLogical == nil {
			return
		}

		// Left side should be comparison
		leftBinary := validateBinaryExpr(t, rootLogical.Left, LESS)
		if leftBinary == nil {
			return
		}

		// Left comparison operands should be 1 and 2
		validateLiteralExpr(t, leftBinary.Left, 1.0)
		validateLiteralExpr(t, leftBinary.Right, 2.0)

		// Right side should be comparison
		rightBinary := validateBinaryExpr(t, rootLogical.Right, GREATER)
		if rightBinary == nil {
			return
		}

		// Right comparison operands should be 3 and 4
		validateLiteralExpr(t, rightBinary.Left, 3.0)
		validateLiteralExpr(t, rightBinary.Right, 4.0)
	})

	t.Run("Unary has higher precedence than multiplication", func(t *testing.T) {
		// "-2 * 3" should parse as "(-2) * 3"
		input := "-2 * 3;"
		lox := NewTestGLox()
		scanner := NewScanner(lox, input)
		tokens := scanner.scanTokens()

		parser := NewParser(lox, tokens)
		statements, err := parser.parse()

		if err != nil {
			t.Errorf("Parse error for %s: %v", input, err)
			return
		}

		exprStmt, ok := statements[0].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt, got %T", statements[0])
			return
		}

		// Root should be multiplication
		rootBinary := validateBinaryExpr(t, exprStmt.expression, STAR)
		if rootBinary == nil {
			return
		}

		// Left side should be unary minus
		unaryExpr, ok := rootBinary.Left.(*UnaryExpr)
		if !ok {
			t.Errorf("Expected UnaryExpr, got %T", rootBinary.Left)
			return
		}

		if unaryExpr.Operator.token_type != MINUS {
			t.Errorf("Expected MINUS operator, got %v", unaryExpr.Operator.token_type)
			return
		}

		// Unary operand should be literal 2
		validateLiteralExpr(t, unaryExpr.Right, 2.0)

		// Right side should be literal 3
		validateLiteralExpr(t, rootBinary.Right, 3.0)
	})

	t.Run("Grouping overrides precedence", func(t *testing.T) {
		// "(1 + 2) * 3" should parse as "(1 + 2) * 3" (grouping forces addition first)
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

		exprStmt, ok := statements[0].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt, got %T", statements[0])
			return
		}

		// Root should be multiplication
		rootBinary := validateBinaryExpr(t, exprStmt.expression, STAR)
		if rootBinary == nil {
			return
		}

		// Left side should be grouping
		groupingExpr, ok := rootBinary.Left.(*GroupingExpr)
		if !ok {
			t.Errorf("Expected GroupingExpr, got %T", rootBinary.Left)
			return
		}

		// Grouped expression should be addition
		groupedBinary := validateBinaryExpr(t, groupingExpr.Expression, PLUS)
		if groupedBinary == nil {
			return
		}

		// Addition operands should be 1 and 2
		validateLiteralExpr(t, groupedBinary.Left, 1.0)
		validateLiteralExpr(t, groupedBinary.Right, 2.0)

		// Right side should be literal 3
		validateLiteralExpr(t, rootBinary.Right, 3.0)
	})
}
