package main

import (
	"time"
)

// Implementation of built-in functions provided by the interpreter 

type clockFn struct {}

func (c clockFn) arity() int {
	return 0
}

func (c clockFn) call(i *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixMilli()), nil 
}