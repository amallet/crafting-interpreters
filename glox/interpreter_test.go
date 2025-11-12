package main

import (
	"reflect"
	"testing"
)

// Test helper to create a new interpreter with test environment
func createTestInterpreter() *Interpreter {
	testGLox := NewTestGLox()
	interpreter := NewInterpreter(testGLox)
	return interpreter
}

// Test helper to evaluate an expression and return the result
func evaluateExpression(t *testing.T, interpreter *Interpreter, expr Expr) any {
	result, err := interpreter.evaluate(expr)
	if err != nil {
		t.Fatalf("Unexpected error evaluating expression: %v", err)
	}
	return result
}

// Test helper to evaluate an expression and expect an error
func evaluateExpressionWithError(t *testing.T, interpreter *Interpreter, expr Expr) error {
	_, err := interpreter.evaluate(expr)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}
	return err
}

// Test helper to execute a statement
func executeStatement(t *testing.T, interpreter *Interpreter, stmt Stmt) {
	err := interpreter.execute(stmt)
	if err != nil {
		t.Fatalf("Unexpected error executing statement: %v", err)
	}
}

// Test helper to execute a statement and expect an error
func executeStatementWithError(t *testing.T, interpreter *Interpreter, stmt Stmt) error {
	err := interpreter.execute(stmt)
	if err == nil {
		t.Fatalf("Expected error but got none")
	}
	return err
}

// Test helper to check if two values are equal
func assertEqual(t *testing.T, expected, actual any, testName string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: Expected %v (%T), got %v (%T)", testName, expected, expected, actual, actual)
	}
}

// Test helper to check if a value is truthy
func assertTruthy(t *testing.T, value any, testName string) {
	if !isTruthy(value) {
		t.Errorf("%s: Expected truthy value, got %v", testName, value)
	}
}

// Test helper to check if a value is falsy
func assertFalsy(t *testing.T, value any, testName string) {
	if isTruthy(value) {
		t.Errorf("%s: Expected falsy value, got %v", testName, value)
	}
}

// ============================================================================
// LITERAL EXPRESSION TESTS
// ============================================================================

func TestLiteralExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name:     "Number literal - integer",
			expr:     &LiteralExpr{Value: 42.0},
			expected: 42.0,
		},
		{
			name:     "Number literal - decimal",
			expr:     &LiteralExpr{Value: 3.14},
			expected: 3.14,
		},
		{
			name:     "String literal",
			expr:     &LiteralExpr{Value: "hello"},
			expected: "hello",
		},
		{
			name:     "Boolean literal - true",
			expr:     &LiteralExpr{Value: true},
			expected: true,
		},
		{
			name:     "Boolean literal - false",
			expr:     &LiteralExpr{Value: false},
			expected: false,
		},
		{
			name:     "Nil literal",
			expr:     &LiteralExpr{Value: nil},
			expected: nil,
		},
		{
			name:     "Empty string",
			expr:     &LiteralExpr{Value: ""},
			expected: "",
		},
		{
			name:     "Zero number",
			expr:     &LiteralExpr{Value: 0.0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateExpression(t, interpreter, tt.expr)
			assertEqual(t, tt.expected, result, tt.name)
		})
	}
}

// ============================================================================
// UNARY EXPRESSION TESTS
// ============================================================================

func TestUnaryExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("Negation operator", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "Negate positive number",
				expr:     &UnaryExpr{Operator: createOperatorToken(MINUS, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: -5.0,
			},
			{
				name:     "Negate negative number",
				expr:     &UnaryExpr{Operator: createOperatorToken(MINUS, 1), Right: &LiteralExpr{Value: -3.0}},
				expected: 3.0,
			},
			{
				name:     "Negate zero",
				expr:     &UnaryExpr{Operator: createOperatorToken(MINUS, 1), Right: &LiteralExpr{Value: 0.0}},
				expected: 0.0,
			},
			{
				name:     "Negate decimal",
				expr:     &UnaryExpr{Operator: createOperatorToken(MINUS, 1), Right: &LiteralExpr{Value: 2.5}},
				expected: -2.5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Logical NOT operator", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "NOT true",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: true}},
				expected: false,
			},
			{
				name:     "NOT false",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: false}},
				expected: true,
			},
			{
				name:     "NOT nil",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: nil}},
				expected: true,
			},
			{
				name:     "NOT number (truthy)",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: 42.0}},
				expected: false,
			},
			{
				name:     "NOT zero (truthy)",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: 0.0}},
				expected: false,
			},
			{
				name:     "NOT string (truthy)",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: "hello"}},
				expected: false,
			},
			{
				name:     "NOT empty string (truthy)",
				expr:     &UnaryExpr{Operator: createOperatorToken(BANG, 1), Right: &LiteralExpr{Value: ""}},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Unary expression errors", func(t *testing.T) {
		t.Run("Negate non-number", func(t *testing.T) {
			expr := &UnaryExpr{Operator: createOperatorToken(MINUS, 1), Right: &LiteralExpr{Value: "hello"}}
			err := evaluateExpressionWithError(t, interpreter, expr)

			runtimeErr, ok := err.(RuntimeError)
			if !ok {
				t.Fatalf("Expected RuntimeError, got %T", err)
			}
			if runtimeErr.message != "operand to operator - must be a number" {
				t.Errorf("Expected 'operand to operator - must be a number', got '%s'", runtimeErr.message)
			}
		})
	})
}

