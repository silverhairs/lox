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
	VisitBinary(exp *Binary) interface{}
	VisitUnary(exp *Unary) interface{}
	VisitGrouping(exp *Grouping) interface{}
	VisitLiteral(exp *Literal) interface{}
	VisitTernary(exp *Ternary) interface{}
}

type Expression interface {
	String() string
	Type() ExpType
	Accept(Visitor) interface{}
}

func parenthesize(name ExpType, value string) string {
	var out bytes.Buffer

	out.WriteString("(" + string(name) + " ")
	out.WriteString(value)
	out.WriteString(" )")

	return out.String()
}
