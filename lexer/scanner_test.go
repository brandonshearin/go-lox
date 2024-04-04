package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexemeLength1(t *testing.T) {
	assert := assert.New(t)

	s := NewScanner("()")
	tokens := s.ScanTokens()

	assert.Len(tokens, 3, "tokens should be length 3")
	assert.Equal(LEFT_PAREN, tokens[0].tokenType)
	assert.Equal(RIGHT_PAREN, tokens[1].tokenType)
	assert.Equal(EOF, tokens[2].tokenType)
}

func TestLexemeLength2(t *testing.T) {
	assert := assert.New(t)

	s := NewScanner("!=<=>")
	tokens := s.ScanTokens()

	assert.Len(tokens, 4, "tokens should be length 2")
	assert.Equal(BANG_EQUAL, tokens[0].tokenType)
	assert.Equal(LESS_EQUAL, tokens[1].tokenType)
	assert.Equal(GREATER, tokens[2].tokenType)
	assert.Equal(EOF, tokens[3].tokenType)
}
