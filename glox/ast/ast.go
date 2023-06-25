package ast

import (
	"bytes"
	"fmt"
)

type ExpType string

const (
	BINARY_EXP  ExpType = "binary"
	UNARY_EXP   ExpType = "unary"
	GROUP_EXP   ExpType = "group"
	LITERAL_EXP ExpType = "literal"
)

type Expression interface {
	String() string
	Describe() string
}

func parenthesize(name ExpType, expressions ...Expression) string {
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
