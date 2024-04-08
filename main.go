package main

import (
	"fmt"
	"os"

	"github.com/brandonshearin/go-lox/lexer"
	"github.com/brandonshearin/go-lox/parser"
)

func main_old() {

	lox := lexer.NewLox()
	// a file was provided. os.Args[0] is the program name, os.Args[1] is the first argument
	if len(os.Args) == 2 {
		err := lox.RunFile(os.Args[1])
		if err != nil {
			fmt.Printf("there was an error running %s: %s \n", os.Args[1], err.Error())
		}
	} else if len(os.Args) == 1 {
		lox.RunPrompt()
	}

}

func main() {
	expr := &parser.BinaryExpr{
		LeftExpr: &parser.UnaryExpr{
			Operator: parser.Operator{
				TokenType: lexer.MINUS,
				Lexeme:    "-",
				Literal:   "-",
				Line:      1,
			},
			Expr: &parser.LiteralExpr{
				Literal: "123",
			},
		},
		Operator: parser.Operator{
			TokenType: lexer.STAR,
			Lexeme:    "*",
			Literal:   "*",
			Line:      1,
		},
		RightExpr: &parser.GroupingExpr{
			Expr: &parser.LiteralExpr{
				Literal: "45.67",
			},
		},
	}

	ast := parser.ASTPrinter{}

	fmt.Println(ast.Print(expr))
}
