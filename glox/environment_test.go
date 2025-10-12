package main

import (
	"fmt"
	"testing"
)

// ============================================================================
// ENVIRONMENT TESTS
// ============================================================================

// Test helper to create a new environment
func createTestEnvironment(enclosing *Environment) *Environment {
	return NewEnvironment(enclosing)
}

// Test helper to check if an error occurred
func assertError(t *testing.T, err error, testName string) {
	if err == nil {
		t.Errorf("%s: Expected error but got none", testName)
	}
}

// Test helper to check if no error occurred
func assertNoError(t *testing.T, err error, testName string) {
	if err != nil {
		t.Errorf("%s: Expected no error but got: %v", testName, err)
	}
}

// ============================================================================
// ENVIRONMENT CREATION TESTS
// ============================================================================

func TestEnvironmentCreation(t *testing.T) {
	t.Run("Create root environment", func(t *testing.T) {
		env := createTestEnvironment(nil)

		if env == nil {
			t.Fatal("Expected environment to be created, got nil")
		}
		if env.enclosing != nil {
			t.Error("Expected root environment to have no enclosing environment")
		}
		if env.values == nil {
			t.Error("Expected environment to have initialized values map")
		}
		if len(env.values) != 0 {
			t.Error("Expected new environment to have empty values map")
		}
	})

	t.Run("Create nested environment", func(t *testing.T) {
		parent := createTestEnvironment(nil)
		child := createTestEnvironment(parent)

		if child == nil {
			t.Fatal("Expected child environment to be created, got nil")
		}
		if child.enclosing != parent {
			t.Error("Expected child environment to have correct parent reference")
		}
		if child.values == nil {
			t.Error("Expected child environment to have initialized values map")
		}
	})

	t.Run("Create deeply nested environment", func(t *testing.T) {
		root := createTestEnvironment(nil)
		level1 := createTestEnvironment(root)
		level2 := createTestEnvironment(level1)
		level3 := createTestEnvironment(level2)

		// Verify the chain
		if level3.enclosing != level2 {
			t.Error("Expected level3 to have level2 as parent")
		}
		if level2.enclosing != level1 {
			t.Error("Expected level2 to have level1 as parent")
		}
		if level1.enclosing != root {
			t.Error("Expected level1 to have root as parent")
		}
		if root.enclosing != nil {
			t.Error("Expected root to have no parent")
		}
	})
}

// ============================================================================
// VARIABLE DEFINITION TESTS
// ============================================================================

func TestVariableDefinition(t *testing.T) {
	t.Run("Define variable in root environment", func(t *testing.T) {
		env := createTestEnvironment(nil)

		env.defineVarValue("x", 42.0)

		// Check that variable was stored
		if value, exists := env.values["x"]; !exists {
			t.Error("Expected variable 'x' to be defined")
		} else {
			assertEqual(t, 42.0, value, "Variable value")
		}
	})

	t.Run("Define multiple variables", func(t *testing.T) {
		env := createTestEnvironment(nil)

		env.defineVarValue("x", 42.0)
		env.defineVarValue("y", "hello")
		env.defineVarValue("z", true)
		env.defineVarValue("w", nil)

		// Check all variables
		assertEqual(t, 42.0, env.values["x"], "Variable x")
		assertEqual(t, "hello", env.values["y"], "Variable y")
		assertEqual(t, true, env.values["z"], "Variable z")
		assertEqual(t, nil, env.values["w"], "Variable w")

		if len(env.values) != 4 {
			t.Errorf("Expected 4 variables, got %d", len(env.values))
		}
	})

	t.Run("Redefine variable", func(t *testing.T) {
		env := createTestEnvironment(nil)

		// Define variable
		env.defineVarValue("x", 42.0)
		assertEqual(t, 42.0, env.values["x"], "Initial value")

		// Redefine with new value
		env.defineVarValue("x", "new value")
		assertEqual(t, "new value", env.values["x"], "Redefined value")

		if len(env.values) != 1 {
			t.Errorf("Expected 1 variable after redefinition, got %d", len(env.values))
		}
	})

	t.Run("Define variable with different types", func(t *testing.T) {
		env := createTestEnvironment(nil)

		tests := []struct {
			name  string
			value any
		}{
			{"number", 42.0},
			{"string", "hello"},
			{"boolean", true},
			{"nil", nil},
			{"decimal", 3.14},
			{"negative", -1.0},
			{"empty_string", ""},
		}

		for _, tt := range tests {
			env.defineVarValue(tt.name, tt.value)
			assertEqual(t, tt.value, env.values[tt.name], tt.name)
		}
	})
}

