package ast

type StmtVisitor interface {
	VisitPrintStmt(*PrintSmt) any
	VisitExprStmt(*ExpressionStmt) any
}

type Statement interface {
	Accept(StmtVisitor) any
}

type PrintSmt struct {
	Exp Expression
}

func NewPrintSmt(exp Expression) *PrintSmt {
	return &PrintSmt{Exp: exp}
}

func (smt *PrintSmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(smt)
}

type ExpressionStmt struct {
	Exp Expression
}

func NewExprStmt(exp Expression) *ExpressionStmt {
	return &ExpressionStmt{Exp: exp}
}

func (smt *ExpressionStmt) Accept(v StmtVisitor) any {
	return v.VisitExprStmt(smt)
}
