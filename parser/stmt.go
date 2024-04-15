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

type BlockStmt struct {
	Stmts []Stmt
}

func (b *BlockStmt) Statement() {}
func (v *BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlockStmt(v)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *IfStmt) Statement() {}
func (i *IfStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitIfStmt(i)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w *WhileStmt) Statement() {}
func (w *WhileStmt) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhileStmt(w)
}

type FunctionStmt struct {
	Name   lexer.Token
	Params []lexer.Token
	Body   []Stmt
}

func (f *FunctionStmt) Statement()                       {}
func (f *FunctionStmt) Accept(visitor StmtVisitor) error { return visitor.VisitFunctionStmt(f) }
