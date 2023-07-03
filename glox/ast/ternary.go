package ast

import (
	"bytes"
	"glox/token"
)

type Ternary struct {
	Condition     Expression  // The expression being evaluated
	LeftOperator  token.Token // Typicaly a `?`
	True          Expression  // Expression executed when the condition is true
	RightOperator token.Token // Typically `:`
	False         Expression  // Expression executed when condition is false
}

func (exp *Ternary) String() string {
	var out bytes.Buffer

	out.WriteString(exp.Condition.String())
	out.WriteString(" " + exp.LeftOperator.Lexeme + " ")
	out.WriteString("(" + exp.True.String() + ")")
	out.WriteString(" " + exp.RightOperator.Lexeme + " ")
	out.WriteString("(" + exp.False.String() + ")")

	return parenthesize(exp.Type(), out.String())
}

func (exp *Ternary) Type() ExpType {
	return TERNARY_EXP
}

func NewTernaryConditional(condition Expression, leftOperator token.Token, positive Expression, rightOperator token.Token, negative Expression) *Ternary {
	return &Ternary{
		Condition:     condition,
		LeftOperator:  leftOperator,
		RightOperator: rightOperator,
		True:          positive,
		False:         negative,
	}
}

func (exp *Ternary) Accept(v Visitor) any {
	return v.VisitTernary(exp)
}
