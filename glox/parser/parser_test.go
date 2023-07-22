package parser

import (
	"glox/ast"
	"glox/lexer"
	"glox/token"
	"testing"
)

func TestParseTernary(t *testing.T) {
	code := ` 15 > 1 ? "abc" : "123";`

	lxr := lexer.New(code)
	prsr := New(lxr.Tokenize())

	program, err := prsr.Parse()
	if err != nil {
		t.Fatalf("Parsing errors caught: %q", err.Error())
	}

	if len(program) != 1 {
		t.Fatalf("program has wrong number of statements. expected=%d got=%d", 1, len(program))
	}

	smt, isOk := program[0].(*ast.ExpressionStmt)
	if !isOk {
		t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
	}

	expr := smt.Exp

	ternary, isTernary := expr.(*ast.Ternary)
	if !isTernary {
		t.Fatalf("result is not *ast.Ternary. got=%T", expr)
	}

	if ternary.ThenOperator.Type != token.QUESTION_MARK {
		t.Fatalf("exp.LeftOperator has the wrong token. expected='token.QUESTION_MARK' but got=%q", ternary.ThenOperator.Type)
	}

	if ternary.OrElseOperator.Type != token.COLON {
		t.Fatalf("exp.RightOperator has the wrong token. expected='token.COLON' but got=%q", ternary.ThenOperator.Type)
	}

	condition, isBinary := ternary.Condition.(*ast.Binary)
	if !isBinary {
		t.Fatalf("exp.Condition is not *ast.Binary. got=%T", ternary.Condition)
	}

	testLiteral(condition.Left, "15", t)
	if condition.Operator.Type != token.GREATER {
		t.Fatalf("exp.Operator.Type is not token.GREATER. got=%q", string(condition.Operator.Type))
	}

	testLiteral(condition.Right, "1", t)
	testLiteral(ternary.Then, "abc", t)
	testLiteral(ternary.OrElse, "123", t)

}

func TestParseUnary(t *testing.T) {
	code := `!true`
	lxr := lexer.New(code)
	prsr := New(lxr.Tokenize())

	program, err := prsr.Parse()
	if err != nil {
		t.Fatalf("Parsing errors caught: %v", err.Error())
	}

	if len(program) != 1 {
		t.Fatalf("program has wrong number of statements. expected=%d got=%d", 1, len(program))
	}

	smt, isOk := program[0].(*ast.ExpressionStmt)
	if !isOk {
		t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
	}

	expr := smt.Exp

	unary, isUnary := expr.(*ast.Unary)
	if !isUnary {
		t.Fatalf("result is not *ast.Unary. got=%T", expr)
	}

	if unary.Operator.Type != token.BANG {
		t.Fatalf("exp.Operator.Type is not token.BANG. got=%q", string(unary.Operator.Type))
	}

	testLiteral(unary.Right, "true", t)
}

func TestParseBinary(t *testing.T) {
	code := `5 + 10;`
	lxr := lexer.New(code)
	prsr := New(lxr.Tokenize())

	program, err := prsr.Parse()
	if err != nil {
		t.Fatalf("Parsing errors caught: %v", err.Error())
	}

	if len(program) != 1 {
		t.Fatalf("program has wrong number of statements. expected=%d got=%d", 1, len(program))
	}

	smt, isOk := program[0].(*ast.ExpressionStmt)
	if !isOk {
		t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
	}

	expr := smt.Exp

	binary, isBinary := expr.(*ast.Binary)
	if !isBinary {
		t.Fatalf("parsed expression is not *ast.Binary. got=%T", expr)
	}

	testLiteral(binary.Left, "5", t)
	if binary.Operator.Type != token.PLUS {
		t.Fatalf("exp.Operator.Type is not token.PLUS. got=%q", string(binary.Operator.Type))
	}

	testLiteral(binary.Right, "10", t)

}

func testLiteral(exp ast.Expression, expectedValue any, t *testing.T) {
	isLiteral, literal := assertLiteral(exp, ast.NewLiteralExpression(expectedValue))
	if !isLiteral {
		t.Fatalf("result.Left is not *ast.Literal. got=%T", exp)
	}
	if literal.Value != expectedValue {
		t.Fatalf("literal.Value wrong. expected=5 got=%q", literal.Value)
	}

}

func assertLiteral(exp ast.Expression, expected *ast.Literal) (bool, *ast.Literal) {
	literal, isLiteral := exp.(*ast.Literal)
	if !isLiteral || literal.String() != exp.String() {
		return false, nil
	}
	return true, expected
}
