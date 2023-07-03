package ast

import (
	"bytes"
)

type ExpType string

const (
	BINARY_EXP  ExpType = "binary"
	UNARY_EXP   ExpType = "unary"
	GROUP_EXP   ExpType = "group"
	LITERAL_EXP ExpType = "literal"
	TERNARY_EXP ExpType = "ternary"
)

type Visitor interface {
	VisitBinary(exp *Binary) any
	VisitUnary(exp *Unary) any
	VisitGrouping(exp *Grouping) any
	VisitLiteral(exp *Literal) any
	VisitTernary(exp *Ternary) any
}

type Expression interface {
	String() string
	Type() ExpType
	Accept(Visitor) any
}

func parenthesize(name ExpType, value string) string {
	var out bytes.Buffer

	out.WriteString("(" + string(name) + " ")
	out.WriteString(value)
	out.WriteString(" )")

	return out.String()
}
