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
	assert.Equal(LEFT_PAREN, tokens[0].TokenType)
	assert.Equal(RIGHT_PAREN, tokens[1].TokenType)
	assert.Equal(EOF, tokens[2].TokenType)
}

func TestLexemeLength2(t *testing.T) {
	assert := assert.New(t)

	s := NewScanner("!=<=>")
	tokens := s.ScanTokens()

	assert.Len(tokens, 4, "tokens should be length 2")
	assert.Equal(BANG_EQUAL, tokens[0].TokenType)
	assert.Equal(LESS_EQUAL, tokens[1].TokenType)
	assert.Equal(GREATER, tokens[2].TokenType)
	assert.Equal(EOF, tokens[3].TokenType)
}

type ScannerTestCase struct {
	ID             int
	Source         string
	Lines          int
	NumTokens      int
	TokenTypes     []TokenType
	IsNegativeCase bool
	ErrorMsgs      []string
}

var testCases = []ScannerTestCase{
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
	{
		ID:         6,
		Source:     "for while",
		Lines:      1,
		NumTokens:  3,
		TokenTypes: []TokenType{FOR, WHILE, EOF},
	},
	// negative cases
	{
		ID:             7,
		Source:         "@", // illegal token
		IsNegativeCase: true,
		ErrorMsgs:      []string{"unexpected character @"},
	},
	{
		ID:             8,
		Source:         "@$", // multiple illegal tokens --> multiple error msgs
		IsNegativeCase: true,
		ErrorMsgs:      []string{"unexpected character @", "unexpected character $"},
	},
	{
		ID:             9,
		Source:         "for @ while if", // illegal token, but lexer can still produce tokens that are valid
		Lines:          1,
		NumTokens:      2,
		IsNegativeCase: true,
		ErrorMsgs:      []string{"unexpected character @"},
		TokenTypes:     []TokenType{FOR, WHILE, IF, EOF},
	},
	{
		ID:             10,
		Source:         "\"this is an unterminated string", // unterminated strings
		Lines:          1,
		NumTokens:      1,
		IsNegativeCase: true,
		ErrorMsgs:      []string{"unterminated string"},
		TokenTypes:     []TokenType{EOF},
	},
}

func TestLongerLexemes(t *testing.T) {
	assert := assert.New(t)

	for _, testCase := range testCases {
		s := NewScanner(testCase.Source)
		tokens := s.ScanTokens()

		if testCase.IsNegativeCase {
			// iterate through collected errors and do a substring match
			for idx, lexerError := range s.Errors {
				assert.Contains(lexerError, testCase.ErrorMsgs[idx], "test case ID: %d", testCase.ID)
			}
			// we can still assert that the lexer emitted tokens if present on the test case
			for idx, expectedType := range testCase.TokenTypes {
				assert.Equal(expectedType, tokens[idx].TokenType, "test case %d failed. wanted %s token type, got %s token type", testCase.ID, expectedType, tokens[idx].TokenType)
			}
		} else {
			// positive assertions
			assert.Len(tokens, testCase.NumTokens, fmt.Sprintf("test case %d failed", testCase.ID))
			assert.Equal(s.line, testCase.Lines, fmt.Sprintf("test case %d failed", testCase.ID))

			for idx, expectedType := range testCase.TokenTypes {
				assert.Equal(expectedType, tokens[idx].TokenType, "test case %d failed. wanted %s token type, got %s token type", testCase.ID, expectedType, tokens[idx].TokenType)
			}
		}

	}
}

var stringCases = []ScannerTestCase{
	{
		ID:        1,
		Source:    `"im an unterminated string`,
		ErrorMsgs: []string{"[line 1] Error unterminated string"},
	},
	{
		ID:        2,
		Source:    `"a'`, // mismatched quotes
		ErrorMsgs: []string{"[line 1] Error unterminated string"},
	},
	{
		ID:        3,
		Source:    "\n \"hello", // line 2
		ErrorMsgs: []string{"[line 2] Error unterminated string"},
	},
	{
		ID:        4,
		Source:    "@ \n \"hello", // errors on multiple lines
		ErrorMsgs: []string{"[line 1] Error unexpected character @", "[line 2] Error unterminated string"},
	},
}

func TestStringHandling(t *testing.T) {
	assert := assert.New(t)

	for _, testCase := range stringCases {
		s := NewScanner(testCase.Source)
		_ = s.ScanTokens()

		for idx, lexerError := range s.Errors {
			assert.Contains(lexerError, testCase.ErrorMsgs[idx], "test case ID: %d", testCase.ID)
		}

	}

}
