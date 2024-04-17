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

func TestVariableDecl(t *testing.T) {
	// define a number
	source := "var a = 1;"
	tokens := lexer.NewScanner(source).ScanTokens()
	stmts := ast.NewParser(tokens).Parse()

	i := NewInterpreter()

	err := i.execute(stmts[0])
	assert.Nil(t, err)

	value, err := i.Environment.Get(lexer.Token{Lexeme: "a"})
	assert.Nil(t, err)
	assert.IsType(t, float64(0), value)
	assert.Equal(t, float64(1), value)

	// define a string
	source = "var b = \"hello\";"
	tokens = lexer.NewScanner(source).ScanTokens()
	stmts = ast.NewParser(tokens).Parse()

	i = NewInterpreter()

	err = i.execute(stmts[0])
	assert.Nil(t, err)

	value, err = i.Environment.Get(lexer.Token{Lexeme: "b"})
	assert.Nil(t, err)
	assert.IsType(t, string(""), value)
	assert.Equal(t, "hello", value)

	// define a variable from an arithmetic expression
	source = "var c = 1 + 2;"
	tokens = lexer.NewScanner(source).ScanTokens()
	stmts = ast.NewParser(tokens).Parse()

	i = NewInterpreter()

	err = i.execute(stmts[0])
	assert.Nil(t, err)

	value, err = i.Environment.Get(lexer.Token{Lexeme: "c"})
	assert.Nil(t, err)
	assert.IsType(t, float64(0), value)
	assert.Equal(t, float64(3), value)

}

// TODO: parser is spinning forever on the scope test case
func TestAssignmentExpr(t *testing.T) {
	// assign a number to a variable
	source := "var a; a = 1234;"
	tokens := lexer.NewScanner(source).ScanTokens()
	stmts := ast.NewParser(tokens).Parse()

	i := NewInterpreter()

	err := i.Interpret(stmts)
	assert.Nil(t, err)

	value, _ := i.Environment.Get(lexer.Token{Lexeme: "a"})
	assert.IsType(t, float64(0), value)
	assert.Equal(t, float64(1234), value)

	// reassignment
	source = "var a \"before\"; a = \"after\"; print a;"
	tokens = lexer.NewScanner(source).ScanTokens()
	stmts = ast.NewParser(tokens).Parse()

	i = NewInterpreter()

	err = i.Interpret(stmts)
	assert.Nil(t, err)
	output := i.Output.String()
	assert.Equal(t, "after\n", output)

	// scope
	source = "{var a \"first\"; print a;} {var a \"second\"; print a;}"
	tokens = lexer.NewScanner(source).ScanTokens()
	p := ast.NewParser(tokens)
	stmts = p.Parse()

	i = NewInterpreter()

	err = i.Interpret(stmts)
	assert.Nil(t, err)
	output = i.Output.String()
	assert.Equal(t, "first\n second\n", output)

}
