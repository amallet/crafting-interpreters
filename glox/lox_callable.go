package main 

// The LoxCallable interface needs to be implemented by anything that can be 
// called from a Lox program ie functions and (class) methods. 
type LoxCallable interface {
	arity() int 
	call(i *Interpreter, arguments []any) (any, error)
}