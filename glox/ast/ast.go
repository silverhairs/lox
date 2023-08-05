package ast

import (
	"bytes"
	"fmt"
	"glox/token"
)

type ExpType string

const (
	BINARY_EXP      ExpType = "binary"
	UNARY_EXP       ExpType = "unary"
	GROUP_EXP       ExpType = "group"
	LITERAL_EXP     ExpType = "literal"
	TERNARY_EXP     ExpType = "ternary"
	VARIABLE_EXP    ExpType = "variable"
	ASSIGNMENT_EXP  ExpType = "assignment"
	LOGICAL_OR_EXP  ExpType = "logical_or"
	LOGICAL_AND_EXP ExpType = "logical_and"
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
	VisitVariable(exp *Variable) any
	VisitAssignment(exp *Assignment) any
	VisitLogical(exp *Logical) any
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
	Condition      Expression  // The conditional expression being evaluated
	ThenOperator   token.Token // Typicaly a `?`
	Then           Expression  // Expression executed when the condition is true
	OrElseOperator token.Token // Typically `:`
	OrElse         Expression  // Expression executed when condition is false
}

func (exp *Ternary) String() string {
	var out bytes.Buffer

	out.WriteString(exp.Condition.String())
	out.WriteString(" " + exp.ThenOperator.Lexeme + " ")
	out.WriteString("(" + exp.Then.String() + ")")
	out.WriteString(" " + exp.OrElseOperator.Lexeme + " ")
	out.WriteString("(" + exp.OrElse.String() + ")")

	return parenthesize(exp.Type(), out.String())
}

func (exp *Ternary) Type() ExpType {
	return TERNARY_EXP
}

func NewTernaryConditional(condition Expression, thenOp token.Token, then Expression, orElseOp token.Token, orElse Expression) *Ternary {
	return &Ternary{
		Condition:      condition,
		ThenOperator:   thenOp,
		OrElseOperator: orElseOp,
		Then:           then,
		OrElse:         orElse,
	}
}

func (exp *Ternary) Accept(v Visitor) any {
	return v.VisitTernary(exp)
}

type Variable struct {
	Name token.Token
}

func NewVariable(name token.Token) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) Type() ExpType {
	return VARIABLE_EXP
}

func (v *Variable) String() string {
	return parenthesize(v.Type(), v.Name.Lexeme)
}

func (v *Variable) Accept(visitor Visitor) any {
	return visitor.VisitVariable(v)
}

type Assignment struct {
	Name  token.Token
	Value Expression
}

func NewAssignment(name token.Token, value Expression) *Assignment {
	return &Assignment{Name: name, Value: value}
}

func (exp *Assignment) Type() ExpType {
	return ASSIGNMENT_EXP
}

func (exp *Assignment) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(exp.Name.Lexeme)
	out.WriteString(exp.Value.String())
	out.WriteString(")")
	return parenthesize(exp.Type(), out.String())
}

func (exp *Assignment) Accept(v Visitor) any {
	return v.VisitAssignment(exp)
}

type Logical struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func NewLogical(left Expression, op token.Token, right Expression) *Logical {
	if !(op.Type == token.AND || op.Type == token.OR) {
		panic(fmt.Sprintf("token '%v' cannot be used for logical expressions.", op))
	}
	return &Logical{Left: left, Operator: op, Right: right}
}

func (exp *Logical) Type() ExpType {
	if exp.Operator.Type == token.AND {
		return LOGICAL_AND_EXP
	} else if exp.Operator.Type == token.OR {
		return LOGICAL_OR_EXP
	}

	panic(fmt.Sprintf("token '%v' cannot be used for logical expressions.", exp.Operator))
}

func (exp *Logical) String() string {
	var out bytes.Buffer

	out.WriteString(exp.Left.String() + " " + string(exp.Operator.Type) + " ")
	out.WriteString(exp.Right.String())
	return parenthesize(exp.Type(), out.String())
}

func (exp *Logical) Accept(v Visitor) any {
	return v.VisitLogical(exp)
}

func parenthesize(name ExpType, value string) string {
	var out bytes.Buffer

	out.WriteString("(" + string(name) + " ")
	out.WriteString(value)
	out.WriteString(" )")

	return out.String()
}
