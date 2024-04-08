package lexer

import (
	"strconv"
	"unicode"
)

type ErrorHandler interface {
	HandleError(line int, message string)
	Report(line int, where string, message string)
}

type Scanner struct {
	// stores the original source code
	source string
	tokens []Token

	start   int
	current int
	line    int

	reservedWords map[string]TokenType

	errorHandler ErrorHandler
}

func NewScanner(source string) *Scanner {
	reservedWords := map[string]TokenType{
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

	return &Scanner{
		source:        source,
		tokens:        make([]Token, 0),
		start:         0,
		current:       0,
		line:          1,
		reservedWords: reservedWords,
		// TODO: not sure what the best way to handle errors so hook into Lox for now
		errorHandler: NewLox(),
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		// we are at the beginning of next lexeme
		s.start = s.current
		s.scanToken()
	}

	// add an EOF marker
	s.tokens = append(s.tokens, *NewToken(EOF, "", nil, s.line))

	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current == len(s.source)
}

// in each turn of the loop, we scan a single token
func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(LEFT_PAREN)
	case ")":
		s.addToken(RIGHT_PAREN)
	case "{":
		s.addToken(LEFT_BRACE)
	case "}":
		s.addToken(RIGHT_BRACE)
	case ",":
		s.addToken(COMMA)
	case ".":
		s.addToken(DOT)
	case "-":
		s.addToken(MINUS)
	case "+":
		s.addToken(PLUS)
	case ":":
		s.addToken(SEMICOLON)
	case "*":
		s.addToken(STAR)
	case "!":
		matches := s.match("=")
		if matches {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case "=":
		matches := s.match("=")
		if matches {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case "<":
		matches := s.match("=")
		if matches {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case ">":
		matches := s.match("=")
		if matches {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case "/":
		if matches := s.match("/"); matches {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case " ":
	case "\r":
	case "\t":
	case "\n":
		s.line += 1
	case "\"":
		s.eatString()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.errorHandler.HandleError(s.line, "unexpected character")
		}
	}
}

// peek just looks at current token, doesn't consume
func (s *Scanner) peek() string {
	if s.isAtEnd() {
		return "\\0"
	} else {
		return string(s.source[s.current])
	}
}

// peekNext looks at next token, doesnt consume
func (s *Scanner) peekNext() string {
	if s.current+1 >= len(s.source) {
		return "\\0"
	}
	return string(s.source[s.current+1])
}

// advance consumes the current character
func (s *Scanner) advance() string {
	nextChar := s.source[s.current]
	s.current += 1
	return string(nextChar)
}

func (s *Scanner) addToken(tt TokenType) {
	s.addTokenWithLiteral(tt, nil)
}

func (s *Scanner) addTokenWithLiteral(tt TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(tt, text, literal, s.line))
}

func (s *Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}

	if string(s.source[s.current]) != expected {
		return false
	}

	s.current += 1
	return true
}

func (s *Scanner) eatString() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errorHandler.HandleError(s.line, "unterminated string")
		return
	}

	// the closing quote
	s.advance()

	// trim the quotes
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func isDigit(c string) bool {
	if len(c) != 1 {
		return false
	}

	runeValue := []rune(c)[0]

	return unicode.IsDigit(runeValue)
}

func isAlpha(c string) bool {
	if len(c) != 1 {
		return false
	}

	runeValue := []rune(c)[0]
	return unicode.IsLetter(runeValue)
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// look for fractional digits
	if s.peek() == "." && isDigit(s.peekNext()) {
		// consume the `.`
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	if float, err := strconv.ParseFloat(s.source[s.start:s.current], 32); err != nil {
		s.errorHandler.HandleError(s.line, "there was an error parsing the number")
	} else {
		s.addTokenWithLiteral(NUMBER, float)
	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	ident := s.source[s.start:s.current]
	if tt, ok := s.reservedWords[ident]; ok {
		s.addToken(tt)
	} else {
		s.addTokenWithLiteral(IDENTIFIER, ident)
	}
}
