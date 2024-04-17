package parser

import "fmt"

// implements LoxCallable
type LoxFunction struct {
	Declaration FunctionStmt
	Closure     Environment
}

func NewLoxFunction(decl FunctionStmt, closure Environment) *LoxFunction {
	return &LoxFunction{
		Declaration: decl,
		Closure:     closure,
	}
}

func (s *LoxFunction) Call(interpreter Interpreter, arguments []any) any {
	env := NewEnvironment(&s.Closure)

	for i, param := range s.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	if err := interpreter.executeBlock(s.Declaration.Body, env); err != nil {
		// the return value will bubble up the call stack as a RuntimeError
		if ret, ok := err.(*Return); ok {
			return ret.Value
		}
	}

	return nil
}

func (s *LoxFunction) Arity() int {
	return len(s.Declaration.Params)
}

func (s *LoxFunction) toString() string {
	return fmt.Sprintf("<fn %s >", s.Declaration.Name.Lexeme)
}
