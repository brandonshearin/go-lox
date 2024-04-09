package parser

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

type ExpressionStatement struct {
	Expr Expr
}

func (p *ExpressionStatement) Statement() {}
func (p *ExpressionStatement) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpressionStmt(p)
}
