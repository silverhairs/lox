package ast

import (
	"bytes"
	"glox/token"
)

type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (exp *Binary) String() string {
	var out bytes.Buffer

	out.WriteString(exp.Left.String())
	out.WriteString(" " + exp.Operator.Lexeme + " ")
	out.WriteString(exp.Right.String())

	return out.String()
}

func NewBinaryExpression(left Expression, operator token.Token, right Expression) *Binary {
	return &Binary{left, operator, right}
}

func (exp *Binary) Describe() string {
	return parenthesize(BINARY_EXP, exp)
}
