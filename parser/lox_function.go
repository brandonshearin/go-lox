package parser

import "fmt"

// implements LoxCallable
type LoxFunction struct {
	Declaration FunctionStmt
}

func NewLoxFunction(decl FunctionStmt) *LoxFunction {
	return &LoxFunction{
		Declaration: decl,
	}
}

func (s *LoxFunction) Call(interpreter Interpreter, arguments []any) any {
	env := NewEnvironment(&interpreter.Environment)

	for i, param := range s.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	if err := interpreter.executeBlock(s.Declaration.Body, env); err != nil {
		return err
	}

	return nil
}

func (s *LoxFunction) Arity() int {
	return len(s.Declaration.Params)
}

func (s *LoxFunction) toString() string {
	return fmt.Sprintf("<fn %s >", s.Declaration.Name.Lexeme)
}
