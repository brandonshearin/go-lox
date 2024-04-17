package parser

import (
	"fmt"

	"github.com/brandonshearin/go-lox/lexer"
)

type Environment struct {
	Values    map[string]any
	Enclosing *Environment // implements scope
}

// factory for the global scope
func NewGlobalEnvironment() *Environment {
	return &Environment{
		Values:    map[string]any{},
		Enclosing: nil,
	}
}

// factory for local scopes
func NewEnvironment(e *Environment) *Environment {
	return &Environment{
		Values:    map[string]any{},
		Enclosing: e,
	}
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name lexer.Token) (any, RuntimeError) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	// name wasn't found in current scope, check one level up
	if e.Enclosing != nil {
		return e.Get(name)
	} else {
		// once we reach the global scope, return a runtime error if name never found
		return nil, &RuntimeErrorImpl{
			Token:   name,
			Message: fmt.Sprintf("undefined variable '%s'.", name.Lexeme),
		}
	}

}

func (e *Environment) Assign(name lexer.Token, value any) RuntimeError {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	} else {
		return &RuntimeErrorImpl{
			Token:   name,
			Message: fmt.Sprintf("undefined variable '%s'.", name.Lexeme),
		}
	}

}