// ============================================================================
// BINARY EXPRESSION TESTS
// ============================================================================

func TestBinaryExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("Arithmetic operators", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "Addition - numbers",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 2.0}, Operator: createOperatorToken(PLUS, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: 5.0,
			},
			{
				name:     "Addition - strings",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: "hello"}, Operator: createOperatorToken(PLUS, 1), Right: &LiteralExpr{Value: " world"}},
				expected: "hello world",
			},
			{
				name:     "Subtraction",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(MINUS, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: 2.0,
			},
			{
				name:     "Multiplication",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 4.0}, Operator: createOperatorToken(STAR, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: 12.0,
			},
			{
				name:     "Division",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 15.0}, Operator: createOperatorToken(SLASH, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: 5.0,
			},
			{
				name:     "Division with decimal result",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 7.0}, Operator: createOperatorToken(SLASH, 1), Right: &LiteralExpr{Value: 2.0}},
				expected: 3.5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Comparison operators", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "Greater than - true",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(GREATER, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: true,
			},
			{
				name:     "Greater than - false",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 3.0}, Operator: createOperatorToken(GREATER, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: false,
			},
			{
				name:     "Greater than or equal - equal",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(GREATER_EQUAL, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: true,
			},
			{
				name:     "Less than - true",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 3.0}, Operator: createOperatorToken(LESS, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: true,
			},
			{
				name:     "Less than or equal - equal",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(LESS_EQUAL, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Equality operators", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "Equal - numbers",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(EQUAL_EQUAL, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: true,
			},
			{
				name:     "Equal - different numbers",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(EQUAL_EQUAL, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: false,
			},
			{
				name:     "Equal - strings",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: "hello"}, Operator: createOperatorToken(EQUAL_EQUAL, 1), Right: &LiteralExpr{Value: "hello"}},
				expected: true,
			},
			{
				name:     "Equal - different strings",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: "hello"}, Operator: createOperatorToken(EQUAL_EQUAL, 1), Right: &LiteralExpr{Value: "world"}},
				expected: false,
			},
			{
				name:     "Equal - booleans",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: true}, Operator: createOperatorToken(EQUAL_EQUAL, 1), Right: &LiteralExpr{Value: true}},
				expected: true,
			},
			{
				name:     "Equal - nil values",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: nil}, Operator: createOperatorToken(EQUAL_EQUAL, 1), Right: &LiteralExpr{Value: nil}},
				expected: true,
			},
			{
				name:     "Not equal - numbers",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(BANG_EQUAL, 1), Right: &LiteralExpr{Value: 3.0}},
				expected: true,
			},
			{
				name:     "Not equal - same numbers",
				expr:     &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(BANG_EQUAL, 1), Right: &LiteralExpr{Value: 5.0}},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Binary expression errors", func(t *testing.T) {
		t.Run("Addition - mixed types", func(t *testing.T) {
			expr := &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(PLUS, 1), Right: &LiteralExpr{Value: "hello"}}
			err := evaluateExpressionWithError(t, interpreter, expr)

			runtimeErr, ok := err.(RuntimeError)
			if !ok {
				t.Fatalf("Expected RuntimeError, got %T", err)
			}
			if runtimeErr.message != "operands to operator + must be numbers/strings" {
				t.Errorf("Expected 'operands to operator + must be numbers/strings', got '%s'", runtimeErr.message)
			}
		})

		t.Run("Division by zero", func(t *testing.T) {
			expr := &BinaryExpr{Left: &LiteralExpr{Value: 5.0}, Operator: createOperatorToken(SLASH, 1), Right: &LiteralExpr{Value: 0.0}}
			err := evaluateExpressionWithError(t, interpreter, expr)

			runtimeErr, ok := err.(RuntimeError)
			if !ok {
				t.Fatalf("Expected RuntimeError, got %T", err)
			}
			if runtimeErr.message != "illegal operation: division by zero" {
				t.Errorf("Expected 'illegal operation: division by zero', got '%s'", runtimeErr.message)
			}
		})

		t.Run("Comparison with non-numbers", func(t *testing.T) {
			expr := &BinaryExpr{Left: &LiteralExpr{Value: "hello"}, Operator: createOperatorToken(GREATER, 1), Right: &LiteralExpr{Value: "world"}}
			err := evaluateExpressionWithError(t, interpreter, expr)

			runtimeErr, ok := err.(RuntimeError)
			if !ok {
				t.Fatalf("Expected RuntimeError, got %T", err)
			}
			if runtimeErr.message != "operands to operator > must be numbers" {
				t.Errorf("Expected 'operands to operator > must be numbers', got '%s'", runtimeErr.message)
			}
		})
	})
}

