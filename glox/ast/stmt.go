package ast

import "glox/token"

type StmtVisitor interface {
	VisitPrintStmt(*PrintStmt) any
	VisitExprStmt(*ExpressionStmt) any
	VisitLetStmt(*LetStmt) any
	VisitBlockStmt(*BlockStmt) any
	VisitIfStmt(*IfStmt) any
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

type BlockStmt struct {
	Stmts []Statement
}

func NewBlockStmt(stmts []Statement) *BlockStmt {
	return &BlockStmt{Stmts: stmts}
}

func (stmt *BlockStmt) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(stmt)
}

type IfStmt struct {
	Condition Expression
	Then      Statement
	OrElse    Statement
}

func NewIfStmt(cond Expression, then Statement, orelse Statement) *IfStmt {
	return &IfStmt{Condition: cond, Then: then, OrElse: orelse}
}

func (stmt *IfStmt) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(stmt)
}
