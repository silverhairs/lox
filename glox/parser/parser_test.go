package parser

import (
	"glox/ast"
	"glox/lexer"
	"glox/token"
	"strings"
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
	code := `(5+10;`
	lxr := lexer.New(code)
	tokens, err := lxr.Tokenize()
	if err != nil {
		t.Fatalf("Scanning failed with exception='%v'", err.Error())
	}
	prsr := New(tokens)
	_, err = prsr.Parse()
	if err == nil {
		t.Fatalf("Parsing should have caught an error on code='%s'", code)
	}

	chunk := "expected ')' after expression"
	got := err.Error()
	if !strings.Contains(got, chunk) {
		t.Fatalf("exception message wrong. expected to contain='%s' but got='%s'", chunk, got)
	}
}

func TestParseVariable(t *testing.T) {
	tests := []struct {
		code string
		want *ast.Variable
	}{
		{
			code: "maybe_12;",
			want: ast.NewVariable(token.Token{Type: token.IDENTIFIER, Lexeme: "maybe_12", Line: 1}),
		},

		{
			code: "random_value;",
			want: ast.NewVariable(token.Token{Type: token.IDENTIFIER, Lexeme: "random_value", Line: 1}),
		},
	}

	for _, test := range tests {
		lxr := lexer.New(test.code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := New(tokens)
		program, err := prsr.Parse()
		if err != nil {
			t.Fatalf("code:'%v'\tParsing errors caught: %v", test.code, err.Error())
		}
		if len(program) != 1 {
			t.Fatalf("code:'%v'\tprogram has wrong number of statements. expected=%d got=%d", test.code, 1, len(program))
		}
		smt, isOk := program[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("code:'%v'\tprogram[0] is not *ast.ExpressionStmt. got=%T", test.code, program[0])
		}

		exp, isOk := smt.Exp.(*ast.Variable)
		if !isOk {
			t.Fatalf("code:'%v'\t smt.Exp is not *ast.Variable. got=%T", test.code, smt.Exp)
		}
		if passed := testVariable(exp, test.want, t); !passed {
			t.Errorf("testVariable failed for '%s'", test.code)
			t.FailNow()
		}
	}
}

func ParseAssignment(t *testing.T) {
	tests := []struct {
		code string
		want *ast.Assignment
	}{
		{
			code: "maybe_12 = 12;",
			want: ast.NewAssignment(
				token.Token{Type: token.IDENTIFIER, Lexeme: "maybe_12", Line: 1},
				ast.NewLiteralExpression(12),
			),
		},
		{
			code: `the_name = "John";`,
			want: ast.NewAssignment(
				token.Token{Type: token.IDENTIFIER, Lexeme: "the_name", Line: 1},
				ast.NewLiteralExpression("John"),
			),
		},
		{
			code: `can_drink = 12 >= 18? true: false;`,
			want: ast.NewAssignment(
				token.Token{Type: token.IDENTIFIER, Lexeme: "can_drink", Line: 1},
				ast.NewTernaryConditional(
					ast.NewBinaryExpression(
						ast.NewLiteralExpression(12),
						token.Token{Type: token.GREATER_EQ, Lexeme: ">=", Line: 1},
						ast.NewLiteralExpression(18),
					),
					token.Token{Type: token.QUESTION_MARK, Lexeme: "?", Line: 1},
					ast.NewLiteralExpression(true),
					token.Token{Type: token.COLON, Lexeme: ":", Line: 1},
					ast.NewLiteralExpression(false),
				),
			),
		},
		{
			code: `negative = -1;`,
			want: ast.NewAssignment(
				token.Token{Type: token.IDENTIFIER, Lexeme: "negative", Line: 1},
				ast.NewUnaryExpression(
					token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
					ast.NewLiteralExpression(1),
				),
			),
		},
		{
			code: `twelve = 6*2;`,
			want: ast.NewAssignment(
				token.Token{Type: token.IDENTIFIER, Lexeme: "twelve", Line: 1},
				ast.NewBinaryExpression(
					ast.NewLiteralExpression(6),
					token.Token{Type: token.ASTERISK, Lexeme: "*", Line: 1},
					ast.NewLiteralExpression(2),
				),
			),
		},
		{
			code: `var_reference = twelve;`,
			want: ast.NewAssignment(
				token.Token{Type: token.IDENTIFIER, Lexeme: "var_reference", Line: 1},
				ast.NewVariable(token.Token{Type: token.IDENTIFIER, Lexeme: "twelve", Line: 1}),
			),
		},
	}

	for _, test := range tests {
		code := test.code
		lxr := lexer.New(code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("code='%s'\tScanning failed with exception='%v'", code, err.Error())
		}
		prsr := New(tokens)
		program, err := prsr.Parse()
		if err != nil {
			t.Fatalf("code='%s'\tParsing errors caught: %v", code, err.Error())
		}
		if len(program) != 1 {
			t.Fatalf("code='%s'\tprogram has wrong number of statements. expected=%d got=%d", code, 1, len(program))
		}
		smt, isOk := program[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("code='%s'\tprogram[0] is not *ast.ExpressionStmt. got=%T", code, program[0])
		}
		expr, isOk := smt.Exp.(*ast.Assignment)
		if !isOk {
			t.Fatalf("code='%s'\t smt.Exp is not *ast.Assignment. got=%T", code, smt.Exp)
		}

		if passed := testAssignment(expr, test.want, t); !passed {
			t.Errorf("testAssignment failed for '%s'", code)
			t.FailNow()
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

func testVariable(exp ast.Expression, expected *ast.Variable, t *testing.T) bool {
	variable, isOk := exp.(*ast.Variable)
	if !isOk {
		t.Errorf("exp is not a *ast.Variable. got='%T'", exp)
		return false
	}

	if variable.Name.Lexeme != expected.Name.Lexeme {
		want, got := expected.Name.Lexeme, variable.Name.Lexeme
		t.Errorf("wrong value for variable.Name. expected='%v' got='%v'", want, got)
		return false
	}
	return true
}

func testAssignment(exp ast.Expression, expected *ast.Assignment, t *testing.T) bool {
	assign, isOk := exp.(*ast.Assignment)
	if !isOk {
		t.Errorf("exp is not a *ast.Assignment. got='%T'", exp)
		return false
	}

	if assign.Name.Lexeme != expected.Name.Lexeme {
		want, got := expected.Name.Lexeme, assign.Name.Lexeme
		t.Errorf("wrong value for assign.Name. expected='%v' got='%v'", want, got)
		return false
	}

	if assign.Value.String() != expected.Value.String() {
		want, got := expected.Value, assign.Value
		t.Errorf("wrong value for assign.Value. expected='%v' got='%v'", want, got)
		return false
	}

	return true
}
