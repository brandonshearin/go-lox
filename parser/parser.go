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

func (p *Parser) comparison() Expr {
	if p.match(lexer.STRING) {
		return &LiteralExpr{
			Literal: p.previous().Lexeme,
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