// ============================================================================
// LOGICAL EXPRESSION TESTS
// ============================================================================

func TestLogicalExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("AND operator", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "true AND true",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: true}, Operator: createKeywordToken(AND, 1), Right: &LiteralExpr{Value: true}},
				expected: true,
			},
			{
				name:     "true AND false",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: true}, Operator: createKeywordToken(AND, 1), Right: &LiteralExpr{Value: false}},
				expected: false,
			},
			{
				name:     "false AND true",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: false}, Operator: createKeywordToken(AND, 1), Right: &LiteralExpr{Value: true}},
				expected: false,
			},
			{
				name:     "false AND false",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: false}, Operator: createKeywordToken(AND, 1), Right: &LiteralExpr{Value: false}},
				expected: false,
			},
			{
				name:     "truthy AND truthy",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: 42.0}, Operator: createKeywordToken(AND, 1), Right: &LiteralExpr{Value: "hello"}},
				expected: "hello",
			},
			{
				name:     "falsy AND truthy",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: nil}, Operator: createKeywordToken(AND, 1), Right: &LiteralExpr{Value: "hello"}},
				expected: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("OR operator", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name:     "true OR true",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: true}, Operator: createKeywordToken(OR, 1), Right: &LiteralExpr{Value: true}},
				expected: true,
			},
			{
				name:     "true OR false",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: true}, Operator: createKeywordToken(OR, 1), Right: &LiteralExpr{Value: false}},
				expected: true,
			},
			{
				name:     "false OR true",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: false}, Operator: createKeywordToken(OR, 1), Right: &LiteralExpr{Value: true}},
				expected: true,
			},
			{
				name:     "false OR false",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: false}, Operator: createKeywordToken(OR, 1), Right: &LiteralExpr{Value: false}},
				expected: false,
			},
			{
				name:     "truthy OR truthy",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: 42.0}, Operator: createKeywordToken(OR, 1), Right: &LiteralExpr{Value: "hello"}},
				expected: 42.0,
			},
			{
				name:     "falsy OR truthy",
				expr:     &LogicalExpr{Left: &LiteralExpr{Value: nil}, Operator: createKeywordToken(OR, 1), Right: &LiteralExpr{Value: "hello"}},
				expected: "hello",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})
}

// ============================================================================
// GROUPING EXPRESSION TESTS
// ============================================================================

func TestGroupingExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("Simple grouping", func(t *testing.T) {
		expr := &GroupingExpr{Expression: &LiteralExpr{Value: 42.0}}
		result := evaluateExpression(t, interpreter, expr)
		assertEqual(t, 42.0, result, "Simple grouping")
	})

	t.Run("Nested grouping", func(t *testing.T) {
		inner := &GroupingExpr{Expression: &LiteralExpr{Value: 42.0}}
		outer := &GroupingExpr{Expression: inner}
		result := evaluateExpression(t, interpreter, outer)
		assertEqual(t, 42.0, result, "Nested grouping")
	})

	t.Run("Grouping with binary expression", func(t *testing.T) {
		// (2 + 3) * 4
		grouped := &GroupingExpr{
			Expression: &BinaryExpr{
				Left:     &LiteralExpr{Value: 2.0},
				Operator: createOperatorToken(PLUS, 1),
				Right:    &LiteralExpr{Value: 3.0},
			},
		}
		expr := &BinaryExpr{
			Left:     grouped,
			Operator: createOperatorToken(STAR, 1),
			Right:    &LiteralExpr{Value: 4.0},
		}
		result := evaluateExpression(t, interpreter, expr)
		assertEqual(t, 20.0, result, "Grouping with binary expression")
	})
}

// ============================================================================
// VARIABLE EXPRESSION TESTS
// ============================================================================

func TestVariableExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("Undefined variable", func(t *testing.T) {
		expr := &VariableExpr{variable: createIdentifierToken("undefined", 1)}
		err := evaluateExpressionWithError(t, interpreter, expr)

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Undefined variable 'undefined'" {
			t.Errorf("Expected 'Undefined variable 'undefined'', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Defined variable", func(t *testing.T) {
		// First define the variable
		varStmt := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: 42.0},
		}
		executeStatement(t, interpreter, varStmt)

		// Then evaluate the variable
		expr := &VariableExpr{variable: createIdentifierToken("x", 2)}
		result := evaluateExpression(t, interpreter, expr)
		assertEqual(t, 42.0, result, "Defined variable")
	})

	t.Run("Variable with different types", func(t *testing.T) {
		tests := []struct {
			name    string
			value   any
			varName string
		}{
			{"Number variable", 3.14, "num"},
			{"String variable", "hello", "str"},
			{"Boolean variable", true, "flag"},
			{"Nil variable", nil, "empty"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Define variable
				varStmt := &VarStmt{
					variable:    createIdentifierToken(tt.varName, 1),
					initializer: &LiteralExpr{Value: tt.value},
				}
				executeStatement(t, interpreter, varStmt)

				// Evaluate variable
				expr := &VariableExpr{variable: createIdentifierToken(tt.varName, 2)}
				result := evaluateExpression(t, interpreter, expr)
				assertEqual(t, tt.value, result, tt.name)
			})
		}
	})
}

// ============================================================================
// ASSIGNMENT EXPRESSION TESTS
// ============================================================================

func TestAssignmentExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("Assign to undefined variable", func(t *testing.T) {
		expr := &AssignExpr{
			variable: createIdentifierToken("x", 1),
			value:    &LiteralExpr{Value: 42.0},
		}
		err := evaluateExpressionWithError(t, interpreter, expr)

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Undefined variable 'x'" {
			t.Errorf("Expected 'Undefined variable 'x'', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Assign to defined variable", func(t *testing.T) {
		// First define the variable
		varStmt := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: 10.0},
		}
		executeStatement(t, interpreter, varStmt)

		// Then assign a new value
		expr := &AssignExpr{
			variable: createIdentifierToken("x", 2),
			value:    &LiteralExpr{Value: 42.0},
		}
		result := evaluateExpression(t, interpreter, expr)
		assertEqual(t, 42.0, result, "Assignment result")

		// Verify the variable was updated
		varExpr := &VariableExpr{variable: createIdentifierToken("x", 3)}
		varResult := evaluateExpression(t, interpreter, varExpr)
		assertEqual(t, 42.0, varResult, "Variable after assignment")
	})

	t.Run("Assignment returns the assigned value", func(t *testing.T) {
		// Define variable
		varStmt := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: 10.0},
		}
		executeStatement(t, interpreter, varStmt)

		// Assignment should return the assigned value
		expr := &AssignExpr{
			variable: createIdentifierToken("x", 2),
			value:    &LiteralExpr{Value: 99.0},
		}
		result := evaluateExpression(t, interpreter, expr)
		assertEqual(t, 99.0, result, "Assignment return value")
	})

	t.Run("Chained assignment", func(t *testing.T) {
		// Define variables
		varStmt1 := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: 10.0},
		}
		varStmt2 := &VarStmt{
			variable:    createIdentifierToken("y", 2),
			initializer: &LiteralExpr{Value: 20.0},
		}
		executeStatement(t, interpreter, varStmt1)
		executeStatement(t, interpreter, varStmt2)

		// Chain assignment: x = y = 42
		innerAssign := &AssignExpr{
			variable: createIdentifierToken("y", 3),
			value:    &LiteralExpr{Value: 42.0},
		}
		outerAssign := &AssignExpr{
			variable: createIdentifierToken("x", 3),
			value:    innerAssign,
		}
		result := evaluateExpression(t, interpreter, outerAssign)
		assertEqual(t, 42.0, result, "Chained assignment result")

		// Verify both variables were updated
		xExpr := &VariableExpr{variable: createIdentifierToken("x", 4)}
		yExpr := &VariableExpr{variable: createIdentifierToken("y", 4)}
		xResult := evaluateExpression(t, interpreter, xExpr)
		yResult := evaluateExpression(t, interpreter, yExpr)
		assertEqual(t, 42.0, xResult, "x after chained assignment")
		assertEqual(t, 42.0, yResult, "y after chained assignment")
	})
}

// ============================================================================
// TRUTHINESS TESTS
// ============================================================================

func TestTruthiness(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"nil is falsy", nil, false},
		{"false is falsy", false, false},
		{"true is truthy", true, true},
		{"zero is truthy", 0.0, true},
		{"positive number is truthy", 42.0, true},
		{"negative number is truthy", -1.0, true},
		{"empty string is truthy", "", true},
		{"non-empty string is truthy", "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTruthy(tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v to be %v, got %v", tt.value, tt.expected, result)
			}
		})
	}
}

// ============================================================================
// EQUALITY TESTS
// ============================================================================

func TestEquality(t *testing.T) {
	tests := []struct {
		name     string
		left     any
		right    any
		expected bool
	}{
		{"nil == nil", nil, nil, true},
		{"nil == false", nil, false, false},
		{"true == true", true, true, true},
		{"true == false", true, false, false},
		{"false == false", false, false, true},
		{"1 == 1", 1.0, 1.0, true},
		{"1 == 2", 1.0, 2.0, false},
		{"1.0 == 1", 1.0, 1.0, true},
		{"\"hello\" == \"hello\"", "hello", "hello", true},
		{"\"hello\" == \"world\"", "hello", "world", false},
		{"\"\" == \"\"", "", "", true},
		{"\"\" == nil", "", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEqual(tt.left, tt.right)
			if result != tt.expected {
				t.Errorf("Expected %v == %v to be %v, got %v", tt.left, tt.right, tt.expected, result)
			}
		})
	}
}

