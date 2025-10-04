package main

import (
	"fmt"
)

// The Environment type maps variables to their value within a given
// scope (ie lexical block). An environment may have a pointer to a parent environment, 
// represents the enclosing scope
type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing, make(map[string]any)}
}

func (e *Environment) defineVarValue(name string, value any) {
	e.values[name] = value
}

func (e *Environment) getVarValue(name Token) (any, error) {
	// Try to retrieve value in current environment, if it exists; if not, fall back to
	// enclosing environment, if there is one
	if value, ok := e.values[name.lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.getVarValue(name)
	}

	return nil, RuntimeError{name,
		fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}

func (e *Environment) assignVarValue(name Token, value any) error {
	// Try to assign value in current environment, if it exists; if not, fall back to
	// enclosing environment, if there is one
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assignVarValue(name, value)
	}

	return RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}
