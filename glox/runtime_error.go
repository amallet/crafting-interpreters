package main

type RuntimeError struct {
	token Token
	message string 
}

func (e RuntimeError) Error() string {
	return e.message 
}