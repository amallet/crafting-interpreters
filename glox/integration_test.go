package main

import (
	"os"
	"strings"
	"testing"
)

// ============================================================================
// INTEGRATION TESTS - PHASE 2.3
// ============================================================================
// These tests verify end-to-end functionality through complete GLox programs
// covering multi-statement programs, complex control flow, and error propagation.

func TestIntegrationCompletePrograms(t *testing.T) {
	t.Run("Simple arithmetic program", func(t *testing.T) {
		program := `
var a = 10;
var b = 20;
var sum = a + b;
var product = a * b;
print sum;
print product;
`

		expected := []string{"30", "200"}
		runProgramAndCheckOutput(t, program, expected, "Simple arithmetic program")
	})

	t.Run("Complex expression evaluation", func(t *testing.T) {
		program := `
var x = 5;
var y = 3;
var result = (x + y) * 2 - 1;
print result;
var comparison = x > y and y < 10;
print comparison;
`

		expected := []string{"15", "true"}
		runProgramAndCheckOutput(t, program, expected, "Complex expression evaluation")
	})

	t.Run("String manipulation", func(t *testing.T) {
		program := `
var greeting = "Hello";
var name = "World";
var message = greeting + " " + name + "!";
print message;
var length_check = greeting == "Hello";
print length_check;
`

		expected := []string{"Hello World!", "true"}
		runProgramAndCheckOutput(t, program, expected, "String manipulation")
	})

	t.Run("Boolean logic", func(t *testing.T) {
		program := `
var a = true;
var b = false;
var c = true;

print a;
print b;
print c;
print !a;
print !b;
`

		expected := []string{"true", "false", "true", "false", "true"}
		runProgramAndCheckOutput(t, program, expected, "Boolean logic")
	})
}

func TestIntegrationMultiStatementPrograms(t *testing.T) {
	t.Run("Variable declarations and assignments", func(t *testing.T) {
		program := `
var x = 1;
var y = 2;
var z = 3;

print x;
print y;
print z;

x = x + 1;
y = y * 2;
z = z - 1;

print x;
print y;
print z;
`

		expected := []string{"1", "2", "3", "2", "4", "2"}
		runProgramAndCheckOutput(t, program, expected, "Variable declarations and assignments")
	})

	t.Run("Mixed statement types", func(t *testing.T) {
		program := `
var counter = 0;
print "Starting counter:";
print counter;

counter = counter + 1;
print "After increment:";
print counter;

var message = "Counter is now";
print message;
print counter;
`

		expected := []string{"Starting counter:", "0", "After increment:", "1", "Counter is now", "1"}
		runProgramAndCheckOutput(t, program, expected, "Mixed statement types")
	})

	t.Run("Expression statements", func(t *testing.T) {
		program := `
var x = 5;
x + 3;  // Expression statement
print x;  // Should still be 5
x = x + 3;  // Assignment statement
print x;  // Should be 8
`

		expected := []string{"5", "8"}
		runProgramAndCheckOutput(t, program, expected, "Expression statements")
	})
}

func TestIntegrationComplexControlFlow(t *testing.T) {
	t.Run("Nested if statements", func(t *testing.T) {
		program := `
var x = 10;
var y = 5;

if (x > y) {
  print "x is greater than y";
  if (x > 15) {
    print "x is also greater than 15";
  } else {
    print "x is not greater than 15";
  }
} else {
  print "x is not greater than y";
}
`

		expected := []string{"x is greater than y", "x is not greater than 15"}
		runProgramAndCheckOutput(t, program, expected, "Nested if statements")
	})

	t.Run("If-else chains", func(t *testing.T) {
		program := `
var score = 85;

if (score >= 90) {
  print "Grade: A";
} else if (score >= 80) {
  print "Grade: B";
} else if (score >= 70) {
  print "Grade: C";
} else {
  print "Grade: F";
}
`

		expected := []string{"Grade: B"}
		runProgramAndCheckOutput(t, program, expected, "If-else chains")
	})

	t.Run("While loop with counter", func(t *testing.T) {
		program := `
var i = 1;
while (i <= 5) {
  print i;
  i = i + 1;
}
print "Loop finished";
`

		expected := []string{"1", "2", "3", "4", "5", "Loop finished"}
		runProgramAndCheckOutput(t, program, expected, "While loop with counter")
	})

	t.Run("Nested loops", func(t *testing.T) {
		program := `
var outer = 1;
while (outer <= 2) {
  var inner = 1;
  while (inner <= 2) {
    print outer * 10 + inner;
    inner = inner + 1;
  }
  outer = outer + 1;
}
`

		expected := []string{"11", "12", "21", "22"}
		runProgramAndCheckOutput(t, program, expected, "Nested loops")
	})

	t.Run("For loop desugaring", func(t *testing.T) {
		program := `
for (var i = 1; i <= 3; i = i + 1) {
  print i;
}
print "For loop finished";
`

		expected := []string{"1", "2", "3", "For loop finished"}
		runProgramAndCheckOutput(t, program, expected, "For loop desugaring")
	})

	t.Run("Complex control flow", func(t *testing.T) {
		program := `
var x = 0;
var y = 0;

while (x < 3) {
  if (x == 1) {
    y = y + 10;
  } else {
    y = y + 1;
  }
  x = x + 1;
}

print "Final x:";
print x;
print "Final y:";
print y;
`

		expected := []string{"Final x:", "3", "Final y:", "12"}
		runProgramAndCheckOutput(t, program, expected, "Complex control flow")
	})
}

