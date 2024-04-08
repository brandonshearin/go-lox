package parser

import "strings"

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
}

// ASTPrinter implements the visitor interface
type ASTPrinter struct{}

func (a *ASTPrinter) Print(expr Expr) string {
	val, _ := expr.Accept(a)
	return val.(string)
}

func (a *ASTPrinter) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.LeftExpr, expr.RightExpr), nil
}

func (a *ASTPrinter) VisitUnaryExpr(expr *UnaryExpr) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Expr), nil
}

func (a *ASTPrinter) VisitGroupingExpr(expr *GroupingExpr) (any, error) {
	return a.parenthesize("group", expr.Expr), nil
}

func (a *ASTPrinter) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	if expr.Value == "" {
		return "nil", nil
	}

	return expr.Value.(string), nil
}

func (a *ASTPrinter) parenthesize(name string, expr ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range expr {
		builder.WriteString(" ")
		str, _ := expr.Accept(a)
		builder.WriteString(str.(string))
	}
	builder.WriteString(")")

	return builder.String()
}
