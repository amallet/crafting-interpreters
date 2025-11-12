package main

import (
	"testing"
)

func TestParserVariableDeclaration(t *testing.T) {
	tests := []struct {
		input   string
		varName string
		hasInit bool
	}{
		{"var x;", "x", false},
		{"var y = 123;", "y", true},
		{"var name = \"hello\";", "name", true},
		{"var flag = true;", "flag", true},
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

			varStmt, ok := statements[0].(*VarStmt)
			if !ok {
				t.Errorf("Expected VarStmt, got %T", statements[0])
				return
			}

			if varStmt.variable.lexeme != test.varName {
				t.Errorf("Expected variable name %s, got %s", test.varName, varStmt.variable.lexeme)
			}

			if (varStmt.initializer != nil) != test.hasInit {
				t.Errorf("Expected hasInit %v, got %v", test.hasInit, varStmt.initializer != nil)
			}
		})
	}
}

func TestParserPrintStatement(t *testing.T) {
	input := "print 123;"
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

	printStmt, ok := statements[0].(*PrintStmt)
	if !ok {
		t.Errorf("Expected PrintStmt, got %T", statements[0])
		return
	}

	literalExpr, ok := printStmt.expression.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr, got %T", printStmt.expression)
		return
	}

	if literalExpr.Value != 123.0 {
		t.Errorf("Expected 123.0, got %v", literalExpr.Value)
	}
}

func TestParserIfStatement(t *testing.T) {
	input := "if (true) print 1;"
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

	ifStmt, ok := statements[0].(*IfStmt)
	if !ok {
		t.Errorf("Expected IfStmt, got %T", statements[0])
		return
	}

	// Check condition
	condition, ok := ifStmt.condition.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr condition, got %T", ifStmt.condition)
		return
	}

	if condition.Value != true {
		t.Errorf("Expected true condition, got %v", condition.Value)
	}

	// Check then branch
	thenPrintStmt, ok := ifStmt.thenBranch.(*PrintStmt)
	if !ok {
		t.Errorf("Expected PrintStmt then branch, got %T", ifStmt.thenBranch)
		return
	}

	thenExpr, ok := thenPrintStmt.expression.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr in then branch, got %T", thenPrintStmt.expression)
		return
	}

	if thenExpr.Value != 1.0 {
		t.Errorf("Expected 1.0 in then branch, got %v", thenExpr.Value)
	}

	// Check else branch (should be nil)
	if ifStmt.elseBranch != nil {
		t.Errorf("Expected nil else branch, got %T", ifStmt.elseBranch)
	}
}

func TestParserIfElseStatement(t *testing.T) {
	input := "if (false) print 1; else print 2;"
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

	ifStmt, ok := statements[0].(*IfStmt)
	if !ok {
		t.Errorf("Expected IfStmt, got %T", statements[0])
		return
	}

	// Check else branch exists
	if ifStmt.elseBranch == nil {
		t.Error("Expected else branch, got nil")
		return
	}

	elsePrintStmt, ok := ifStmt.elseBranch.(*PrintStmt)
	if !ok {
		t.Errorf("Expected PrintStmt else branch, got %T", ifStmt.elseBranch)
		return
	}

	elseExpr, ok := elsePrintStmt.expression.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr in else branch, got %T", elsePrintStmt.expression)
		return
	}

	if elseExpr.Value != 2.0 {
		t.Errorf("Expected 2.0 in else branch, got %v", elseExpr.Value)
	}
}

func TestParserWhileStatement(t *testing.T) {
	input := "while (true) print 1;"
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

	whileStmt, ok := statements[0].(*WhileStmt)
	if !ok {
		t.Errorf("Expected WhileStmt, got %T", statements[0])
		return
	}

	// Check condition
	condition, ok := whileStmt.condition.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr condition, got %T", whileStmt.condition)
		return
	}

	if condition.Value != true {
		t.Errorf("Expected true condition, got %v", condition.Value)
	}

	// Check body
	bodyPrintStmt, ok := whileStmt.body.(*PrintStmt)
	if !ok {
		t.Errorf("Expected PrintStmt body, got %T", whileStmt.body)
		return
	}

	bodyExpr, ok := bodyPrintStmt.expression.(*LiteralExpr)
	if !ok {
		t.Errorf("Expected LiteralExpr in body, got %T", bodyPrintStmt.expression)
		return
	}

	if bodyExpr.Value != 1.0 {
		t.Errorf("Expected 1.0 in body, got %v", bodyExpr.Value)
	}
}

