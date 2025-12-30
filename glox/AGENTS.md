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

### Four-Stage Pipeline

GLox processes Lox source code through four distinct stages:

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

3. **Resolver (resolver.go)** - Static analysis
   - Performs semantic analysis before interpretation
   - Resolves variable references to their correct scopes (local vs global)
   - Stores resolution information (scope distance) in interpreter's `locals` map
   - Detects static errors:
     - Reading a local variable in its own initializer
     - Returning from top-level code (outside functions)
     - Variable redeclaration in the same local scope
     - Unused local variables (variables declared but never read or assigned to)
   - Enables efficient variable lookup during interpretation

4. **Interpreter (interpreter.go)** - Execution
   - Tree-walk interpreter using visitor pattern
   - Evaluates AST nodes and executes statements
   - Manages runtime state through Environment chains
   - Uses resolver's scope distance information for efficient variable access

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
- `PropGetExpr` - Property access on instances: `instance.property`
- `PropSetExpr` - Property assignment on instances: `instance.property = value`
- `ThisExpr` - `this` keyword expression (refers to the current instance in a method)

**Statement nodes (stmt.go)**: Implement `Stmt` interface with `Accept(StmtVisitor)` method
- `ExpressionStmt` - Expression statements (discards result)
- `PrintStmt` - Print statements (outputs to stdout)
- `VarStmt` - Variable declarations with optional initializer
- `BlockStmt` - Scoped blocks containing multiple statements
- `IfStmt` - Conditional branching with optional else
- `WhileStmt` + `IfStmt` - Loop statements
- `FunctionStmt` - Function declarations with name, parameters, and body
- `ReturnStmt` - Return statements (with optional return value)
- `ClassStmt` - Class declarations with name and methods

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

### Variable Resolution (resolver.go)

The `Resolver` performs static analysis to resolve variable references:
- Traverses the AST before interpretation to determine variable scope distances
- Stores resolution information in interpreter's `locals` map (maps expressions to scope distances)
- Enables efficient variable lookup: local variables use `getAt(distance, name)` instead of walking the environment chain
- Detects static semantic errors:
  - **Self-reference in initializer**: `var a = a;` - can't read local variable in its own initializer
  - **Top-level return**: `return value;` at top level - can't return from top-level code
  - **Variable redeclaration**: `var a = 1; var a = 2;` in same scope - already a variable with this name
  - **Unused local variables**: `var unused = 1;` - local variables that are declared but never read or assigned to
- Tracks function context to validate return statements are only inside functions
- Manages scope stack to track variable declaration/definition status (prevents reading uninitialized variables)
- **`this` keyword handling**: When resolving a class declaration, the resolver injects `this` into the method scope using `injectThis()`, allowing methods to reference the instance via `this`
- **`this` resolution**: The `VisitThisExpr` method resolves `this` as a local variable, which will be bound to the instance when the method is called
- **Initializer detection**: The resolver tracks when it's inside an `init()` method using `functionTypeInitializer` to enforce that initializers cannot return values


### Function Calling (lox_callable.go, lox_function.go)

The `LoxCallable` interface abstracts callable entities (functions and built-ins):
- `arity()` - Returns the number of parameters the callable expects
- `call(interpreter, arguments)` - Executes the callable with given arguments

`LoxFunction` implements `LoxCallable` and wraps user-defined functions:
- Functions capture their enclosing environment (closure) at declaration time
- Each function call creates a new environment chained to the captured closure (not the global environment)
- Parameters are bound to argument values in the function's environment
- Return values are propagated via `ReturnValue` error wrapper to unwind the call stack
- Functions can be stored in variables and passed as values
- Functions can access and modify variables from their enclosing scope even after that scope exits
- **Method binding**: When a method is accessed on an instance, `LoxInstance.get()` calls `bind()` to create a bound method before returning it, enabling methods to access instance state


**Closures**: Functions capture variables from their lexical environment. This enables:
- Functions returning other functions that maintain independent state
- Function factories (e.g., `makeAdder(n)` returning a function that adds `n`)
- Multiple closures with independent captured variables
- Closures persisting after their enclosing scope exits

Built-in functions (builtin_fns.go):
- `clock()` - Returns current Unix time in milliseconds (0 parameters)

### Classes and Instances (lox_class.go, lox_instance.go)

**LoxClass** represents a class definition:
- Classes are callable (implement `LoxCallable` interface) and can be instantiated
- Calling a class (e.g., `ClassName()`) creates a new `LoxInstance`
- Classes can be stored in variables and passed around as values
- Classes have a name that is used for string representation
- **Single inheritance**: a class may optionally declare a superclass using `class Sub < Super { ... }`

