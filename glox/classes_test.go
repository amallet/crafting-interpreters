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
		runProgramAndExpectError(t, program, "Undefined property name undefined", "Get undefined property")
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
