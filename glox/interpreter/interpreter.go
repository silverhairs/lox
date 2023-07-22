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

func (i *Interpreter) Interpret(exp ast.Expression) any {
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
		num, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return -*num

	}

	return nil
}

func (i *Interpreter) VisitBinary(exp *ast.Binary) any {
	left := i.evaluate(exp.Left)
	right := i.evaluate(exp.Right)

	switch exp.Operator.Type {
	case token.GREATER:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return *leftNum > *rightNum
	case token.GREATER_EQ:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return *leftNum >= *rightNum
	case token.LESS:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return *leftNum < *rightNum
	case token.LESS_EQ:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return *leftNum <= *rightNum
	case token.EQ_EQ:
		return isEqual(left, right)
	case token.BANG_EQ:
		return !isEqual(left, right)
	case token.MINUS:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return *leftNum - *rightNum
	case token.PLUS:
		if leftNum, isLFloat := left.(float64); isLFloat {
			rightNum, isRFloat := right.(float64)
			if isRFloat {
				return leftNum + rightNum
			} else if rightVal, isRightStr := right.(string); isRightStr {
				return fmt.Sprintf("%v%s", leftNum, rightVal)
			}
		} else if leftVal, isLeftStr := left.(string); isLeftStr {
			if rightVal, isRightStr := right.(string); isRightStr {
				return leftVal + rightVal
			} else if rightNum, isRightNum := right.(float64); isRightNum {
				return fmt.Sprintf("%s%v", leftVal, rightNum)
			}
		}

		return exception.Runtime(exp.Operator, "unsupported operands. This operation can only be performed with numbers and strings.")

	case token.SLASH:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}

		if *rightNum == 0 {
			return exception.Runtime(exp.Operator, "division by zero")
		}
		return *leftNum / *rightNum
	case token.ASTERISK:
		leftNum, err := checkOperand(exp.Operator, left)
		if err != nil {
			return err
		}
		rightNum, err := checkOperand(exp.Operator, right)
		if err != nil {
			return err
		}
		return *leftNum * *rightNum
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