func TestIntegrationVariableScoping(t *testing.T) {
	t.Run("Basic variable scoping", func(t *testing.T) {
		program := `
var x = "global";
{
  var y = "local";
  print x; // Should print "global"
  print y; // Should print "local"
}
print x; // Should print "global"
// y is not accessible here
`

		expected := []string{"global", "local", "global"}
		runProgramAndCheckOutput(t, program, expected, "Basic variable scoping")
	})

	t.Run("Variable shadowing", func(t *testing.T) {
		program := `
var x = "global";
{
  var x = "local";
  print x; // Should print "local"
}
print x; // Should print "global"
`

		expected := []string{"local", "global"}
		runProgramAndCheckOutput(t, program, expected, "Variable shadowing")
	})

	t.Run("Assignment to shadowed variable", func(t *testing.T) {
		program := `
var x = "global";
{
  var x = "local";
  x = "modified local";
  print x; // Should print "modified local"
}
print x; // Should print "global"
`

		expected := []string{"modified local", "global"}
		runProgramAndCheckOutput(t, program, expected, "Assignment to shadowed variable")
	})

	t.Run("Assignment to parent variable", func(t *testing.T) {
		program := `
var x = "global";
{
  x = "modified global";
  print x; // Should print "modified global"
}
print x; // Should print "modified global"
`

		expected := []string{"modified global", "modified global"}
		runProgramAndCheckOutput(t, program, expected, "Assignment to parent variable")
	})

	t.Run("Deep nesting", func(t *testing.T) {
		program := `
var a = "level1";
{
  var b = "level2";
  {
    var c = "level3";
    print a; // Should print "level1"
    print b; // Should print "level2"
    print c; // Should print "level3"
  }
  print a; // Should print "level1"
  print b; // Should print "level2"
  // c is not accessible here
}
print a; // Should print "level1"
// b and c are not accessible here
`

		expected := []string{"level1", "level2", "level3", "level1", "level2", "level1"}
		runProgramAndCheckOutput(t, program, expected, "Deep nesting")
	})

	t.Run("Variable shadowing at multiple levels", func(t *testing.T) {
		program := `
var x = "level1";
{
  var x = "level2";
  {
    var x = "level3";
    print x; // Should print "level3"
  }
  print x; // Should print "level2"
}
print x; // Should print "level1"
`

		expected := []string{"level3", "level2", "level1"}
		runProgramAndCheckOutput(t, program, expected, "Variable shadowing at multiple levels")
	})

	t.Run("Assignment through multiple levels", func(t *testing.T) {
		program := `
var x = "original";
{
  {
    x = "modified from level3";
    print x; // Should print "modified from level3"
  }
  print x; // Should print "modified from level3"
}
print x; // Should print "modified from level3"
`

		expected := []string{"modified from level3", "modified from level3", "modified from level3"}
		runProgramAndCheckOutput(t, program, expected, "Assignment through multiple levels")
	})

	t.Run("Mixed variable access patterns", func(t *testing.T) {
		program := `
var global = "global value";
var shared = "level1 shared";
{
  var local2 = "level2 local";
  var shared = "level2 shared";
  {
    var local3 = "level3 local";
    print global; // Should print "global value"
    print shared; // Should print "level2 shared"
    print local2; // Should print "level2 local"
    print local3; // Should print "level3 local"
  }
  print global; // Should print "global value"
  print shared; // Should print "level2 shared"
  print local2; // Should print "level2 local"
  // local3 is not accessible here
}
print global; // Should print "global value"
print shared; // Should print "level1 shared"
// local2 and local3 are not accessible here
`

		expected := []string{
			"global value", "level2 shared", "level2 local", "level3 local",
			"global value", "level2 shared", "level2 local",
			"global value", "level1 shared",
		}
		runProgramAndCheckOutput(t, program, expected, "Mixed variable access patterns")
	})

	t.Run("Undefined variable access", func(t *testing.T) {
		program := `
print undefined; // Should cause runtime error
`

		runProgramAndExpectError(t, program, "Undefined variable 'undefined'", "Undefined variable access")
	})

	t.Run("Undefined variable assignment", func(t *testing.T) {
		program := `
undefined = "value"; // Should cause runtime error
`

		runProgramAndExpectError(t, program, "Undefined variable 'undefined'", "Undefined variable assignment")
	})

	t.Run("Variable redefinition", func(t *testing.T) {
		program := `
var x = "first";
var x = "second";
print x; // Should print "second"
`

		expected := []string{"second"}
		runProgramAndCheckOutput(t, program, expected, "Variable redefinition")
	})

	t.Run("Different variable types", func(t *testing.T) {
		program := `
var num = 42;
var str = "hello";
var bool = true;
var nil_var = nil;

print num;
print str;
print bool;
print nil_var;
`

		expected := []string{"42", "hello", "true", "<nil>"}
		runProgramAndCheckOutput(t, program, expected, "Different variable types")
	})

	t.Run("Variable type changes", func(t *testing.T) {
		program := `
var x = 42;
print x;
x = "now a string";
print x;
x = true;
print x;
x = nil;
print x;
`

		expected := []string{"42", "now a string", "true", "<nil>"}
		runProgramAndCheckOutput(t, program, expected, "Variable type changes")
	})

	t.Run("Assignment in nested scopes", func(t *testing.T) {
		program := `
var x = "original";
{
  var y = "local";
  x = "modified by local";
  y = "modified local";
  print x; // Should print "modified by local"
  print y; // Should print "modified local"
}
print x; // Should print "modified by local"
// y is not accessible here
`

		expected := []string{"modified by local", "modified local", "modified by local"}
		runProgramAndCheckOutput(t, program, expected, "Assignment in nested scopes")
	})

	t.Run("Shadowing with assignment", func(t *testing.T) {
		program := `
var x = "global";
{
  var x = "local";
  x = "modified local";
  print x; // Should print "modified local"
}
{
  x = "modified global";
  print x; // Should print "modified global"
}
print x; // Should print "modified global"
`

		expected := []string{"modified local", "modified global", "modified global"}
		runProgramAndCheckOutput(t, program, expected, "Shadowing with assignment")
	})

	t.Run("Loop variable scoping", func(t *testing.T) {
		program := `
var i = 100;  // This should be shadowed

for (var i = 1; i <= 2; i = i + 1) {
  print i;
}

print i;  // Should print 100, not 3
`

		expected := []string{"1", "2", "100"}
		runProgramAndCheckOutput(t, program, expected, "Loop variable scoping")
	})
}

