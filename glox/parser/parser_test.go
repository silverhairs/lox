package parser

import (
	"glox/ast"
	"glox/lexer"
	"glox/token"
	"glox/utils"
	"testing"
)

func TestParseTernary(t *testing.T) {
	code := ` 15 > 1 ? "abc" : "123";`

	lxr := lexer.New(code)
	prsr := New(lxr.Tokenize())

	result, errs := prsr.Parse()
	if len(errs) > 0 {
		messages := utils.Map[error, string](errs, func(e error) string { return e.Error() })
		t.Fatalf("Parsing errors caught: %v", messages)
	}

	exp, isTernary := result.(*ast.Ternary)
	if !isTernary {
		t.Fatalf("result is not *ast.Ternary. got=%T", result)
	}

	if exp.ThenOperator.Type != token.QUESTION_MARK {
		t.Fatalf("exp.LeftOperator has the wrong token. expected='token.QUESTION_MARK' but got=%q", exp.ThenOperator.Type)
	}

	if exp.OrElseOperator.Type != token.COLON {
		t.Fatalf("exp.RightOperator has the wrong token. expected='token.COLON' but got=%q", exp.ThenOperator.Type)
	}

	condition, isBinary := exp.Condition.(*ast.Binary)
	if !isBinary {
		t.Fatalf("exp.Condition is not *ast.Binary. got=%T", exp.Condition)
	}

	testLiteral(condition.Left, "15", t)
	if condition.Operator.Type != token.GREATER {
		t.Fatalf("exp.Operator.Type is not token.GREATER. got=%q", string(condition.Operator.Type))
	}

	testLiteral(condition.Right, "1", t)
	testLiteral(exp.Then, "abc", t)
	testLiteral(exp.OrElse, "123", t)

}

func TestParseUnary(t *testing.T) {
	code := `!true`
	lxr := lexer.New(code)
	prsr := New(lxr.Tokenize())

	result, errs := prsr.Parse()
	if len(errs) > 0 {
		messages := utils.Map[error, string](errs, func(e error) string { return e.Error() })
		t.Fatalf("Parsing errors caught: %v", messages)
	}

	exp, isUnary := result.(*ast.Unary)
	if !isUnary {
		t.Fatalf("result is not *ast.Unary. got=%T", result)
	}

	if exp.Operator.Type != token.BANG {
		t.Fatalf("exp.Operator.Type is not token.BANG. got=%q", string(exp.Operator.Type))
	}

	testLiteral(exp.Right, "true", t)
}

func TestParseBinary(t *testing.T) {
	code := `5 + 10;`
	lxr := lexer.New(code)
	prsr := New(lxr.Tokenize())

	result, errs := prsr.Parse()
	if len(errs) > 0 {
		messages := utils.Map[error, string](errs, func(e error) string { return e.Error() })
		t.Fatalf("Parsing errors caught: %v", messages)
	}

	exp, isBinary := result.(*ast.Binary)
	if !isBinary {
		t.Fatalf("parsed expression is not *ast.Binary. got=%T", result)
	}

	testLiteral(exp.Left, "5", t)
	if exp.Operator.Type != token.PLUS {
		t.Fatalf("exp.Operator.Type is not token.PLUS. got=%q", string(exp.Operator.Type))
	}

	testLiteral(exp.Right, "10", t)

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
