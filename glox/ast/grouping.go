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

	return out.String()
}

func (exp *Grouping) Describe() string {
	return parenthesize(GROUP_EXP, exp)
}

func NewGroupingExp(exp Expression) *Grouping {
	return &Grouping{exp}
}
