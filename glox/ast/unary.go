package ast

import (
	"bytes"
	"glox/token"
)

type Unary struct {
	Operator token.Token
	Right    Expression[any]
}

func (exp *Unary) String() string {
	var out bytes.Buffer

	out.WriteString("( " + exp.Operator.Lexeme + " ) ")
	out.WriteString(exp.Right.String())

	return out.String()
}

// FIXME: Generic should not be `any`. This is a workaround
// due to the limitation in Go's typesystem.
func (exp *Unary) Accept(visitor Vistor[any]) any {
	return visitor.VisitUnary(exp)
}

func NewUnaryExpression(operator token.Token, right Exp) *Unary {
	return &Unary{operator, right}
}
