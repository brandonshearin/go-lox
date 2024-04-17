package interpreter

type LoxCallable interface {
	Call(interpreter Interpreter, arguments []any) any
	Arity() int
}
