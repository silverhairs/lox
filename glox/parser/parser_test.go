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
	tokens, err := lxr.Tokenize()
	if err != nil {
		t.Fatalf("Scanning failed with exception='%v'", err.Error())
	}
	prsr := New(tokens)

	program, err := prsr.Parse()
	if err != nil {
		t.Fatalf("Parsing errors caught: %q", err.Error())
	}

	if len(program) != 1 {
		t.Fatalf("program has wrong number of statements. expected=%d got=%d", 1, len(program))
	}

	stmt, isOk := program[0].(*ast.ExpressionStmt)
	if !isOk {
		t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
	}

	expr := stmt.Exp

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
	tests := []struct {
		code string
		want *ast.Unary
	}{
		{
			code: `!true;`,
			want: ast.NewUnaryExpression(
				token.Token{Type: token.BANG, Lexeme: "!", Line: 1},
				ast.NewLiteralExpression(true),
			),
		},
		{
			code: "-1;",
			want: ast.NewUnaryExpression(
				token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
				ast.NewLiteralExpression(1),
			),
		},
		{
			code: `!false;`,
			want: ast.NewUnaryExpression(
				token.Token{Type: token.BANG, Lexeme: "!", Line: 1},
				ast.NewLiteralExpression(false),
			),
		},
	}
	for _, test := range tests {
		code := test.code
		lxr := lexer.New(code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := New(tokens)

		program, err := prsr.Parse()
		if err != nil {
			t.Fatalf("Parsing errors caught: %v", err.Error())
		}

		if len(program) != 1 {
			t.Fatalf("program has wrong number of statements. expected=%d got=%d", 1, len(program))
		}

		stmt, isOk := program[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
		}

		expr := stmt.Exp

		if passed := testUnary(expr, test.want, t); !passed {
			t.Errorf("testUnary failed for '%s'", code)
			t.FailNow()
		}
	}

}

