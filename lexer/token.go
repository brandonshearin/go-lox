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
	return fmt.Sprintf("%d %s %s", t.tokenType, t.lexeme, t.literal)
}
