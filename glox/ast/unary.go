package ast

import (
	"bytes"
	"glox/token"
)

type Unary struct {
	Operator token.Token
	Right    Expression
}

func (exp *Unary) String() string {
	var out bytes.Buffer

	out.WriteString("( " + exp.Operator.Lexeme + " ) ")
	out.WriteString(exp.Right.String())

	return out.String()
}

func (exp *Unary) Describe() string {
	return parenthesize(UNARY_EXP, exp)
}

func NewUnaryExpression(operator token.Token, right Expression) *Unary {
	return &Unary{operator, right}
}