func TestParserBlockStatement(t *testing.T) {
	input := "{ var x = 1; print x; }"
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

	blockStmt, ok := statements[0].(*BlockStmt)
	if !ok {
		t.Errorf("Expected BlockStmt, got %T", statements[0])
		return
	}

	if len(blockStmt.statements) != 2 {
		t.Errorf("Expected 2 statements in block, got %d", len(blockStmt.statements))
		return
	}

	// Check first statement (var declaration)
	varStmt, ok := blockStmt.statements[0].(*VarStmt)
	if !ok {
		t.Errorf("Expected VarStmt as first statement, got %T", blockStmt.statements[0])
		return
	}

	if varStmt.variable.lexeme != "x" {
		t.Errorf("Expected variable name 'x', got %s", varStmt.variable.lexeme)
	}

	// Check second statement (print)
	printStmt, ok := blockStmt.statements[1].(*PrintStmt)
	if !ok {
		t.Errorf("Expected PrintStmt as second statement, got %T", blockStmt.statements[1])
		return
	}

	variableExpr, ok := printStmt.expression.(*VariableExpr)
	if !ok {
		t.Errorf("Expected VariableExpr in print, got %T", printStmt.expression)
		return
	}

	if variableExpr.variable.lexeme != "x" {
		t.Errorf("Expected variable name 'x' in print, got %s", variableExpr.variable.lexeme)
	}
}

func TestParserFunctionDeclaration(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedName   string
		expectedParams int
		expectedBody   int
	}{
		{
			name:           "Function with no parameters",
			input:          "fun foo() { return 1; }",
			expectedName:   "foo",
			expectedParams: 0,
			expectedBody:   1,
		},
		{
			name:           "Function with one parameter",
			input:          "fun greet(name) { print name; }",
			expectedName:   "greet",
			expectedParams: 1,
			expectedBody:   1,
		},
		{
			name:           "Function with multiple parameters",
			input:          "fun add(a, b) { return a + b; }",
			expectedName:   "add",
			expectedParams: 2,
			expectedBody:   1,
		},
		{
			name:           "Function with three parameters",
			input:          "fun triple(a, b, c) { return a + b + c; }",
			expectedName:   "triple",
			expectedParams: 3,
			expectedBody:   1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

			funcStmt, ok := statements[0].(*FunctionStmt)
			if !ok {
				t.Errorf("Expected FunctionStmt, got %T", statements[0])
				return
			}

			if funcStmt.functionName.lexeme != test.expectedName {
				t.Errorf("Expected function name %s, got %s", test.expectedName, funcStmt.functionName.lexeme)
			}

			if len(funcStmt.params) != test.expectedParams {
				t.Errorf("Expected %d parameters, got %d", test.expectedParams, len(funcStmt.params))
			}

			if len(funcStmt.body) != test.expectedBody {
				t.Errorf("Expected %d statements in body, got %d", test.expectedBody, len(funcStmt.body))
			}
		})
	}
}

func TestParserReturnStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasValue bool
	}{
		{
			name:     "Return with value",
			input:    "return 42;",
			hasValue: true,
		},
		{
			name:     "Return without value",
			input:    "return;",
			hasValue: false,
		},
		{
			name:     "Return with expression",
			input:    "return x + 1;",
			hasValue: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

			returnStmt, ok := statements[0].(*ReturnStmt)
			if !ok {
				t.Errorf("Expected ReturnStmt, got %T", statements[0])
				return
			}

			if (returnStmt.returnValue != nil) != test.hasValue {
				t.Errorf("Expected hasValue %v, got %v", test.hasValue, returnStmt.returnValue != nil)
			}
		})
	}
}
