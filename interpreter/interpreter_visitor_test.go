package interpreter

import (
	"testing"

	"github.com/brandonshearin/go-lox/lexer"
	ast "github.com/brandonshearin/go-lox/parser"
	"github.com/stretchr/testify/assert"
)

// todo: maybe the best way to test the interpreter is to test specific methods, ie evalute() and execute(), or visitIfStmt() etc...

func TestInterpreter(t *testing.T) {

	source := "a = 1;"
	tokens := lexer.NewScanner(source).ScanTokens()
	stmts := ast.NewParser(tokens).Parse()

	interpreter := NewInterpreter()

	err := interpreter.Interpret(stmts)

	assert.Nil(t, err)
}
