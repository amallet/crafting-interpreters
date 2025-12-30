package main

import "testing"

// ============================================================================
// INHERITANCE + SUPER TESTS
// ============================================================================

func TestInheritanceAndSuper(t *testing.T) {
	t.Run("Inherited method lookup", func(t *testing.T) {
		program := `
class A {
    method() {
        return "A";
    }
}

class B < A {}

print B().method();
`
		expected := []string{"A"}
		runProgramAndCheckOutput(t, program, expected, "Inherited method lookup")
	})

	t.Run("Override method and call super", func(t *testing.T) {
		program := `
class Doughnut {
    cook() {
        print "Fry";
    }
}

class BostonCream < Doughnut {
    cook() {
        super.cook();
        print "Pipe";
    }
}

BostonCream().cook();
`
		expected := []string{"Fry", "Pipe"}
		runProgramAndCheckOutput(t, program, expected, "Override method and call super")
	})

	t.Run("Inherited initializer runs when constructing subclass", func(t *testing.T) {
		program := `
class A {
    init() {
        this.x = 1;
    }
}

class B < A {}

var b = B();
print b.x;
`
		expected := []string{"1"}
		runProgramAndCheckOutput(t, program, expected, "Inherited initializer runs when constructing subclass")
	})

	t.Run("Subclass initializer can call super.init()", func(t *testing.T) {
		program := `
class A {
    init() {
        this.x = 1;
    }
}

class B < A {
    init() {
        super.init();
        this.y = 2;
    }
}

var b = B();
print b.x;
print b.y;
`
		expected := []string{"1", "2"}
		runProgramAndCheckOutput(t, program, expected, "Subclass initializer can call super.init()")
	})

	t.Run("Can't inherit from self", func(t *testing.T) {
		program := `
class Foo < Foo {}
`
		runProgramAndExpectError(t, program, "A class can't inherit from itself", "Can't inherit from self")
	})

	t.Run("Superclass must be a class", func(t *testing.T) {
		program := `
var NotAClass = "nope";
class Child < NotAClass {}
`
		runProgramAndExpectError(t, program, "Not a class.", "Superclass must be a class")
	})

	t.Run("Super method must exist on superclass", func(t *testing.T) {
		program := `
class A {}
class B < A {
    test() {
        super.nope();
    }
}

B().test();
`
		runProgramAndExpectError(t, program, "Undefined property nope.", "Super method must exist on superclass")
	})

	t.Run("Deep hierarchy: super resolves to immediate superclass at each level", func(t *testing.T) {
		program := `
class A {
    greet() { return "A"; }
}

class B < A {
    greet() { return "B"; }
    parentGreet() { return super.greet(); } // should be A.greet()
}

class C < B {
    greet() { return "C"; }
    parentGreet() { return super.greet(); } // should be B.greet()
}

print B().parentGreet();
print C().parentGreet();
`
		expected := []string{"A", "B"}
		runProgramAndCheckOutput(t, program, expected, "Deep hierarchy: super resolves to immediate superclass at each level")
	})

	t.Run("Deep hierarchy: nested super calls chain correctly", func(t *testing.T) {
		program := `
class A {
    say() { return "A"; }
}

class B < A {
    say() { return "B(" + super.say() + ")"; }
}

class C < B {
    say() { return "C(" + super.say() + ")"; }
}

print C().say();
`
		expected := []string{"C(B(A))"}
		runProgramAndCheckOutput(t, program, expected, "Deep hierarchy: nested super calls chain correctly")
	})

	t.Run("Deep hierarchy: inherited method keeps its own super binding (lexical)", func(t *testing.T) {
		program := `
class A {
    foo() { return "A.foo"; }
}

class B < A {
    bar() { return "B.bar -> " + super.foo(); } // super should be A
}

class C < B {} // inherits bar() from B

print C().bar();
`
		expected := []string{"B.bar -> A.foo"}
		runProgramAndCheckOutput(t, program, expected, "Deep hierarchy: inherited method keeps its own super binding (lexical)")
	})

	t.Run("Can't use super at top level", func(t *testing.T) {
		program := `
super.foo();
`
		runProgramAndExpectError(t, program, "Can't use 'super' outside a class", "Can't use super at top level")
	})

	t.Run("Can't use super in a class with no superclass", func(t *testing.T) {
		program := `
class A {
    method() {
        super.foo();
    }
}

A().method();
`
		runProgramAndExpectError(t, program, "Can't use 'super' in a class with no superclass", "Can't use super in a class with no superclass")
	})

	t.Run("Can't use super in a top-level function", func(t *testing.T) {
		program := `
fun f() {
    super.foo();
}

f();
`
		runProgramAndExpectError(t, program, "Can't use 'super' outside a class", "Can't use super in a top-level function")
	})
}