func TestParseBinary(t *testing.T) {
	tests := []struct {
		code string
		want *ast.Binary
	}{
		{
			code: "5+10;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(5),
				token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
				ast.NewLiteralExpression(10),
			),
		},
		{
			code: "5==12;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(5),
				token.Token{Type: token.EQ_EQ, Lexeme: "==", Line: 1},
				ast.NewLiteralExpression(12),
			),
		},
		{
			code: "true != false;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(true),
				token.Token{Type: token.BANG_EQ, Lexeme: "!=", Line: 1},
				ast.NewLiteralExpression(false),
			),
		},
		{
			code: "13 > 90;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(13),
				token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},

		{
			code: "13 >= 90;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(13),
				token.Token{Type: token.GREATER_EQ, Lexeme: ">=", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: "87 < 90;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(87),
				token.Token{Type: token.LESS, Lexeme: "<", Line: 1},
				ast.NewLiteralExpression(90),
			),
		}, {
			code: "87 <= 90;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(87),
				token.Token{Type: token.LESS_EQ, Lexeme: "<=", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: " 12 * 90;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(12),
				token.Token{Type: token.ASTERISK, Lexeme: "*", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: " 12 / 90;",
			want: ast.NewBinaryExpression(
				ast.NewLiteralExpression(12),
				token.Token{Type: token.SLASH, Lexeme: "/", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
	}

	for _, test := range tests {
		code := test.code
		lxr := lexer.New(code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := New(tokens)

		program, err := prsr.Parse()
		if err != nil {
			t.Fatalf("Parsing errors caught: %v", err.Error())
		}

		if len(program) != 1 {
			t.Fatalf("program has wrong number of statements. expected=%d got=%d", 1, len(program))
		}

		stmt, isOk := program[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
		}

		expr := stmt.Exp
		if passed := testBinary(expr, test.want, t); !passed {
			t.Errorf("testBinary failed for '%s'", code)
			t.FailNow()
		}
	}
}

func TestParseGrouping(t *testing.T) {
	tests := []struct {
		code string
		exp  ast.Expression
	}{
		{
			code: "(5+10);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(5),
				token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
				ast.NewLiteralExpression(10),
			),
		},
		{
			code: "(5==12);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(5),
				token.Token{Type: token.EQ_EQ, Lexeme: "==", Line: 1},
				ast.NewLiteralExpression(12),
			),
		},
		{
			code: "(true != false);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(true),
				token.Token{Type: token.BANG_EQ, Lexeme: "!=", Line: 1},
				ast.NewLiteralExpression(false),
			),
		},
		{
			code: "(13 > 90);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(13),
				token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},

		{
			code: "(13 >= 90);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(13),
				token.Token{Type: token.GREATER_EQ, Lexeme: ">=", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: "(87 < 90);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(87),
				token.Token{Type: token.LESS, Lexeme: "<", Line: 1},
				ast.NewLiteralExpression(90),
			),
		}, {
			code: "(87 <= 90);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(87),
				token.Token{Type: token.LESS_EQ, Lexeme: "<=", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: "(12 * 90);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(12),
				token.Token{Type: token.ASTERISK, Lexeme: "*", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: "(12 / 90);",
			exp: ast.NewBinaryExpression(
				ast.NewLiteralExpression(12),
				token.Token{Type: token.SLASH, Lexeme: "/", Line: 1},
				ast.NewLiteralExpression(90),
			),
		},
		{
			code: `(!true);`,
			exp: ast.NewUnaryExpression(
				token.Token{Type: token.BANG, Lexeme: "!", Line: 1},
				ast.NewLiteralExpression(true),
			),
		},
		{
			code: "(-1);",
			exp: ast.NewUnaryExpression(
				token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
				ast.NewLiteralExpression(1),
			),
		},
		{
			code: `(!false);`,
			exp: ast.NewUnaryExpression(
				token.Token{Type: token.BANG, Lexeme: "!", Line: 1},
				ast.NewLiteralExpression(false),
			),
		},
	}
	for _, test := range tests {
		code := test.code
		lxr := lexer.New(code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := New(tokens)
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

		exp := smt.Exp
		group, isGroup := exp.(*ast.Grouping)
		if !isGroup {
			t.Fatalf("exp is not *ast.Grouping. got=%T", exp)
		}
		if group.Exp.String() != test.exp.String() {
			t.Fatalf("wrong group.Exp expected='%v'. got='%v'", test.exp, group.Exp)
		}
	}
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

func testBinary(exp ast.Expression, expected *ast.Binary, t *testing.T) bool {
	binary, isOk := exp.(*ast.Binary)
	if !isOk {
		t.Errorf("exp is not a *ast.Binary. got='%T'", exp)
		return false
	}

	if binary.Left.String() != expected.Left.String() {
		want, got := expected.Left, binary.Left
		t.Errorf("wrong value for binary.Left. expected='%v' got='%v'", want, got)
		return false
	}

	if binary.Operator.Lexeme != expected.Operator.Lexeme {
		want, got := expected.Operator.Lexeme, binary.Operator.Lexeme
		t.Errorf("wrong value for binary.Operator. expected='%s' got='%s'", want, got)
		return false
	}

	if binary.Right.String() != expected.Right.String() {
		want, got := expected.Right, binary.Right
		t.Errorf("wrong value for binary.Right. expected='%v' got='%v'", want, got)
		return false
	}
	return true
}

func testUnary(exp ast.Expression, expected *ast.Unary, t *testing.T) bool {
	unary, isOk := exp.(*ast.Unary)
	if !isOk {
		t.Errorf("exp is not a *ast.Unary. got='%T'", exp)
		return false
	}

	if unary.Operator.Lexeme != expected.Operator.Lexeme {
		want, got := expected.Operator, unary.Operator
		t.Errorf("wrong value for unary.Operator. expected='%v' got='%v'", want, got)
		return false
	}

	if unary.Right.String() != expected.Right.String() {
		want, got := expected.Right, unary.Right
		t.Errorf("wrong value for unary.Right. expected='%v' got='%v'", want, got)
		return false
	}
	return true
}
