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

func (p *printer) VisitBinary(binary *Binary) any {
	return binary.String()
}

func (p *printer) VisitUnary(unary *Unary) any {
	return unary.String()
}

func (p *printer) VisitGrouping(grouping *Grouping) any {
	return grouping.String()
}

func (p *printer) VisitLiteral(literal *Literal) any {
	return literal.String()
}

func (p *printer) VisitTernary(ternary *Ternary) any {
	return ternary.String()
}