func TestIntegrationErrorPropagation(t *testing.T) {
	t.Run("Runtime error in expression", func(t *testing.T) {
		program := `
var x = 10;
var y = 0;
var result = x / y;  // Division by zero
`

		runProgramAndExpectError(t, program, "illegal operation: division by zero", "Runtime error in expression")
	})

	t.Run("Runtime error in statement", func(t *testing.T) {
		program := `
var x = 10;
print x;
var y = undefined;  // Undefined variable
print y;
`

		runProgramAndExpectError(t, program, "Undefined variable 'undefined'", "Runtime error in statement")
	})

	t.Run("Type error in binary operation", func(t *testing.T) {
		program := `
var x = "hello";
var y = 5;
var error = x - y;   // String - number should fail
`

		runProgramAndExpectError(t, program, "operands to operator - must be numbers", "Type error in binary operation")
	})

	t.Run("Error in nested expression", func(t *testing.T) {
		program := `
var x = 10;
var y = 0;
var result = (x + 5) / (y + 0);  // Division by zero in nested expression
`

		runProgramAndExpectError(t, program, "illegal operation: division by zero", "Error in nested expression")
	})

	t.Run("Error in control flow", func(t *testing.T) {
		program := `
var x = 10;
if (x / 0 > 5) {  // Division by zero in condition
  print "This won't print";
}
`

		runProgramAndExpectError(t, program, "illegal operation: division by zero", "Error in control flow")
	})

	t.Run("Error in loop condition", func(t *testing.T) {
		program := `
var x = 10;
while (x / 0 > 5) {  // Division by zero in loop condition
  print "This won't print";
}
`

		runProgramAndExpectError(t, program, "illegal operation: division by zero", "Error in loop condition")
	})
}

