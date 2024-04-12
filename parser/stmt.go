package parser

import "github.com/brandonshearin/go-lox/lexer"

type Stmt interface {
	Statement()
	Accept(visitor StmtVisitor) error
}

type PrintStmt struct {
	Expr Expr
}

func (p *PrintStmt) Statement() {}
func (p *PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitPrintStmt(p)
}

type ExpressionStmt struct {
	Expr Expr
}

func (p *ExpressionStmt) Statement() {}
func (p *ExpressionStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(p)
}

type VariableDeclarationStmt struct {
	Name        lexer.Token
	Initializer Expr
}

func (v *VariableDeclarationStmt) Statement() {}
func (v *VariableDeclarationStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitVariableDeclStmt(v)
}
