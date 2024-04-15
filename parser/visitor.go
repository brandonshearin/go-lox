package parser

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (any, error)
	VisitUnaryExpr(expr *UnaryExpr) (any, error)
	VisitGroupingExpr(expr *GroupingExpr) (any, error)
	VisitLiteralExpr(expr *LiteralExpr) (any, error)
	VisitVariableExpr(expr *VariableExpr) (any, error)
	VisitAssignExpr(expr *AssignExpr) (any, error)
	VisitLogicalExpr(expr *LogicalExpr) (any, error)
	VisitCallExpr(expr *CallExpr) (any, error)
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt) error
	VisitExpressionStmt(stmt *ExpressionStmt) error
	VisitVariableDeclStmt(stmt *VariableDeclarationStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
}
