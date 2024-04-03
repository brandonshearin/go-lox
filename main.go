package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	lox := NewLox()
	// a file was provided. os.Args[0] is the program name, os.Args[1] is the first argument
	if len(os.Args) == 2 {
		err := lox.runFile(os.Args[1])
		if err != nil {
			fmt.Printf("there was an error running %s: %s \n", os.Args[1], err.Error())
		}
	} else if len(os.Args) == 1 {
		lox.runPrompt()
	}

}

func NewLox() *Lox {
	return &Lox{
		hadError: false,
	}
}

type Lox struct {
	hadError bool
}

func (l *Lox) runFile(filename string) error {
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

func (l *Lox) runPrompt() {
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
	// Create a new scanner from the source string
	scanner := bufio.NewScanner(strings.NewReader(source))

	// Set the split function for the scanning operation. lets break on words
	scanner.Split(bufio.ScanWords)

	// Create a slice to hold the tokens
	var tokens []string

	// Loop over the lines in the input and append them to the tokens slice
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}

	// for now, just print the tokens
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func (l *Lox) handleError(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	l.hadError = true
	fmt.Println("[line ", line, "] Error ", where, ": ", message)
}
