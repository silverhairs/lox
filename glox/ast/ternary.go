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

	return out.String()
}

func (exp *Ternary) Describe() string {
	return parenthesize(TERNARY_EXP, exp)
}
