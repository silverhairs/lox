package ast

import (
	"bytes"
	"glox/token"
)

type Binary struct {
	Left     Expression[any]
	Operator token.Token
	Right    Expression[any]
}

func (exp *Binary) String() string {
	var out bytes.Buffer

	out.WriteString(exp.Left.String())
	out.WriteString(" " + exp.Operator.Lexeme + " ")
	out.WriteString(exp.Right.String())

	return out.String()
}

// FIXME: Generic should not be `any`. This is a workaround
// due to the limitation in Go's typesystem.
func (exp *Binary) Accept(visitor Vistor[any]) any {
	return visitor.VisitBinary(exp)
}

func NewBinaryExpression(left Exp, operator token.Token, right Exp) *Binary {
	return &Binary{left, operator, right}
}
