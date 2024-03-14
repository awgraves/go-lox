package expressions

type Visitor interface {
	VisitBinary(Binary) (interface{}, error)
	VisitGrouping(Grouping) (interface{}, error)
	VisitLiteral(Literal) (interface{}, error)
	VisitUnary(Unary) (interface{}, error)
	VisitVariable(Variable) (interface{}, error)
	VisitAssign(Assign) (interface{}, error)
	VisitLogical(Logical) (interface{}, error)
	VisitCall(Call) (interface{}, error)
}
