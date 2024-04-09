package parser

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt) error
	VisitExpressionStmt(stmt *ExpressionStatement) error
}
