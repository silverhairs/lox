package ast

import "fmt"

type Literal struct {
	Value any
}

func (exp *Literal) String() string {
	return parenthesize(exp.Type(), fmt.Sprintf(" %v ", exp.Value))
}

func (exp *Literal) Type() ExpType {
	return LITERAL_EXP
}

func NewLiteralExpression(val any) *Literal {
	return &Literal{val}
}
