package lexer

type ErrorHandler interface {
	handleError(line int, message string)
}

type Scanner struct {
	// stores the original source code
	source string
	tokens []Token

	start   int
	current int
	line    int

	errorHandler ErrorHandler
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
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
	case "/":
		s.addToken(SLASH)
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
	default:
		s.errorHandler.handleError(s.line, "unexpected character")
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
