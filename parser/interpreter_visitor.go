package parser

import (
	"fmt"
	"time"

	"github.com/brandonshearin/go-lox/lexer"
)

// Interpreter implements `ExprVisitor` interface and `StmtVisitor` interface
type Interpreter struct {
	Environment Environment
}

type Clock struct{}

func (c *Clock) Arity() int { return 0 }
func (c *Clock) Call(i Interpreter, arguments []any) any {
	return time.Now().UnixMilli() / 1000
}
func (c *Clock) toString() string { return "<native fn>" }

func NewInterpreter() *Interpreter {
	globals := NewGlobalEnvironment()
	globals.Define("clock", Clock{})
	return &Interpreter{
		Environment: *globals,
	}
}

func (s *Interpreter) Interpret(stmts []Stmt) *RuntimeErrorImpl {
	for _, stmt := range stmts {
		if err := s.execute(stmt); err != nil {
			e := err.(*RuntimeErrorImpl)
			return e
		}
	}
	return nil
}

// ExprVisitor implementation below ----------------------------------------------------------------
func (s *Interpreter) evaluate(expr Expr) (any, RuntimeError) {
	return expr.Accept(s)
}

func (s *Interpreter) VisitLiteralExpr(expr *LiteralExpr) (any, RuntimeError) {
	return expr.Value, nil
}

func (s *Interpreter) VisitGroupingExpr(expr *GroupingExpr) (any, RuntimeError) {
	return s.evaluate(expr.Expr)
}

func (s *Interpreter) VisitUnaryExpr(expr *UnaryExpr) (any, RuntimeError) {
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

func (s *Interpreter) VisitBinaryExpr(expr *BinaryExpr) (any, RuntimeError) {
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

			return nil, &RuntimeErrorImpl{
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

func (s *Interpreter) VisitVariableExpr(expr *VariableExpr) (any, RuntimeError) {
	return s.Environment.Get(expr.Name)
}

func (s *Interpreter) VisitAssignExpr(expr *AssignExpr) (any, RuntimeError) {
	if value, err := s.evaluate(expr.Value); err != nil {
		return nil, err
	} else if err := s.Environment.Assign(expr.Name, value); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

func (s *Interpreter) VisitLogicalExpr(expr *LogicalExpr) (any, RuntimeError) {
	if left, err := s.evaluate(expr.Left); err != nil {
		return nil, err
	} else {
		if expr.Operator.TokenType == lexer.OR {
			if isTruthy(left) {
				return left, nil
			}
		} else {
			if !isTruthy(left) {
				return left, nil
			}
		}
	}

	return s.evaluate(expr.Right)
}

func (s *Interpreter) VisitCallExpr(expr *CallExpr) (any, RuntimeError) {
	if callee, err := s.evaluate(expr.Callee); err != nil {
		return nil, err
	} else {
		args := []any{}
		for _, arg := range expr.Arguments {
			if arg, err := s.evaluate(arg); err != nil {
				return nil, err
			} else {
				args = append(args, arg)
			}
		}

		if c, ok := callee.(LoxCallable); !ok {
			return nil, &RuntimeErrorImpl{
				Token:   expr.Paren,
				Message: fmt.Sprintf("can only call functions and classes."),
			}
		} else if c.Arity() != len(args) {
			return nil, &RuntimeErrorImpl{
				Token:   expr.Paren,
				Message: fmt.Sprintf("expected %d arguments, got %d", c.Arity, len(args)),
			}
		} else {
			return c.Call(*s, args), nil
		}

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

func checkNumberOperand(operator lexer.Token, operand any) RuntimeError {
	if _, ok := operand.(float64); ok {
		return nil
	}

	return &RuntimeErrorImpl{
		Token:   operator,
		Message: "operand must be a number.",
	}
}

func checkNumberOperands(operator lexer.Token, left, right any) RuntimeError {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}

	return &RuntimeErrorImpl{
		Token:   operator,
		Message: "operands must be numbers.",
	}
}

type RuntimeError interface {
	RuntimeError()
}
type RuntimeErrorImpl struct {
	Token   lexer.Token
	Message string
}

func (e *RuntimeErrorImpl) Error() string {
	return fmt.Sprintf("Operator: %s, Message: %s", e.Token.Lexeme, e.Message)
}

func (e *RuntimeErrorImpl) RuntimeError() {}

// Return implements the RuntimeError interface because it will be thrown/caught following a pattern similar to runtime errors
type Return struct {
	Value any
}

func (e *Return) RuntimeError() {}

// StmtVisitor implementation below ----------------------------------------------------------------
func (s *Interpreter) execute(stmt Stmt) RuntimeError {
	return stmt.Accept(s)
}

// todo: this feels like it wont work
func (s *Interpreter) VisitWhileStmt(stmt *WhileStmt) RuntimeError {

	if val, err := s.evaluate(stmt.Condition); err != nil {
		return err
	} else {
		for isTruthy(val) {
			if err := s.execute(stmt.Body); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Interpreter) VisitReturnStmt(stmt *ReturnStmt) RuntimeError {
	if stmt.Value != nil {
		if value, err := s.evaluate(stmt.Value); err != nil {
			return err
		} else {
			return &Return{
				Value: value,
			}
		}
	}

	return nil
}

func (s *Interpreter) VisitIfStmt(stmt *IfStmt) RuntimeError {
	if val, err := s.evaluate(stmt.Condition); err != nil {
		return err
	} else if isTruthy(val) {
		if err := s.execute(stmt.ThenBranch); err != nil {
			return err
		}
	} else if stmt.ElseBranch != nil {
		if err := s.execute(stmt.ElseBranch); err != nil {
			return err
		}
	}

	return nil
}

func (s *Interpreter) VisitPrintStmt(stmt *PrintStmt) RuntimeError {
	if val, err := s.evaluate(stmt.Expr); err != nil {
		return err
	} else {
		fmt.Println(val)
	}
	return nil
}

func (s *Interpreter) VisitBlockStmt(stmt *BlockStmt) RuntimeError {
	s.executeBlock(stmt.Stmts, NewEnvironment(&s.Environment))
	return nil
}

func (s *Interpreter) executeBlock(stmts []Stmt, env *Environment) RuntimeError {
	prev := s.Environment
	defer func() { s.Environment = prev }()

	// before executing these statements, replace the interpreters environment with the new
	s.Environment = *env
	for _, stmt := range stmts {
		// TODO: RuntimeError handling??
		if err := s.execute(stmt); err != nil {
			return err
		}
	}

	return nil

}

func (s *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) RuntimeError {
	if _, err := s.evaluate(stmt.Expr); err != nil {
		return err
	}
	return nil
}

func (s *Interpreter) VisitVariableDeclStmt(stmt *VariableDeclarationStmt) RuntimeError {

	if stmt.Initializer != nil {
		if value, err := s.evaluate(stmt.Initializer); err != nil {
			// return fmt.Errorf("error evaluating initializer expression: %v", err)
			return &RuntimeErrorImpl{
				Message: fmt.Sprintf("error evaluating initializer expression: %v", err),
			}
		} else {
			s.Environment.Define(stmt.Name.Lexeme, value)
		}
	} else {
		s.Environment.Define(stmt.Name.Lexeme, nil)
	}

	return nil
}

func (s *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) RuntimeError {
	function := NewLoxFunction(*stmt, s.Environment)

	s.Environment.Define(stmt.Name.Lexeme, function)

	return nil
}
