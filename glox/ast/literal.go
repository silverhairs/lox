package ast

import "fmt"

type Literal struct {
	Value any
}

func (exp *Literal) String() string {
	return fmt.Sprintf("%v", exp.Value)
}

func (exp *Literal) Describe() string {
	return parenthesize(LITERAL_EXP, exp)
}

func NewLiteralExpression(val any) *Literal {
	return &Literal{val}
}
