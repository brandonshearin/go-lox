package parser

import (
	"testing"

	"github.com/brandonshearin/go-lox/lexer"
	"github.com/stretchr/testify/assert"
)

func TestEqualityExpr(t *testing.T) {
	// a
	p := NewParser([]lexer.Token{
		{TokenType: lexer.STRING, Lexeme: "a", Literal: "a", Line: 1},
		{TokenType: lexer.EOF, Line: 1},
	})

	exprAST := p.expression()

	visitor := ASTPrinter{}
	prettyPrinted := visitor.Print(exprAST)

	assert.Equal(t, prettyPrinted, "a")

	// a == b
	p = NewParser([]lexer.Token{
		{TokenType: lexer.STRING, Lexeme: "a", Literal: "a", Line: 1},
		{TokenType: lexer.EQUAL_EQUAL, Lexeme: "==", Literal: "==", Line: 1},
		{TokenType: lexer.STRING, Lexeme: "b", Literal: "b", Line: 1},
		{TokenType: lexer.EOF, Line: 1},
	})

	exprAST = p.expression()

	visitor = ASTPrinter{}
	prettyPrinted = visitor.Print(exprAST)

	assert.Equal(t, prettyPrinted, "(== a b)")

	// a != b
	p = NewParser([]lexer.Token{
		{TokenType: lexer.STRING, Lexeme: "a", Literal: "a", Line: 1},
		{TokenType: lexer.BANG_EQUAL, Lexeme: "!=", Literal: "!=", Line: 1},
		{TokenType: lexer.STRING, Lexeme: "b", Literal: "b", Line: 1},
		{TokenType: lexer.EOF, Line: 1},
	})

	exprAST = p.expression()

	visitor = ASTPrinter{}
	prettyPrinted = visitor.Print(exprAST)

	assert.Equal(t, prettyPrinted, "(!= a b)")

	// a == b == c
	p = NewParser([]lexer.Token{
		{TokenType: lexer.STRING, Lexeme: "a", Literal: "a", Line: 1},
		{TokenType: lexer.EQUAL_EQUAL, Lexeme: "==", Literal: "==", Line: 1},
		{TokenType: lexer.STRING, Lexeme: "b", Literal: "b", Line: 1},
		{TokenType: lexer.EQUAL_EQUAL, Lexeme: "==", Literal: "==", Line: 1},
		{TokenType: lexer.STRING, Lexeme: "c", Literal: "c", Line: 1},
		{TokenType: lexer.EOF, Line: 1},
	})

	exprAST = p.expression()

	visitor = ASTPrinter{}
	prettyPrinted = visitor.Print(exprAST)

	assert.Equal(t, prettyPrinted, "(== (== a b) c)")
}

func TestComparisonExpr(t *testing.T) {
	// a == b > c
	p := NewParser([]lexer.Token{
		{TokenType: lexer.STRING, Lexeme: "a", Literal: "a", Line: 1},
		{TokenType: lexer.EQUAL_EQUAL, Lexeme: "==", Literal: "==", Line: 1},
		{TokenType: lexer.STRING, Lexeme: "b", Literal: "b", Line: 1},
		{TokenType: lexer.GREATER, Lexeme: ">", Literal: ">", Line: 1},
		{TokenType: lexer.STRING, Lexeme: "c", Literal: "c", Line: 1},
		{TokenType: lexer.EOF, Line: 1},
	})

	exprAST := p.expression()

	visitor := ASTPrinter{}
	prettyPrinted := visitor.Print(exprAST)

	assert.Equal(t, "(== a (> b c))", prettyPrinted)
}
