package statements

type Visitor interface {
	VisitExpressionStmt(ExpStmt) error
	VisitFunctionStmt(FunctionStmt) error
	VisitPrintStmt(PrintStmt) error
	VisitVarStmt(VarStmt) error
	VisitBlock(Block) error
	VisitIfStmt(IfStmt) error
	VisitWhileStmt(WhileStmt) error
}
