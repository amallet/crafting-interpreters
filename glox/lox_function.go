package main

// LoxFunction implements the LoxCallable interface, and wraps the code that
// needs to be interpreted to execute a Lox function.
type LoxFunction struct {
	declaration *FunctionStmt
	closure *Environment 
}

// Execute the actual function that's wrapped by the enclosing LoxFunction
func (lf *LoxFunction) call(interpreter *Interpreter, arguments []any) (any, error) {
	
	// Create the scope within which the function will execute, as a child of the closure/scope
	// associated with the function 
	env := NewEnvironment(lf.closure)

	// Bind parameters to their values within the function's scope 
	for i, param := range lf.declaration.params {
		env.defineVarValue(param.lexeme, arguments[i])
	}

	// Execute the function's code
	err := interpreter.executeBlock(lf.declaration.body, env)
	if err != nil {
		// If the error returned is of type ReturnValue, then it's not 
		// really an error, but a  wrapper for the actual return value of the 
		// function code that was just interpreted, so that's what's returned
		if retValue, ok := err.(*ReturnValue); ok {
			return retValue.value, nil 
		} else { 
			return nil, err
		}
	}
	return nil, nil
}

// Implement arity() to comply with LoxCallable interface
func (lf *LoxFunction) arity() int {
	return len(lf.declaration.params)
}

// bind() binds the 'this' variable for the given function instance to the 
// supplied class instance,  
func (lf *LoxFunction) bind(li *LoxInstance) *LoxFunction {
	env := NewEnvironment(lf.closure)
	env.defineVarValue("this", li)
	return &LoxFunction{lf.declaration, env}
}
