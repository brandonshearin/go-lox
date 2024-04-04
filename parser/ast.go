package parser

import (
	"strings"

	"github.com/brandonshearin/go-lox/lexer"
)

type Expr interface {
	Expression()
	Accept(visitor Visitor) string
}

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) string
	VisitUnaryExpr(expr *UnaryExpr) string
	VisitGroupingExpr(expr *GroupingExpr) string
	VisitLiteralExpr(expr *LiteralExpr) string
}

type ASTPrinter struct{}

func (a *ASTPrinter) Print(expr Expr) string {
	return expr.Accept(a)
}

func (a *ASTPrinter) VisitBinaryExpr(expr *BinaryExpr) string {
	return a.parenthesize(expr.Operator.Token.Lexeme, expr.LeftExpr, expr.RightExpr)
}

func (a *ASTPrinter) VisitUnaryExpr(expr *UnaryExpr) string {
	return a.parenthesize(expr.Operator.Token.Lexeme, expr.Expr)
}

func (a *ASTPrinter) VisitGroupingExpr(expr *GroupingExpr) string {
	return a.parenthesize("group", expr.Expr)
}

func (a *ASTPrinter) VisitLiteralExpr(expr *LiteralExpr) string {
	if expr.Literal == "" {
		return "nil"
	}
	return expr.Literal
}

func (a *ASTPrinter) parenthesize(name string, expr ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range expr {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(a))
	}
	builder.WriteString(")")

	return builder.String()
}

type LiteralExpr struct {
	Literal string
}

func (l *LiteralExpr) Expression() {}

func (l *LiteralExpr) Accept(visitor Visitor) string {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	// TODO: type of `operator` could maybe be `lexer.Token`?  or the local def for `Operator` could be improved
	Operator Operator
	Expr     Expr
}

func (u *UnaryExpr) Expression() {}

func (u *UnaryExpr) Accept(visitor Visitor) string {
	return visitor.VisitUnaryExpr(u)
}

type BinaryExpr struct {
	LeftExpr  Expr
	Operator  Operator
	RightExpr Expr
}

func (b *BinaryExpr) Expression() {}

func (b *BinaryExpr) Accept(visitor Visitor) string {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expr Expr
}

func (g *GroupingExpr) Expression() {}

func (g *GroupingExpr) Accept(visitor Visitor) string {
	return visitor.VisitGroupingExpr(g)
}

type Operator struct {
	Token lexer.Token
}
