package ast

import "fmt"

type Literal struct {
	Value interface{}
}

func (exp *Literal) String() string {
	return parenthesize(exp.Type(), fmt.Sprintf(" %v ", exp.Value))
}

func (exp *Literal) Type() ExpType {
	return LITERAL_EXP
}

func NewLiteralExpression(val interface{}) *Literal {
	return &Literal{val}
}

func (exp *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteral(exp)
}
