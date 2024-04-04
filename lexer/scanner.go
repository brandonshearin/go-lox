package lexer

import "fmt"

type Scanner struct {
	// stores the original source code
	source string
	tokens []Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
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
		break
	default:
		fmt.Printf("unexpected character: %s \n", c)
	}
}

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
