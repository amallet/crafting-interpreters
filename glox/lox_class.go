package main

type LoxClass struct {
	name    string
	superclass *LoxClass 
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, superclass: superclass, methods: methods}
}

// call() is invoked on a LoxClass to construct a new instance of the class
func (lc *LoxClass) call(i *Interpreter, arguments []any) (any, error) {

	// construct the instance
	instance := NewLoxInstance(i, lc)

	// If class has an init() function, call it to do any instance initialization needed
	if initializer := lc.findMethod("init"); initializer != nil {
		if _, err := initializer.bindThis(instance).call(i, arguments); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (lc *LoxClass) arity() int {
	// Arity is determined by number of parameters the init function takes,
	// if there is an init function
	if initFn := lc.findMethod("init"); initFn != nil {
		return initFn.arity()
	}

	return 0
}

func (lc *LoxClass) findMethod(methodName string) *LoxFunction {
	if method, ok := lc.methods[methodName]; ok {
		return method
	}

	// Couldn't find method on class itself, look in superclass 
	if lc.superclass != nil {
		return lc.superclass.findMethod(methodName)
	}

	return nil
}

func (lc *LoxClass) String() string {
	return lc.name
}
