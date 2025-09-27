package main

import (
	"fmt"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing, make(map[string]any)}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) getValue(name Token) (any, error) {
	// Try to retrieve value in current environment, if it exists; if not, fall back to
	// enclosing environment, if there is one
	if value, ok := e.values[name.lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.getValue(name)
	}

	return nil, RuntimeError{name,
		fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}

func (e *Environment) assign(name Token, value any) error {
	// Try to assign value in current environment, if it exists; if not, fall back to
	// enclosing environment, if there is one
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}