**LoxInstance** represents an instance of a class:
- Each instance has a reference to its class (`klass`)
- Instances have a `fields` map that stores property values
- Properties are accessed via `instance.property` (PropGetExpr)
- Properties are set via `instance.property = value` (PropSetExpr)
- Getting an undefined property raises a runtime error: "Undefined property name X"
- Setting a property on a non-instance raises a runtime error: "Only instances have fields"
- Multiple instances of the same class have independent field storage
- Methods are defined in class declarations: `class ClassName { methodName() { ... } }`
- Methods are accessed via `instance.method` and are automatically bound to the instance

**`super` keyword (single inheritance)**:
- `super.method` refers to the *immediate* superclass of the current class
- `super.method()` evaluates by locating the method on the superclass and binding it to the current instance (`this`)
- `super` is only valid inside **subclass** method bodies

**Getter Functions**:
- Getter functions are methods with no parameter list that execute immediately when accessed
- Syntax: `class ClassName { propertyName { <code that returns a value> }` (no parentheses after the name)
- Getters are accessed like properties: `instance.propertyName` (not `instance.propertyName()`)
- When a property is accessed, if no field exists with that name, the class is checked for a getter with that name
- If a getter is found, it is executed immediately and its return value is used
- Getters can use `this`, call other methods, access fields, and use all language features
- Fields shadow getters: if an instance field exists with the same name, the field is returned instead of calling the getter
- Getters must be defined inside a class (cannot be top-level functions)
- Getters cannot have the same name as regular methods in the same class (duplicate method names are not allowed)
- The `init()` method cannot be a getter (it must have a parameter list, even if empty)

**Initializers (init() method)**:
- The `init()` method is a special constructor method that is called automatically when a class is instantiated
- Syntax: `class ClassName { init(parameters) { ... } }`
- The `init()` method can accept optional parameters that are passed during instantiation (e.g., `ClassName(arg1, arg2)`)
- The `init()` method always returns the initialized instance (cannot return any other value)
- If `init()` has an explicit `return` statement without a value, it still returns the instance
- The `init()` method can be called explicitly as a regular method on an existing instance (e.g., `instance.init(args)`)
- The resolver detects and reports errors when `init()` attempts to return a value (e.g., `return "value"` is not allowed)
- The `init()` method can access and modify instance fields using `this`, call other methods, and use all language features


### Runtime Abstraction (runtime.go)

The `LoxRuntime` interface abstracts error reporting to allow:
- `GLox` for production use (glox.go)
- `TestGLox` for testing (test_utils.go)
- Unified error handling across scanner, parser, and interpreter

### Error Handling

**Parse errors**: Reported with line numbers and token context via `parseError()`
- Parser synchronizes to next statement boundary on error
- Allows REPL to continue after syntax errors

**Resolver errors**: Static semantic errors detected during resolution phase
- Reported via `parseError()` with line numbers and token context
- Errors include:
  - Reading local variable in its own initializer
  - Returning from top-level code
  - Variable redeclaration in same scope
  - Unused local variables (declared but never read or assigned to)
  - Using `this` outside of a class (if validation is implemented)
  - Using `super` outside a class
  - Using `super` in a class with no superclass
  - A class can't inherit from itself (`class A < A`)
  - Returning a value from an initializer (`init()` method cannot return values)
  - Duplicate method names in the same class (including getters and cross-type duplicates)
  - Getter functions defined outside of a class
- Execution stops if resolver errors are detected

**Runtime errors**: Wrapped in `RuntimeError` type with token location
- Type checking failures (e.g., adding string to number)
- Undefined variable access (including `this` when used outside a method context)
- Division by zero
- Undefined property access on instances
- Property access/assignment on non-instances (classes, nil, numbers, strings, etc.)
- Calling undefined methods on instances
- Method arity mismatches (wrong number of arguments)
- Initializer arity mismatches (wrong number of arguments when instantiating a class with `init()`)
- Invalid inheritance target: superclass must be a class (`Not a class.`)
- `super` method lookup failures (`Undefined property <name>.`)

## Project Structure

