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
}

func main() {
	/*
	lox := GLox {}
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		lox.runFile(os.Args[1])
	} else {
		lox.runPrompt()
	}
		*/
	expression := &Binary { 
		&Unary { Token { MINUS, "-", nil, 1}, 
				&Literal {123} },
		Token{ STAR, "*", nil, 1 },
		&Grouping { &Literal {45.67} } }

	printer := astPrinter{}
	fmt.Printf("%s \n", printer.print(expression))

	expression = &Binary {
		&Grouping {
			&Binary {
				&Literal {1},
				Token {PLUS, "+", nil, 1},
				&Literal {2},
			},
		},
		Token{STAR, "*", nil, 1},
		&Grouping {
			&Binary {
				&Literal {3},
				Token {MINUS, "-", nil, 1},
				&Literal {4},
			},
		},
	}
	
	rpn_printer := rpnPrinter{}
	fmt.Printf("%s \n", rpn_printer.print(expression))
}

func (l GLox) runFile(file string) {
	if data, err := os.ReadFile(file); err != nil {
		log.Fatal(err)
	} else {
		l.run(string(data))
	}

	if l.hadError {
		os.Exit(65)
	}

}

func (l GLox) runPrompt() {
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

func (l GLox) run(source string) {
	scanner := NewScanner(l, source)
	tokens := scanner.scanTokens()
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}

func (l GLox) error(line int, message string) {
	l.report(line, "", message)
}

func (l GLox) report(line int, where string, message string) {
	fmt.Fprint(os.Stderr, "[line " + strconv.Itoa(line) + "] Error" + where + ": " + message + "\n")
	l.hadError = true 
}