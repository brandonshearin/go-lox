package parser

import (
	"fmt"

	"github.com/brandonshearin/go-lox/lexer"
)

type Environment struct {
	Values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		Values: map[string]any{},
	}
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name lexer.Token) (any, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	return nil, &RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("undefined variable '%s'.", name.Lexeme),
	}
}
