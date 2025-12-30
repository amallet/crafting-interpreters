package main

import (
	"fmt"
)

// The Environment type maps variables to their value within a given
// scope (ie lexical block). An environment may have a pointer to a parent environment, 
// which represents the enclosing scope
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

func (e *Environment) getVarValue(varToken Token) (any, error) {
	// Try to retrieve value in current environment, if it exists; if not, fall back to
	// enclosing environment, if there is one
	if value, ok := e.values[varToken.lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.getVarValue(varToken)
	}

	return nil, RuntimeError{varToken,
		fmt.Sprintf("Undefined variable '%s'", varToken.lexeme)}
}

func (e *Environment) assignVarValue(varToken Token, value any) error {
	// Try to assign value in current environment, if it exists; if not, fall back to
	// enclosing environment, if there is one
	if _, ok := e.values[varToken.lexeme]; ok {
		e.values[varToken.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assignVarValue(varToken, value)
	}

	return RuntimeError{varToken, fmt.Sprintf("Undefined variable '%s'", varToken.lexeme)}
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e 
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env 
}

func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) assignAt(distance int, varToken Token, value any) {
	e.ancestor(distance).values[varToken.lexeme] = value 
}
