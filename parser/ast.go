package parser

import (
	"github.com/brandonshearin/go-lox/lexer"
)

type Expr interface {
	Expression()
	Accept(visitor Visitor) string
}

type LiteralExpr struct {
	Value     any
	IsBoolean bool
	IsNil     bool
}

func (l *LiteralExpr) Expression() {}

func (l *LiteralExpr) Accept(visitor Visitor) string {
	return visitor.VisitLiteralExpr(l).(string)
}

type UnaryExpr struct {
	// TODO: type of `operator` could maybe be `lexer.Token`?  or the local def for `Operator` could be improved
	Operator Operator
	Expr     Expr
}

func (u *UnaryExpr) Expression() {}

func (u *UnaryExpr) Accept(visitor Visitor) string {
	return visitor.VisitUnaryExpr(u).(string)
}

type BinaryExpr struct {
	LeftExpr  Expr
	Operator  Operator
	RightExpr Expr
}

func (b *BinaryExpr) Expression() {}

func (b *BinaryExpr) Accept(visitor Visitor) string {
	return visitor.VisitBinaryExpr(b).(string)
}

type GroupingExpr struct {
	Expr Expr
}

func (g *GroupingExpr) Expression() {}

func (g *GroupingExpr) Accept(visitor Visitor) string {
	return visitor.VisitGroupingExpr(g).(string)
}

// TODO: is this jank?
type Operator lexer.Token