func TestIntegrationComplexScenarios(t *testing.T) {
	t.Run("Fibonacci sequence", func(t *testing.T) {
		program := `
var a = 0;
var b = 1;
var count = 0;

while (count < 6) {
  print a;
  var temp = a + b;
  a = b;
  b = temp;
  count = count + 1;
}
`

		expected := []string{"0", "1", "1", "2", "3", "5"}
		runProgramAndCheckOutput(t, program, expected, "Fibonacci sequence")
	})

	// Note: Prime number checker removed due to complex modulo logic
	// This would require proper modulo operator implementation

	// Note: String palindrome checker removed - would require string indexing
	// which is not implemented in GLox yet

	t.Run("Nested scoping with complex logic", func(t *testing.T) {
		program := `
var global_counter = 0;

{
  var local_counter = 0;
  
  for (var i = 1; i <= 3; i = i + 1) {
    local_counter = local_counter + i;
    global_counter = global_counter + i;
    
    {
      var inner_counter = local_counter * 2;
      print "Iteration";
      print i;
      print "inner";
      print inner_counter;
    }
  }
  
  print "Local counter:";
  print local_counter;
}

print "Global counter:";
print global_counter;
`

		expected := []string{
			"Iteration", "1", "inner", "2",
			"Iteration", "2", "inner", "6",
			"Iteration", "3", "inner", "12",
			"Local counter:", "6",
			"Global counter:", "6",
		}
		runProgramAndCheckOutput(t, program, expected, "Nested scoping with complex logic")
	})
}

func TestIntegrationEdgeCases(t *testing.T) {
	t.Run("Empty program", func(t *testing.T) {
		program := ``
		expected := []string{}
		runProgramAndCheckOutput(t, program, expected, "Empty program")
	})

	t.Run("Single statement", func(t *testing.T) {
		program := `print "Hello, World!";`
		expected := []string{"Hello, World!"}
		runProgramAndCheckOutput(t, program, expected, "Single statement")
	})

	t.Run("Multiple empty blocks", func(t *testing.T) {
		program := `
{
  {
    {
      print "Deep nesting";
    }
  }
}
`

		expected := []string{"Deep nesting"}
		runProgramAndCheckOutput(t, program, expected, "Multiple empty blocks")
	})

	t.Run("Variable redefinition in same scope", func(t *testing.T) {
		program := `
var x = 1;
var x = 2;
print x;
`

		expected := []string{"2"}
		runProgramAndCheckOutput(t, program, expected, "Variable redefinition in same scope")
	})

	t.Run("Complex expression with all operators", func(t *testing.T) {
		program := `
var a = 10;
var b = 3;
var c = 2;

var result = (a + b) * c - a / b;
print result;
var comparison = a > b;
print comparison;
var equality = c == 2;
print equality;
`

		expected := []string{"22.666666666666668", "true", "true"}
		runProgramAndCheckOutput(t, program, expected, "Complex expression with all operators")
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func runProgramAndCheckOutput(t *testing.T, program string, expected []string, testName string) {
	t.Helper()

	// Create a new GLox instance
	glox := &GLox{}
	glox.interpreter = NewInterpreter(glox)

	// Capture output by intercepting the print statements
	var output []string
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the program
	glox.run(program, false)

	// Restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	capturedOutput := string(buf[:n])

	// Parse output lines
	lines := strings.Split(strings.TrimSpace(capturedOutput), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			output = append(output, strings.TrimSpace(line))
		}
	}

	// Check for errors
	if glox.hadError || glox.hadRuntimeError {
		t.Errorf("%s: Program execution failed with errors", testName)
		return
	}

	// Check output
	if len(output) != len(expected) {
		t.Errorf("%s: Expected %d output lines, got %d", testName, len(expected), len(output))
		t.Errorf("Expected: %v", expected)
		t.Errorf("Got: %v", output)
		return
	}

	for i, expectedLine := range expected {
		if i >= len(output) {
			t.Errorf("%s: Missing output line %d, expected: %s", testName, i, expectedLine)
			continue
		}
		if output[i] != expectedLine {
			t.Errorf("%s: Output line %d mismatch. Expected: %s, Got: %s", testName, i, expectedLine, output[i])
		}
	}
}

func runProgramAndExpectError(t *testing.T, program string, expectedError string, testName string) {
	t.Helper()

	// Create a new GLox instance
	glox := &GLox{}
	glox.interpreter = NewInterpreter(glox)

	// Capture stderr to check for error messages
	originalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Run the program and expect an error
	glox.run(program, false)

	// Restore stderr
	w.Close()
	os.Stderr = originalStderr

	// Read captured error output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	capturedError := string(buf[:n])

	// Check if we got the expected error
	if !glox.hadRuntimeError {
		t.Errorf("%s: Expected runtime error but got none", testName)
		return
	}

	// Check if the error message contains what we expect
	if !strings.Contains(capturedError, expectedError) {
		t.Errorf("%s: Expected error message to contain '%s', got: %s", testName, expectedError, capturedError)
	}
}
