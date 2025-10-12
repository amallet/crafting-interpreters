package main

import (
	"testing"
)

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestParserErrorRecovery(t *testing.T) {
	// Test that parser reports errors for invalid syntax
	input := "var x = ;"
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	parser := NewParser(lox, tokens)
	statements, err := parser.parse()

	// Should have an error
	if !lox.hadError {
		t.Error("Expected parse error for invalid syntax")
	}

	// Should return an error
	if err == nil {
		t.Error("Expected parse error to be returned")
	}

	// Should not produce any statements due to error
	if len(statements) != 0 {
		t.Errorf("Expected 0 statements due to error, got %d", len(statements))
	}
}

func TestParserInvalidSyntax(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		// Variable declaration errors
		{
			name:        "Missing variable name",
			input:       "var = 123;",
			expectError: true,
			description: "Variable declaration without identifier",
		},
		{
			name:        "Missing assignment value",
			input:       "var x = ;",
			expectError: true,
			description: "Variable declaration without value",
		},
		{
			name:        "Missing semicolon after var",
			input:       "var x = 123",
			expectError: true,
			description: "Variable declaration without semicolon",
		},

		// Expression errors
		{
			name:        "Missing right operand",
			input:       "1 + ;",
			expectError: true,
			description: "Binary expression missing right operand",
		},
		{
			name:        "Missing left operand",
			input:       "+ 2;",
			expectError: true,
			description: "Binary expression missing left operand",
		},
		{
			name:        "Missing unary operand",
			input:       "!;",
			expectError: true,
			description: "Unary expression missing operand",
		},
		{
			name:        "Missing unary operand minus",
			input:       "-;",
			expectError: true,
			description: "Unary minus missing operand",
		},

		// Grouping errors
		{
			name:        "Unclosed parenthesis",
			input:       "(1 + 2;",
			expectError: true,
			description: "Unclosed grouping expression",
		},
		{
			name:        "Empty parentheses",
			input:       "();",
			expectError: true,
			description: "Empty grouping expression",
		},
		{
			name:        "Unopened parenthesis",
			input:       "1 + 2);",
			expectError: true,
			description: "Unopened closing parenthesis",
		},

		// Statement errors
		{
			name:        "Missing print expression",
			input:       "print ;",
			expectError: true,
			description: "Print statement without expression",
		},
		{
			name:        "Missing print semicolon",
			input:       "print 123",
			expectError: true,
			description: "Print statement without semicolon",
		},

		// If statement errors
		{
			name:        "Missing if condition",
			input:       "if () print 1;",
			expectError: true,
			description: "If statement without condition",
		},
		{
			name:        "Missing if opening parenthesis",
			input:       "if true print 1;",
			expectError: true,
			description: "If statement without opening parenthesis",
		},
		{
			name:        "Missing if closing parenthesis",
			input:       "if (true print 1;",
			expectError: true,
			description: "If statement without closing parenthesis",
		},
		{
			name:        "Missing if body",
			input:       "if (true) ;",
			expectError: true,
			description: "If statement without body",
		},

		// While statement errors
		{
			name:        "Missing while condition",
			input:       "while () print 1;",
			expectError: true,
			description: "While statement without condition",
		},
		{
			name:        "Missing while opening parenthesis",
			input:       "while true print 1;",
			expectError: true,
			description: "While statement without opening parenthesis",
		},
		{
			name:        "Missing while closing parenthesis",
			input:       "while (true print 1;",
			expectError: true,
			description: "While statement without closing parenthesis",
		},
		{
			name:        "Missing while body",
			input:       "while (true) ;",
			expectError: true,
			description: "While statement without body",
		},

		// Block statement errors
		{
			name:        "Unclosed block",
			input:       "{ var x = 1;",
			expectError: true,
			description: "Block statement without closing brace",
		},
		{
			name:        "Unopened block",
			input:       "var x = 1; }",
			expectError: true,
			description: "Closing brace without opening brace",
		},

		// Assignment errors
		{
			name:        "Assignment to literal",
			input:       "123 = 456;",
			expectError: true,
			description: "Assignment to literal value",
		},
		{
			name:        "Assignment without value",
			input:       "x = ;",
			expectError: true,
			description: "Assignment without value",
		},
		{
			name:        "Assignment without target",
			input:       "= 123;",
			expectError: true,
			description: "Assignment without target variable",
		},

		// Logical expression errors
		{
			name:        "Missing right operand for AND",
			input:       "true and ;",
			expectError: true,
			description: "Logical AND missing right operand",
		},
		{
			name:        "Missing left operand for AND",
			input:       "and false;",
			expectError: true,
			description: "Logical AND missing left operand",
		},
		{
			name:        "Missing right operand for OR",
			input:       "true or ;",
			expectError: true,
			description: "Logical OR missing right operand",
		},
		{
			name:        "Missing left operand for OR",
			input:       "or false;",
			expectError: true,
			description: "Logical OR missing left operand",
		},

		// Complex error cases
		{
			name:        "Multiple errors in one statement",
			input:       "var = + ;",
			expectError: true,
			description: "Multiple syntax errors in one statement",
		},
		{
			name:        "Nested grouping errors",
			input:       "((1 + 2;",
			expectError: true,
			description: "Nested unclosed grouping",
		},
		{
			name:        "Mixed valid and invalid statements",
			input:       "var x = 1; var = 2; var y = 3;",
			expectError: true,
			description: "Mixed valid and invalid statements",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			parser := NewParser(lox, tokens)
			statements, err := parser.parse()

			if test.expectError {
				// Should have an error
				if !lox.hadError {
					t.Errorf("Expected parse error for: %s", test.description)
				}

				// Should return an error
				if err == nil {
					t.Errorf("Expected parse error to be returned for: %s", test.description)
				}

				// Should not produce any statements due to error
				if len(statements) != 0 {
					t.Errorf("Expected 0 statements due to error, got %d for: %s", len(statements), test.description)
				}

				// Should have at least one error message
				if len(lox.errors) == 0 {
					t.Errorf("Expected at least one error message for: %s", test.description)
				}
			} else {
				// Should not have an error
				if lox.hadError {
					t.Errorf("Unexpected parse error for: %s. Errors: %v", test.description, lox.errors)
				}

				// Should not return an error
				if err != nil {
					t.Errorf("Unexpected parse error returned for: %s. Error: %v", test.description, err)
				}
			}
		})
	}
}

func TestParserErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedErrors []string // Substrings that should appear in error messages
	}{
		{
			name:  "Variable declaration without identifier",
			input: "var = 123;",
			expectedErrors: []string{
				"Expect variable name",
			},
		},
		{
			name:  "Missing assignment value",
			input: "var x = ;",
			expectedErrors: []string{
				"Expected expression",
			},
		},
		{
			name:  "Missing semicolon",
			input: "var x = 123",
			expectedErrors: []string{
				"Expect ';' after variable declaration",
			},
		},
		{
			name:  "Unclosed parenthesis",
			input: "(1 + 2;",
			expectedErrors: []string{
				"Expect ')' after expression",
			},
		},
		{
			name:  "Missing print expression",
			input: "print ;",
			expectedErrors: []string{
				"Expected expression",
			},
		},
		{
			name:  "Missing if condition",
			input: "if () print 1;",
			expectedErrors: []string{
				"Expected expression",
			},
		},
		{
			name:  "Missing if opening parenthesis",
			input: "if true print 1;",
			expectedErrors: []string{
				"Expect '(' after 'if",
			},
		},
		{
			name:  "Missing while opening parenthesis",
			input: "while true print 1;",
			expectedErrors: []string{
				"Expect '(' after while condition",
			},
		},
		{
			name:  "Assignment to literal",
			input: "123 = 456;",
			expectedErrors: []string{
				"Invalid assignment target",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			parser := NewParser(lox, tokens)
			_, _ = parser.parse()

			// Should have an error
			if !lox.hadError {
				t.Errorf("Expected parse error for: %s", test.name)
				return
			}

			// Check that expected error messages are present
			for _, expectedError := range test.expectedErrors {
				found := false
				for _, actualError := range lox.errors {
					if contains(actualError, expectedError) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message containing '%s' not found. Actual errors: %v", expectedError, lox.errors)
				}
			}
		})
	}
}
