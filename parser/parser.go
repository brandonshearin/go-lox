package parser

import (
	"errors"
	"fmt"

	"github.com/brandonshearin/go-lox/lexer"
)

type Parser struct {
	Tokens  []lexer.Token
	Current int
	Errors  []string
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		Current: 0,
		Tokens:  tokens,
	}
}

func (p *Parser) Parse() []Stmt {
	stmts := []Stmt{}

	for !p.isAtEnd() {
		stmts = append(stmts, p.statement())
	}

	return stmts
}

func (p *Parser) statement() Stmt {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()

	p.consume(lexer.SEMICOLON, "expect ';' after expression.")

	return &PrintStmt{
		Expr: value,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()

	p.consume(lexer.SEMICOLON, "expect ';' after expression.")

	return &ExpressionStatement{
		Expr: expr,
	}
}

// expression → equality ;
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()

		expr = &BinaryExpr{
			LeftExpr:  expr,
			RightExpr: right,
			Operator:  Operator(operator),
		}
	}

	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

		expr = &BinaryExpr{
			LeftExpr:  expr,
			Operator:  Operator(operator),
			RightExpr: right,
		}
	}

	return expr
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{
			LeftExpr:  expr,
			Operator:  Operator(operator),
			RightExpr: right,
		}
	}
	return expr
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{
			LeftExpr:  expr,
			Operator:  Operator(operator),
			RightExpr: right,
		}
	}

	return expr
}

// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{
			Operator: Operator(operator),
			Expr:     right,
		}
	}

	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() Expr {
	if p.match(lexer.NUMBER, lexer.STRING) {
		return &LiteralExpr{
			Value: p.previous().Literal,
		}
	}

	// our lexer doesn't support storing raw boolean values in the tokens it emits, so account for that here instead
	if p.match(lexer.TRUE) {
		return &LiteralExpr{
			Value:     true,
			IsBoolean: true,
		}
	}

	if p.match(lexer.FALSE) {
		return &LiteralExpr{
			Value:     false,
			IsBoolean: true,
		}
	}

	if p.match(lexer.NIL) {
		return &LiteralExpr{
			Value: p.previous().Literal,
			IsNil: true,
		}
	}

	if p.match(lexer.LEFT_PAREN) {
		expr := p.expression()
		p.consume(lexer.RIGHT_PAREN, "Expect '(' after expression)")
		return &GroupingExpr{
			Expr: expr,
		}
	}

	// if we reach here, throw a syntax error
	p.handleError(p.peek(), "Expect expression.")

	return nil
}

// --------------------- parser machinery
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType lexer.TokenType, message string) lexer.Token {
	if p.check(tokenType) {
		return p.advance()
	} else {
		p.handleError(p.peek(), message)
		return lexer.Token{}
	}
}

var (
	ErrParse = errors.New("parse error")
)

func (p *Parser) handleError(token lexer.Token, message string) error {
	if token.TokenType == lexer.EOF {

		msg := formatErrorMessage(token.Line, " at end", message)
		p.Errors = append(p.Errors, msg)
	} else {
		msg := formatErrorMessage(token.Line, fmt.Sprintf(" at %s", token.Lexeme), message)
		p.Errors = append(p.Errors, msg)
	}
	return ErrParse
}

func formatErrorMessage(line int, where string, message string) string {
	return fmt.Sprintf("[line %d] Error %s: %s", line, where, message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == lexer.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case lexer.CLASS:
		case lexer.FUN:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) check(tokenType lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.Current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == lexer.EOF
}

func (p *Parser) peek() lexer.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() lexer.Token {
	return p.Tokens[p.Current-1]
}
