package ast

import "bytes"

type Grouping struct {
	Exp Expression[any]
}

func (exp *Grouping) String() string {
	var out bytes.Buffer

	out.WriteString("( ")
	out.WriteString(exp.Exp.String())
	out.WriteString(" )")

	return out.String()
}

// FIXME: Generic should not be `any`. This is a workaround
// due to the limitation in Go's typesystem.
func (exp *Grouping) Accept(visitor Vistor[any]) any {
	return visitor.VisitGrouping(exp)
}

func NewGroupingExp(exp Exp) *Grouping {
	return &Grouping{exp}
}
