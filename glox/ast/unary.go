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

	return parenthesize(exp.Type(), out.String())
}

func (exp *Unary) Type() ExpType {
	return UNARY_EXP
}

func NewUnaryExpression(operator token.Token, right Expression) *Unary {
	return &Unary{operator, right}
}