```
glox/
├── glox.go              # Main entry point and REPL
├── scanner.go           # Lexical analysis (tokenizer)
├── parser.go            # Syntax analysis (parser)
├── resolver.go          # Static analysis and variable resolution
├── interpreter.go       # Tree-walk interpreter
├── expr.go              # Expression AST nodes
├── stmt.go              # Statement AST nodes
├── environment.go       # Variable environment management
├── token.go             # Token representation
├── token_type.go        # Token type definitions
├── runtime.go           # Runtime error handling
├── runtime_error.go     # Runtime error types
├── return_value.go      # Return value propagation mechanism
├── lox_callable.go      # Callable interface for functions and classes
├── lox_function.go      # User-defined function implementation
├── lox_class.go         # Class definition and instantiation
├── lox_instance.go      # Instance property management
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
- `resolver_test.go` - Variable resolution, static error detection
- `interpreter_test.go` - Expression evaluation, statement execution
- `environment_test.go` - Variable scoping behavior
- `classes_test.go` - Class definition, instantiation, and property access
- `inheritance_test.go` - Single inheritance and `super` semantics (including deep hierarchies and misuse errors)
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
- Functions capture their enclosing environment at declaration time (closures)
- Function calls create new environments chained to the captured closure (enables closure behavior)
- Return values are propagated using `ReturnValue` error wrapper to unwind the call stack
- Maximum 255 parameters and 255 arguments are enforced by the parser
- Variable resolution happens before interpretation for efficient lookup and static error detection
- Resolver stores scope distances in interpreter's `locals` map for O(1) local variable access
- Unused variable detection: Resolver tracks variable usage and reports errors for unused local variables and function parameters
- Classes are callable: Classes implement `LoxCallable` interface and can be instantiated by calling them like functions
- Instances have independent field storage: Each instance maintains its own `fields` map, allowing multiple instances of the same class to have different property values
- Property access is dynamic: Properties are accessed and set at runtime, with no compile-time checking
- Methods are bound to instances: When a method is accessed on an instance (e.g., `instance.method`), `LoxInstance.get()` calls `LoxFunction.bind()` to create a bound method with `this` bound to that instance
- Method binding: The `bind()` method creates a new `LoxFunction` with an environment that has `this` defined as the instance, enabling methods to access instance state
- Getter execution: When a property is accessed on an instance, `LoxInstance.get()` first checks for a field, then checks for a method. If the method is a getter (`isGetter` flag), it executes immediately and returns the result instead of returning the function
- Duplicate method detection: The resolver tracks method names during class resolution and reports an error if duplicate method names are found (applies to regular methods, getters, and cross-type duplicates)
- `this` keyword: The resolver injects `this` into the scope when resolving methods within a class declaration, and the interpreter resolves `this` as a variable lookup
- `this` in closures: When methods are stored in closures, `this` is correctly captured and refers to the original instance even when the method is called later
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
- **Closures**: Functions capture variables from their enclosing lexical scope
- **Classes**: Class declarations with methods (syntax: `class ClassName { method() {} }`)
- **Single inheritance**: Class declarations may specify a superclass (e.g., `class Sub < Super { ... }`)
- **Instances**: Class instantiation via constructor call (e.g., `var obj = ClassName()`)
- **Initializers**: The `init()` method acts as a constructor that is automatically called during instantiation (e.g., `class Point { init(x, y) { this.x = x; this.y = y; } }`)
- **Properties**: Dynamic property access (`instance.field`) and assignment (`instance.field = value`)
- **Methods**: Methods defined on classes can be called on instances (e.g., `instance.method()`)
- **Getter Functions**: Getter functions are methods with no parameter list that execute immediately when accessed as properties (e.g., `class Foo { value { return 42; } }` accessed as `instance.value`)
- **This Keyword**: The `this` keyword refers to the current instance within a method, enabling methods to access instance fields and call other methods
- **Super Keyword**: The `super` keyword refers to the superclass inside subclass methods (syntax: `super.method` / `super.method()`)
- **Print Statement**: Built-in `print` statement for output
- **Built-in Functions**: `clock()` function for getting current time
- **Comments**: Single-line comments with `//`

### Grammar and Precedence

Parser implements this grammar (listed by increasing precedence):

```
program        → declaration* EOF
declaration    → classDecl | funDecl | varDecl | statement
classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" method* "}"
method         → IDENTIFIER ("(" parameters? ")")? block
funDecl        → "fun" function
function       → IDENTIFIER "(" parameters? ")" block
parameters     → IDENTIFIER ( "," IDENTIFIER )*
varDecl        → "var" IDENTIFIER ("=" expression)? ";"
statement      → exprStmt | ifStmt | printStmt | whileStmt | forStmt | returnStmt | block
returnStmt     → "return" expression? ";"
block          → "{" declaration* "}"
expression     → assignmentOrValue
assignmentOrValue    → (call ".")? IDENTIFIER "=" assignment  | logic_or
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )*
logic_or       → logic_and ( "or" logic_and )*
logic_and      → equality ( "and" equality )*
equality       → comparison ( ( "!=" | "==" ) comparison )*
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
term           → factor ( ( "-" | "+" ) factor )*
factor         → unary ( ( "/" | "*" ) unary )*
unary          → ( "!" | "-" ) unary | call
call           → primary ( "(" arguments? ")" )*
arguments      → expression ( "," expression )*
primary        → "true" | "false" | "nil" | "this" | NUMBER | STRING | "(" expression ")" | IDENTIFIER | "super" "." IDENTIFIER
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

**Closures**
```lox
// Function factory - returns a function with captured state
fun makeCounter() {
    var i = 0;
    fun count() {
        i = i + 1;
        return i;
    }
    return count;
}