// ============================================================================
// STATEMENT EXECUTION TESTS
// ============================================================================

func TestStatementExecution(t *testing.T) {
	t.Run("Variable declaration", func(t *testing.T) {
		interpreter := createTestInterpreter()

		tests := []struct {
			name     string
			stmt     Stmt
			expected any
		}{
			{
				name: "Variable with initializer",
				stmt: &VarStmt{
					variable:    createIdentifierToken("x", 1),
					initializer: &LiteralExpr{Value: 42.0},
				},
				expected: 42.0,
			},
			{
				name: "Variable without initializer",
				stmt: &VarStmt{
					variable:    createIdentifierToken("y", 1),
					initializer: nil,
				},
				expected: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				executeStatement(t, interpreter, tt.stmt)

				// Verify variable was defined
				varStmt, ok := tt.stmt.(*VarStmt)
				if !ok {
					t.Fatalf("Expected VarStmt")
				}
				expr := &VariableExpr{variable: varStmt.variable}
				result := evaluateExpression(t, interpreter, expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Expression statement", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Expression statements don't produce a value, just execute the expression
		stmt := &ExpressionStmt{
			expression: &BinaryExpr{
				Left:     &LiteralExpr{Value: 2.0},
				Operator: createOperatorToken(PLUS, 1),
				Right:    &LiteralExpr{Value: 3.0},
			},
		}
		executeStatement(t, interpreter, stmt)
		// No assertion needed - just ensuring no error occurs
	})

	t.Run("Print statement", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Print statements don't return a value, they output to stdout
		// We can't easily test the output in unit tests, so we just ensure no error
		stmt := &PrintStmt{
			expression: &LiteralExpr{Value: "Hello, World!"},
		}
		executeStatement(t, interpreter, stmt)
		// No assertion needed - just ensuring no error occurs
	})

	t.Run("If statement", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define a variable to track which branch was executed
		varStmt := &VarStmt{
			variable:    createIdentifierToken("result", 1),
			initializer: &LiteralExpr{Value: "none"},
		}
		executeStatement(t, interpreter, varStmt)

		t.Run("True condition - then branch", func(t *testing.T) {
			// Reset result
			assignStmt := &ExpressionStmt{
				expression: &AssignExpr{
					variable: createIdentifierToken("result", 2),
					value:    &LiteralExpr{Value: "none"},
				},
			}
			executeStatement(t, interpreter, assignStmt)

			// If statement with true condition
			ifStmt := &IfStmt{
				condition: &LiteralExpr{Value: true},
				thenBranch: &ExpressionStmt{
					expression: &AssignExpr{
						variable: createIdentifierToken("result", 3),
						value:    &LiteralExpr{Value: "then"},
					},
				},
				elseBranch: nil,
			}
			executeStatement(t, interpreter, ifStmt)

			// Check result
			expr := &VariableExpr{variable: createIdentifierToken("result", 4)}
			result := evaluateExpression(t, interpreter, expr)
			assertEqual(t, "then", result, "If true condition")
		})

		t.Run("False condition - else branch", func(t *testing.T) {
			// Reset result
			assignStmt := &ExpressionStmt{
				expression: &AssignExpr{
					variable: createIdentifierToken("result", 5),
					value:    &LiteralExpr{Value: "none"},
				},
			}
			executeStatement(t, interpreter, assignStmt)

			// If statement with false condition and else branch
			ifStmt := &IfStmt{
				condition: &LiteralExpr{Value: false},
				thenBranch: &ExpressionStmt{
					expression: &AssignExpr{
						variable: createIdentifierToken("result", 6),
						value:    &LiteralExpr{Value: "then"},
					},
				},
				elseBranch: &ExpressionStmt{
					expression: &AssignExpr{
						variable: createIdentifierToken("result", 6),
						value:    &LiteralExpr{Value: "else"},
					},
				},
			}
			executeStatement(t, interpreter, ifStmt)

			// Check result
			expr := &VariableExpr{variable: createIdentifierToken("result", 7)}
			result := evaluateExpression(t, interpreter, expr)
			assertEqual(t, "else", result, "If false condition with else")
		})

		t.Run("False condition - no else branch", func(t *testing.T) {
			// Reset result
			assignStmt := &ExpressionStmt{
				expression: &AssignExpr{
					variable: createIdentifierToken("result", 8),
					value:    &LiteralExpr{Value: "none"},
				},
			}
			executeStatement(t, interpreter, assignStmt)

			// If statement with false condition and no else branch
			ifStmt := &IfStmt{
				condition: &LiteralExpr{Value: false},
				thenBranch: &ExpressionStmt{
					expression: &AssignExpr{
						variable: createIdentifierToken("result", 9),
						value:    &LiteralExpr{Value: "then"},
					},
				},
				elseBranch: nil,
			}
			executeStatement(t, interpreter, ifStmt)

			// Check result should remain unchanged
			expr := &VariableExpr{variable: createIdentifierToken("result", 10)}
			result := evaluateExpression(t, interpreter, expr)
			assertEqual(t, "none", result, "If false condition without else")
		})
	})

	t.Run("While statement", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define counter variable
		varStmt := &VarStmt{
			variable:    createIdentifierToken("counter", 1),
			initializer: &LiteralExpr{Value: 0.0},
		}
		executeStatement(t, interpreter, varStmt)

		// While loop that increments counter 3 times
		whileStmt := &WhileStmt{
			condition: &BinaryExpr{
				Left:     &VariableExpr{variable: createIdentifierToken("counter", 2)},
				Operator: createOperatorToken(LESS, 2),
				Right:    &LiteralExpr{Value: 3.0},
			},
			body: &ExpressionStmt{
				expression: &AssignExpr{
					variable: createIdentifierToken("counter", 2),
					value: &BinaryExpr{
						Left:     &VariableExpr{variable: createIdentifierToken("counter", 2)},
						Operator: createOperatorToken(PLUS, 2),
						Right:    &LiteralExpr{Value: 1.0},
					},
				},
			},
		}
		executeStatement(t, interpreter, whileStmt)

		// Check final counter value
		expr := &VariableExpr{variable: createIdentifierToken("counter", 3)}
		result := evaluateExpression(t, interpreter, expr)
		assertEqual(t, 3.0, result, "While loop counter")
	})

	t.Run("Block statement", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define outer variable
		outerVarStmt := &VarStmt{
			variable:    createIdentifierToken("outer", 1),
			initializer: &LiteralExpr{Value: "outer"},
		}
		executeStatement(t, interpreter, outerVarStmt)

		// Block with inner variable and assignment
		blockStmt := &BlockStmt{
			statements: []Stmt{
				&VarStmt{
					variable:    createIdentifierToken("inner", 2),
					initializer: &LiteralExpr{Value: "inner"},
				},
				&ExpressionStmt{
					expression: &AssignExpr{
						variable: createIdentifierToken("outer", 3),
						value:    &LiteralExpr{Value: "modified"},
					},
				},
			},
		}
		executeStatement(t, interpreter, blockStmt)

		// Check that outer variable was modified
		outerExpr := &VariableExpr{variable: createIdentifierToken("outer", 4)}
		outerResult := evaluateExpression(t, interpreter, outerExpr)
		assertEqual(t, "modified", outerResult, "Outer variable after block")

		// Check that inner variable is not accessible (should cause error)
		innerExpr := &VariableExpr{variable: createIdentifierToken("inner", 5)}
		err := evaluateExpressionWithError(t, interpreter, innerExpr)

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Undefined variable 'inner'" {
			t.Errorf("Expected 'Undefined variable 'inner'', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Nested blocks", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define variable in outer scope
		outerVarStmt := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: "outer"},
		}
		executeStatement(t, interpreter, outerVarStmt)

		// Nested blocks with variable shadowing
		innerBlock := &BlockStmt{
			statements: []Stmt{
				&VarStmt{
					variable:    createIdentifierToken("x", 3),
					initializer: &LiteralExpr{Value: "inner"},
				},
			},
		}
		outerBlock := &BlockStmt{
			statements: []Stmt{
				&VarStmt{
					variable:    createIdentifierToken("x", 2),
					initializer: &LiteralExpr{Value: "middle"},
				},
				innerBlock,
			},
		}
		executeStatement(t, interpreter, outerBlock)

		// Check that outer variable is still accessible and unchanged
		outerExpr := &VariableExpr{variable: createIdentifierToken("x", 4)}
		outerResult := evaluateExpression(t, interpreter, outerExpr)
		assertEqual(t, "outer", outerResult, "Outer variable after nested blocks")
	})
}

