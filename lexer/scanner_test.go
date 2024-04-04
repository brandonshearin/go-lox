package lexer

import (
	"fmt"
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

type ScannerTestCase struct {
	ID         int
	Source     string
	Lines      int
	NumTokens  int
	TokenTypes []TokenType
}

func TestLongerLexemes(t *testing.T) {
	assert := assert.New(t)

	testCases := []ScannerTestCase{
		{
			ID:         1,
			Source:     "// this is a comment",
			Lines:      1,
			NumTokens:  1,
			TokenTypes: []TokenType{EOF},
		},
		{
			ID: 2,
			Source: `// this is a comment 
	(( )) {} // grouping some things`,
			Lines:      2,
			NumTokens:  7,
			TokenTypes: []TokenType{LEFT_PAREN, LEFT_PAREN, RIGHT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, EOF},
		},
		{
			ID:         3,
			Source:     `!*+-/=<> <= == // some operators`,
			Lines:      1,
			NumTokens:  11,
			TokenTypes: []TokenType{BANG, STAR, PLUS, MINUS, SLASH, EQUAL, LESS, GREATER, LESS_EQUAL, EQUAL_EQUAL, EOF},
		},
		{
			ID:         4,
			Source:     "\"hello world\" \"whats your name?\"",
			Lines:      1,
			NumTokens:  3,
			TokenTypes: []TokenType{STRING, STRING, EOF},
		},
		{
			ID:         5,
			Source:     "123456 1.43",
			Lines:      1,
			NumTokens:  3,
			TokenTypes: []TokenType{NUMBER, NUMBER, EOF},
		},
	}

	for _, testCase := range testCases {
		s := NewScanner(testCase.Source)
		tokens := s.ScanTokens()

		assert.Len(tokens, testCase.NumTokens, fmt.Sprintf("test case %d failed", testCase.ID))
		assert.Equal(s.line, testCase.Lines, fmt.Sprintf("test case %d failed", testCase.ID))

		for idx, expectedType := range testCase.TokenTypes {
			assert.Equal(expectedType, tokens[idx].tokenType, "test case %d failed. wanted %s token type, got %s token type", testCase.ID, expectedType, tokens[idx].tokenType)
		}
	}
}
