package ast

import "bytes"

type Grouping struct {
	Exp Expression
}

func (exp *Grouping) String() string {
	var out bytes.Buffer

	out.WriteString("( ")
	out.WriteString(exp.Exp.String())
	out.WriteString(" )")

	return parenthesize(exp.Type(), out.String())
}

func (exp *Grouping) Type() ExpType {
	return GROUP_EXP
}

func NewGroupingExp(exp Expression) *Grouping {
	return &Grouping{exp}
}
