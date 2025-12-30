package main

import (
	"fmt"
)

type LoxInstance struct {
	interpreter *Interpreter
	class       *LoxClass // class that this is an instance of
	fields      map[string]any
}

func NewLoxInstance(interpreter *Interpreter, class *LoxClass) *LoxInstance {
	return &LoxInstance{
		interpreter: interpreter,
		class:       class,
		fields:      make(map[string]any),
	}
}

func (li *LoxInstance) get(token Token) (any, error) {
	// First look for instance fields with matching name
	if value, ok := li.fields[token.lexeme]; ok {
		return value, nil
	}

	// No instance field matches, look for matching method on class, and
	// bind it to this instance
	if method := li.class.findMethod(token.lexeme); method != nil {
		boundMethod := method.bindThis(li)
		// If the method is a getter function, execute it immediately, to generate
		// the return value from the getter
		if method.declaration.isGetter {
			return boundMethod.call(li.interpreter, nil)
		}

		// Otherwise, just return the function itself
		return boundMethod, nil
	}

	return nil, RuntimeError{token, fmt.Sprintf("undefined property name %s", token.lexeme)}
}

func (li *LoxInstance) set(token Token, value any) {
	li.fields[token.lexeme] = value
}

func (li *LoxInstance) String() string {
	return fmt.Sprintf("Instance of class %s", li.class.name)
}
