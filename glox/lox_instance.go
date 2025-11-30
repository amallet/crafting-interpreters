package main

import (
	"fmt"
)

type LoxInstance struct {
	class *LoxClass // class that this is an instance of 
	fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class: class,
		fields: make(map[string]any),
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
		boundMethod := method.bind(li)
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