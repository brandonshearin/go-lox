package lexer

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func NewToken(t TokenType, lexeme string, literal any, line int) *Token {
	return &Token{t, lexeme, literal, line}
}

func (t *Token) String() string {
	return fmt.Sprintf("token type: %s, lexeme: %s", t.tokenType.String(), t.lexeme)
}