// ============================================================================
// VARIABLE RETRIEVAL TESTS
// ============================================================================

func TestVariableRetrieval(t *testing.T) {
	t.Run("Retrieve defined variable", func(t *testing.T) {
		env := createTestEnvironment(nil)
		env.defineVarValue("x", 42.0)

		token := createIdentifierToken("x", 1)
		value, err := env.getVarValue(token)

		assertNoError(t, err, "Retrieve defined variable")
		assertEqual(t, 42.0, value, "Variable value")
	})

	t.Run("Retrieve undefined variable", func(t *testing.T) {
		env := createTestEnvironment(nil)

		token := createIdentifierToken("undefined", 1)
		_, err := env.getVarValue(token)

		assertError(t, err, "Retrieve undefined variable")

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Undefined variable 'undefined'" {
			t.Errorf("Expected 'Undefined variable 'undefined'', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Retrieve variable from parent environment", func(t *testing.T) {
		parent := createTestEnvironment(nil)
		parent.defineVarValue("x", 42.0)

		child := createTestEnvironment(parent)

		token := createIdentifierToken("x", 1)
		value, err := child.getVarValue(token)

		assertNoError(t, err, "Retrieve variable from parent")
		assertEqual(t, 42.0, value, "Variable value from parent")
	})

	t.Run("Retrieve variable from grandparent environment", func(t *testing.T) {
		grandparent := createTestEnvironment(nil)
		grandparent.defineVarValue("x", 42.0)

		parent := createTestEnvironment(grandparent)
		child := createTestEnvironment(parent)

		token := createIdentifierToken("x", 1)
		value, err := child.getVarValue(token)

		assertNoError(t, err, "Retrieve variable from grandparent")
		assertEqual(t, 42.0, value, "Variable value from grandparent")
	})

	t.Run("Variable shadowing - child shadows parent", func(t *testing.T) {
		parent := createTestEnvironment(nil)
		parent.defineVarValue("x", "parent value")

		child := createTestEnvironment(parent)
		child.defineVarValue("x", "child value")

		token := createIdentifierToken("x", 1)
		value, err := child.getVarValue(token)

		assertNoError(t, err, "Retrieve shadowed variable")
		assertEqual(t, "child value", value, "Shadowed variable value")
	})

	t.Run("Variable shadowing - grandchild shadows grandparent", func(t *testing.T) {
		grandparent := createTestEnvironment(nil)
		grandparent.defineVarValue("x", "grandparent value")

		parent := createTestEnvironment(grandparent)
		child := createTestEnvironment(parent)
		child.defineVarValue("x", "child value")

		token := createIdentifierToken("x", 1)
		value, err := child.getVarValue(token)

		assertNoError(t, err, "Retrieve shadowed variable")
		assertEqual(t, "child value", value, "Shadowed variable value")
	})

	t.Run("Retrieve different variables from different scopes", func(t *testing.T) {
		parent := createTestEnvironment(nil)
		parent.defineVarValue("parent_var", "parent")

		child := createTestEnvironment(parent)
		child.defineVarValue("child_var", "child")

		// Retrieve parent variable
		parentToken := createIdentifierToken("parent_var", 1)
		parentValue, err := child.getVarValue(parentToken)
		assertNoError(t, err, "Retrieve parent variable")
		assertEqual(t, "parent", parentValue, "Parent variable value")

		// Retrieve child variable
		childToken := createIdentifierToken("child_var", 2)
		childValue, err := child.getVarValue(childToken)
		assertNoError(t, err, "Retrieve child variable")
		assertEqual(t, "child", childValue, "Child variable value")
	})

	t.Run("Retrieve variable with different types", func(t *testing.T) {
		env := createTestEnvironment(nil)

		tests := []struct {
			name  string
			value any
		}{
			{"number", 42.0},
			{"string", "hello"},
			{"boolean", true},
			{"nil", nil},
			{"decimal", 3.14},
			{"negative", -1.0},
			{"empty_string", ""},
		}

		for _, tt := range tests {
			env.defineVarValue(tt.name, tt.value)

			token := createIdentifierToken(tt.name, 1)
			value, err := env.getVarValue(token)

			assertNoError(t, err, "Retrieve "+tt.name)
			assertEqual(t, tt.value, value, tt.name+" value")
		}
	})
}

// ============================================================================
// VARIABLE ASSIGNMENT TESTS
// ============================================================================

func TestVariableAssignment(t *testing.T) {
	t.Run("Assign to defined variable", func(t *testing.T) {
		env := createTestEnvironment(nil)
		env.defineVarValue("x", 42.0)

		token := createIdentifierToken("x", 1)
		err := env.assignVarValue(token, "new value")

		assertNoError(t, err, "Assign to defined variable")
		assertEqual(t, "new value", env.values["x"], "Assigned value")
	})

	t.Run("Assign to undefined variable", func(t *testing.T) {
		env := createTestEnvironment(nil)

		token := createIdentifierToken("undefined", 1)
		err := env.assignVarValue(token, 42.0)

		assertError(t, err, "Assign to undefined variable")

		runtimeErr, ok := err.(RuntimeError)
		if !ok {
			t.Fatalf("Expected RuntimeError, got %T", err)
		}
		if runtimeErr.message != "Undefined variable 'undefined'" {
			t.Errorf("Expected 'Undefined variable 'undefined'', got '%s'", runtimeErr.message)
		}
	})

	t.Run("Assign to variable in parent environment", func(t *testing.T) {
		parent := createTestEnvironment(nil)
		parent.defineVarValue("x", "original")

		child := createTestEnvironment(parent)

		token := createIdentifierToken("x", 1)
		err := child.assignVarValue(token, "modified")

		assertNoError(t, err, "Assign to parent variable")
		assertEqual(t, "modified", parent.values["x"], "Parent variable modified")

		// Child should not have a local copy - assignment modifies parent
		if _, exists := child.values["x"]; exists {
			t.Error("Child should not have local copy of parent variable")
		}
	})

	t.Run("Assign to variable in grandparent environment", func(t *testing.T) {
		grandparent := createTestEnvironment(nil)
		grandparent.defineVarValue("x", "original")

		parent := createTestEnvironment(grandparent)
		child := createTestEnvironment(parent)

		token := createIdentifierToken("x", 1)
		err := child.assignVarValue(token, "modified")

		assertNoError(t, err, "Assign to grandparent variable")
		assertEqual(t, "modified", grandparent.values["x"], "Grandparent variable modified")
	})

	t.Run("Assign to shadowed variable", func(t *testing.T) {
		parent := createTestEnvironment(nil)
		parent.defineVarValue("x", "parent value")

		child := createTestEnvironment(parent)
		child.defineVarValue("x", "child value")

		token := createIdentifierToken("x", 1)
		err := child.assignVarValue(token, "new child value")

		assertNoError(t, err, "Assign to shadowed variable")
		assertEqual(t, "new child value", child.values["x"], "Child variable modified")
		assertEqual(t, "parent value", parent.values["x"], "Parent variable unchanged")
	})

	t.Run("Assign different types to same variable", func(t *testing.T) {
		env := createTestEnvironment(nil)
		env.defineVarValue("x", 42.0)

		token := createIdentifierToken("x", 1)

		tests := []struct {
			name  string
			value any
		}{
			{"string", "hello"},
			{"boolean", true},
			{"nil", nil},
			{"decimal", 3.14},
			{"negative", -1.0},
		}

		for _, tt := range tests {
			err := env.assignVarValue(token, tt.value)
			assertNoError(t, err, "Assign "+tt.name)
			assertEqual(t, tt.value, env.values["x"], tt.name+" value")
		}
	})

	t.Run("Assign to multiple variables", func(t *testing.T) {
		env := createTestEnvironment(nil)
		env.defineVarValue("x", 1.0)
		env.defineVarValue("y", 2.0)
		env.defineVarValue("z", 3.0)

		// Assign to all variables
		xToken := createIdentifierToken("x", 1)
		yToken := createIdentifierToken("y", 2)
		zToken := createIdentifierToken("z", 3)

		err := env.assignVarValue(xToken, "new x")
		assertNoError(t, err, "Assign to x")

		err = env.assignVarValue(yToken, "new y")
		assertNoError(t, err, "Assign to y")

		err = env.assignVarValue(zToken, "new z")
		assertNoError(t, err, "Assign to z")

		// Verify all assignments
		assertEqual(t, "new x", env.values["x"], "Variable x")
		assertEqual(t, "new y", env.values["y"], "Variable y")
		assertEqual(t, "new z", env.values["z"], "Variable z")
	})
}

// ============================================================================
// SCOPE NESTING TESTS
// ============================================================================

func TestScopeNesting(t *testing.T) {
	t.Run("Three-level nesting", func(t *testing.T) {
		level1 := createTestEnvironment(nil)
		level1.defineVarValue("a", "level1")

		level2 := createTestEnvironment(level1)
		level2.defineVarValue("b", "level2")

		level3 := createTestEnvironment(level2)
		level3.defineVarValue("c", "level3")

		// Test retrieval from each level
		aToken := createIdentifierToken("a", 1)
		bToken := createIdentifierToken("b", 2)
		cToken := createIdentifierToken("c", 3)

		// From level3, should access all variables
		value, err := level3.getVarValue(aToken)
		assertNoError(t, err, "Access level1 variable from level3")
		assertEqual(t, "level1", value, "Level1 variable value")

		value, err = level3.getVarValue(bToken)
		assertNoError(t, err, "Access level2 variable from level3")
		assertEqual(t, "level2", value, "Level2 variable value")

		value, err = level3.getVarValue(cToken)
		assertNoError(t, err, "Access level3 variable from level3")
		assertEqual(t, "level3", value, "Level3 variable value")

		// From level2, should access level1 and level2 variables
		value, err = level2.getVarValue(aToken)
		assertNoError(t, err, "Access level1 variable from level2")
		assertEqual(t, "level1", value, "Level1 variable value")

		value, err = level2.getVarValue(bToken)
		assertNoError(t, err, "Access level2 variable from level2")
		assertEqual(t, "level2", value, "Level2 variable value")

		_, err = level2.getVarValue(cToken)
		assertError(t, err, "Access level3 variable from level2")
	})

	t.Run("Variable shadowing at multiple levels", func(t *testing.T) {
		level1 := createTestEnvironment(nil)
		level1.defineVarValue("x", "level1")

		level2 := createTestEnvironment(level1)
		level2.defineVarValue("x", "level2")

		level3 := createTestEnvironment(level2)
		level3.defineVarValue("x", "level3")

		token := createIdentifierToken("x", 1)

		// Each level should see its own shadowed version
		value, err := level1.getVarValue(token)
		assertNoError(t, err, "Access x from level1")
		assertEqual(t, "level1", value, "Level1 x value")

		value, err = level2.getVarValue(token)
		assertNoError(t, err, "Access x from level2")
		assertEqual(t, "level2", value, "Level2 x value")

		value, err = level3.getVarValue(token)
		assertNoError(t, err, "Access x from level3")
		assertEqual(t, "level3", value, "Level3 x value")
	})

	t.Run("Assignment through multiple levels", func(t *testing.T) {
		level1 := createTestEnvironment(nil)
		level1.defineVarValue("x", "original")

		level2 := createTestEnvironment(level1)
		level3 := createTestEnvironment(level2)

		token := createIdentifierToken("x", 1)

		// Assign from level3, should modify level1
		err := level3.assignVarValue(token, "modified from level3")
		assertNoError(t, err, "Assign from level3")
		assertEqual(t, "modified from level3", level1.values["x"], "Level1 variable modified")

		// Assign from level2, should modify level1
		err = level2.assignVarValue(token, "modified from level2")
		assertNoError(t, err, "Assign from level2")
		assertEqual(t, "modified from level2", level1.values["x"], "Level1 variable modified again")
	})

	t.Run("Mixed variable access patterns", func(t *testing.T) {
		level1 := createTestEnvironment(nil)
		level1.defineVarValue("global", "global value")
		level1.defineVarValue("shared", "level1 shared")

		level2 := createTestEnvironment(level1)
		level2.defineVarValue("local2", "level2 local")
		level2.defineVarValue("shared", "level2 shared")

		level3 := createTestEnvironment(level2)
		level3.defineVarValue("local3", "level3 local")

		// Test various access patterns
		globalToken := createIdentifierToken("global", 1)
		sharedToken := createIdentifierToken("shared", 2)
		local2Token := createIdentifierToken("local2", 3)
		local3Token := createIdentifierToken("local3", 4)

		// From level3
		value, err := level3.getVarValue(globalToken)
		assertNoError(t, err, "Access global from level3")
		assertEqual(t, "global value", value, "Global value")

		value, err = level3.getVarValue(sharedToken)
		assertNoError(t, err, "Access shared from level3")
		assertEqual(t, "level2 shared", value, "Shared value (level2)")

		value, err = level3.getVarValue(local2Token)
		assertNoError(t, err, "Access local2 from level3")
		assertEqual(t, "level2 local", value, "Local2 value")

		value, err = level3.getVarValue(local3Token)
		assertNoError(t, err, "Access local3 from level3")
		assertEqual(t, "level3 local", value, "Local3 value")
	})
}

// ============================================================================
// EDGE CASES AND ERROR SCENARIOS
// ============================================================================

func TestEnvironmentEdgeCases(t *testing.T) {

	t.Run("Nil value handling", func(t *testing.T) {
		env := createTestEnvironment(nil)

		env.defineVarValue("nil_var", nil)

		token := createIdentifierToken("nil_var", 1)
		value, err := env.getVarValue(token)

		assertNoError(t, err, "Retrieve nil variable")
		assertEqual(t, nil, value, "Nil variable value")
	})

	t.Run("Large number of variables", func(t *testing.T) {
		env := createTestEnvironment(nil)

		// Define many variables
		for i := 0; i < 1000; i++ {
			env.defineVarValue("var"+fmt.Sprintf("%d", i), i)
		}

		if len(env.values) != 1000 {
			t.Errorf("Expected 1000 variables, got %d", len(env.values))
		}

		// Test retrieval of some variables
		token := createIdentifierToken("var0", 1)
		value, err := env.getVarValue(token)
		assertNoError(t, err, "Retrieve first variable")
		assertEqual(t, 0, value, "First variable value")

		token = createIdentifierToken("var999", 2)
		value, err = env.getVarValue(token)
		assertNoError(t, err, "Retrieve last variable")
		assertEqual(t, 999, value, "Last variable value")
	})

	t.Run("Deep nesting", func(t *testing.T) {
		// Create a deep chain of environments
		current := createTestEnvironment(nil)
		current.defineVarValue("deep_var", "deep value")

		for i := 0; i < 100; i++ {
			next := createTestEnvironment(current)
			current = next
		}

		// Try to access the variable from the deepest level
		token := createIdentifierToken("deep_var", 1)
		value, err := current.getVarValue(token)

		assertNoError(t, err, "Retrieve variable from deep nesting")
		assertEqual(t, "deep value", value, "Deep variable value")
	})

	t.Run("Circular reference prevention", func(t *testing.T) {
		// This test ensures we don't accidentally create circular references
		env1 := createTestEnvironment(nil)
		env2 := createTestEnvironment(env1)

		// Verify the chain is correct
		if env2.enclosing != env1 {
			t.Error("Expected env2 to have env1 as parent")
		}
		if env1.enclosing != nil {
			t.Error("Expected env1 to have no parent")
		}

		// Verify we can traverse the chain
		env1.defineVarValue("x", "value")
		token := createIdentifierToken("x", 1)
		value, err := env2.getVarValue(token)

		assertNoError(t, err, "Retrieve variable through chain")
		assertEqual(t, "value", value, "Variable value through chain")
	})
}
