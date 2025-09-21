package main 

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)


type GLox struct {
	hadError bool
	hadRuntimeError bool 
	interpreter *Interpreter 
}

func main() {
	
	lox := GLox {}
	lox.interpreter = NewInterpreter(&lox)
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		lox.runFile(os.Args[1])
	} else {
		lox.runPrompt()
	}
}

func (l *GLox) runFile(file string) {
	if data, err := os.ReadFile(file); err != nil {
		log.Fatal(err)
	} else {
		l.run(string(data))
	}

	if l.hadError {
		os.Exit(65)
	}
	if l.hadRuntimeError {
		os.Exit(70);
	}
}

func (l *GLox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
        fmt.Print("> ")
        if scanner.Scan() {
            line := scanner.Text()
			l.run(line)
        } else {
            break
        }
    }

}

func (l *GLox) run(source string) {
	scanner := NewScanner(l, source)
	tokens := scanner.scanTokens()
	parser := NewParser(l, tokens)
	expression, _ := parser.parse()

	if l.hadError {
		return 
	}

	l.interpreter.interpret(expression)
	//fmt.Printf("%s\n", (&astPrinter{}).print(expression))
}

func (l *GLox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *GLox) parseError(token Token, message string) {
	if (token.token_type == EOF) {
		l.report(token.line, " at end", message)
	} else {
		l.report(token.line, fmt.Sprintf(" at %s ", token.lexeme), message)
	}
}

func (l *GLox) runtimeError(err error) {
	runtime_err, _ := err.(RuntimeError)
	fmt.Fprintf(os.Stderr,"[line %d] %s\n", runtime_err.token.line, runtime_err.Error())
	l.hadRuntimeError = true 
}

func (l *GLox) report(line int, where string, message string) {
	fmt.Fprint(os.Stderr, "[line " + strconv.Itoa(line) + "] Error" + where + ": " + message + "\n")
	l.hadError = true 
}