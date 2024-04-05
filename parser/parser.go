package parser

import "github.com/brandonshearin/go-lox/lexer"

type Parser struct {
	Tokens  []lexer.Token
	Current int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		Current: 0,
		Tokens:  tokens,
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
			Literal: p.previous().Lexeme,
		}
	}

	if p.match(lexer.TRUE, lexer.FALSE) {
		return &LiteralExpr{
			Literal:   p.previous().Lexeme,
			IsBoolean: true,
		}
	}

	if p.match(lexer.NIL) {
		return &LiteralExpr{
			Literal: p.previous().Lexeme,
			IsNil:   true,
		}
	}

	if p.match(lexer.LEFT_PAREN) {
		expr := p.expression()
		p.consume(lexer.RIGHT_PAREN, "Expect '(' after expression)")
		return &GroupingExpr{
			Expr: expr,
		}
	}

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
