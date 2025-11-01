package main

type ReturnValue struct {
	value any 
}

func (*ReturnValue) Error() string {
	return ""
}