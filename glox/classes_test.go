package main

import (
	"testing"
)

// ============================================================================
// CLASS TESTS
// ============================================================================
// Tests for class definition, instantiation, and property access

func TestClassDefinition(t *testing.T) {
	t.Run("Define simple class", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
instance.test = true;
print instance.test;
`
		expected := []string{"true"}
		runProgramAndCheckOutput(t, program, expected, "Define simple class")
	})

	t.Run("Define class with methods", func(t *testing.T) {
		program := `
class Foo {
    method1() {}
    method2() {}
}
var instance = Foo();
instance.test = true;
print instance.test;
`
		expected := []string{"true"}
		runProgramAndCheckOutput(t, program, expected, "Define class with methods")
	})

	t.Run("Class can be stored in variable", func(t *testing.T) {
		program := `
class Foo {}
var Bar = Foo;
var instance = Bar();
instance.test = true;
print instance.test;
`
		expected := []string{"true"}
		runProgramAndCheckOutput(t, program, expected, "Class can be stored in variable")
	})

	t.Run("Multiple class definitions", func(t *testing.T) {
		program := `
class Foo {}
class Bar {}
var foo = Foo();
var bar = Bar();
foo.test = true;
bar.test = true;
print foo.test;
print bar.test;
`
		expected := []string{"true", "true"}
		runProgramAndCheckOutput(t, program, expected, "Multiple class definitions")
	})
}

func TestClassInstantiation(t *testing.T) {
	t.Run("Instance stored in variable", func(t *testing.T) {
		program := `
class Foo {}
var obj = Foo();
obj.test = true;
print obj.test;
`
		expected := []string{"true"}
		runProgramAndCheckOutput(t, program, expected, "Instance stored in variable")
	})

	t.Run("Instance can be reassigned", func(t *testing.T) {
		program := `
class Foo {}
var obj = Foo();
obj.x = 1;
obj = Foo();
obj.x = 2;
print obj.x;
`
		expected := []string{"2"}
		runProgramAndCheckOutput(t, program, expected, "Instance can be reassigned")
	})
}

func TestPropertyAccess(t *testing.T) {

	t.Run("Get undefined property", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
print instance.undefined;
`
		runProgramAndExpectError(t, program, "undefined property name undefined", "Get undefined property")
	})

	t.Run("Get property with different types", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
instance.num = 42;
instance.str = "hello";
instance.bool = true;
instance.none = nil;
print instance.num;
print instance.str;
print instance.bool;
print instance.none;
`
		expected := []string{"42", "hello", "true", "<nil>"}
		runProgramAndCheckOutput(t, program, expected, "Get property with different types")
	})

	t.Run("Chained property access", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
instance.nested = Foo();
instance.nested.field = "nested value";
print instance.nested.field;
`
		expected := []string{"nested value"}
		runProgramAndCheckOutput(t, program, expected, "Chained property access")
	})
}

