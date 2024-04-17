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
		stmts = append(stmts, p.declaration())
	}

	return stmts
}

// declaration    → varDecl | statement ;
func (p *Parser) declaration() Stmt {
	var stmt Stmt

	if p.match(lexer.VAR) {
		stmt = p.varDeclaration()
	} else if p.match(lexer.FUN) {
		stmt = p.functionDeclaration("function")
	} else {
		stmt = p.statement()
	}
	// TODO: whats the best way to handle errors
	if len(p.Errors) > 0 {
		p.synchronize()
		return stmt
	}

	return stmt
}

func (p *Parser) functionDeclaration(kind string) Stmt {
	name := p.consume(lexer.IDENTIFIER, fmt.Sprintf("expect %s name.", kind))

	p.consume(lexer.LEFT_PAREN, fmt.Sprintf("expect '(' after %s name.", kind))

	params := []lexer.Token{}

	for {
		params = append(params, p.consume(lexer.IDENTIFIER, "expect parameter name."))

		if !p.match(lexer.COMMA) {
			break
		}
	}

	p.consume(lexer.RIGHT_PAREN, "expect ')' after parameters")

	p.consume(lexer.LEFT_BRACE, fmt.Sprintf("expect '{' before %s body.}", kind))

	body := p.block()

	return &FunctionStmt{
		Name:   name,
		Params: params,
		Body:   body,
	}

}

// varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() Stmt {
	name := p.consume(lexer.IDENTIFIER, "expect variable name.")

	var initializer Expr
	if p.match(lexer.EQUAL) {
		initializer = p.expression()
	}

	p.consume(lexer.SEMICOLON, "expect ';' after variable declaration.")

	return &VariableDeclarationStmt{
		Name:        name,
		Initializer: initializer,
	}
}

func (p *Parser) statement() Stmt {
	if p.match(lexer.IF) {
		return p.ifStatement()
	}

	if p.match(lexer.PRINT) {
		return p.printStatement()
	}

	if p.match(lexer.WHILE) {
		return p.whileStatement()
	}

	if p.match(lexer.FOR) {
		return p.forStatement()
	}

	if p.match(lexer.LEFT_BRACE) {
		return &BlockStmt{
			Stmts: p.block(),
		}
	}

	return p.expressionStatement()
}

func (p *Parser) ifStatement() Stmt {
	p.consume(lexer.LEFT_PAREN, "expect '(' after 'if'.")
	condition := p.expression()
	p.consume(lexer.RIGHT_PAREN, "expect ')' after if condition")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(lexer.ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()

	p.consume(lexer.SEMICOLON, "expect ';' after expression.")

	return &PrintStmt{
		Expr: value,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(lexer.LEFT_PAREN, "expected '(' after while")
	condition := p.expression()
	p.consume(lexer.RIGHT_PAREN, "expected ')' after while condition")

	body := p.statement()

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume(lexer.LEFT_PAREN, "expect '(' after 'for'.")

	// initializer
	var initializer Stmt
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(lexer.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(lexer.SEMICOLON, "expect ';' after loop condition")

	var increment Expr
	if !p.check(lexer.RIGHT_PAREN) {
		increment = p.expression()
	}

	p.consume(lexer.RIGHT_PAREN, "expect ')' after for clauses.")

	body := p.statement()

	// if our loop contains an increment expression, then we append it to the original body so that it executes after the original body stmts
	if increment != nil {
		body = &BlockStmt{
			Stmts: []Stmt{body, &ExpressionStmt{Expr: increment}},
		}
	}

	// if condition is omitted, jam in `true` for an infinite loop`
	if condition == nil {
		condition = &LiteralExpr{
			Value:     true,
			IsBoolean: true,
		}
	}

	// if there is an initializer, it runs once before the entire loop
	if initializer != nil {
		body = &BlockStmt{
			Stmts: []Stmt{initializer, body},
		}
	}

	return body

}

func (p *Parser) block() []Stmt {
	stmts := []Stmt{}

	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}

	p.consume(lexer.RIGHT_BRACE, "expect '}' after block.")
	return stmts
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()

	p.consume(lexer.SEMICOLON, "expect ';' after expression.")

	return &ExpressionStmt{
		Expr: expr,
	}
}

// expression → assignment ;
func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment → IDENTIFIER "=" assignment | equality ;
func (p *Parser) assignment() Expr {
	// expr holds the l-value of the assignment.
	expr := p.or()

	// after parsing the l-value, if an = operator exists, then pass the r-value of the assignment
	if p.match(lexer.EQUAL) {
		equalsTok := p.previous()
		value := p.assignment()

		// the only valid l-value type we accept right now is a variable. ie expressions like `a = "hello"; b = 1234;`
		if variableExpr, ok := expr.(*VariableExpr); ok {
			name := variableExpr.Name

			return &AssignExpr{
				Name:  name,
				Value: value,
			}
		} else {
			p.handleError(equalsTok, "invalid assignment target")
		}
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(lexer.OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{
			Left:     expr,
			Operator: Operator(operator),
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(lexer.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicalExpr{
			Left:     expr,
			Operator: Operator(operator),
			Right:    right,
		}
	}

	return expr
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

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(lexer.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	args := []Expr{}

	if !p.check(lexer.RIGHT_PAREN) {
		for {
			args = append(args, p.expression())

			if !p.match(lexer.COMMA) {
				break
			}
		}
	}

	paren := p.consume(lexer.RIGHT_PAREN, "expect ')' after arguments.")

	return &CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: args,
	}
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
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

	if p.match(lexer.IDENTIFIER) {
		return &VariableExpr{
			Name: p.previous(),
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
		msg := formatErrorMessage(token.Line, fmt.Sprintf("at %s", token.Lexeme), message)
		p.Errors = append(p.Errors, msg)
	}
	return ErrParse
}

func formatErrorMessage(line int, where string, message string) string {
	return fmt.Sprintf("[line %d] Error %s: %s", line, where, message)
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		if p.previous().TokenType == lexer.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		// return once we find the next statement
		case lexer.CLASS,
			lexer.FUN,
			lexer.VAR,
			lexer.FOR,
			lexer.IF,
			lexer.WHILE,
			lexer.PRINT,
			lexer.RETURN:
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
