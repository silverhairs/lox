package ast

import "fmt"

type Literal struct {
	Value any
}

func (exp *Literal) String() string {
	return fmt.Sprintf("%v", exp.Value)
}

// FIXME: Generic should not be `any`. This is a workaround
// due to the limitation in Go's typesystem.
func (exp *Literal) Accept(visitor Vistor[any]) any {
	return visitor.VisitLiteral(exp)
}

func NewLiteralExpression(val any) *Literal {
	return &Literal{val}
}
