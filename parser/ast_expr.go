package parser

import "github.com/brandonshearin/go-lox/lexer"

type Expr interface {
	Expression()
	Accept(visitor ExprVisitor) (any, error)
}

type LiteralExpr struct {
	Value     any
	IsBoolean bool
	IsNil     bool
}

func (l *LiteralExpr) Expression() {}

func (l *LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	// TODO: type of `operator` could maybe be `lexer.Token`?  or the local def for `Operator` could be improved
	Operator Operator
	Expr     Expr
}

func (u *UnaryExpr) Expression() {}

func (u *UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

type BinaryExpr struct {
	LeftExpr  Expr
	Operator  Operator
	RightExpr Expr
}

func (b *BinaryExpr) Expression() {}

func (b *BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expr Expr
}

func (g *GroupingExpr) Expression() {}

func (g *GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

// TODO: is this jank?
type Operator lexer.Token

type VariableExpr struct {
	Name lexer.Token
}

func (v *VariableExpr) Expression() {}
func (v *VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(v)
}

type AssignExpr struct {
	Name  lexer.Token
	Value Expr
}

func (a *AssignExpr) Expression() {}
func (a *AssignExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssignExpr(a)
}

// logical operators 'and' and 'or'
type LogicalExpr struct {
	Operator Operator
	Left     Expr
	Right    Expr
}

func (l *LogicalExpr) Expression() {}
func (l *LogicalExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

type CallExpr struct {
	Callee    Expr
	Paren     lexer.Token
	Arguments []Expr
}

func (c *CallExpr) Expression()                             {}
func (c *CallExpr) Accept(visitor ExprVisitor) (any, error) { return visitor.VisitCallExpr(c) }
