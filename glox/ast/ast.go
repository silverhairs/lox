package ast

import (
	"bytes"
)

const (
	BINARY_EXP  = "binary"
	UNARY_EXP   = "unary"
	GROUP_EXP   = "group"
	LITERAL_EXP = "literal"
)

type Exp Expression[any]

type Expression[R any] interface {
	String() string
	Accept(v Vistor[R]) R
}

type Vistor[T any] interface {
	VisitBinary(exp *Binary) T
	VisitUnary(exp *Unary) T
	VisitLiteral(exp *Literal) T
	VisitGrouping(exp *Grouping) T
}

func parenthesize(name string, expressions ...Expression[any]) string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(name)

	for _, exp := range expressions {
		out.WriteString(" ")
		out.WriteString(exp.String())
	}
	out.WriteString(")")

	return out.String()
}

type PrettyPrinter struct {
}

func (p *PrettyPrinter) VisitBinary(exp *Binary) string {
	return parenthesize(BINARY_EXP, exp)
}

func (p *PrettyPrinter) VisitGrouping(exp *Grouping) string {
	return parenthesize(GROUP_EXP, exp)
}

func (p *PrettyPrinter) VisitUnary(exp *Unary) string {
	return parenthesize(UNARY_EXP, exp)
}

func (p *PrettyPrinter) VisitLiteral(exp *Literal) string {
	return parenthesize(LITERAL_EXP, exp)
}

func (p *PrettyPrinter) Format(exp Exp) string {
	text := exp.Accept(p)
}
