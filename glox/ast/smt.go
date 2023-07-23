package ast

import "glox/token"

type StmtVisitor interface {
	VisitPrintStmt(*PrintSmt) any
	VisitExprStmt(*ExpressionStmt) any
	VisitLetSmt(*LetSmt) any
}

type Statement interface {
	Accept(StmtVisitor) any
}

type PrintSmt struct {
	Exp Expression
}

type LetSmt struct {
	Name  token.Token
	Value Expression
}

func NewLetSmt(name token.Token, value Expression) *LetSmt {
	return &LetSmt{Name: name, Value: value}
}

func (smt *LetSmt) Accept(v StmtVisitor) any {
	return v.VisitLetSmt(smt)
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
