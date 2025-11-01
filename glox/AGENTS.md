# CLAUDE.md

This file provides guidance to AI agents when working with code in this repository.

## Project Overview

GLox is a Go implementation of the Lox programming language from "Crafting Interpreters" by Robert Nystrom. It implements a tree-walk interpreter featuring a complete lexical analyzer (scanner), recursive descent parser, and expression/statement evaluator. The codebase follows the book's structure but adapts the design patterns to idiomatic Go.

**Prerequisites**: Go 1.24.1 or later

## Development Commands

### Building
```bash
go build -o glox
```

### Testing
```bash
# Run all tests
go test

# Run specific test suite
go test -run TestScanner
go test -run TestParser
go test -run TestInterpreter
go test -run TestIntegration

# Verbose output
go test -v

# Run single test
go test -run TestParserEmptyInput
```

### Running
```bash
# Interactive REPL mode
./glox

# Execute a Lox file
./glox script.lox
```

## Architecture

### Three-Stage Pipeline

GLox processes Lox source code through three distinct stages:

1. **Scanner (scanner.go)** - Lexical analysis
   - Converts source text into tokens
   - Handles keywords, identifiers, literals (strings/numbers), operators
   - Uses rune-based scanning for proper Unicode support
   - Keywords are resolved via `reservedKeyWordMap`

2. **Parser (parser.go)** - Syntax analysis
   - Recursive descent parser implementing Lox grammar
   - Builds Abstract Syntax Tree (AST) from tokens
   - Handles operator precedence automatically through grammar structure
   - **For loops are desugared**: Parser transforms `for` loops into `while` loops with initialization blocks (see parser.go:214-290)
   - Error recovery via `synchronize()` method for better REPL experience

3. **Interpreter (interpreter.go)** - Execution
   - Tree-walk interpreter using visitor pattern
   - Evaluates AST nodes and executes statements
   - Manages runtime state through Environment chains

### AST Node Types

**Expression nodes (expr.go)**: Implement `Expr` interface with `Accept(ExprVisitor)` method
- `BinaryExpr` - Binary operations: +, -, *, /, ==, !=, <, >, <=, >=
- `UnaryExpr` - Unary operations: !, - (negation)
- `LiteralExpr` - Numbers (float64), strings, booleans, nil
- `GroupingExpr` - Parenthesized expressions
- `VariableExpr` - Variable references
- `AssignExpr` - Variable assignments (returns assigned value)
- `LogicalExpr` - Short-circuiting logical operations: and, or
- `CallExpr` - Function calls with callee expression and argument list

**Statement nodes (stmt.go)**: Implement `Stmt` interface with `Accept(StmtVisitor)` method
- `ExpressionStmt` - Expression statements (discards result)
- `PrintStmt` - Print statements (outputs to stdout)
- `VarStmt` - Variable declarations with optional initializer
- `BlockStmt` - Scoped blocks containing multiple statements
- `IfStmt` - Conditional branching with optional else
- `WhileStmt` - Loop statements
- `FunctionStmt` - Function declarations with name, parameters, and body
- `ReturnStmt` - Return statements (with optional return value)

### Visitor Pattern Implementation

Go interfaces are used to implement the visitor pattern:
- `ExprVisitor` interface defines visit methods for each expression type
- `StmtVisitor` interface defines visit methods for each statement type
- `Interpreter` implements both visitor interfaces
- Each AST node calls `visitor.Visit<NodeType>()` in its `Accept()` method

### Variable Scoping (environment.go)

The `Environment` type implements lexical scoping through a chain of environments:
- Each `Environment` has a map of variable names to values
- Pointer to enclosing (parent) environment for nested scopes
- Variable lookup walks up the environment chain until variable is found
- Variable assignment searches current scope first, then parent scopes
- Block statements create new child environments (see interpreter.go:117-133)
- Function calls create new child environments with parameter bindings

### Function Calling (lox_callable.go, lox_function.go)

The `LoxCallable` interface abstracts callable entities (functions and built-ins):
- `arity()` - Returns the number of parameters the callable expects
- `call(interpreter, arguments)` - Executes the callable with given arguments

`LoxFunction` implements `LoxCallable` and wraps user-defined functions:
- Each function creates a new environment chained to the global environment
- Parameters are bound to argument values in the function's environment
- Return values are propagated via `ReturnValue` error wrapper to unwind the call stack
- Functions can be stored in variables and passed as values

Built-in functions (builtin_fns.go):
- `clock()` - Returns current Unix time in milliseconds (0 parameters)

### Runtime Abstraction (runtime.go)

The `LoxRuntime` interface abstracts error reporting to allow:
- `GLox` for production use (glox.go)
- `TestGLox` for testing (test_utils.go)
- Unified error handling across scanner, parser, and interpreter

### Error Handling

**Parse errors**: Reported with line numbers and token context via `parseError()`
- Parser synchronizes to next statement boundary on error
- Allows REPL to continue after syntax errors

**Runtime errors**: Wrapped in `RuntimeError` type with token location
- Type checking failures (e.g., adding string to number)
- Undefined variable access
- Division by zero

## Project Structure

