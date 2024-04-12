package parser

import (
	"testing"

	"github.com/brandonshearin/go-lox/lexer"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter(t *testing.T) {

	source := "var a = 1; print a;"
	tokens := lexer.NewScanner(source).ScanTokens()
	stmts := NewParser(tokens).Parse()

	interpreter := NewInterpreter()

	// stmts := []Stmt{
	// 	&VariableDeclarationStmt{
	// 		Name: *lexer.NewToken(lexer.IDENTIFIER, "hello", "hello", 1),
	// 		Initializer: &LiteralExpr{
	// 			Value: 1,
	// 		},
	// 	},
	// }

	err := interpreter.Interpret(stmts)

	assert.Nil(t, err)
}
