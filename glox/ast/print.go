package ast

import (
	"fmt"
)

type printer struct{}

func NewPrinter() *printer {
	return &printer{}
}

func (p *printer) Print(exp Expression) string {
	if str, isOk := exp.Accept(p).(string); isOk {
		return str
	}

	panic(fmt.Sprintf("%T %+v cannot accept *ASTPrinter", exp, exp))
}

func (p *printer) VisitBinary(binary *Binary) interface{} {
	return binary.String()
}

func (p *printer) VisitUnary(unary *Unary) interface{} {
	return unary.String()
}

func (p *printer) VisitGrouping(grouping *Grouping) interface{} {
	return grouping.String()
}

func (p *printer) VisitLiteral(literal *Literal) interface{} {
	return literal.String()
}

func (p *printer) VisitTernary(ternary *Ternary) interface{} {
	return ternary.String()
}
