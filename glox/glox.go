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
		l.run(string(data), false)
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
			l.run(line, true)
			l.hadError = false; // reset error state
        } else {
            break
        }
    }

}

func (l *GLox) run(source string, in_repl bool) {
	// Tokenize input 
	scanner := NewScanner(l, source)
	tokens := scanner.scanTokens()

	// Parse tokens into valid statements
	parser := NewParser(l, tokens)
	statements, _ := parser.parse()
	if l.hadError { // bail out if parsing failed 
		return 
	}

	// Do some static analysis to resolve variables to the right scopes/closures
	resolver := NewResolver(l, l.interpreter)
	resolver.resolveStmts(statements)
	if l.hadError {
		return 
	}

	// Interpret the parsed statements
	results := l.interpreter.interpret(statements)

	// If in REPL mode, also print the results of any expressions that were 
	// entered 
	if in_repl && len(results) > 0 {
		for _, result := range(results) {
			fmt.Printf("%v\n", result)
		}
	}
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