// ============================================================================
// COMPLEX EXPRESSION TESTS
// ============================================================================

func TestComplexExpressionEvaluation(t *testing.T) {
	interpreter := createTestInterpreter()

	t.Run("Operator precedence", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name: "Multiplication before addition",
				expr: &BinaryExpr{
					Left:     &LiteralExpr{Value: 2.0},
					Operator: createOperatorToken(PLUS, 1),
					Right: &BinaryExpr{
						Left:     &LiteralExpr{Value: 3.0},
						Operator: createOperatorToken(STAR, 1),
						Right:    &LiteralExpr{Value: 4.0},
					},
				},
				expected: 14.0, // 2 + (3 * 4) = 2 + 12 = 14
			},
			{
				name: "Division before subtraction",
				expr: &BinaryExpr{
					Left: &BinaryExpr{
						Left:     &LiteralExpr{Value: 10.0},
						Operator: createOperatorToken(SLASH, 1),
						Right:    &LiteralExpr{Value: 2.0},
					},
					Operator: createOperatorToken(MINUS, 1),
					Right:    &LiteralExpr{Value: 1.0},
				},
				expected: 4.0, // (10 / 2) - 1 = 5 - 1 = 4
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Complex logical expressions", func(t *testing.T) {
		tests := []struct {
			name     string
			expr     Expr
			expected any
		}{
			{
				name: "AND with short-circuit",
				expr: &LogicalExpr{
					Left:     &LiteralExpr{Value: false},
					Operator: createKeywordToken(AND, 1),
					Right:    &LiteralExpr{Value: "should not evaluate"},
				},
				expected: false,
			},
			{
				name: "OR with short-circuit",
				expr: &LogicalExpr{
					Left:     &LiteralExpr{Value: true},
					Operator: createKeywordToken(OR, 1),
					Right:    &LiteralExpr{Value: "should not evaluate"},
				},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := evaluateExpression(t, interpreter, tt.expr)
				assertEqual(t, tt.expected, result, tt.name)
			})
		}
	})

	t.Run("Mixed type operations", func(t *testing.T) {
		t.Run("String concatenation", func(t *testing.T) {
			expr := &BinaryExpr{
				Left:     &LiteralExpr{Value: "Hello"},
				Operator: createOperatorToken(PLUS, 1),
				Right:    &LiteralExpr{Value: " World"},
			}
			result := evaluateExpression(t, interpreter, expr)
			assertEqual(t, "Hello World", result, "String concatenation")
		})

		t.Run("Number to string concatenation", func(t *testing.T) {
			expr := &BinaryExpr{
				Left:     &LiteralExpr{Value: "The answer is "},
				Operator: createOperatorToken(PLUS, 1),
				Right:    &LiteralExpr{Value: 42.0},
			}
			err := evaluateExpressionWithError(t, interpreter, expr)

			runtimeErr, ok := err.(RuntimeError)
			if !ok {
				t.Fatalf("Expected RuntimeError, got %T", err)
			}
			if runtimeErr.message != "operands to operator + must be numbers/strings" {
				t.Errorf("Expected 'operands to operator + must be numbers/strings', got '%s'", runtimeErr.message)
			}
		})
	})
}

