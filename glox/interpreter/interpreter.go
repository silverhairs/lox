package interpreter

import (
	"fmt"
	"glox/ast"
	"glox/exception"
	"glox/token"
	"math"
)

type Interpreter struct{}

func New() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Eval(exp ast.Expression) any {
	return exp.Accept(i)
}

func (i *Interpreter) VisitLiteral(exp *ast.Literal) any {
	return exp.Value
}

func (i Interpreter) VisitGrouping(exp *ast.Grouping) any {
	return i.evaluate(exp.Exp)
}

func (i *Interpreter) VisitUnary(exp *ast.Unary) any {
	right := i.evaluate(exp.Right)

	switch exp.Operator.Type {
	case token.BANG:
		return !isTruthy(right)
	case token.MINUS:
		num := panicWhenOperandIsNotNumber(exp.Operator, right)
		return -num

	}

	return nil
}

func (i *Interpreter) VisitBinary(exp *ast.Binary) any {
	left := i.evaluate(exp.Left)
	right := i.evaluate(exp.Right)

	switch exp.Operator.Type {
	case token.GREATER:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum > rightNum
	case token.GREATER_EQ:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum >= rightNum
	case token.LESS:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum < rightNum
	case token.LESS_EQ:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum <= rightNum
	case token.EQ_EQ:
		return isEqual(left, right)
	case token.BANG_EQ:
		return !isEqual(left, right)
	case token.MINUS:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum - rightNum
	case token.PLUS:
		var err error
		if leftNum, isLFloat := left.(float64); isLFloat {
			rightNum, isRFloat := right.(float64)
			if isRFloat {
				return leftNum + rightNum
			}
		} else if leftVal, isLeftStr := left.(string); isLeftStr {
			if rightVal, isRightStr := right.(string); isRightStr {
				return leftVal + rightVal
			}
		}
		err = exception.Runtime(exp.Operator, "Both operands must be eihter numbers or strings.")
		panic(err)

	case token.SLASH:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum / rightNum
	case token.ASTERISK:
		leftNum := panicWhenOperandIsNotNumber(exp.Operator, left)
		rightNum := panicWhenOperandIsNotNumber(exp.Operator, right)
		return leftNum * rightNum
	}

	return nil
}

func (i *Interpreter) VisitTernary(exp *ast.Ternary) any {
	condition := i.evaluate(exp.Condition)
	then := i.evaluate(exp.Then)
	orElse := i.evaluate(exp.OrElse)

	if isTruthy(condition) {
		return then
	} else {
		return orElse
	}

}

func (i *Interpreter) evaluate(exp ast.Expression) any {
	return exp.Accept(i)
}

// Only `nil` and `false` are falsey, everything else is truthy.
func isTruthy(object any) bool {
	if object == nil {
		return false
	}

	if val, isBool := object.(bool); isBool {
		return val
	}

	return true
}

func isEqual(l any, r any) bool {
	lNum, isLOk := l.(float64)
	rNum, isROk := r.(float64)
	if isLOk && isROk {
		if math.IsNaN(lNum) && math.IsNaN(rNum) {
			return true
		}

		return lNum == rNum
	}

	return l == r
}

func checkOperand(operator token.Token, operand any) (*float64, error) {
	num, isNum := operand.(float64)
	if !isNum {
		return nil, exception.Runtime(
			operator,
			fmt.Sprintf("Operator %q only accepts number operands.", operator.Lexeme),
		)
	}

	return &num, nil
}

func panicWhenOperandIsNotNumber(operator token.Token, operand any) float64 {
	num, err := checkOperand(operator, operand)
	if err != nil {
		panic(err)
	} else {
		return *num
	}
}
