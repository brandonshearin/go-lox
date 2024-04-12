package parser

import (
	"fmt"

	"github.com/brandonshearin/go-lox/lexer"
)

// Interpreter implements `ExprVisitor` interface and `StmtVisitor` interface
type Interpreter struct {
	Environment Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Environment: *NewGlobalEnvironment(),
	}
}

func (s *Interpreter) Interpret(stmts []Stmt) *RuntimeError {
	for _, stmt := range stmts {
		if err := s.execute(stmt); err != nil {
			e := err.(*RuntimeError)
			return e
		}
	}
	return nil
}

// ExprVisitor implementation below ----------------------------------------------------------------
func (s *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.Accept(s)
}

func (s *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (s *Interpreter) VisitGroupingExpr(expr *GroupingExpr) (any, error) {
	return s.evaluate(expr.Expr)
}

func (s *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (any, error) {
	if right, err := s.evaluate(expr.Expr); err != nil {
		return nil, err
	} else {
		switch expr.Operator.TokenType {
		case lexer.BANG:
			return !isTruthy(right), nil
		case lexer.MINUS:
			if err := checkNumberOperand(lexer.Token(expr.Operator), right); err != nil {
				return nil, err
			} else {
				return -right.(float64), nil
			}
		}
	}

	return nil, nil
}

func (s *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (any, error) {
	if left, err := s.evaluate(expr.LeftExpr); err != nil {
		return nil, err
	} else if right, err := s.evaluate(expr.RightExpr); err != nil {
		return nil, err
	} else {
		switch expr.Operator.TokenType {
		case lexer.MINUS:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) - right.(float64), nil
		case lexer.SLASH:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) / right.(float64), nil
		case lexer.STAR:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) * right.(float64), nil
		case lexer.PLUS:
			leftNum, leftIsNumber := left.(float64)
			rightNum, rightIsNumber := right.(float64)

			if leftIsNumber && rightIsNumber {
				return leftNum + rightNum, nil
			}

			leftStr, leftIsString := left.(string)
			rightStr, rightIsString := right.(string)

			if leftIsString && rightIsString {
				return leftStr + rightStr, nil
			}

			return nil, &RuntimeError{
				Token:   lexer.Token(expr.Operator),
				Message: "operands must be two numbers or two strings.",
			}
		case lexer.GREATER:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) > right.(float64), nil
		case lexer.GREATER_EQUAL:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) >= right.(float64), nil
		case lexer.LESS:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) < right.(float64), nil
		case lexer.LESS_EQUAL:
			if err := checkNumberOperands(lexer.Token(expr.Operator), left, right); err != nil {
				return nil, err
			}
			return left.(float64) <= right.(float64), nil
		case lexer.BANG_EQUAL:
			return !isEqual(left, right), nil
		case lexer.EQUAL_EQUAL:
			return isEqual(left, right), nil
		}
	}

	return nil, nil
}

func (s *Interpreter) VisitVariableExpr(expr *VariableExpr) (any, error) {
	return s.Environment.Get(expr.Name)
}

func (s *Interpreter) VisitAssignExpr(expr *AssignExpr) (any, error) {
	if value, err := s.evaluate(expr.Value); err != nil {
		return nil, err
	} else if err := s.Environment.Assign(expr.Name, value); err != nil {
		return nil, err
	} else {
		return value, nil
	}
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

// StmtVisitor implementation below ----------------------------------------------------------------
func (s *Interpreter) execute(stmt Stmt) error {
	return stmt.Accept(s)
}

func (s *Interpreter) VisitPrintStmt(stmt *PrintStmt) error {
	if val, err := s.evaluate(stmt.Expr); err != nil {
		return err
	} else {
		fmt.Println(val)
	}
	return nil
}

func (s *Interpreter) VisitBlockStmt(stmt *BlockStmt) error {
	s.executeBlock(stmt.Stmts, NewEnvironment(&s.Environment))
	return nil
}

func (s *Interpreter) executeBlock(stmts []Stmt, env *Environment) error {
	prev := s.Environment
	defer func() { s.Environment = prev }()

	// before executing these statements, replace the interpreters environment with the new
	s.Environment = *env
	for _, stmt := range stmts {
		// TODO: error handling??
		if err := s.execute(stmt); err != nil {
			return err
		}
	}

	return nil

}

func (s *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) error {
	if _, err := s.evaluate(stmt.Expr); err != nil {
		return err
	}
	return nil
}

func (s *Interpreter) VisitVariableDeclStmt(stmt *VariableDeclarationStmt) error {

	if stmt.Initializer != nil {
		if value, err := s.evaluate(stmt.Initializer); err != nil {
			return fmt.Errorf("error evaluating initializer expression: %v", err)
		} else {
			s.Environment.Define(stmt.Name.Lexeme, value)
		}
	} else {
		s.Environment.Define(stmt.Name.Lexeme, nil)
	}

	return nil
}
