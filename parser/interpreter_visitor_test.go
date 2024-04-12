package parser

import (
	"testing"

	"github.com/brandonshearin/go-lox/lexer"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {

	source := "a = 1;"
	tokens := lexer.NewScanner(source).ScanTokens()
	stmts := NewParser(tokens).Parse()

	interpreter := NewInterpreter()

	err := interpreter.Interpret(stmts)

	assert.Nil(t, err)
}