func TestPropertyAssignment(t *testing.T) {
	t.Run("Update existing property", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
instance.x = 10;
instance.x = 20;
print instance.x;
`
		expected := []string{"20"}
		runProgramAndCheckOutput(t, program, expected, "Update existing property")
	})

	t.Run("Set property with expression", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
instance.x = 5 + 3;
print instance.x;
`
		expected := []string{"8"}
		runProgramAndCheckOutput(t, program, expected, "Set property with expression")
	})

	t.Run("Set property returns value", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
var result = instance.x = 42;
print result;
print instance.x;
`
		expected := []string{"42", "42"}
		runProgramAndCheckOutput(t, program, expected, "Set property returns value")
	})
}

func TestInstanceIndependence(t *testing.T) {
	t.Run("Instances from different classes are independent", func(t *testing.T) {
		program := `
class Foo {}
class Bar {}
var foo = Foo();
var bar = Bar();
foo.x = "foo value";
bar.x = "bar value";
print foo.x;
print bar.x;
`
		expected := []string{"foo value", "bar value"}
		runProgramAndCheckOutput(t, program, expected, "Instances from different classes are independent")
	})

	t.Run("Instance properties persist after assignment", func(t *testing.T) {
		program := `
class Foo {}
var instance = Foo();
instance.x = 10;
var copy = instance;
copy.x = 20;
print instance.x;
print copy.x;
`
		expected := []string{"20", "20"}
		runProgramAndCheckOutput(t, program, expected, "Instance properties persist after assignment")
	})
}

func TestClassErrors(t *testing.T) {
	t.Run("Get property on non-instance", func(t *testing.T) {
		program := `
class Foo {}
var notInstance = Foo;
print notInstance.field;
`
		runProgramAndExpectError(t, program, "Only instances have properties", "Get property on non-instance")
	})

	t.Run("Set property on non-instance", func(t *testing.T) {
		program := `
class Foo {}
var notInstance = Foo;
notInstance.field = "value";
`
		runProgramAndExpectError(t, program, "Only instances have fields", "Set property on non-instance")
	})

	t.Run("Get property on nil", func(t *testing.T) {
		program := `
var nilValue = nil;
print nilValue.field;
`
		runProgramAndExpectError(t, program, "Only instances have properties", "Get property on nil")
	})

	t.Run("Set property on nil", func(t *testing.T) {
		program := `
var nilValue = nil;
nilValue.field = "value";
`
		runProgramAndExpectError(t, program, "Only instances have fields", "Set property on nil")
	})

	t.Run("Get property on number", func(t *testing.T) {
		program := `
var num = 42;
print num.field;
`
		runProgramAndExpectError(t, program, "Only instances have properties", "Get property on number")
	})

	t.Run("Set property on string", func(t *testing.T) {
		program := `
var str = "hello";
str.field = "value";
`
		runProgramAndExpectError(t, program, "Only instances have fields", "Set property on string")
	})
}

func TestClassInExpressions(t *testing.T) {
	t.Run("Class in loop", func(t *testing.T) {
		program := `
class Foo {}
var i = 0;
while (i < 3) {
    var instance = Foo();
    instance.count = i;
    print instance.count;
    i = i + 1;
}
`
		expected := []string{"0", "1", "2"}
		runProgramAndCheckOutput(t, program, expected, "Class in loop")
	})

	t.Run("Class as function argument", func(t *testing.T) {
		program := `
class Foo {}
fun setProperty(instance) {
    instance.x = 100;
}
var instance = Foo();
setProperty(instance);
print instance.x;
`
		expected := []string{"100"}
		runProgramAndCheckOutput(t, program, expected, "Class as function argument")
	})
}

func TestComplexClassScenarios(t *testing.T) {
	t.Run("Nested instances", func(t *testing.T) {
		program := `
class Container {}
class Item {}
var container = Container();
container.item = Item();
container.item.name = "nested";
print container.item.name;
`
		expected := []string{"nested"}
		runProgramAndCheckOutput(t, program, expected, "Nested instances")
	})

	t.Run("Instance in closure", func(t *testing.T) {
		program := `
class Counter {}
fun makeCounter() {
    var counter = Counter();
    counter.value = 0;
    fun increment() {
        counter.value = counter.value + 1;
        return counter.value;
    }
    return increment;
}
var inc = makeCounter();
print inc();
print inc();
`
		expected := []string{"1", "2"}
		runProgramAndCheckOutput(t, program, expected, "Instance in closure")
	})

	t.Run("Class name shadowing", func(t *testing.T) {
		program := `
class Foo {}
{
    class Foo {}
    var inner = Foo();
    inner.x = 1;
    print inner.x;
}
var outer = Foo();
outer.x = 2;
print outer.x;
`
		expected := []string{"1", "2"}
		runProgramAndCheckOutput(t, program, expected, "Class name shadowing")
	})

	t.Run("Property access in expression", func(t *testing.T) {
		program := `
class Point {}
var point = Point();
point.x = 10;
point.y = 20;
var sum = point.x + point.y;
print sum;
`
		expected := []string{"30"}
		runProgramAndCheckOutput(t, program, expected, "Property access in expression")
	})

	t.Run("Property assignment in expression", func(t *testing.T) {
		program := `
class Counter {}
var counter = Counter();
counter.value = 0;
counter.value = counter.value + 1;
print counter.value;
`
		expected := []string{"1"}
		runProgramAndCheckOutput(t, program, expected, "Property assignment in expression")
	})
}

func TestMethodDefinition(t *testing.T) {
	t.Run("Simple method definition and call", func(t *testing.T) {
		program := `
class Foo {
    greet() {
        print "Hello";
    }
}
var foo = Foo();
foo.greet();
`
		expected := []string{"Hello"}
		runProgramAndCheckOutput(t, program, expected, "Simple method definition and call")
	})

	t.Run("Method with empty body returns nil", func(t *testing.T) {
		program := `
class Foo {
    bar() {}
}
var foo = Foo();
print foo.bar();
`
		expected := []string{"<nil>"}
		runProgramAndCheckOutput(t, program, expected, "Method with empty body returns nil")
	})

	t.Run("Method with return value", func(t *testing.T) {
		program := `
class Foo {
    getValue() {
        return 42;
    }
}
var foo = Foo();
print foo.getValue();
`
		expected := []string{"42"}
		runProgramAndCheckOutput(t, program, expected, "Method with return value")
	})

	t.Run("Method with multiple parameters", func(t *testing.T) {
		program := `
class Math {
    sum(a, b, c) {
        return a + b + c;
    }
}
var math = Math();
print math.sum(1, 2, 3);
`
		expected := []string{"6"}
		runProgramAndCheckOutput(t, program, expected, "Method with multiple parameters")
	})

	t.Run("Multiple methods on same class", func(t *testing.T) {
		program := `
class Foo {
    method1() {
        return "first";
    }
    method2() {
        return "second";
    }
    method3() {
        return "third";
    }
}
var foo = Foo();
print foo.method1();
print foo.method2();
print foo.method3();
`
		expected := []string{"first", "second", "third"}
		runProgramAndCheckOutput(t, program, expected, "Multiple methods on same class")
	})
}

func TestMethodCalling(t *testing.T) {
	t.Run("Method call with different argument types", func(t *testing.T) {
		program := `
class Printer {
    printValue(value) {
        print value;
    }
}
var printer = Printer();
printer.printValue(42);
printer.printValue("hello");
printer.printValue(true);
`
		expected := []string{"42", "hello", "true"}
		runProgramAndCheckOutput(t, program, expected, "Method call with different argument types")
	})

	t.Run("Method call in expression", func(t *testing.T) {
		program := `
class Calculator {
    double(n) {
        return n * 2;
    }
}
var calc = Calculator();
var result = calc.double(5) + calc.double(3);
print result;
`
		expected := []string{"16"}
		runProgramAndCheckOutput(t, program, expected, "Method call in expression")
	})

	t.Run("Nested method calls", func(t *testing.T) {
		program := `
class Math {
    add(a, b) {
        return a + b;
    }
    multiply(a, b) {
        return a * b;
    }
}
var math = Math();
var result = math.add(math.multiply(2, 3), math.multiply(4, 5));
print result;
`
		expected := []string{"26"}
		runProgramAndCheckOutput(t, program, expected, "Nested method calls")
	})
}

func TestMethodAsValue(t *testing.T) {
	t.Run("Method stored in variable", func(t *testing.T) {
		program := `
class Foo {
    greet() {
        print "Hello";
    }
}
var foo = Foo();
var method = foo.greet;
method();
`
		expected := []string{"Hello"}
		runProgramAndCheckOutput(t, program, expected, "Method stored in variable")
	})

	t.Run("Method passed as argument", func(t *testing.T) {
		program := `
class Foo {
    say(message) {
        print message;
    }
}
fun callMethod(fn) {
    fn("called");
}
var foo = Foo();
callMethod(foo.say);
`
		expected := []string{"called"}
		runProgramAndCheckOutput(t, program, expected, "Method passed as argument")
	})

	t.Run("Method returned from function", func(t *testing.T) {
		program := `
class Counter {
    count() {
        print "counting";
    }
}
fun getCounter() {
    var counter = Counter();
    return counter.count;
}
var method = getCounter();
method();
`
		expected := []string{"counting"}
		runProgramAndCheckOutput(t, program, expected, "Method returned from function")
	})
}

func TestMethodErrors(t *testing.T) {
	t.Run("Call undefined method", func(t *testing.T) {
		program := `
class Foo {}
var foo = Foo();
foo.unknown();
`
		runProgramAndExpectError(t, program, "undefined property name unknown", "Call undefined method")
	})

	t.Run("Call method with too few arguments", func(t *testing.T) {
		program := `
class Foo {
    method(a, b) {
        return a + b;
    }
}
var foo = Foo();
foo.method(1);
`
		runProgramAndExpectError(t, program, "Expected 2 arguments but got 1", "Call method with too few arguments")
	})

	t.Run("Call method with too many arguments", func(t *testing.T) {
		program := `
class Foo {
    method(a) {
        return a;
    }
}
var foo = Foo();
foo.method(1, 2);
`
		runProgramAndExpectError(t, program, "Expected 1 arguments but got 2", "Call method with too many arguments")
	})

	t.Run("Call method on non-instance", func(t *testing.T) {
		program := `
class Foo {
    method() {}
}
var notInstance = Foo;
notInstance.method();
`
		runProgramAndExpectError(t, program, "Only instances have properties", "Call method on non-instance")
	})
}

func TestMethodVsField(t *testing.T) {
	t.Run("Field shadows method", func(t *testing.T) {
		program := `
class Foo {
    value() {
        return "method";
    }
}
var foo = Foo();
foo.value = "field";
print foo.value;
`
		expected := []string{"field"}
		runProgramAndCheckOutput(t, program, expected, "Field shadows method")
	})

	t.Run("Method accessible when field not set", func(t *testing.T) {
		program := `
class Foo {
    value() {
        return "method";
    }
}
var foo = Foo();
var method = foo.value;
print method();
`
		expected := []string{"method"}
		runProgramAndCheckOutput(t, program, expected, "Method accessible when field not set")
	})
}

func TestMethodComplexScenarios(t *testing.T) {
	t.Run("Method in loop", func(t *testing.T) {
		program := `
class Counter {
    increment() {
        return 1;
    }
}
var counter = Counter();
var sum = 0;
var i = 0;
while (i < 3) {
    sum = sum + counter.increment();
    i = i + 1;
}
print sum;
`
		expected := []string{"3"}
		runProgramAndCheckOutput(t, program, expected, "Method in loop")
	})

	t.Run("Method in closure", func(t *testing.T) {
		program := `
class Counter {
    getValue() {
        return 42;
    }
}
fun makeGetter() {
    var counter = Counter();
    fun getter() {
        return counter.getValue();
    }
    return getter;
}
var getter = makeGetter();
print getter();
`
		expected := []string{"42"}
		runProgramAndCheckOutput(t, program, expected, "Method in closure")
	})

	t.Run("Method on different instances", func(t *testing.T) {
		program := `
class Person {
    getName() {
        return "Person";
    }
}
var person1 = Person();
var person2 = Person();
print person1.getName();
print person2.getName();
`
		expected := []string{"Person", "Person"}
		runProgramAndCheckOutput(t, program, expected, "Method on different instances")
	})

	t.Run("Method with conditional logic", func(t *testing.T) {
		program := `
class Number {
    isPositive(n) {
        if (n > 0) {
            return true;
        }
        return false;
    }
}
var number = Number();
print number.isPositive(5);
print number.isPositive(-3);
`
		expected := []string{"true", "false"}
		runProgramAndCheckOutput(t, program, expected, "Method with conditional logic")
	})

	t.Run("Method with local variables", func(t *testing.T) {
		program := `
class Calculator {
    compute(a, b) {
        var sum = a + b;
        var product = a * b;
        return sum + product;
    }
}
var calc = Calculator();
print calc.compute(2, 3);
`
		expected := []string{"11"}
		runProgramAndCheckOutput(t, program, expected, "Method with local variables")
	})
}

func TestThisKeyword(t *testing.T) {
	t.Run("Access instance field with this", func(t *testing.T) {
		program := `
class Person {
    getName() {
        return this.name;
    }
}
var person = Person();
person.name = "Alice";
print person.getName();
`
		expected := []string{"Alice"}
		runProgramAndCheckOutput(t, program, expected, "Access instance field with this")
	})

	t.Run("Set instance field with this", func(t *testing.T) {
		program := `
class Counter {
    setValue(n) {
        this.value = n;
    }
    getValue() {
        return this.value;
    }
}
var counter = Counter();
counter.setValue(42);
print counter.getValue();
`
		expected := []string{"42"}
		runProgramAndCheckOutput(t, program, expected, "Set instance field with this")
	})

	t.Run("Call method with this", func(t *testing.T) {
		program := `
class Foo {
    bar() {
        return "bar";
    }
    baz() {
        return this.bar();
    }
}
var foo = Foo();
print foo.baz();
`
		expected := []string{"bar"}
		runProgramAndCheckOutput(t, program, expected, "Call method with this")
	})

	t.Run("Return this from method", func(t *testing.T) {
		program := `
class Foo {
    getSelf() {
        return this;
    }
    getName() {
        return "Foo";
    }
}
var foo = Foo();
var self = foo.getSelf();
print self.getName();
`
		expected := []string{"Foo"}
		runProgramAndCheckOutput(t, program, expected, "Return this from method")
	})

	t.Run("Method chaining with this", func(t *testing.T) {
		program := `
class Builder {
    add(n) {
        this.value = this.value + n;
        return this;
    }
    getValue() {
        return this.value;
    }
}
var builder = Builder();
builder.value = 0;
builder.add(1).add(2).add(3);
print builder.getValue();
`
		expected := []string{"6"}
		runProgramAndCheckOutput(t, program, expected, "Method chaining with this")
	})

	t.Run("This in different instances", func(t *testing.T) {
		program := `
class Person {
    getName() {
        return this.name;
    }
}
var person1 = Person();
person1.name = "Alice";
var person2 = Person();
person2.name = "Bob";
print person1.getName();
print person2.getName();
`
		expected := []string{"Alice", "Bob"}
		runProgramAndCheckOutput(t, program, expected, "This in different instances")
	})

	t.Run("This in nested method calls", func(t *testing.T) {
		program := `
class Math {
    double(n) {
        return n * 2;
    }
    quadruple(n) {
        return this.double(this.double(n));
    }
}
var math = Math();
print math.quadruple(5);
`
		expected := []string{"20"}
		runProgramAndCheckOutput(t, program, expected, "This in nested method calls")
	})

	t.Run("This with property access", func(t *testing.T) {
		program := `
class Point {
    getX() {
        return this.x;
    }
    getY() {
        return this.y;
    }
}
var point = Point();
point.x = 10;
point.y = 20;
print point.getX();
print point.getY();
`
		expected := []string{"10", "20"}
		runProgramAndCheckOutput(t, program, expected, "This with property access")
	})

	t.Run("This in conditional logic", func(t *testing.T) {
		program := `
class Number {
    isPositive() {
        if (this.value > 0) {
            return true;
        }
        return false;
    }
}
var num = Number();
num.value = 5;
print num.isPositive();
num.value = -3;
print num.isPositive();
`
		expected := []string{"true", "false"}
		runProgramAndCheckOutput(t, program, expected, "This in conditional logic")
	})

	t.Run("This in loop", func(t *testing.T) {
		program := `
class Counter {
    increment() {
        this.value = this.value + 1;
    }
    getValue() {
        return this.value;
    }
}
var counter = Counter();
counter.value = 0;
var i = 0;
while (i < 3) {
    counter.increment();
    i = i + 1;
}
print counter.getValue();
`
		expected := []string{"3"}
		runProgramAndCheckOutput(t, program, expected, "This in loop")
	})

	t.Run("This in closure", func(t *testing.T) {
		program := `
class Foo {
    getClosure() {
        fun closure() {
            return this.name;
        }
        return closure;
    }
}
var foo = Foo();
foo.name = "Foo";
var closure = foo.getClosure();
print closure();
`
		expected := []string{"Foo"}
		runProgramAndCheckOutput(t, program, expected, "This in closure")
	})

	t.Run("This in nested closure", func(t *testing.T) {
		program := `
class Foo {
    getClosure() {
        fun f() {
            fun g() {
                fun h() {
                    return this.name;
                }
                return h;
            }
            return g;
        }
        return f;
    }
}
var foo = Foo();
foo.name = "Foo";
var closure = foo.getClosure();
print closure()()();
`
		expected := []string{"Foo"}
		runProgramAndCheckOutput(t, program, expected, "This in nested closure")
	})

	t.Run("This with recursive method", func(t *testing.T) {
		program := `
class Math {
    factorial(n) {
        if (n <= 1) {
            return 1;
        }
        return n * this.factorial(n - 1);
    }
}
var math = Math();
print math.factorial(5);
`
		expected := []string{"120"}
		runProgramAndCheckOutput(t, program, expected, "This with recursive method")
	})

	t.Run("This passed as argument", func(t *testing.T) {
		program := `
class Person {
    getName() {
        return this.name;
    }
}
fun printName(person) {
    print person.getName();
}
var person = Person();
person.name = "Alice";
printName(person);
`
		expected := []string{"Alice"}
		runProgramAndCheckOutput(t, program, expected, "This passed as argument")
	})

	t.Run("This stored in variable", func(t *testing.T) {
		program := `
class Foo {
    getSelf() {
        return this;
    }
    getName() {
        return "Foo";
    }
}
var foo = Foo();
var self = foo.getSelf();
print self.getName();
`
		expected := []string{"Foo"}
		runProgramAndCheckOutput(t, program, expected, "This stored in variable")
	})
}

func TestThisKeywordErrors(t *testing.T) {
	t.Run("This at top level", func(t *testing.T) {
		program := `this;`
		// Note: This may be a resolver error ("Can't use 'this' outside of a class")
		// or a runtime error ("Undefined variable 'this'") depending on implementation
		runProgramAndExpectError(t, program, "this", "This at top level")
	})

	t.Run("This in top-level function", func(t *testing.T) {
		program := `
fun foo() {
    this;
}
foo();
`
		runProgramAndExpectError(t, program, "this", "This in top-level function")
	})

	t.Run("This in nested function outside class", func(t *testing.T) {
		program := `
fun outer() {
    fun inner() {
        this;
    }
    inner();
}
outer();
`
		runProgramAndExpectError(t, program, "this", "This in nested function outside class")
	})
}

func TestInitMethod(t *testing.T) {
	t.Run("Init with no parameters", func(t *testing.T) {
		program := `
class Foo {
    init() {
        this.value = 42;
    }
}
var foo = Foo();
print foo.value;
`
		expected := []string{"42"}
		runProgramAndCheckOutput(t, program, expected, "Init with no parameters")
	})

	t.Run("Init with multiple parameters", func(t *testing.T) {
		program := `
class Person {
    init(name, age, city) {
        this.name = name;
        this.age = age;
        this.city = city;
    }
}
var person = Person("Alice", 30, "NYC");
print person.name;
print person.age;
print person.city;
`
		expected := []string{"Alice", "30", "NYC"}
		runProgramAndCheckOutput(t, program, expected, "Init with multiple parameters")
	})

	t.Run("Init returns instance implicitly", func(t *testing.T) {
		program := `
class Foo {
    init() {
        this.value = 42;
    }
}
var foo = Foo();
print foo.value;
`
		expected := []string{"42"}
		runProgramAndCheckOutput(t, program, expected, "Init returns instance implicitly")
	})

	t.Run("Init with explicit return (no value)", func(t *testing.T) {
		program := `
class Foo {
    init() {
        this.value = 42;
        return;
    }
}
var foo = Foo();
print foo.value;
`
		expected := []string{"42"}
		runProgramAndCheckOutput(t, program, expected, "Init with explicit return (no value)")
	})

	t.Run("Init calling methods", func(t *testing.T) {
		program := `
class Counter {
    init(start) {
        this.value = start;
        this.increment();
    }
    increment() {
        this.value = this.value + 1;
    }
}
var counter = Counter(5);
print counter.value;
`
		expected := []string{"6"}
		runProgramAndCheckOutput(t, program, expected, "Init calling methods")
	})

	t.Run("Init with this keyword", func(t *testing.T) {
		program := `
class Person {
    init(name) {
        this.name = name;
        this.greeting = "Hello, " + this.name;
    }
}
var person = Person("Alice");
print person.greeting;
`
		expected := []string{"Hello, Alice"}
		runProgramAndCheckOutput(t, program, expected, "Init with this keyword")
	})

	t.Run("Multiple instances with different init arguments", func(t *testing.T) {
		program := `
class Point {
    init(x, y) {
        this.x = x;
        this.y = y;
    }
}
var p1 = Point(1, 2);
var p2 = Point(10, 20);
print p1.x;
print p1.y;
print p2.x;
print p2.y;
`
		expected := []string{"1", "2", "10", "20"}
		runProgramAndCheckOutput(t, program, expected, "Multiple instances with different init arguments")
	})

	t.Run("Init can be called explicitly as method", func(t *testing.T) {
		program := `
class Foo {
    init(arg) {
        this.field = arg;
        return;
    }
}
var foo = Foo("one");
foo.field = "field";
var foo2 = foo.init("two");
print foo2;
print foo.field;
`
		expected := []string{"Instance of class Foo", "two"}
		runProgramAndCheckOutput(t, program, expected, "Init can be called explicitly as method")
	})

	t.Run("Init accessing global variables", func(t *testing.T) {
		program := `
var defaultName = "Default";
class Person {
    init(name) {
        if (name == nil) {
            this.name = defaultName;
        } else {
            this.name = name;
        }
    }
}
var person1 = Person("Alice");
var person2 = Person(nil);
print person1.name;
print person2.name;
`
		expected := []string{"Alice", "Default"}
		runProgramAndCheckOutput(t, program, expected, "Init accessing global variables")
	})
}

func TestInitMethodErrors(t *testing.T) {
	t.Run("Init returning a value", func(t *testing.T) {
		program := `
class Foo {
    init() {
        return "result";
    }
}
`
		runProgramAndExpectError(t, program, "can't return a value from an initializer", "Init returning a value")
	})

	t.Run("Init with missing arguments", func(t *testing.T) {
		program := `
class Foo {
    init(a, b) {
        this.sum = a + b;
    }
}
var foo = Foo(1);
`
		runProgramAndExpectError(t, program, "Expected 2 arguments but got 1", "Init with missing arguments")
	})

	t.Run("Init with too many arguments", func(t *testing.T) {
		program := `
class Foo {
    init(a) {
        this.value = a;
    }
}
var foo = Foo(1, 2);
`
		runProgramAndExpectError(t, program, "Expected 1 arguments but got 2", "Init with too many arguments")
	})

	t.Run("Init with zero arguments when required", func(t *testing.T) {
		program := `
class Foo {
    init(a) {
        this.value = a;
    }
}
var foo = Foo();
`
		runProgramAndExpectError(t, program, "Expected 1 arguments but got 0", "Init with zero arguments when required")
	})

	t.Run("Init returning in conditional", func(t *testing.T) {
		program := `
class Foo {
    init() {
        if (true) {
            return "value";
        }
    }
}
`
		runProgramAndExpectError(t, program, "can't return a value from an initializer", "Init returning in conditional")
	})

}
