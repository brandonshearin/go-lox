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

func TestVarDecl(t *testing.T) {
	source := "var hello = 1;"
	tokens := lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p := NewParser(tokens)
	ast := p.declaration()

	assert.IsType(t, &VariableDeclarationStmt{}, ast)
	assert.Equal(t, "hello", ast.(*VariableDeclarationStmt).Name.Lexeme)

	// TODO: need to test how the synchronize() feature of the parser begins at the next statement.
	source = "var hello = 1 \n var world = 2;"
	tokens = lexer.NewScanner(source).ScanTokens()

	// parse tokens into expression AST
	p = NewParser(tokens)
	stmts := p.Parse()

	assert.IsType(t, &VariableDeclarationStmt{}, stmts[0])
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
