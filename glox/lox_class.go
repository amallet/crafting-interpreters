package main

type LoxClass struct {
	name string 
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (lc *LoxClass) call(i *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(lc)
	return instance, nil 
}

func (lc *LoxClass) arity() int {
	return 0 
}

func (lc *LoxClass) String() string {
	return lc.name 
}