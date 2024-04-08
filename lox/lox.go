package lox

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/brandonshearin/go-lox/lexer"
	"github.com/brandonshearin/go-lox/parser"
)

func NewLox() *Lox {
	return &Lox{
		hadError:        false,
		hadRuntimeError: false,
	}
}

type Lox struct {
	hadError        bool
	hadRuntimeError bool
}

func (l *Lox) RunFile(filename string) error {
	// ReadFile reads the file named by filename and returns the contents.
	// A successful call returns err == nil, not err == EOF.
	data, err := os.ReadFile(filename)
	if err != nil {
		// Return an empty slice and the error
		return err
	}

	sourceCode := string(data)

	l.run(sourceCode)

	return nil
}

func (l *Lox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	// loops until user `control C's``
	for {
		fmt.Print("> ")

		// Wait for the user to input something and press Enter
		scanner.Scan()
		input := scanner.Text()

		// Trim the input to remove any leading or trailing whitespace
		trimmedInput := strings.TrimSpace(input)

		l.run(trimmedInput)
	}
}

func (l *Lox) run(source string) {

	s := lexer.NewScanner(source)

	tokens := s.ScanTokens()

	if l.hadError {
		fmt.Println(">>> lexical error occurred")
		fmt.Println(s.Errors)
		return
	}

	p := parser.NewParser(tokens)
	ast := p.Parse()

	if l.hadError {
		fmt.Println(">>> syntax error occurred")
		fmt.Println(s.Errors)
		return
	}

	interpreter := parser.Interpreter{}
	if value, err := interpreter.Interpret(ast); err != nil {
		l.HandleRuntimeError(*err)
	} else {
		fmt.Println(">>> ast interpreter successfully", value)
	}

	// Create a new scanner from the source string
	// scanner := bufio.NewScanner(strings.NewReader(source))

	// // Set the split function for the scanning operation. lets break on words
	// scanner.Split(bufio.ScanWords)

	// // Create a slice to hold the tokens
	// var tokens []string

	// // Loop over the lines in the input and append them to the tokens slice
	// for scanner.Scan() {
	// 	tokens = append(tokens, scanner.Text())
	// }

	// // for now, just print the tokens
	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }
}

// TODO: maybe collect error messages in a slice on the Lox struct for test assertions
func (l *Lox) HandleError(line int, message string) {
	l.Report(line, "", message)
}

func (l *Lox) Report(line int, where string, message string) {
	l.hadError = true
	fmt.Println("[line ", line, "] Error ", where, ": ", message)
}

func (l *Lox) HandleRuntimeError(e parser.RuntimeError) {
	fmt.Println(e.Message, "\n[line ], ", e.Token.Line, "]")
	l.hadRuntimeError = true
}
