package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	assert := assert.New(t)

	s := NewScanner("(")
	tokens := s.ScanTokens()

	assert.Len(tokens, 2, "tokens should be length 2")
	assert.Equal(LEFT_PAREN, tokens[0].tokenType)
	assert.Equal(EOF, tokens[1].tokenType)
}
