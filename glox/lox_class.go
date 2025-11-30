package main

type LoxClass struct {
	name string 
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, methods: methods}
}

// call() is invoked on a LoxClass to construct a new instance of the class
func (lc *LoxClass) call(i *Interpreter, arguments []any) (any, error) {
	
	// construct the instance
	instance := NewLoxInstance(lc)

	// If class has an init() function, call it to do any instance initialization needed
	if initializer := lc.findMethod("init"); initializer != nil {
		if _, err := initializer.bind(instance).call(i, arguments); err != nil {
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

	return nil 
}

func (lc *LoxClass) String() string {
	return lc.name 
}