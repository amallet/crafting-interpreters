package main 

import (
	"strconv"
	"unicode"
)

var reservedKeyWordMap = map[string]TokenType {
    "and":    AND,
    "class":  CLASS,
    "else":   ELSE,
    "false":  FALSE,
    "for":    FOR,
    "fun":    FUN,
    "if":     IF,
    "nil":    NIL,
    "or":     OR,
    "print":  PRINT,
    "return": RETURN,
    "super":  SUPER,
    "this":   THIS,
    "true":   TRUE,
    "var":    VAR,
    "while":  WHILE,
}

type Scanner struct {
	lox GLox
	source string
	source_runes []rune
	tokens []Token
	start int 
	current int 
	line int
}

func NewScanner(lox GLox, source string) Scanner {
	s := Scanner {
		lox: lox,
		source: source,
		source_runes: []rune(source),
		line: 1,
	}
	return s 
}

func (s *Scanner) scanTokens() []Token {
	for (!s.isAtEnd()) {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{ EOF, "", nil, s.line })
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch (c) {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// Just scan through comment to the end, then discard it
			for (s.peek() != '\n' && !s.isAtEnd()) {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
		
	case ' ':
		// NOP
	case '\r':
		// NOP
	case '\t':
		// NOP
	case '\n':
		s.line += 1 

	case '"':
		s.scanString()
	
	default:
		if s.isDigit(c) {
			s.scanNumber()
		} else if s.isAlpha(c) {
			s.scanIdentifier()
		} else {
			s.lox.error(s.line, "Unexpected character")
		}
	}
}

func (s *Scanner) scanIdentifier() {
	// Scan to the end of the identifier 
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start : s.current]
	
	// Identifier can either be a reserved keyword or a variable name
	if tokenType, exists := reservedKeyWordMap[text]; exists {
		s.addToken(tokenType)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *Scanner) scanNumber() {
	// Consume as many digits as possible
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Handle decimal point - must be followed by at least one number
	if (s.peek() == '.' && s.isDigit(s.peekNext())) {
		s.advance()

		// Consume digits after decimal point
		for (s.isDigit(s.peek())) {
			s.advance()
		}
	}

	// Store all numbers as float64
	value, _ := strconv.ParseFloat(string(s.source_runes[s.start : s.current]), 64)
	s.addLiteralToken(NUMBER, value)
}

func (s *Scanner) scanString() {
	// Are at start of string, look for the end of the string, terminated by
	// another double-quote
	for (s.peek() != '"' && !s.isAtEnd()) {
		if (s.peek() == '\n') { // multi-line strings are ok
			s.line++
		}
		s.advance()
	}

	if (s.isAtEnd()) {
		s.lox.error(s.line, "Unterminated string")
		return 
	}

	s.advance() // found terminating double quote, consume it

	// Extract text between the starting and ending double quotes
	value := string(s.source_runes[s.start + 1 : s.current - 1])

	s.addLiteralToken(STRING, value)
}

func (s *Scanner) match(expected rune) bool {
	if (s.isAtEnd()) {
		return false 
	}

	if (s.source_runes[s.current] != expected) {
		return false 
	}

	s.current++
	return true 
}

func (s *Scanner) peek() rune {
	if (s.isAtEnd()) {
		return '\x00'
	}
	return s.source_runes[s.current]
}

func (s *Scanner) peekNext() rune {
	if (s.current + 1 >= len(s.source_runes)) {
		return '\x00'
	}
	return s.source_runes[s.current + 1]
}

func (s *Scanner) isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (s *Scanner) isAlphaNumeric(r rune) bool {
	return unicode.IsDigit(r) || s.isAlpha(r)
}

func (s *Scanner) isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source_runes)
}

func (s *Scanner) advance() rune {
	r := s.source_runes[s.current]
	s.current++
	return r 
}

func (s *Scanner) addToken(tokenType TokenType) {
	text := string(s.source_runes[s.start : s.current])
	s.tokens = append(s.tokens, Token{ tokenType, text, nil, s.line })
}

func (s *Scanner) addLiteralToken(tokenType TokenType, literal any) {
	text := string(s.source_runes[s.start : s.current])
	s.tokens = append(s.tokens, Token{ tokenType, text, literal, s.line })
}