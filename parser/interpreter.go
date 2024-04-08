package parser

import (
	"fmt"

	"github.com/brandonshearin/go-lox/lexer"
)

type Interpreter struct {
}

func (s *Interpreter) VisitLiteralExpr(expr *LiteralExpr) any {
	return expr.Value
}

func (s *Interpreter) VisitGroupingExpr(expr *GroupingExpr) any {
	return s.evaluate(expr.Expr)
}

func (s *Interpreter) VisitUnaryExpr(expr *UnaryExpr) any {
	right := s.evaluate(expr.Expr)

	switch expr.Operator.TokenType {
	case lexer.BANG:
		return !isTruthy(right)
	case lexer.MINUS:
		checkNumberOperand(lexer.Token(expr.Operator), right)
		return -right.(float64)
	}

	return nil
}

func (s *Interpreter) VisitBinaryExpr(expr *BinaryExpr) any {
	left := s.evaluate(expr.LeftExpr)
	right := s.evaluate(expr.RightExpr)

	switch expr.Operator.TokenType {
	case lexer.MINUS:
		// if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
		// 	return nil, err
		// }
		return left.(float64) - right.(float64)
	case lexer.SLASH:
		checkNumberOperands(lexer.Token(expr.Operator), left, right)
		return left.(float64) / right.(float64)
	case lexer.STAR:
		checkNumberOperands(lexer.Token(expr.Operator), left, right)
		return left.(float64) * right.(float64)
	case lexer.PLUS:
		leftNum, leftIsNumber := left.(float64)
		rightNum, rightIsNumber := right.(float64)

		if leftIsNumber && rightIsNumber {
			return leftNum + rightNum
		}

		leftStr, leftIsString := left.(string)
		rightStr, rightIsString := right.(string)

		if leftIsString && rightIsString {
			return leftStr + rightStr
		}
	case lexer.GREATER:
		checkNumberOperands(lexer.Token(expr.Operator), left, right)
		return left.(float64) > right.(float64)
	case lexer.GREATER_EQUAL:
		checkNumberOperands(lexer.Token(expr.Operator), left, right)
		return left.(float64) >= right.(float64)
	case lexer.LESS:
		checkNumberOperands(lexer.Token(expr.Operator), left, right)
		return left.(float64) < right.(float64)
	case lexer.LESS_EQUAL:
		checkNumberOperands(lexer.Token(expr.Operator), left, right)
		return left.(float64) <= right.(float64)
	case lexer.BANG_EQUAL:
		return !isEqual(left, right)
	case lexer.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	return nil
}

func (s *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(s)
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	if _, ok := obj.(bool); ok {
		return obj.(bool)
	} else {
		return true
	}
}

func isEqual(left, right any) bool {
	if left == nil && right == nil {
		return true
	}

	if left == nil {
		return false
	}

	return left == right
}

func checkNumberOperand(operator lexer.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}

	return &RuntimeError{
		Token:   operator,
		Message: "operand must be a number.",
	}
}

func checkNumberOperands(operator lexer.Token, left, right any) error {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}

	return &RuntimeError{
		Token:   operator,
		Message: "operands must be numbers.",
	}
}

type RuntimeError struct {
	Token   lexer.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Operator: %s, Message: %s", e.Token.Lexeme, e.Message)
}
