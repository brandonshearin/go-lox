package main

import (
	"fmt"
	"os"

	"github.com/brandonshearin/go-lox/lox"
)

func main() {

	lox := lox.NewLox()
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

// func main_old() {
// 	expr := &parser.BinaryExpr{
// 		LeftExpr: &parser.UnaryExpr{
// 			Operator: parser.Operator{
// 				TokenType: lexer.MINUS,
// 				Lexeme:    "-",
// 				Literal:   "-",
// 				Line:      1,
// 			},
// 			Expr: &parser.LiteralExpr{
// 				Value: "123",
// 			},
// 		},
// 		Operator: parser.Operator{
// 			TokenType: lexer.STAR,
// 			Lexeme:    "*",
// 			Literal:   "*",
// 			Line:      1,
// 		},
// 		RightExpr: &parser.GroupingExpr{
// 			Expr: &parser.LiteralExpr{
// 				Value: "45.67",
// 			},
// 		},
// 	}

// 	ast := parser.ASTPrinter{}

// 	fmt.Println(ast.Print(expr))
// }
