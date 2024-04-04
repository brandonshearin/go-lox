package main

import (
	"fmt"
	"os"

	"github.com/brandonshearin/go-lox/lexer"
)

func main() {

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
