package main

import (
	"fmt"
	"testing"
)

func TestScannerEmptyInput(t *testing.T) {
	lox := NewTestGLox()
	scanner := NewScanner(lox, "")
	tokens := scanner.scanTokens()

	expected := []Token{createEOFToken(1)}
	compareTokens(t, expected, tokens, "Empty input")
}

func TestScannerSingleCharacterTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
		lexeme   string
	}{
		{"(", LEFT_PAREN, "("},
		{")", RIGHT_PAREN, ")"},
		{"{", LEFT_BRACE, "{"},
		{"}", RIGHT_BRACE, "}"},
		{",", COMMA, ","},
		{".", DOT, "."},
		{"-", MINUS, "-"},
		{"+", PLUS, "+"},
		{";", SEMICOLON, ";"},
		{"/", SLASH, "/"},
		{"*", STAR, "*"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			expected := []Token{
				createToken(test.expected, test.lexeme, nil, 1),
				createEOFToken(1),
			}
			compareTokens(t, expected, tokens, test.input)
		})
	}
}

func TestScannerTwoCharacterTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
		lexeme   string
	}{
		{"!=", BANG_EQUAL, "!="},
		{"==", EQUAL_EQUAL, "=="},
		{">=", GREATER_EQUAL, ">="},
		{"<=", LESS_EQUAL, "<="},
		{"!", BANG, "!"},
		{"=", EQUAL, "="},
		{">", GREATER, ">"},
		{"<", LESS, "<"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			expected := []Token{
				createToken(test.expected, test.lexeme, nil, 1),
				createEOFToken(1),
			}
			compareTokens(t, expected, tokens, test.input)
		})
	}
}

func TestScannerNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"0", 0.0},
		{"123", 123.0},
		{"123.45", 123.45},
		{"0.5", 0.5},
		{"42.0", 42.0},
		{"999.999", 999.999},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			// Get non-EOF tokens
			nonEOFTokens := getNonEOFTokens(tokens)
			if len(nonEOFTokens) != 1 {
				t.Errorf("Expected 1 non-EOF token, got %d", len(nonEOFTokens))
				return
			}

			token := nonEOFTokens[0]
			if token.token_type != NUMBER {
				t.Errorf("Expected NUMBER token, got %v", token.token_type)
			}

			if token.literal != test.expected {
				t.Errorf("Expected literal %v, got %v", test.expected, token.literal)
			}
		})
	}
}

func TestScannerNegativeNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			"-0",
			[]Token{
				createOperatorToken(MINUS, 1),
				createNumberToken(0.0, 1),
			},
		},
		{
			"-123",
			[]Token{
				createOperatorToken(MINUS, 1),
				createNumberToken(123.0, 1),
			},
		},
		{
			"-123.45",
			[]Token{
				createOperatorToken(MINUS, 1),
				createNumberToken(123.45, 1),
			},
		},
		{
			"-0.5",
			[]Token{
				createOperatorToken(MINUS, 1),
				createNumberToken(0.5, 1),
			},
		},
		{
			"-42.0",
			[]Token{
				createOperatorToken(MINUS, 1),
				createNumberToken(42.0, 1),
			},
		},
		{
			"-999.999",
			[]Token{
				createOperatorToken(MINUS, 1),
				createNumberToken(999.999, 1),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			// Get non-EOF tokens
			nonEOFTokens := getNonEOFTokens(tokens)
			if len(nonEOFTokens) != len(test.expected) {
				t.Errorf("Expected %d non-EOF tokens, got %d", len(test.expected), len(nonEOFTokens))
				return
			}

			// Check each token
			for i, expectedToken := range test.expected {
				actualToken := nonEOFTokens[i]

				// Check token type
				if actualToken.token_type != expectedToken.token_type {
					t.Errorf("Token %d type mismatch: expected %v, got %v",
						i, expectedToken.token_type, actualToken.token_type)
					continue
				}

				// Check line number
				if actualToken.line != expectedToken.line {
					t.Errorf("Token %d line mismatch: expected %d, got %d",
						i, expectedToken.line, actualToken.line)
					continue
				}

				// Check literal value (for numbers)
				if expectedToken.token_type == NUMBER {
					if actualToken.literal != expectedToken.literal {
						t.Errorf("Token %d literal mismatch: expected %v, got %v",
							i, expectedToken.literal, actualToken.literal)
					}
				}
			}
		})
	}
}

func TestScannerStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		lexeme   string
	}{
		{`"hello"`, "hello", `"hello"`},
		{`"world"`, "world", `"world"`},
		{`""`, "", `""`},
		{`"hello world"`, "hello world", `"hello world"`},
		{`"123"`, "123", `"123"`},
		{`"special chars !@#$%^&*()"`, "special chars !@#$%^&*()", `"special chars !@#$%^&*()"`},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			expected := []Token{
				createStringToken(test.expected, 1),
				createEOFToken(1),
			}
			compareTokens(t, expected, tokens, test.input)
		})
	}
}

func TestScannerMultiLineStrings(t *testing.T) {
	input := `"hello
world"`
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	expected := []Token{
		createStringToken("hello\nworld", 2), // Line number should be 2
		createEOFToken(2),
	}
	compareTokens(t, expected, tokens, "Multi-line string")
}

func TestScannerIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"world", "world"},
		{"_underscore", "_underscore"},
		{"camelCase", "camelCase"},
		{"snake_case", "snake_case"},
		{"var123", "var123"},
		{"a", "a"},
		{"A", "A"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			expected := []Token{
				createIdentifierToken(test.expected, 1),
				createEOFToken(1),
			}
			compareTokens(t, expected, tokens, test.input)
		})
	}
}

func TestScannerKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"var", VAR},
		{"print", PRINT},
		{"if", IF},
		{"else", ELSE},
		{"while", WHILE},
		{"for", FOR},
		{"true", TRUE},
		{"false", FALSE},
		{"nil", NIL},
		{"and", AND},
		{"or", OR},
		{"class", CLASS},
		{"fun", FUN},
		{"return", RETURN},
		{"super", SUPER},
		{"this", THIS},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			expected := []Token{
				createKeywordToken(test.expected, 1),
				createEOFToken(1),
			}
			compareTokens(t, expected, tokens, test.input)
		})
	}
}

func TestScannerComments(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			"// This is a comment",
			[]Token{createEOFToken(1)},
		},
		{
			"123 // This is a comment",
			[]Token{
				createNumberToken(123.0, 1),
				createEOFToken(1),
			},
		},
		{
			"// Comment\n123",
			[]Token{
				createNumberToken(123.0, 2),
				createEOFToken(2),
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("comment_%d", i), func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			compareTokens(t, test.expected, tokens, test.input)
		})
	}
}

func TestScannerWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			"   ",
			[]Token{createEOFToken(1)},
		},
		{
			"\t\t",
			[]Token{createEOFToken(1)},
		},
		{
			"\r\r",
			[]Token{createEOFToken(1)},
		},
		{
			" \t \r ",
			[]Token{createEOFToken(1)},
		},
		{
			"123   456",
			[]Token{
				createNumberToken(123.0, 1),
				createNumberToken(456.0, 1),
				createEOFToken(1),
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("whitespace_%d", i), func(t *testing.T) {
			lox := NewTestGLox()
			scanner := NewScanner(lox, test.input)
			tokens := scanner.scanTokens()

			compareTokens(t, test.expected, tokens, test.input)
		})
	}
}

func TestScannerNewlines(t *testing.T) {
	input := "123\n456\n789"
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	expected := []Token{
		createNumberToken(123.0, 1),
		createNumberToken(456.0, 2),
		createNumberToken(789.0, 3),
		createEOFToken(3),
	}
	compareTokens(t, expected, tokens, "Newlines")
}

func TestScannerComplexExpression(t *testing.T) {
	input := "var x = 123 + 456;"
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	expected := []Token{
		createKeywordToken(VAR, 1),
		createIdentifierToken("x", 1),
		createOperatorToken(EQUAL, 1),
		createNumberToken(123.0, 1),
		createOperatorToken(PLUS, 1),
		createNumberToken(456.0, 1),
		createOperatorToken(SEMICOLON, 1),
		createEOFToken(1),
	}
	compareTokens(t, expected, tokens, "Complex expression")
}

func TestScannerUnterminatedString(t *testing.T) {
	input := `"hello world`
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	// Should have an error
	if !lox.hadError {
		t.Error("Expected error for unterminated string")
	}

	if len(lox.errors) == 0 {
		t.Error("Expected error message for unterminated string")
	}

	// Should still produce EOF token
	expected := []Token{createEOFToken(1)}
	compareTokens(t, expected, tokens, "Unterminated string")
}

func TestScannerUnexpectedCharacter(t *testing.T) {
	input := "@"
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	// Should have an error
	if !lox.hadError {
		t.Error("Expected error for unexpected character")
	}

	if len(lox.errors) == 0 {
		t.Error("Expected error message for unexpected character")
	}

	// Should still produce EOF token
	expected := []Token{createEOFToken(1)}
	compareTokens(t, expected, tokens, "Unexpected character")
}

func TestScannerMixedContent(t *testing.T) {
	input := `var name = "John";
print name + " Doe";`
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	expected := []Token{
		createKeywordToken(VAR, 1),
		createIdentifierToken("name", 1),
		createOperatorToken(EQUAL, 1),
		createStringToken("John", 1),
		createOperatorToken(SEMICOLON, 1),
		createKeywordToken(PRINT, 2),
		createIdentifierToken("name", 2),
		createOperatorToken(PLUS, 2),
		createStringToken(" Doe", 2),
		createOperatorToken(SEMICOLON, 2),
		createEOFToken(2),
	}
	compareTokens(t, expected, tokens, "Mixed content")
}

func TestScannerLineNumberTracking(t *testing.T) {
	input := `123
456
789`
	lox := NewTestGLox()
	scanner := NewScanner(lox, input)
	tokens := scanner.scanTokens()

	expected := []Token{
		createNumberToken(123.0, 1),
		createNumberToken(456.0, 2),
		createNumberToken(789.0, 3),
		createEOFToken(3),
	}
	compareTokens(t, expected, tokens, "Line number tracking")
}
