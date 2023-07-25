package ast

import "glox/token"

type StmtVisitor interface {
	VisitPrintStmt(*PrintStmt) any
	VisitExprStmt(*ExpressionStmt) any
	VisitLetStmt(*LetStmt) any
}

type Statement interface {
	Accept(StmtVisitor) any
}

type PrintStmt struct {
	Exp Expression
}

type LetStmt struct {
	Name  token.Token
	Value Expression
}

func NewLetStmt(name token.Token, value Expression) *LetStmt {
	return &LetStmt{Name: name, Value: value}
}

func (stmt *LetStmt) Accept(v StmtVisitor) any {
	return v.VisitLetStmt(stmt)
}

func NewPrintStmt(exp Expression) *PrintStmt {
	return &PrintStmt{Exp: exp}
}

func (stmt *PrintStmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(stmt)
}

type ExpressionStmt struct {
	Exp Expression
}

func NewExprStmt(exp Expression) *ExpressionStmt {
	return &ExpressionStmt{Exp: exp}
}

func (stmt *ExpressionStmt) Accept(v StmtVisitor) any {
	return v.VisitExprStmt(stmt)
}
