package parser

import (
	"fmt"
	"testing"

	"github.com/brandonshearin/go-lox/lexer"
	"github.com/stretchr/testify/assert"
)

func TestEqualityExpr(t *testing.T) {
	// "a" EOF
	p := NewParser([]lexer.Token{
		{TokenType: lexer.STRING, Lexeme: "a", Literal: "a", Line: 1},
		{TokenType: lexer.EOF, Line: 1},
	})

	exprAST := p.expression()

	visitor := ASTPrinter{}
	prettyPrinted := visitor.Print(exprAST)

	assert.Equal(t, prettyPrinted, "a")

	// "a" == "b"
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

	// "a" != "b"
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

	// "a" == "b" == "c"
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
	// "a" == "b" > "c"
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

// TODO:
func TestErrorHandling(t *testing.T) {
	p := NewParser([]lexer.Token{
		{TokenType: lexer.GREATER, Lexeme: ">", Literal: ">", Line: 1},
		{TokenType: lexer.EOF, Lexeme: "\\0", Literal: "\\0", Line: 1},
	})

	_ = p.Parse()

	// assert.Len(t, errs, 1)

}

type TestCase struct {
	ID                     int
	Source                 string
	ExpectedRepresentation string
}

var testCases = []TestCase{
	// equality expr, precedence level is equal
	{
		ID:                     1,
		Source:                 "1 == 2 == 3 == 4",
		ExpectedRepresentation: "(== (== (== 1.00 2.00) 3.00) 4.00)",
	},
	// equality expr, precedence level is equal
	{
		ID:                     2,
		Source:                 "1 == 2 != 3 == 4",
		ExpectedRepresentation: "(== (!= (== 1.00 2.00) 3.00) 4.00)",
	},
	//  comparison expr > equality expr
	{
		ID:                     3,
		Source:                 "1 == 2 >= 3 > 4",
		ExpectedRepresentation: "(== 1.00 (> (>= 2.00 3.00) 4.00))",
	},
	//  comparison expr > equality expr
	{
		ID:                     4,
		Source:                 "1 == 2 >= 3 != 4",
		ExpectedRepresentation: "(!= (== 1.00 (>= 2.00 3.00)) 4.00)",
	},
	// comparison expr, precedence level is equal
	{
		ID:                     5,
		Source:                 "1 > 2 >= 3 < 4 <= 5",
		ExpectedRepresentation: "(<= (< (>= (> 1.00 2.00) 3.00) 4.00) 5.00)",
	},
	// term expr, precdence level is equal
	{
		ID:                     6,
		Source:                 "1 - 2 + 3",
		ExpectedRepresentation: "(+ (- 1.00 2.00) 3.00)",
	},
	//  term exprs > comparison exprs > equality exprs
	{
		ID:                     7,
		Source:                 "true == false > 6 - 1",
		ExpectedRepresentation: "(== true (> false (- 6.00 1.00)))",
	},
	//  factor exprs > term exprs > comparison exprs > equality exprs
	{
		ID:                     8,
		Source:                 "true == false > 6 - 2 * 5",
		ExpectedRepresentation: "(== true (> false (- 6.00 (* 2.00 5.00))))",
	},
	//  factor exprs > term exprs > comparison exprs > equality exprs
	{
		ID:                     9,
		Source:                 "5 * 2 - 6 > false != true",
		ExpectedRepresentation: "(!= (> (- (* 5.00 2.00) 6.00) false) true)",
	},
}

func TestPrecedence(t *testing.T) {
	for _, testCase := range testCases {
		source := testCase.Source
		// assuming no lexical errors
		tokens := lexer.NewScanner(source).ScanTokens()

		// parse tokens into expression AST
		p := NewParser(tokens)
		exprAST := p.expression()

		printer := ASTPrinter{}
		stringifiedAST := printer.Print(exprAST)

		assert.Equal(t, testCase.ExpectedRepresentation, stringifiedAST, fmt.Sprintf("test case %d failed", testCase.ID))
	}
}

func TestLiterals(t *testing.T) {
	// variable expressions----------------------------------------------------------------
	source := "a == b"
	tokens := lexer.NewScanner(source).ScanTokens()

	p := NewParser(tokens)
	exprAST := p.expression()

	assert.IsType(t, &VariableExpr{}, exprAST.(*BinaryExpr).LeftExpr)
	assert.IsType(t, &VariableExpr{}, exprAST.(*BinaryExpr).RightExpr)

	// grouping expressions----------------------------------------------------------------
	source = "(1 + 2)"
	tokens = lexer.NewScanner(source).ScanTokens()

	p = NewParser(tokens)
	exprAST = p.expression()

	assert.IsType(t, &GroupingExpr{}, exprAST)
	assert.IsType(t, &BinaryExpr{}, exprAST.(*GroupingExpr).Expr)

	source = "(a)"
	tokens = lexer.NewScanner(source).ScanTokens()

	p = NewParser(tokens)
	exprAST = p.expression()

	assert.IsType(t, &GroupingExpr{}, exprAST)
	assert.IsType(t, &VariableExpr{}, exprAST.(*GroupingExpr).Expr)
	assert.Equal(t, "a", exprAST.(*GroupingExpr).Expr.(*VariableExpr).Name.Lexeme)

	// number literals -------------------------------- --------------------------------
	source = "1"
	tokens = lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p = NewParser(tokens)
	exprAST = p.expression()

	assert.IsType(t, &LiteralExpr{}, exprAST)
	assert.IsType(t, float64(0), exprAST.(*LiteralExpr).Value)
	assert.Equal(t, float64(1), exprAST.(*LiteralExpr).Value)

	// boolean literals -------------------------------- --------------------------------
	source = "true"
	tokens = lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p = NewParser(tokens)
	exprAST = p.expression()

	assert.IsType(t, &LiteralExpr{}, exprAST)
	assert.IsType(t, bool(false), exprAST.(*LiteralExpr).Value)
	assert.Equal(t, true, exprAST.(*LiteralExpr).Value)

	source = "false"
	tokens = lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p = NewParser(tokens)
	exprAST = p.expression()

	assert.IsType(t, &LiteralExpr{}, exprAST)
	assert.IsType(t, bool(false), exprAST.(*LiteralExpr).Value)
	assert.Equal(t, false, exprAST.(*LiteralExpr).Value)

	// string literals -------------------------------- --------------------------------
	source = "\"hello world\""
	tokens = lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p = NewParser(tokens)
	exprAST = p.expression()

	assert.IsType(t, &LiteralExpr{}, exprAST)
	assert.IsType(t, string(""), exprAST.(*LiteralExpr).Value)
	assert.Equal(t, "hello world", exprAST.(*LiteralExpr).Value)
}

func getStmtFromSource(source string) Stmt {
	tokens := lexer.NewScanner(source).ScanTokens()
	p := NewParser(tokens)
	ast := p.declaration()

	return ast
}

func TestVarDecl(t *testing.T) {
	// literal initializer
	ast := getStmtFromSource("var hello = 1;")

	assert.IsType(t, &VariableDeclarationStmt{}, ast)
	varDecl := ast.(*VariableDeclarationStmt)

	assert.Equal(t, "hello", varDecl.Name.Lexeme)
	assert.IsType(t, &LiteralExpr{}, varDecl.Initializer)

	assert.Equal(t, float64(1), varDecl.Initializer.(*LiteralExpr).Value)

	// variable initializer
	ast = getStmtFromSource("var second = a;")

	assert.IsType(t, &VariableDeclarationStmt{}, ast)
	varDecl = ast.(*VariableDeclarationStmt)

	assert.Equal(t, "second", varDecl.Name.Lexeme)
	assert.IsType(t, &VariableExpr{}, varDecl.Initializer)

	assert.Equal(t, "a", varDecl.Initializer.(*VariableExpr).Name.Lexeme)

	// call expression initializer
	ast = getStmtFromSource("var third = foobar();")

	assert.IsType(t, &VariableDeclarationStmt{}, ast)
	varDecl = ast.(*VariableDeclarationStmt)

	assert.Equal(t, "third", varDecl.Name.Lexeme)
	assert.IsType(t, &CallExpr{}, varDecl.Initializer)

	assert.Equal(t, "foobar", varDecl.Initializer.(*CallExpr).Callee.(*VariableExpr).Name.Lexeme)

	// logical expression initializer
	ast = getStmtFromSource("var fourth = a and b;")

	assert.IsType(t, &VariableDeclarationStmt{}, ast)
	varDecl = ast.(*VariableDeclarationStmt)

	assert.Equal(t, "fourth", varDecl.Name.Lexeme)
	assert.IsType(t, &LogicalExpr{}, varDecl.Initializer)
}

func TestFunctionDecl(t *testing.T) {
	// simple function declaration
	source := "fun foobar(){}"
	tokens := lexer.NewScanner(source).ScanTokens()
	p := NewParser(tokens)
	ast := p.declaration()

	assert.Empty(t, p.Errors)
	assert.IsType(t, &FunctionStmt{}, ast)
	assert.Equal(t, "foobar", ast.(*FunctionStmt).Name.Lexeme)

	// function decl with one param
	source = "fun foobar(a){}"
	tokens = lexer.NewScanner(source).ScanTokens()
	p = NewParser(tokens)
	ast = p.declaration()

	assert.Empty(t, p.Errors)
	assert.IsType(t, &FunctionStmt{}, ast)

	functionDecl := ast.(*FunctionStmt)
	assert.Len(t, functionDecl.Params, 1)
	assert.Equal(t, "a", functionDecl.Params[0].Lexeme)

	// function decl with multiple params
	source = "fun foobar(a,b){}"
	tokens = lexer.NewScanner(source).ScanTokens()
	p = NewParser(tokens)
	ast = p.declaration()

	assert.Empty(t, p.Errors)
	assert.IsType(t, &FunctionStmt{}, ast)

	functionDecl = ast.(*FunctionStmt)
	assert.Len(t, functionDecl.Params, 2)
	assert.Equal(t, "a", functionDecl.Params[0].Lexeme)
	assert.Equal(t, "b", functionDecl.Params[1].Lexeme)

	// function decl with body
	source = `fun foobar(a){
		var b = a;
	}`
	tokens = lexer.NewScanner(source).ScanTokens()
	p = NewParser(tokens)
	ast = p.declaration()

	assert.Empty(t, p.Errors)
	assert.IsType(t, &FunctionStmt{}, ast)

	functionDecl = ast.(*FunctionStmt)
	assert.Len(t, functionDecl.Params, 1)
	assert.Equal(t, "a", functionDecl.Params[0].Lexeme)
	assert.Len(t, functionDecl.Body, 1)

	// error case: no identifier
	source = "fun (){}"
	tokens = lexer.NewScanner(source).ScanTokens()
	p = NewParser(tokens)
	_ = p.declaration()

	assert.NotEmpty(t, p.Errors)
	assert.Len(t, p.Errors, 1)

	// error case: no identifier, no parens
	source = "fun {}"
	tokens = lexer.NewScanner(source).ScanTokens()
	p = NewParser(tokens)
	_ = p.declaration()

	assert.NotEmpty(t, p.Errors)
	assert.Len(t, p.Errors, 3)

	// error case: no identifier, no parens
	source = "fun "
	tokens = lexer.NewScanner(source).ScanTokens()
	p = NewParser(tokens)
	_ = p.declaration()

	assert.NotEmpty(t, p.Errors)
	assert.Len(t, p.Errors, 5)

}

func TestSynchronize(t *testing.T) {
	source := "var = 1; var = 2; var = 3;"
	tokens := lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p := NewParser(tokens)
	stmts := p.Parse()

	assert.Len(t, stmts, 3)
	assert.Len(t, p.Errors, 3)

	for _, stmt := range stmts {
		assert.IsType(t, &VariableDeclarationStmt{}, stmt, "stmt should be a variable declaration with an empty token for the identifier")
	}

}

func TestAssignExpr(t *testing.T) {
	// l-value is a `VariableExpr`
	source := "a = 1;"
	tokens := lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p := NewParser(tokens)
	ast := p.declaration()

	assert.IsType(t, &AssignExpr{}, ast.(*ExpressionStmt).Expr)
	assert.Equal(t, "a", ast.(*ExpressionStmt).Expr.(*AssignExpr).Name.Lexeme)
}
