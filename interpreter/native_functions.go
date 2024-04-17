package interpreter

import "time"

type Clock struct{}

func (c *Clock) Arity() int { return 0 }
func (c *Clock) Call(i Interpreter, arguments []any) any {
	return time.Now().UnixMilli() / 1000
}

func (c *Clock) toString() string { return "<native fn>" }
