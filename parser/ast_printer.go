package parser

import "strings"

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
	VisitLiteralExpr(expr *LiteralExpr) any
}

// ASTPrinter implements the visitor interface
type ASTPrinter struct{}

func (a *ASTPrinter) Print(expr Expr) string {
	return expr.Accept(a)
}

func (a *ASTPrinter) VisitBinaryExpr(expr *BinaryExpr) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.LeftExpr, expr.RightExpr)
}

func (a *ASTPrinter) VisitUnaryExpr(expr *UnaryExpr) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Expr)
}

func (a *ASTPrinter) VisitGroupingExpr(expr *GroupingExpr) any {
	return a.parenthesize("group", expr.Expr)
}

func (a *ASTPrinter) VisitLiteralExpr(expr *LiteralExpr) any {
	if expr.Value == "" {
		return "nil"
	}

	return expr.Value.(string)
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
