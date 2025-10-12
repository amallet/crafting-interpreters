package main

import (
	"fmt"
	"reflect"
	"testing"
)

// TestGLox is a test-friendly version of GLox that captures errors instead of printing them
type TestGLox struct {
	hadError        bool
	hadRuntimeError bool
	errors          []string
	runtimeErrors   []string
}

func NewTestGLox() *TestGLox {
	return &TestGLox{
		errors:        make([]string, 0),
		runtimeErrors: make([]string, 0),
	}
}

func (l *TestGLox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *TestGLox) parseError(token Token, message string) {
	if token.token_type == EOF {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, fmt.Sprintf(" at %s ", token.lexeme), message)
	}
}

func (l *TestGLox) runtimeError(err error) {
	runtime_err, _ := err.(RuntimeError)
	errorMsg := fmt.Sprintf("[line %d] %s", runtime_err.token.line, runtime_err.Error())
	l.runtimeErrors = append(l.runtimeErrors, errorMsg)
	l.hadRuntimeError = true
}

func (l *TestGLox) report(line int, where string, message string) {
	errorMsg := fmt.Sprintf("[line %d] Error%s: %s", line, where, message)
	l.errors = append(l.errors, errorMsg)
	l.hadError = true
}

// Helper functions for creating test tokens
func createToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		token_type: tokenType,
		lexeme:     lexeme,
		literal:    literal,
		line:       line,
	}
}

func createNumberToken(value float64, line int) Token {
	// For testing, we need to match what the scanner actually produces
	// The scanner uses the original source text as the lexeme
	var lexeme string
	if value == float64(int64(value)) {
		// Integer values
		lexeme = fmt.Sprintf("%.0f", value)
	} else {
		// Decimal values - use the original format
		lexeme = fmt.Sprintf("%g", value)
	}
	return createToken(NUMBER, lexeme, value, line)
}

func createStringToken(value string, line int) Token {
	return createToken(STRING, fmt.Sprintf("\"%s\"", value), value, line)
}

func createIdentifierToken(name string, line int) Token {
	return createToken(IDENTIFIER, name, nil, line)
}

func createKeywordToken(keyword TokenType, line int) Token {
	keywordMap := map[TokenType]string{
		VAR:    "var",
		PRINT:  "print",
		IF:     "if",
		ELSE:   "else",
		WHILE:  "while",
		FOR:    "for",
		TRUE:   "true",
		FALSE:  "false",
		NIL:    "nil",
		AND:    "and",
		OR:     "or",
		CLASS:  "class",
		FUN:    "fun",
		RETURN: "return",
		SUPER:  "super",
		THIS:   "this",
	}
	return createToken(keyword, keywordMap[keyword], nil, line)
}

func createOperatorToken(operator TokenType, line int) Token {
	operatorMap := map[TokenType]string{
		LEFT_PAREN:    "(",
		RIGHT_PAREN:   ")",
		LEFT_BRACE:    "{",
		RIGHT_BRACE:   "}",
		COMMA:         ",",
		DOT:           ".",
		MINUS:         "-",
		PLUS:          "+",
		SEMICOLON:     ";",
		SLASH:         "/",
		STAR:          "*",
		BANG:          "!",
		BANG_EQUAL:    "!=",
		EQUAL:         "=",
		EQUAL_EQUAL:   "==",
		GREATER:       ">",
		GREATER_EQUAL: ">=",
		LESS:          "<",
		LESS_EQUAL:    "<=",
	}
	return createToken(operator, operatorMap[operator], nil, line)
}

func createEOFToken(line int) Token {
	return createToken(EOF, "", nil, line)
}

// Helper function to compare token slices
func compareTokens(t *testing.T, expected []Token, actual []Token, testName string) {
	if len(expected) != len(actual) {
		t.Errorf("%s: Expected %d tokens, got %d", testName, len(expected), len(actual))
		return
	}

	for i, exp := range expected {
		act := actual[i]
		if !tokensEqual(exp, act) {
			t.Errorf("%s: Token %d mismatch:\nExpected: %v\nActual:   %v",
				testName, i, exp, act)
		}
	}
}

// Helper function to compare individual tokens
func tokensEqual(expected, actual Token) bool {
	return expected.token_type == actual.token_type &&
		expected.lexeme == actual.lexeme &&
		expected.line == actual.line &&
		reflect.DeepEqual(expected.literal, actual.literal)
}

// Helper function to check if a slice of tokens contains only EOF
func isOnlyEOF(tokens []Token) bool {
	return len(tokens) == 1 && tokens[0].token_type == EOF
}

// Helper function to extract non-EOF tokens
func getNonEOFTokens(tokens []Token) []Token {
	result := make([]Token, 0)
	for _, token := range tokens {
		if token.token_type != EOF {
			result = append(result, token)
		}
	}
	return result
}
