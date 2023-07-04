package ast

import (
	"bytes"
	"fmt"
	"glox/token"
)

type ExpType string

const (
	BINARY_EXP  ExpType = "binary"
	UNARY_EXP   ExpType = "unary"
	GROUP_EXP   ExpType = "group"
	LITERAL_EXP ExpType = "literal"
	TERNARY_EXP ExpType = "ternary"
)

type Expression interface {
	String() string
	Type() ExpType
	Accept(Visitor) any
}

type Visitor interface {
	VisitBinary(exp *Binary) any
	VisitUnary(exp *Unary) any
	VisitGrouping(exp *Grouping) any
	VisitLiteral(exp *Literal) any
	VisitTernary(exp *Ternary) any
}

type Literal struct {
	Value any
}

func (exp *Literal) String() string {
	return parenthesize(exp.Type(), fmt.Sprintf(" %v ", exp.Value))
}

func (exp *Literal) Type() ExpType {
	return LITERAL_EXP
}

func NewLiteralExpression(val any) *Literal {
	return &Literal{val}
}

func (exp *Literal) Accept(v Visitor) any {
	return v.VisitLiteral(exp)
}

type Unary struct {
	Operator token.Token
	Right    Expression
}

func (exp *Unary) String() string {
	var out bytes.Buffer

	out.WriteString("( " + exp.Operator.Lexeme + " ) ")
	out.WriteString(exp.Right.String())

	return parenthesize(exp.Type(), out.String())
}

func (exp *Unary) Type() ExpType {
	return UNARY_EXP
}

func NewUnaryExpression(operator token.Token, right Expression) *Unary {
	return &Unary{operator, right}
}

func (exp *Unary) Accept(v Visitor) any {
	return v.VisitUnary(exp)
}

type Binary struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (exp *Binary) String() string {
	var out bytes.Buffer

	out.WriteString(exp.Left.String())
	out.WriteString(" " + exp.Operator.Lexeme + " ")
	out.WriteString(exp.Right.String())

	return parenthesize(exp.Type(), out.String())
}

func NewBinaryExpression(left Expression, operator token.Token, right Expression) *Binary {
	return &Binary{left, operator, right}
}

func (exp *Binary) Type() ExpType {
	return BINARY_EXP
}

func (exp *Binary) Accept(v Visitor) any {
	return v.VisitBinary(exp)
}

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

func (exp *Grouping) Accept(v Visitor) any {
	return v.VisitGrouping(exp)
}

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

func parenthesize(name ExpType, value string) string {
	var out bytes.Buffer

	out.WriteString("(" + string(name) + " ")
	out.WriteString(value)
	out.WriteString(" )")

	return out.String()
}
