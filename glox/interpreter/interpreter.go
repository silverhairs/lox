package interpreter

import (
	"fmt"
	"glox/ast"
	"glox/env"
	"glox/exception"
	"glox/token"
	"io"
	"math"
)

type Interpreter struct {
	StdOut io.Writer
	StdErr io.Writer
	Env    *env.Environment
}

func New(stderr io.Writer, stdout io.Writer) *Interpreter {
	return &Interpreter{StdOut: stdout, StdErr: stderr, Env: env.Global()}
}

func (i *Interpreter) Interpret(stmts []ast.Statement) any {
	var err error
	for _, stmt := range stmts {
		i.execute(stmt)
	}
	return err
}

func (i *Interpreter) execute(stmt ast.Statement) {
	val := stmt.Accept(i)
	if err, isErr := val.(error); isErr {
		fmt.Fprintf(i.StdErr, "%s\n", err.Error())
	}
}

func (i *Interpreter) VisitLetStmt(stmt *ast.LetStmt) any {
	var val any
	if stmt.Value != nil {
		val = i.evaluate(stmt.Value)
	}
	if err, isErr := val.(error); isErr {
		return err
	}

	i.Env.Define(stmt.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) any {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Then)
	} else if stmt.OrElse != nil {
		i.execute(stmt.OrElse)
	}

	return nil
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExpressionStmt) any {
	val := i.evaluate(stmt.Exp)
	if err, isErr := val.(error); isErr {
		return err
	}
	fmt.Fprintf(i.StdOut, "%v\n", val)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) any {
	val := i.evaluate(stmt.Exp)
	if err, isErr := val.(error); isErr {
		fmt.Fprintf(i.StdErr, "%s\n", err.Error())
	} else {
		fmt.Fprintf(i.StdOut, "%v\n", val)
	}
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) any {
	i.executeBlock(stmt.Stmts, env.New(i.Env))
	return nil
}

func (i *Interpreter) executeBlock(stmts []ast.Statement, env *env.Environment) {
	prev := i.Env
	i.Env = env
	for _, stmt := range stmts {
		i.execute(stmt)
	}
	i.Env = prev
}

func (i *Interpreter) VisitVariable(exp *ast.Variable) any {
	if res := i.Env.Get(exp.Name); res != nil {
		return res
	}

	return exception.Runtime(exp.Name, fmt.Sprintf("tried to access variable '%s' which holds a nil value.", exp.Name.Lexeme))
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

func (i *Interpreter) VisitAssignment(exp *ast.Assignment) any {
	val := i.evaluate(exp.Value)
	if err, isErr := val.(error); isErr {
		return err
	}
	if err := i.Env.Assign(exp.Name, val); err != nil {
		return err
	}
	return val
}

func (i *Interpreter) VisitLogical(exp *ast.Logical) any {
	left := i.evaluate(exp.Left)
	if err, isErr := left.(error); isErr {
		return err
	}

	if exp.Operator.Type == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(exp.Right)
}

func (i *Interpreter) VisitWhile(exp *ast.WhileStmt) any {
	cond := i.evaluate(exp.Condition)
	for isTruthy(cond) {
		if err, isErr := cond.(error); isErr {
			return err
		}
		i.execute(exp.Body)
		cond = i.evaluate(exp.Condition)
	}

	return nil
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
