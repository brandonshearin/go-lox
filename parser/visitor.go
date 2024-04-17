package parser

type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (any, RuntimeError)
	VisitUnaryExpr(expr *UnaryExpr) (any, RuntimeError)
	VisitGroupingExpr(expr *GroupingExpr) (any, RuntimeError)
	VisitLiteralExpr(expr *LiteralExpr) (any, RuntimeError)
	VisitVariableExpr(expr *VariableExpr) (any, RuntimeError)
	VisitAssignExpr(expr *AssignExpr) (any, RuntimeError)
	VisitLogicalExpr(expr *LogicalExpr) (any, RuntimeError)
	VisitCallExpr(expr *CallExpr) (any, RuntimeError)
}

type StmtVisitor interface {
	VisitPrintStmt(stmt *PrintStmt) RuntimeError
	VisitExpressionStmt(stmt *ExpressionStmt) RuntimeError
	VisitVariableDeclStmt(stmt *VariableDeclarationStmt) RuntimeError
	VisitBlockStmt(stmt *BlockStmt) RuntimeError
	VisitIfStmt(stmt *IfStmt) RuntimeError
	VisitWhileStmt(stmt *WhileStmt) RuntimeError
	VisitFunctionStmt(stmt *FunctionStmt) RuntimeError
	VisitReturnStmt(stmt *ReturnStmt) RuntimeError
}