// ============================================================================
// ERROR HANDLING TESTS
// ============================================================================

func TestErrorHandling(t *testing.T) {
	t.Run("Runtime error propagation", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Test that runtime errors in expressions are properly propagated
		expr := &BinaryExpr{
			Left:     &LiteralExpr{Value: 5.0},
			Operator: createOperatorToken(SLASH, 1),
			Right:    &LiteralExpr{Value: 0.0},
		}
		err := evaluateExpressionWithError(t, interpreter, expr)

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "illegal operation: division by zero" {
			t.Errorf("Expected 'illegal operation: division by zero', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Error in statement execution", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Test that runtime errors in statements are properly propagated
		stmt := &ExpressionStmt{
			expression: &BinaryExpr{
				Left:     &LiteralExpr{Value: "hello"},
				Operator: createOperatorToken(MINUS, 1),
				Right:    &LiteralExpr{Value: "world"},
			},
		}
		err := executeStatementWithError(t, interpreter, stmt)

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "operands to operator - must be numbers" {
			t.Errorf("Expected 'operands to operator - must be numbers', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Error in variable assignment", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define variable first
		varStmt := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: 10.0},
		}
		executeStatement(t, interpreter, varStmt)

		// Try to assign with an error in the expression
		assignStmt := &ExpressionStmt{
			expression: &AssignExpr{
				variable: createIdentifierToken("x", 2),
				value: &BinaryExpr{
					Left:     &LiteralExpr{Value: "hello"},
					Operator: createOperatorToken(MINUS, 2),
					Right:    &LiteralExpr{Value: "world"},
				},
			},
		}
		err := executeStatementWithError(t, interpreter, assignStmt)

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "operands to operator - must be numbers" {
			t.Errorf("Expected 'operands to operator - must be numbers', got '%s'", runtimeErr.message)
		}
	})
}

// ============================================================================
// FUNCTION TESTS
// ============================================================================
// These tests focus on internal mechanics and error handling.
// Higher-level function behavior is tested in integration_test.go

func TestFunctionDeclaration(t *testing.T) {
	// Test internal mechanics: function is stored as LoxFunction with correct properties
	t.Run("Function declaration stores LoxFunction in environment", func(t *testing.T) {
		interpreter := createTestInterpreter()

		funcStmt := &FunctionStmt{
			functionName: createIdentifierToken("add", 1),
			params: []Token{
				createIdentifierToken("a", 1),
				createIdentifierToken("b", 1),
			},
			body: []Stmt{
				&ReturnStmt{
					keyword:     createKeywordToken(RETURN, 1),
					returnValue: &LiteralExpr{Value: 42.0},
				},
			},
		}

		executeStatement(t, interpreter, funcStmt)

		// Verify function is stored and has correct arity
		funcExpr := &VariableExpr{variable: createIdentifierToken("add", 2)}
		result := evaluateExpression(t, interpreter, funcExpr)

		funcImpl, ok := result.(*LoxFunction)
		if !ok {
			t.Fatalf("Expected LoxFunction, got %T", result)
		}

		if funcImpl.arity() != 2 {
			t.Errorf("Expected arity 2, got %d", funcImpl.arity())
		}
	})
}

func TestFunctionCallErrors(t *testing.T) {
	t.Run("Call non-callable value", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define a variable that's not a function
		varStmt := &VarStmt{
			variable:    createIdentifierToken("x", 1),
			initializer: &LiteralExpr{Value: 42.0},
		}
		executeStatement(t, interpreter, varStmt)

		// Try to call it
		callExpr := &CallExpr{
			Callee:    &VariableExpr{variable: createIdentifierToken("x", 2)},
			Paren:     createOperatorToken(RIGHT_PAREN, 2),
			Arguments: []Expr{},
		}

		err := evaluateExpressionWithError(t, interpreter, callExpr)
		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Can only call functions and classes." {
			t.Errorf("Expected 'Can only call functions and classes.', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Call function with wrong number of arguments - too few", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define a function that requires 2 arguments
		funcStmt := &FunctionStmt{
			functionName: createIdentifierToken("add", 1),
			params: []Token{
				createIdentifierToken("a", 1),
				createIdentifierToken("b", 1),
			},
			body: []Stmt{
				&ReturnStmt{
					keyword: createKeywordToken(RETURN, 1),
					returnValue: &BinaryExpr{
						Left:     &VariableExpr{variable: createIdentifierToken("a", 1)},
						Operator: createOperatorToken(PLUS, 1),
						Right:    &VariableExpr{variable: createIdentifierToken("b", 1)},
					},
				},
			},
		}
		executeStatement(t, interpreter, funcStmt)

		// Call with only 1 argument
		callExpr := &CallExpr{
			Callee:    &VariableExpr{variable: createIdentifierToken("add", 2)},
			Paren:     createOperatorToken(RIGHT_PAREN, 2),
			Arguments: []Expr{&LiteralExpr{Value: 1.0}},
		}

		err := evaluateExpressionWithError(t, interpreter, callExpr)
		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Expected 2 arguments but got 1" {
			t.Errorf("Expected 'Expected 2 arguments but got 1', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Call function with wrong number of arguments - too many", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Define a function that requires 1 argument
		funcStmt := &FunctionStmt{
			functionName: createIdentifierToken("double", 1),
			params:       []Token{createIdentifierToken("x", 1)},
			body: []Stmt{
				&ReturnStmt{
					keyword: createKeywordToken(RETURN, 1),
					returnValue: &BinaryExpr{
						Left:     &VariableExpr{variable: createIdentifierToken("x", 1)},
						Operator: createOperatorToken(STAR, 1),
						Right:    &LiteralExpr{Value: 2.0},
					},
				},
			},
		}
		executeStatement(t, interpreter, funcStmt)

		// Call with 2 arguments
		callExpr := &CallExpr{
			Callee: &VariableExpr{variable: createIdentifierToken("double", 2)},
			Paren:  createOperatorToken(RIGHT_PAREN, 2),
			Arguments: []Expr{
				&LiteralExpr{Value: 5.0},
				&LiteralExpr{Value: 10.0},
			},
		}

		err := evaluateExpressionWithError(t, interpreter, callExpr)
		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Expected 1 arguments but got 2" {
			t.Errorf("Expected 'Expected 1 arguments but got 2', got '%s'", runtimeErr.message)
		}
	})
}

func TestBuiltinFunction(t *testing.T) {
	t.Run("Call built-in clock function", func(t *testing.T) {
		interpreter := createTestInterpreter()

		// Call clock() - it should return a number
		callExpr := &CallExpr{
			Callee:    &VariableExpr{variable: createIdentifierToken("clock", 1)},
			Paren:     createOperatorToken(RIGHT_PAREN, 1),
			Arguments: []Expr{},
		}

		result := evaluateExpression(t, interpreter, callExpr)

		// Verify it returns a number (float64)
		_, ok := result.(float64)
		if !ok {
			t.Errorf("Expected float64 result from clock(), got %T", result)
		}
	})
}

// Closure behavior is extensively tested in integration_test.go
