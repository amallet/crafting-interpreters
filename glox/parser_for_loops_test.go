package main

import (
	"testing"
)

func TestParserForLoops(t *testing.T) {
	t.Run("Basic for loop with all clauses", func(t *testing.T) {
		// "for (var i = 0; i < 10; i = i + 1) print i;"
		input := "for (var i = 0; i < 10; i = i + 1) print i;"
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

		// For loops are desugared into blocks with initialization and while loop
		blockStmt, ok := statements[0].(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt (desugared for loop), got %T", statements[0])
			return
		}

		if len(blockStmt.statements) != 2 {
			t.Errorf("Expected 2 statements in desugared for loop, got %d", len(blockStmt.statements))
			return
		}

		// First statement should be variable declaration
		varStmt, ok := blockStmt.statements[0].(*VarStmt)
		if !ok {
			t.Errorf("Expected VarStmt as first statement, got %T", blockStmt.statements[0])
			return
		}

		if varStmt.variable.lexeme != "i" {
			t.Errorf("Expected variable name 'i', got %s", varStmt.variable.lexeme)
		}

		// Check variable initializer value
		initLiteral, ok := varStmt.initializer.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as variable initializer, got %T", varStmt.initializer)
			return
		}

		if initLiteral.Value != 0.0 {
			t.Errorf("Expected variable initializer value 0.0, got %v", initLiteral.Value)
		}

		// Second statement should be while loop
		whileStmt, ok := blockStmt.statements[1].(*WhileStmt)
		if !ok {
			t.Errorf("Expected WhileStmt as second statement, got %T", blockStmt.statements[1])
			return
		}

		// While condition should be the loop condition: i < 10
		condition, ok := whileStmt.condition.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr condition, got %T", whileStmt.condition)
			return
		}

		if condition.Operator.token_type != LESS {
			t.Errorf("Expected LESS operator, got %v", condition.Operator.token_type)
		}

		// Check left operand of condition (variable i)
		leftVar, ok := condition.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as left operand, got %T", condition.Left)
			return
		}

		if leftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as left operand, got %s", leftVar.variable.lexeme)
		}

		// Check right operand of condition (literal 10)
		rightLiteral, ok := condition.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as right operand, got %T", condition.Right)
			return
		}

		if rightLiteral.Value != 10.0 {
			t.Errorf("Expected condition right operand 10.0, got %v", rightLiteral.Value)
		}

		// While body should be a block containing the original body and update
		bodyBlock, ok := whileStmt.body.(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt as while body, got %T", whileStmt.body)
			return
		}

		if len(bodyBlock.statements) != 2 {
			t.Errorf("Expected 2 statements in while body, got %d", len(bodyBlock.statements))
			return
		}

		// First statement in body should be the original print statement
		printStmt, ok := bodyBlock.statements[0].(*PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt as first body statement, got %T", bodyBlock.statements[0])
			return
		}

		// Check print expression (variable i)
		printVar, ok := printStmt.expression.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as print expression, got %T", printStmt.expression)
			return
		}

		if printVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' in print statement, got %s", printVar.variable.lexeme)
		}

		// Second statement in body should be the update expression: i = i + 1
		updateStmt, ok := bodyBlock.statements[1].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt as update, got %T", bodyBlock.statements[1])
			return
		}

		updateAssign, ok := updateStmt.expression.(*AssignExpr)
		if !ok {
			t.Errorf("Expected AssignExpr as update, got %T", updateStmt.expression)
			return
		}

		if updateAssign.variable.lexeme != "i" {
			t.Errorf("Expected update variable 'i', got %s", updateAssign.variable.lexeme)
		}

		// Check update expression: i + 1
		updateBinary, ok := updateAssign.value.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr as update value, got %T", updateAssign.value)
			return
		}

		if updateBinary.Operator.token_type != PLUS {
			t.Errorf("Expected PLUS operator in update, got %v", updateBinary.Operator.token_type)
		}

		// Check left operand of update (variable i)
		updateLeftVar, ok := updateBinary.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as update left operand, got %T", updateBinary.Left)
			return
		}

		if updateLeftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as update left operand, got %s", updateLeftVar.variable.lexeme)
		}

		// Check right operand of update (literal 1)
		updateRightLiteral, ok := updateBinary.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as update right operand, got %T", updateBinary.Right)
			return
		}

		if updateRightLiteral.Value != 1.0 {
			t.Errorf("Expected update right operand 1.0, got %v", updateRightLiteral.Value)
		}
	})

	t.Run("For loop with no initializer", func(t *testing.T) {
		// "for (; i < 10; i = i + 1) print i;"
		input := "for (; i < 10; i = i + 1) print i;"
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

		// Should be a while loop directly (no block wrapper)
		whileStmt, ok := statements[0].(*WhileStmt)
		if !ok {
			t.Errorf("Expected WhileStmt (no initializer), got %T", statements[0])
			return
		}

		// Condition should be the loop condition: i < 10
		condition, ok := whileStmt.condition.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr condition, got %T", whileStmt.condition)
			return
		}

		if condition.Operator.token_type != LESS {
			t.Errorf("Expected LESS operator, got %v", condition.Operator.token_type)
		}

		// Check left operand of condition (variable i)
		leftVar, ok := condition.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as left operand, got %T", condition.Left)
			return
		}

		if leftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as left operand, got %s", leftVar.variable.lexeme)
		}

		// Check right operand of condition (literal 10)
		rightLiteral, ok := condition.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as right operand, got %T", condition.Right)
			return
		}

		if rightLiteral.Value != 10.0 {
			t.Errorf("Expected condition right operand 10.0, got %v", rightLiteral.Value)
		}

		// While body should be a block containing the original body and update
		bodyBlock, ok := whileStmt.body.(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt as while body, got %T", whileStmt.body)
			return
		}

		if len(bodyBlock.statements) != 2 {
			t.Errorf("Expected 2 statements in while body, got %d", len(bodyBlock.statements))
			return
		}

		// First statement in body should be the original print statement
		printStmt, ok := bodyBlock.statements[0].(*PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt as first body statement, got %T", bodyBlock.statements[0])
			return
		}

		// Check print expression (variable i)
		printVar, ok := printStmt.expression.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as print expression, got %T", printStmt.expression)
			return
		}

		if printVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' in print statement, got %s", printVar.variable.lexeme)
		}

		// Second statement in body should be the update expression: i = i + 1
		updateStmt, ok := bodyBlock.statements[1].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt as update, got %T", bodyBlock.statements[1])
			return
		}

		updateAssign, ok := updateStmt.expression.(*AssignExpr)
		if !ok {
			t.Errorf("Expected AssignExpr as update, got %T", updateStmt.expression)
			return
		}

		if updateAssign.variable.lexeme != "i" {
			t.Errorf("Expected update variable 'i', got %s", updateAssign.variable.lexeme)
		}

		// Check update expression: i + 1
		updateBinary, ok := updateAssign.value.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr as update value, got %T", updateAssign.value)
			return
		}

		if updateBinary.Operator.token_type != PLUS {
			t.Errorf("Expected PLUS operator in update, got %v", updateBinary.Operator.token_type)
		}

		// Check left operand of update (variable i)
		updateLeftVar, ok := updateBinary.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as update left operand, got %T", updateBinary.Left)
			return
		}

		if updateLeftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as update left operand, got %s", updateLeftVar.variable.lexeme)
		}

		// Check right operand of update (literal 1)
		updateRightLiteral, ok := updateBinary.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as update right operand, got %T", updateBinary.Right)
			return
		}

		if updateRightLiteral.Value != 1.0 {
			t.Errorf("Expected update right operand 1.0, got %v", updateRightLiteral.Value)
		}
	})

	t.Run("For loop with no condition", func(t *testing.T) {
		// "for (var i = 0; ; i = i + 1) print i;"
		input := "for (var i = 0; ; i = i + 1) print i;"
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

		// Should be a block with initialization and while loop
		blockStmt, ok := statements[0].(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt, got %T", statements[0])
			return
		}

		if len(blockStmt.statements) != 2 {
			t.Errorf("Expected 2 statements in desugared for loop, got %d", len(blockStmt.statements))
			return
		}

		// First statement should be variable declaration
		varStmt, ok := blockStmt.statements[0].(*VarStmt)
		if !ok {
			t.Errorf("Expected VarStmt as first statement, got %T", blockStmt.statements[0])
			return
		}

		if varStmt.variable.lexeme != "i" {
			t.Errorf("Expected variable name 'i', got %s", varStmt.variable.lexeme)
		}

		// Check variable initializer value
		initLiteral, ok := varStmt.initializer.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as variable initializer, got %T", varStmt.initializer)
			return
		}

		if initLiteral.Value != 0.0 {
			t.Errorf("Expected variable initializer value 0.0, got %v", initLiteral.Value)
		}

		// Second statement should be while loop with true condition
		whileStmt, ok := blockStmt.statements[1].(*WhileStmt)
		if !ok {
			t.Errorf("Expected WhileStmt, got %T", blockStmt.statements[1])
			return
		}

		// Condition should be true literal (infinite loop)
		condition, ok := whileStmt.condition.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr condition (true), got %T", whileStmt.condition)
			return
		}

		if condition.Value != true {
			t.Errorf("Expected true condition, got %v", condition.Value)
		}

		// While body should be a block containing the original body and update
		bodyBlock, ok := whileStmt.body.(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt as while body, got %T", whileStmt.body)
			return
		}

		if len(bodyBlock.statements) != 2 {
			t.Errorf("Expected 2 statements in while body, got %d", len(bodyBlock.statements))
			return
		}

		// First statement in body should be the original print statement
		printStmt, ok := bodyBlock.statements[0].(*PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt as first body statement, got %T", bodyBlock.statements[0])
			return
		}

		// Check print expression (variable i)
		printVar, ok := printStmt.expression.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as print expression, got %T", printStmt.expression)
			return
		}

		if printVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' in print statement, got %s", printVar.variable.lexeme)
		}

		// Second statement in body should be the update expression: i = i + 1
		updateStmt, ok := bodyBlock.statements[1].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt as update, got %T", bodyBlock.statements[1])
			return
		}

		updateAssign, ok := updateStmt.expression.(*AssignExpr)
		if !ok {
			t.Errorf("Expected AssignExpr as update, got %T", updateStmt.expression)
			return
		}

		if updateAssign.variable.lexeme != "i" {
			t.Errorf("Expected update variable 'i', got %s", updateAssign.variable.lexeme)
		}

		// Check update expression: i + 1
		updateBinary, ok := updateAssign.value.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr as update value, got %T", updateAssign.value)
			return
		}

		if updateBinary.Operator.token_type != PLUS {
			t.Errorf("Expected PLUS operator in update, got %v", updateBinary.Operator.token_type)
		}

		// Check left operand of update (variable i)
		updateLeftVar, ok := updateBinary.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as update left operand, got %T", updateBinary.Left)
			return
		}

		if updateLeftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as update left operand, got %s", updateLeftVar.variable.lexeme)
		}

		// Check right operand of update (literal 1)
		updateRightLiteral, ok := updateBinary.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as update right operand, got %T", updateBinary.Right)
			return
		}

		if updateRightLiteral.Value != 1.0 {
			t.Errorf("Expected update right operand 1.0, got %v", updateRightLiteral.Value)
		}
	})

	t.Run("For loop with no increment", func(t *testing.T) {
		// "for (var i = 0; i < 10;) print i;"
		input := "for (var i = 0; i < 10;) print i;"
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

		// Should be a block with initialization and while loop
		blockStmt, ok := statements[0].(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt, got %T", statements[0])
			return
		}

		if len(blockStmt.statements) != 2 {
			t.Errorf("Expected 2 statements in desugared for loop, got %d", len(blockStmt.statements))
			return
		}

		// First statement should be variable declaration
		varStmt, ok := blockStmt.statements[0].(*VarStmt)
		if !ok {
			t.Errorf("Expected VarStmt as first statement, got %T", blockStmt.statements[0])
			return
		}

		if varStmt.variable.lexeme != "i" {
			t.Errorf("Expected variable name 'i', got %s", varStmt.variable.lexeme)
		}

		// Check variable initializer value
		initLiteral, ok := varStmt.initializer.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as variable initializer, got %T", varStmt.initializer)
			return
		}

		if initLiteral.Value != 0.0 {
			t.Errorf("Expected variable initializer value 0.0, got %v", initLiteral.Value)
		}

		// Second statement should be while loop
		whileStmt, ok := blockStmt.statements[1].(*WhileStmt)
		if !ok {
			t.Errorf("Expected WhileStmt, got %T", blockStmt.statements[1])
			return
		}

		// While condition should be the loop condition: i < 10
		condition, ok := whileStmt.condition.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr condition, got %T", whileStmt.condition)
			return
		}

		if condition.Operator.token_type != LESS {
			t.Errorf("Expected LESS operator, got %v", condition.Operator.token_type)
		}

		// Check left operand of condition (variable i)
		leftVar, ok := condition.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as left operand, got %T", condition.Left)
			return
		}

		if leftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as left operand, got %s", leftVar.variable.lexeme)
		}

		// Check right operand of condition (literal 10)
		rightLiteral, ok := condition.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as right operand, got %T", condition.Right)
			return
		}

		if rightLiteral.Value != 10.0 {
			t.Errorf("Expected condition right operand 10.0, got %v", rightLiteral.Value)
		}

		// Body should be the original statement directly (no block wrapper)
		printStmt, ok := whileStmt.body.(*PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt as while body, got %T", whileStmt.body)
			return
		}

		// Check print expression (variable i)
		printVar, ok := printStmt.expression.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as print expression, got %T", printStmt.expression)
			return
		}

		if printVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' in print statement, got %s", printVar.variable.lexeme)
		}
	})

	t.Run("For loop with expression initializer", func(t *testing.T) {
		// "for (i = 0; i < 10; i = i + 1) print i;"
		input := "for (i = 0; i < 10; i = i + 1) print i;"
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

		// Should be a block with expression statement and while loop
		blockStmt, ok := statements[0].(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt, got %T", statements[0])
			return
		}

		if len(blockStmt.statements) != 2 {
			t.Errorf("Expected 2 statements in desugared for loop, got %d", len(blockStmt.statements))
			return
		}

		// First statement should be expression statement: i = 0
		exprStmt, ok := blockStmt.statements[0].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt as first statement, got %T", blockStmt.statements[0])
			return
		}

		assignExpr, ok := exprStmt.expression.(*AssignExpr)
		if !ok {
			t.Errorf("Expected AssignExpr, got %T", exprStmt.expression)
			return
		}

		if assignExpr.variable.lexeme != "i" {
			t.Errorf("Expected assignment to 'i', got %s", assignExpr.variable.lexeme)
		}

		// Check assignment value (literal 0)
		assignLiteral, ok := assignExpr.value.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as assignment value, got %T", assignExpr.value)
			return
		}

		if assignLiteral.Value != 0.0 {
			t.Errorf("Expected assignment value 0.0, got %v", assignLiteral.Value)
		}

		// Second statement should be while loop
		whileStmt, ok := blockStmt.statements[1].(*WhileStmt)
		if !ok {
			t.Errorf("Expected WhileStmt, got %T", blockStmt.statements[1])
			return
		}

		// While condition should be the loop condition: i < 10
		condition, ok := whileStmt.condition.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr condition, got %T", whileStmt.condition)
			return
		}

		if condition.Operator.token_type != LESS {
			t.Errorf("Expected LESS operator, got %v", condition.Operator.token_type)
		}

		// Check left operand of condition (variable i)
		leftVar, ok := condition.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as left operand, got %T", condition.Left)
			return
		}

		if leftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as left operand, got %s", leftVar.variable.lexeme)
		}

		// Check right operand of condition (literal 10)
		rightLiteral, ok := condition.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as right operand, got %T", condition.Right)
			return
		}

		if rightLiteral.Value != 10.0 {
			t.Errorf("Expected condition right operand 10.0, got %v", rightLiteral.Value)
		}

		// While body should be a block containing the original body and update
		bodyBlock, ok := whileStmt.body.(*BlockStmt)
		if !ok {
			t.Errorf("Expected BlockStmt as while body, got %T", whileStmt.body)
			return
		}

		if len(bodyBlock.statements) != 2 {
			t.Errorf("Expected 2 statements in while body, got %d", len(bodyBlock.statements))
			return
		}

		// First statement in body should be the original print statement
		printStmt, ok := bodyBlock.statements[0].(*PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt as first body statement, got %T", bodyBlock.statements[0])
			return
		}

		// Check print expression (variable i)
		printVar, ok := printStmt.expression.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as print expression, got %T", printStmt.expression)
			return
		}

		if printVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' in print statement, got %s", printVar.variable.lexeme)
		}

		// Second statement in body should be the update expression: i = i + 1
		updateStmt, ok := bodyBlock.statements[1].(*ExpressionStmt)
		if !ok {
			t.Errorf("Expected ExpressionStmt as update, got %T", bodyBlock.statements[1])
			return
		}

		updateAssign, ok := updateStmt.expression.(*AssignExpr)
		if !ok {
			t.Errorf("Expected AssignExpr as update, got %T", updateStmt.expression)
			return
		}

		if updateAssign.variable.lexeme != "i" {
			t.Errorf("Expected update variable 'i', got %s", updateAssign.variable.lexeme)
		}

		// Check update expression: i + 1
		updateBinary, ok := updateAssign.value.(*BinaryExpr)
		if !ok {
			t.Errorf("Expected BinaryExpr as update value, got %T", updateAssign.value)
			return
		}

		if updateBinary.Operator.token_type != PLUS {
			t.Errorf("Expected PLUS operator in update, got %v", updateBinary.Operator.token_type)
		}

		// Check left operand of update (variable i)
		updateLeftVar, ok := updateBinary.Left.(*VariableExpr)
		if !ok {
			t.Errorf("Expected VariableExpr as update left operand, got %T", updateBinary.Left)
			return
		}

		if updateLeftVar.variable.lexeme != "i" {
			t.Errorf("Expected variable 'i' as update left operand, got %s", updateLeftVar.variable.lexeme)
		}

		// Check right operand of update (literal 1)
		updateRightLiteral, ok := updateBinary.Right.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as update right operand, got %T", updateBinary.Right)
			return
		}

		if updateRightLiteral.Value != 1.0 {
			t.Errorf("Expected update right operand 1.0, got %v", updateRightLiteral.Value)
		}
	})

	t.Run("Empty for loop", func(t *testing.T) {
		// "for (;;) print 1;"
		input := "for (;;) print 1;"
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

		// Should be a while loop with true condition
		whileStmt, ok := statements[0].(*WhileStmt)
		if !ok {
			t.Errorf("Expected WhileStmt, got %T", statements[0])
			return
		}

		// Condition should be true literal
		condition, ok := whileStmt.condition.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr condition (true), got %T", whileStmt.condition)
			return
		}

		if condition.Value != true {
			t.Errorf("Expected true condition, got %v", condition.Value)
		}

		// Body should be the original statement directly (no block wrapper)
		printStmt, ok := whileStmt.body.(*PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt as while body, got %T", whileStmt.body)
			return
		}

		// Check print expression (literal 1)
		printLiteral, ok := printStmt.expression.(*LiteralExpr)
		if !ok {
			t.Errorf("Expected LiteralExpr as print expression, got %T", printStmt.expression)
			return
		}

		if printLiteral.Value != 1.0 {
			t.Errorf("Expected print expression value 1.0, got %v", printLiteral.Value)
		}
	})
}

func TestParserForLoopErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "Missing opening parenthesis",
			input:       "for var i = 0; i < 10; i = i + 1) print i;",
			expectError: true,
			description: "For loop without opening parenthesis",
		},
		{
			name:        "Missing closing parenthesis",
			input:       "for (var i = 0; i < 10; i = i + 1 print i;",
			expectError: true,
			description: "For loop without closing parenthesis",
		},
		{
			name:        "Missing first semicolon",
			input:       "for (var i = 0 i < 10; i = i + 1) print i;",
			expectError: true,
			description: "For loop without first semicolon",
		},
		{
			name:        "Missing second semicolon",
			input:       "for (var i = 0; i < 10 i = i + 1) print i;",
			expectError: true,
			description: "For loop without second semicolon",
		},
		{
			name:        "Missing for body",
			input:       "for (var i = 0; i < 10; i = i + 1);",
			expectError: true,
			description: "For loop without body",
		},
		{
			name:        "Invalid initializer",
			input:       "for (var = 0; i < 10; i = i + 1) print i;",
			expectError: true,
			description: "For loop with invalid variable declaration",
		},
		{
			name:        "Invalid condition",
			input:       "for (var i = 0; + ; i = i + 1) print i;",
			expectError: true,
			description: "For loop with invalid condition",
		},
		{
			name:        "Invalid increment",
			input:       "for (var i = 0; i < 10; + ) print i;",
			expectError: true,
			description: "For loop with invalid increment",
		},
		{
			name:        "Missing condition semicolon",
			input:       "for (var i = 0; i < 10 i = i + 1) print i;",
			expectError: true,
			description: "For loop missing semicolon after condition",
		},
		{
			name:        "Missing increment semicolon",
			input:       "for (var i = 0; i < 10; i = i + 1 print i;",
			expectError: true,
			description: "For loop missing semicolon after increment",
		},
		{
			name:        "Nested for loop errors",
			input:       "for (for (;;) print 1; i < 10; i = i + 1) print i;",
			expectError: true,
			description: "For loop with nested for loop in initializer",
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

func TestParserForLoopErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedErrors []string // Substrings that should appear in error messages
	}{
		{
			name:  "Missing opening parenthesis",
			input: "for var i = 0; i < 10; i = i + 1) print i;",
			expectedErrors: []string{
				"Expect '(' after 'for'",
			},
		},
		{
			name:  "Missing closing parenthesis",
			input: "for (var i = 0; i < 10; i = i + 1 print i;",
			expectedErrors: []string{
				"Expect ')' after 'for' clause",
			},
		},
		{
			name:  "Missing first semicolon",
			input: "for (var i = 0 i < 10; i = i + 1) print i;",
			expectedErrors: []string{
				"Expect ';' after variable declaration",
			},
		},
		{
			name:  "Missing second semicolon",
			input: "for (var i = 0; i < 10 i = i + 1) print i;",
			expectedErrors: []string{
				"Expect ';' after loop condition",
			},
		},
		{
			name:  "Invalid initializer",
			input: "for (var = 0; i < 10; i = i + 1) print i;",
			expectedErrors: []string{
				"Expect variable name",
			},
		},
		{
			name:  "Invalid condition",
			input: "for (var i = 0; + ; i = i + 1) print i;",
			expectedErrors: []string{
				"Expected expression",
			},
		},
		{
			name:  "Invalid increment",
			input: "for (var i = 0; i < 10; + ) print i;",
			expectedErrors: []string{
				"Expected expression",
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