```
glox/
├── glox.go              # Main entry point and REPL
├── scanner.go           # Lexical analysis (tokenizer)
├── parser.go            # Syntax analysis (parser)
├── interpreter.go       # Tree-walk interpreter
├── expr.go              # Expression AST nodes
├── stmt.go              # Statement AST nodes
├── environment.go       # Variable environment management
├── token.go             # Token representation
├── token_type.go        # Token type definitions
├── runtime.go           # Runtime error handling
├── runtime_error.go     # Runtime error types
├── return_value.go      # Return value propagation mechanism
├── lox_callable.go      # Callable interface for functions
├── lox_function.go      # User-defined function implementation
├── builtin_fns.go       # Built-in function implementations
├── test_utils.go        # Testing utilities
├── *_test.go           # Comprehensive test suites
└── go.mod              # Go module definition
```

## Testing Strategy

Tests are organized by component:
- `scanner_test.go` - Token recognition, string/number scanning, error cases
- `parser_basic_test.go` - Expression parsing
- `parser_statements_test.go` - Statement parsing
- `parser_precedence_test.go` - Operator precedence verification
- `parser_errors_test.go` - Error handling and recovery
- `parser_for_loops_test.go` - For loop desugaring
- `interpreter_test.go` - Expression evaluation, statement execution
- `environment_test.go` - Variable scoping behavior
- `integration_test.go` - End-to-end program execution

### Test Utilities (test_utils.go)

`TestGLox` implements `LoxRuntime` for testing:
- Captures errors instead of printing to stderr
- Provides `getErrors()` and `clearErrors()` for test assertions
- Use `MakeTokens()` helper to create token slices for parser tests

## Implementation Notes

- All numbers are stored as `float64` internally
- String scanning supports multi-line strings
- Comments are single-line only (// style)
- Exit codes: 64 (usage error), 65 (parse error), 70 (runtime error)
- REPL mode prints expression results automatically
- The interpreter maintains a global environment that persists across REPL inputs
- Assignment is an expression (returns the assigned value) not a statement
- Functions are first-class values that can be stored in variables and passed around
- Function calls create new environments chained to the global environment (not the caller's environment)
- Return values are propagated using `ReturnValue` error wrapper to unwind the call stack
- Maximum 255 parameters and 255 arguments are enforced by the parser
- Go interfaces are used to implement the visitor pattern (vs. classes in Java)
- Go's type system is leveraged for AST node definitions
- Error handling follows Go conventions (returning errors vs. throwing exceptions)

## Lox Language Reference

### Supported Language Features
- **Variables**: Declaration and assignment with `var` keyword
- **Data Types**: Numbers (floating-point), strings, booleans, and `nil`
- **Arithmetic Operations**: `+`, `-`, `*`, `/` with proper precedence
- **Comparison Operations**: `==`, `!=`, `<`, `<=`, `>`, `>=`
- **Logical Operations**: `and`, `or`, `!` (not) with short-circuit evaluation
- **Control Flow**: `if`/`else`, `while` loops, `for` loops (desugared to while)
- **Blocks**: Lexical scoping with `{}` blocks
- **Functions**: Function declarations with parameters and return statements
- **Function Calls**: Call expressions with argument passing and return value handling
- **Print Statement**: Built-in `print` statement for output
- **Built-in Functions**: `clock()` function for getting current time
- **Comments**: Single-line comments with `//`

### Grammar and Precedence

Parser implements this grammar (listed by increasing precedence):

```
program        → declaration* EOF
declaration    → funDecl | varDecl | statement
funDecl        → "fun" function
function       → IDENTIFIER "(" parameters? ")" block
parameters     → IDENTIFIER ( "," IDENTIFIER )*
varDecl        → "var" IDENTIFIER ("=" expression)? ";"
statement      → exprStmt | ifStmt | printStmt | whileStmt | forStmt | returnStmt | block
returnStmt     → "return" expression? ";"
block          → "{" declaration* "}"
expression     → assignment
assignment     → IDENTIFIER "=" assignment | logic_or
logic_or       → logic_and ( "or" logic_and )*
logic_and      → equality ( "and" equality )*
equality       → comparison ( ( "!=" | "==" ) comparison )*
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
term           → factor ( ( "-" | "+" ) factor )*
factor         → unary ( ( "/" | "*" ) unary )*
unary          → ( "!" | "-" ) unary | call
call           → primary ( "(" arguments? ")" )*
arguments      → expression ( "," expression )*
primary        → "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")" | IDENTIFIER
```

Precedence levels (lowest to highest):
1. Assignment (right-associative)
2. Logical OR
3. Logical AND
4. Equality (==, !=)
5. Comparison (<, <=, >, >=)
6. Term (+, -)
7. Factor (*, /)
8. Unary (!, -)
9. Call (function calls)
10. Primary (literals, variables, grouping)

Note: Function calls have higher precedence than unary operators, allowing expressions like `-factorial(5)`.

### Example Programs

**Fibonacci Sequence**
```lox
var a = 0;
var b = 1;
var temp;

for (var i = 0; i < 10; i = i + 1) {
    print a;
    temp = a + b;
    a = b;
    b = temp;
}
```

**Variable Scoping Demo**
```lox
var a = "global a";
var b = "global b";
var c = "global c";
{
    var a = "outer a";
    var b = "outer b";
    {
        var a = "inner a";
        print a;  // "inner a"
        print b;  // "outer b"
        print c;  // "global c"
    }
    print a;      // "outer a"
    print b;      // "outer b"
    print c;      // "global c"
}
print a;          // "global a"
print b;          // "global b"
print c;          // "global c"
```

**Function Definition and Calling**
```lox
fun greet(name) {
    print "Hello, " + name + "!";
}

greet("World");  // "Hello, World!"

fun factorial(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

print factorial(5);  // 120

fun add(a, b) {
    return a + b;
}

var sum = add(3, 4);
print sum;  // 7
```
