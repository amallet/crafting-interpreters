package main

import (
	"testing"
)

// ============================================================================
// RESOLVER ERROR TESTS
// ============================================================================

func TestResolverErrors(t *testing.T) {
	t.Run("Can't read local variable in its own initializer", func(t *testing.T) {
		program := `
var a = "outer";
{
  var a = a; // Error: Can't read local variable in its own initializer
}
`

		runProgramAndExpectError(t, program, "Can't read local variable in its own initializer", "Read variable in own initializer")
	})

	t.Run("Can't return from top-level code", func(t *testing.T) {
		program := `return "wat"; // Error: Can't return from top-level code`

		runProgramAndExpectError(t, program, "Can't return from top-level code", "Return at top level")
	})

	t.Run("Return with value at top level", func(t *testing.T) {
		program := `return 42; // Error: Can't return from top-level code`

		runProgramAndExpectError(t, program, "Can't return from top-level code", "Return with value at top level")
	})

	t.Run("Return without value at top level", func(t *testing.T) {
		program := `return; // Error: Can't return from top-level code`

		runProgramAndExpectError(t, program, "Can't return from top-level code", "Return without value at top level")
	})

	t.Run("Variable redeclaration in same local scope", func(t *testing.T) {
		program := `
{
  var a = "first";
  var a = "second"; // Error: Already a variable with this name in this scope
}
`

		runProgramAndExpectError(t, program, "Already a variable with this name in this scope", "Variable redeclaration in local scope")
	})

	t.Run("Variable redeclaration in nested block", func(t *testing.T) {
		program := `
{
  var a = "first";
  {
    var a = "second"; // This is OK - different scope
    var a = "third";  // Error: Already a variable with this name in this scope
  }
}
`

		runProgramAndExpectError(t, program, "Already a variable with this name in this scope", "Variable redeclaration in nested block")
	})

	t.Run("Variable redeclaration in function parameters", func(t *testing.T) {
		program := `
fun test(a, a) { // Error: Already a variable with this name in this scope
  return a;
}
`

		runProgramAndExpectError(t, program, "Already a variable with this name in this scope", "Variable redeclaration in function parameters")
	})

	t.Run("Variable redeclaration with parameter name", func(t *testing.T) {
		program := `
fun test(a) {
  var a = "local"; // Error: Already a variable with this name in this scope
  return a;
}
`

		runProgramAndExpectError(t, program, "Already a variable with this name in this scope", "Variable redeclaration with parameter name")
	})

	t.Run("Unused local variable in block", func(t *testing.T) {
		program := `
{
  var unused = "value"; // Error: Unused variable
}
`

		runProgramAndExpectError(t, program, "Unused variable", "Unused local variable in block")
	})

	t.Run("Unused local variable in function", func(t *testing.T) {
		program := `
fun test() {
  var unused = "value"; // Error: Unused variable
}
test();
`

		runProgramAndExpectError(t, program, "Unused variable", "Unused local variable in function")
	})

	t.Run("Unused variable in nested block", func(t *testing.T) {
		program := `
{
  var used = "outer";
  {
    var unused = "inner"; // Error: Unused variable
    print used;
  }
}
`

		runProgramAndExpectError(t, program, "Unused variable", "Unused variable in nested block")
	})

	t.Run("Multiple unused variables", func(t *testing.T) {
		program := `
{
  var unused1 = "first"; // Error: Unused variable
  var unused2 = "second"; // Error: Unused variable
}
`

		runProgramAndExpectError(t, program, "Unused variable", "Multiple unused variables")
	})

	t.Run("Unused function parameter", func(t *testing.T) {
		program := `
fun test(param) { // Error: Unused variable
  print "hello";
}
test("arg");
`

		runProgramAndExpectError(t, program, "Unused variable", "Unused function parameter")
	})

	t.Run("Unused variable that shadows global", func(t *testing.T) {
		program := `
var global = "global";
{
  var global = "local"; // Error: Unused variable (even though it shadows)
}
`

		runProgramAndExpectError(t, program, "Unused variable", "Unused variable that shadows global")
	})
}

// ============================================================================
// RESOLVER SCOPE RESOLUTION TESTS
// ============================================================================

func TestResolverScopeResolution(t *testing.T) {
	// Note: Basic scope resolution, nested scopes, and assignment tests are already
	// covered in integration_test.go. These tests focus on resolver-specific resolution
	// behaviors, particularly around function scopes and parameter resolution.

	t.Run("Function parameters resolve correctly", func(t *testing.T) {
		program := `
fun test(param) {
  print param; // Should resolve to parameter
  var local = "local";
  print local; // Should resolve to local variable
}
test("argument");
`

		expected := []string{"argument", "local"}
		runProgramAndCheckOutput(t, program, expected, "Function parameters resolve correctly")
	})
}