var counter = makeCounter();
print counter();  // 1
print counter();  // 2
print counter();  // 3

// Multiple independent closures
var counter1 = makeCounter();
var counter2 = makeCounter();
print counter1();  // 1
print counter1();  // 2
print counter2();  // 1 (independent state)
print counter2();  // 2

// Closure factory with parameters
fun makeAdder(n) {
    fun add(x) {
        return x + n;
    }
    return add;
}

var add5 = makeAdder(5);
var add10 = makeAdder(10);
print add5(3);   // 8
print add10(3);  // 13

// Closures capture variables even after scope exits
fun outer() {
    var message = "captured";
    fun inner() {
        return message;
    }
    return inner;
}

var getMessage = outer();
print getMessage();  // "captured"
```

**Classes with Methods and This**
```lox
// Basic class with methods
class Person {
    getName() {
        return this.name;
    }
    setName(name) {
        this.name = name;
    }
}

var person = Person();
person.setName("Alice");
print person.getName();  // "Alice"

// Method chaining with this
class Counter {
    increment() {
        this.value = this.value + 1;
        return this;  // Return this for chaining
    }
    getValue() {
        return this.value;
    }
}

var counter = Counter();
counter.value = 0;
counter.increment().increment().increment();
print counter.getValue();  // 3

// Methods calling other methods
class Math {
    double(n) {
        return n * 2;
    }
    quadruple(n) {
        return this.double(this.double(n));
    }
}

var math = Math();
print math.quadruple(5);  // 20

// Methods stored in variables (bound methods)
class Greeter {
    greet() {
        return "Hello, " + this.name;
    }
}

var greeter = Greeter();
greeter.name = "World";
var greetFn = greeter.greet;  // Method is bound to instance
print greetFn();  // "Hello, World"

// Methods in closures
class Counter {
    getClosure() {
        fun closure() {
            return this.value;
        }
        return closure;
    }
}

var counter = Counter();
counter.value = 42;
var closure = counter.getClosure();
print closure();  // 42 (this is correctly captured)
```

**Classes with Initializers (init method)**
```lox
// Basic initializer with parameters
class Point {
    init(x, y) {
        this.x = x;
        this.y = y;
    }
}

var point = Point(10, 20);
print point.x;  // 10
print point.y;  // 20

// Initializer with no parameters
class Counter {
    init() {
        this.value = 0;
    }
}

var counter = Counter();
print counter.value;  // 0

// Initializer calling methods
class Person {
    init(name) {
        this.name = name;
        this.greet();  // Can call methods from init
    }
    greet() {
        print "Hello, " + this.name;
    }
}

var person = Person("Alice");  // "Hello, Alice"

// Initializer with computed values
class Rectangle {
    init(width, height) {
        this.width = width;
        this.height = height;
        this.area = width * height;  // Compute derived value
    }
}

var rect = Rectangle(5, 3);
print rect.area;  // 15

// Multiple instances with different init arguments
class Point {
    init(x, y) {
        this.x = x;
        this.y = y;
    }
}

var p1 = Point(1, 2);
var p2 = Point(10, 20);
print p1.x;  // 1
print p2.x;  // 10

// Initializer can be called explicitly as a method
class Foo {
    init(value) {
        this.value = value;
        return;  // Explicit return (still returns instance)
    }
}

var foo = Foo(42);
var result = foo.init(100);  // Can call init() explicitly
print result.value;  // 100
print foo.value;     // 100
```

**Classes with Getter Functions**
```lox
// Basic getter function
class Foo {
    value {
        return 42;
    }
}
var foo = Foo();
print foo.value;  // 42 (getter executes immediately)

// Getter using this
class Person {
    name {
        return this._name;
    }
}
var person = Person();
person._name = "Alice";
print person.name;  // "Alice"

// Getter with computed value
class Rectangle {
    init(w, h) {
        this.width = w;
        this.height = h;
    }
    area {
        return this.width * this.height;
    }
}
var rect = Rectangle(5, 3);
print rect.area;  // 15

// Field shadows getter
class Foo {
    value {
        return "getter";
    }
}
var foo = Foo();
foo.value = "field";
print foo.value;  // "field" (field takes precedence)

// Getter calling other methods
class Calculator {
    result {
        return this.compute();
    }
    compute() {
        return 10 + 20;
    }
}
var calc = Calculator();
print calc.result;  // 30

// Getter vs regular method
class Foo {
    getter {
        return "getter result";
    }
    method() {
        return "method result";
    }
}
var foo = Foo();
print foo.getter;      // "getter result" (executes immediately)
var method = foo.method;
print method();        // "method result" (returns function)
```
