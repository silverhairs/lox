package interpreter

import (
	"glox/ast"
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
		return -right.(float64)
	}

	return nil
}

func (i *Interpreter) VisitBinary(exp *ast.Binary) any {
	left := i.evaluate(exp.Left)
	right := i.evaluate(exp.Right)

	switch exp.Operator.Type {
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQ:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQ:
		return left.(float64) <= right.(float64)
	case token.EQ_EQ:
		return isEqual(left, right)
	case token.BANG_EQ:
		return !isEqual(left, right)
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.PLUS:
		if leftNum, isLFloat := left.(float64); isLFloat {
			rightNum, isRFloat := right.(float64)
			if isRFloat {
				return leftNum + rightNum
			}
		}

		// String concatenation
		if leftStr, isLStr := left.(string); isLStr {
			rightStr, isRStr := right.(string)
			if isRStr {
				return leftStr + rightStr
			}
		}

	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.ASTERISK:
		return left.(float64) * right.(float64)
	}

	return nil
}

func (i *Interpreter) VisitTernary(exp *ast.Ternary) any {
	condition := i.evaluate(exp.Condition)
	then := i.evaluate(exp.Then)
	orElse := i.evaluate(exp.OrElse)

	passes := isTruthy(condition)
	if passes {
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
