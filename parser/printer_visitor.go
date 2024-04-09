package parser

import (
	"fmt"
	"strings"
)

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
	// handle nil
	if expr.IsNil {
		return "nil", nil
	}

	// handle booleans
	if expr.IsBoolean {
		if expr.Value.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	}

	switch expr.Value.(type) {
	case string:
		return expr.Value.(string), nil
	case float64:
		return fmt.Sprintf("%.2f", expr.Value), nil
	}

	return nil, fmt.Errorf("encountered error printing LiteralExpr: %+v ", expr)

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
