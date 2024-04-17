package parser

import "github.com/brandonshearin/go-lox/lexer"

type Stmt interface {
	Statement()
	Accept(visitor StmtVisitor) RuntimeError
}

type PrintStmt struct {
	Expr Expr
}

func (p *PrintStmt) Statement() {}
func (p *PrintStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitPrintStmt(p)
}

type ExpressionStmt struct {
	Expr Expr
}

func (p *ExpressionStmt) Statement() {}
func (p *ExpressionStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitExpressionStmt(p)
}

type VariableDeclarationStmt struct {
	Name        lexer.Token
	Initializer Expr
}

func (v *VariableDeclarationStmt) Statement() {}
func (v *VariableDeclarationStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitVariableDeclStmt(v)
}

type BlockStmt struct {
	Stmts []Stmt
}

func (b *BlockStmt) Statement() {}
func (v *BlockStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitBlockStmt(v)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *IfStmt) Statement() {}
func (i *IfStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitIfStmt(i)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w *WhileStmt) Statement() {}
func (w *WhileStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitWhileStmt(w)
}

type FunctionStmt struct {
	Name   lexer.Token
	Params []lexer.Token
	Body   []Stmt
}

func (f *FunctionStmt) Statement()                              {}
func (f *FunctionStmt) Accept(visitor StmtVisitor) RuntimeError { return visitor.VisitFunctionStmt(f) }

type ReturnStmt struct {
	Keyword lexer.Token
	Value   Expr
}

func (r *ReturnStmt) Statement() {}
func (r *ReturnStmt) Accept(visitor StmtVisitor) RuntimeError {
	return visitor.VisitReturnStmt(r)
}
