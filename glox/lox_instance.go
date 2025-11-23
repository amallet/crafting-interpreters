package main

import (
	"fmt"
)

type LoxInstance struct {
	klass *LoxClass 
	fields map[string]any
}

func NewLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{
		klass: klass,
		fields: make(map[string]any),
	}
}

func (li *LoxInstance) get(token Token) (any, error) {
	if value, ok := li.fields[token.lexeme]; !ok {
		return nil, RuntimeError{token, fmt.Sprintf("Undefined property name %s", token.lexeme)}
	} else {
		return value, nil 
	}
}

func (li *LoxInstance) set(token Token, value any) {
	li.fields[token.lexeme] = value 
}

func (li *LoxInstance) String() string {
	return fmt.Sprintf("Instance of class %s", li.klass.name)
}