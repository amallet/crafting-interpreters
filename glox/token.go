package main

import "fmt"

type Token struct {
	token_type TokenType
	lexeme     string // string representation 
	literal    any // actual value, for numbers and strings
	line       int // line of code where token was found
}

func (t Token) String() string {
	f := fmt.Appendf([]byte{},"%s %s %v",t.token_type.String(), t.lexeme, t.literal)
	return string(f)

}