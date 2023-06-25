package ast

import (
	"bytes"
	"fmt"
)

type ExpName string

const (
	BINARY_EXP  ExpName = "binary"
	UNARY_EXP   ExpName = "unary"
	GROUP_EXP   ExpName = "group"
	LITERAL_EXP ExpName = "literal"
)

type Expression interface {
	String() string
	Describe() string
}

func parenthesize(name ExpName, expressions ...Expression) string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(fmt.Sprintf("%v", name))

	for _, exp := range expressions {
		out.WriteString(" ")
		out.WriteString(exp.String())
	}
	out.WriteString(")")

	return out.String()
}
