package statements

type Visitor interface {
	VisitExpressionStmt(ExpStmt) error
	VisitPrintStmt(PrintStmt) error
	VisitVarStmt(VarStmt) error
	VisitBlock(Block) error
}
