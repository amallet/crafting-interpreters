package main

type LoxClass struct {
	name string 
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, methods: methods}
}

func (lc *LoxClass) call(i *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(lc)
	return instance, nil 
}

func (lc *LoxClass) arity() int {
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