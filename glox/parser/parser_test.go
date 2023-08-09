package parser

import (
	"glox/ast"
	"glox/lexer"
	"glox/token"
	"strings"
	"testing"
)

func TestParseTernary(t *testing.T) {

	tests := []struct {
		code string
		want *ast.Ternary
	}{
		{
			code: ` 15 > 1 ? "abc" : "123";`,
			want: ast.NewTernaryConditional(
				ast.NewBinaryExpression(
					ast.NewLiteralExpression(15),
					token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
					ast.NewLiteralExpression(1),
				),
				token.Token{Type: token.QUESTION_MARK, Lexeme: "?", Line: 1},
				ast.NewLiteralExpression("abc"),
				token.Token{Type: token.COLON, Lexeme: ":", Line: 1},
				ast.NewLiteralExpression("123"),
			),
		},
		{
			code: ` false ? "abc" : 123;`,
			want: ast.NewTernaryConditional(
				ast.NewLiteralExpression(false),
				token.Token{Type: token.QUESTION_MARK, Lexeme: "?", Line: 1},
				ast.NewLiteralExpression("abc"),
				token.Token{Type: token.COLON, Lexeme: ":", Line: 1},
				ast.NewLiteralExpression(123),
			),
		},
	}

	for _, test := range tests {
		lxr := lexer.New(test.code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("`%s` -> failed to tokenize got error='%s'", test.code, err.Error())
		}
		prsr := New(tokens)

		program, err := prsr.Parse()
		if err != nil {
			t.Fatalf("`%s` -> failed to parse got error='%s'", test.code, err.Error())
		}

		if len(program) != 1 {
			t.Fatalf("program has wrong number of statements. want=%d got=%d", 1, len(program))
		}

		stmt, isOk := program[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("program[0] is not *ast.ExpressionStmt. got=%T", program[0])
		}

		expr := stmt.Exp

		if !testTernary(expr, test.want, t) {
			t.Logf("failed on code=`%s`", test.code)
			t.Fail()
		}
	}
}

func TestParseLiteral(t *testing.T) {
	tests := []struct {
		code string
		want ast.Literal
	}{
		{
			code: `"john doe";`,
			want: *ast.NewLiteralExpression("john doe"),
		},
		{
			code: `5;`,
			want: *ast.NewLiteralExpression(5),
		},
		{
			code: `5.9797;`,
			want: *ast.NewLiteralExpression(5.9797),
		},
		{
			code: `false;`,
			want: *ast.NewLiteralExpression(false),
		},
		{
			code: `true;`,
			want: *ast.NewLiteralExpression(true),
		},
		{
			code: `nil;`,
			want: *ast.NewLiteralExpression(nil),
		},
	}

	for _, test := range tests {
		lxr := lexer.New(test.code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("failed to tokenize code='%s' got error='%s'", test.code, err.Error())
		}
		prsr := New(tokens)
		stmts, err := prsr.Parse()
		if err != nil {
			t.Fatalf("failed to parse code='%s' got error='%s'", test.code, err.Error())
		}
		if len(stmts) != 1 {
			t.Fatalf("parsed into wrong number of statements. want=1 got=%d", len(tokens))
		}
		stmt, isOk := stmts[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("`%v` -> stmts[0] is not a *ast.ExpressionStmt. got=%T", test.code, stmts[0])
		}

		if !testLiteral(stmt.Exp, test.want, t) {
			t.Logf("failed on code=`%v`", test.code)
			t.Fail()
		}

	}
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
			t.Fatalf("program has wrong number of statements. want=%d got=%d", 1, len(program))
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
			t.Fatalf("program has wrong number of statements. want=%d got=%d", 1, len(program))
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
			t.Fatalf("program has wrong number of statements. want=%d got=%d", 1, len(program))
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
			t.Fatalf("wrong group.Exp want='%v'. got='%v'", test.exp, group.Exp)
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
		t.Fatalf("exception message wrong. want to contain='%s' but got='%s'", chunk, got)
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
			t.Fatalf("code:'%v'\tprogram has wrong number of statements. want=%d got=%d", test.code, 1, len(program))
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

func TestParseAssignment(t *testing.T) {
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
			t.Fatalf("code='%s'\tprogram has wrong number of statements. want=%d got=%d", code, 1, len(program))
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

func TestParseStatement(t *testing.T) {
	tests := []struct {
		code string
		want ast.Statement
	}{
		{
			code: "var maybe_12 = 12;",
			want: ast.NewLetStmt(
				token.Token{Type: token.IDENTIFIER, Lexeme: "maybe_12", Line: 1},
				ast.NewLiteralExpression(12),
			),
		},
		{
			code: `print "John";`,
			want: ast.NewPrintStmt(
				ast.NewLiteralExpression("John"),
			),
		},
		{
			code: `12;`,
			want: ast.NewExprStmt(
				ast.NewLiteralExpression(12),
			),
		},
		{
			code: `{ 12; }`,
			want: ast.NewBlockStmt(
				[]ast.Statement{
					ast.NewExprStmt(ast.NewLiteralExpression(12)),
				},
			),
		},
		{
			code: `if (12 > 10) { print "yes"; }`,
			want: ast.NewIfStmt(
				ast.NewBinaryExpression(
					ast.NewLiteralExpression(12),
					token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
					ast.NewLiteralExpression(10),
				),
				ast.NewBlockStmt(
					[]ast.Statement{
						ast.NewPrintStmt(ast.NewLiteralExpression("yes")),
					},
				),
				nil,
			),
		},
		{
			code: `if (12 > 10) { print "yes"; } else { print "no"; }`,
			want: ast.NewIfStmt(
				ast.NewBinaryExpression(
					ast.NewLiteralExpression(12),
					token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
					ast.NewLiteralExpression(10),
				),
				ast.NewBlockStmt(
					[]ast.Statement{
						ast.NewPrintStmt(ast.NewLiteralExpression("yes")),
					},
				),
				ast.NewBlockStmt(
					[]ast.Statement{
						ast.NewPrintStmt(ast.NewLiteralExpression("no")),
					},
				),
			),
		},
		{
			code: `for(let i = 0; i < 10; i = i + 1) print i;`,
			want: ast.NewBlockStmt(
				[]ast.Statement{
					ast.NewLetStmt(
						token.Token{Type: token.IDENTIFIER, Lexeme: "i", Line: 1},
						ast.NewLiteralExpression(0),
					),
					ast.NewWhileStmt(
						ast.NewBinaryExpression(
							ast.NewVariable(token.Token{Type: token.IDENTIFIER, Lexeme: "i", Line: 1}),
							token.Token{Type: token.LESS, Lexeme: "<", Line: 1},
							ast.NewLiteralExpression(10),
						),
						ast.NewBlockStmt(
							[]ast.Statement{
								ast.NewPrintStmt(ast.NewVariable(token.Token{Type: token.IDENTIFIER, Lexeme: "i", Line: 1})),
								ast.NewExprStmt(
									ast.NewAssignment(
										token.Token{Type: token.IDENTIFIER, Lexeme: "i", Line: 1},
										ast.NewBinaryExpression(
											ast.NewVariable(token.Token{Type: token.IDENTIFIER, Lexeme: "i", Line: 1}),
											token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
											ast.NewLiteralExpression(1),
										),
									),
								),
							},
						),
					),
				},
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
			t.Fatalf("code='%s'\tprogram has wrong number of statements. want=%d got=%d", code, 1, len(program))
		}

		stmt := program[0]

		if passed := testStmt(stmt, test.want, t); !passed {
			t.Errorf("testStmt failed for '%s'", code)
			t.FailNow()
		}
	}
}

func TestParseLogical(t *testing.T) {
	tests := []struct {
		code string
		want *ast.Logical
	}{
		{
			code: `true or false;`,
			want: ast.NewLogical(
				ast.NewLiteralExpression(true),
				token.Token{Type: token.OR, Lexeme: "or", Line: 1},
				ast.NewLiteralExpression(false),
			),
		},
		{
			code: `true and false;`,
			want: ast.NewLogical(
				ast.NewLiteralExpression(true),
				token.Token{Type: token.AND, Lexeme: "and", Line: 1},
				ast.NewLiteralExpression(false),
			),
		},
		{
			code: `true or false and true;`,
			want: ast.NewLogical(
				ast.NewLiteralExpression(true),
				token.Token{Type: token.OR, Lexeme: "or", Line: 1},
				ast.NewLogical(
					ast.NewLiteralExpression(false),
					token.Token{Type: token.AND, Lexeme: "and", Line: 1},
					ast.NewLiteralExpression(true),
				),
			),
		},
	}
	for _, test := range tests {
		code := test.code
		lxr := lexer.New(code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("failed to tokenize code `%s`", code)
		}
		prsr := New(tokens)
		stmts, err := prsr.Parse()
		if err != nil {
			t.Fatalf("failed to parse code `%s`", code)
		}
		if len(stmts) != 1 {
			t.Fatalf("wrong number of statements. want=1 got=%d", len(stmts))
		}
		stmt, isOk := stmts[0].(*ast.ExpressionStmt)
		if !isOk {
			t.Fatalf("stmts[0] is not a *ast.ExpressionStmt. got=%T", stmts[0])
		}
		if isOk := testLogical(stmt.Exp, test.want, t); !isOk {
			t.Errorf("testLogical failed for '%s'", code)
			t.FailNow()
		}
	}

}

func TestParseWhile(t *testing.T) {
	tests := []struct {
		code string
		want *ast.WhileStmt
	}{
		{
			code: `while(false){print "yes";}`,
			want: ast.NewWhileStmt(
				ast.NewLiteralExpression(false),
				ast.NewBlockStmt(
					[]ast.Statement{
						ast.NewPrintStmt(ast.NewLiteralExpression("yes")),
					},
				),
			),
		},
		{
			code: `while(true)print "yes";`,
			want: ast.NewWhileStmt(
				ast.NewLiteralExpression(true),
				ast.NewPrintStmt(ast.NewLiteralExpression("yes")),
			),
		},
		{
			code: `while(0>1 and 1==1){print "yes";}`,
			want: ast.NewWhileStmt(
				ast.NewLogical(
					ast.NewBinaryExpression(
						ast.NewLiteralExpression(0),
						token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
						ast.NewLiteralExpression(1),
					),
					token.Token{Type: token.AND, Lexeme: "and", Line: 1},
					ast.NewBinaryExpression(
						ast.NewLiteralExpression(1),
						token.Token{Type: token.EQ_EQ, Lexeme: "==", Line: 1},
						ast.NewLiteralExpression(1),
					),
				),
				ast.NewBlockStmt(
					[]ast.Statement{
						ast.NewPrintStmt(ast.NewLiteralExpression("yes")),
					},
				),
			),
		},
		{
			code: `while(0>1 or false)print "yes";`,
			want: ast.NewWhileStmt(
				ast.NewLogical(
					ast.NewBinaryExpression(
						ast.NewLiteralExpression(0),
						token.Token{Type: token.GREATER, Lexeme: ">", Line: 1},
						ast.NewLiteralExpression(1),
					),
					token.Token{Type: token.OR, Lexeme: "or", Line: 1},
					ast.NewLiteralExpression(false),
				),
				ast.NewPrintStmt(ast.NewLiteralExpression("yes")),
			),
		},
	}

	for _, test := range tests {
		lxr := lexer.New(test.code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("failed to tokenize code `%s`", test.code)
		}

		prsr := New(tokens)
		stmts, err := prsr.Parse()
		if err != nil {
			t.Fatalf("failed to parse code `%s`", test.code)
		}
		if len(stmts) != 1 {
			t.Fatalf("wrong number of statements. want=1 got=%d", len(stmts))
		}
		if !testWhile(stmts[0], test.want, t) {
			t.Errorf("testWhile failed for '%s'", test.code)
			t.FailNow()
		}
	}
}

func testLiteral(exp ast.Expression, wantValue any, t *testing.T) bool {
	isLiteral, literal := assertLiteral(exp, ast.NewLiteralExpression(wantValue))
	if !isLiteral {
		t.Errorf("result.Left is not *ast.Literal. got=%T", exp)
		return false
	}
	if literal.Value != wantValue {
		t.Errorf("literal.Value wrong. want=5 got=%q", literal.Value)
		return false
	}
	return true
}

func assertLiteral(exp ast.Expression, want *ast.Literal) (bool, *ast.Literal) {
	literal, isLiteral := exp.(*ast.Literal)
	if !isLiteral || literal.String() != exp.String() {
		return false, nil
	}
	return true, want
}

func testUnary(exp ast.Expression, want *ast.Unary, t *testing.T) bool {
	unary, isOk := exp.(*ast.Unary)
	if !isOk {
		t.Errorf("exp is not a *ast.Unary. got='%T'", exp)
		return false
	}

	if unary.Operator.Lexeme != want.Operator.Lexeme {
		want, got := want.Operator, unary.Operator
		t.Errorf("wrong value for unary.Operator. want='%v' got='%v'", want, got)
		return false
	}

	return testExpression(unary.Right, want.Right, t)
}

func testBinary(exp ast.Expression, want *ast.Binary, t *testing.T) bool {
	binary, isOk := exp.(*ast.Binary)
	if !isOk {
		t.Errorf("exp is not a *ast.Binary. got='%T'", exp)
		return false
	}

	return testExpression(binary.Left, want.Left, t) && testExpression(binary.Right, want.Right, t)
}

func testGrouping(got ast.Expression, want *ast.Grouping, t *testing.T) bool {
	grouping, isOk := got.(*ast.Grouping)
	if !isOk {
		t.Errorf("want=*ast.Grouping. got='%T'", got)
		return false
	}
	return testExpression(grouping.Exp, want.Exp, t)

}

func testTernary(exp ast.Expression, want *ast.Ternary, t *testing.T) bool {
	ternary, isOk := exp.(*ast.Ternary)
	if !isOk {
		t.Errorf("exp is not a *ast.Ternary. got='%T'", exp)
		return false
	}
	return testExpression(ternary.Condition, want.Condition, t) && testExpression(ternary.Then, want.Then, t) && testExpression(ternary.OrElse, want.OrElse, t)
}

func testVariable(exp ast.Expression, want *ast.Variable, t *testing.T) bool {
	variable, isOk := exp.(*ast.Variable)
	if !isOk {
		t.Errorf("exp is not a *ast.Variable. got='%T'", exp)
		return false
	}

	if variable.Name.Lexeme != want.Name.Lexeme {
		want, got := want.Name.Lexeme, variable.Name.Lexeme
		t.Errorf("wrong value for variable.Name. want='%v' got='%v'", want, got)
		return false
	}
	return true
}

func testAssignment(exp ast.Expression, want *ast.Assignment, t *testing.T) bool {
	assign, isOk := exp.(*ast.Assignment)
	if !isOk {
		t.Errorf("exp is not a *ast.Assignment. got='%T'", exp)
		return false
	}

	if assign.Name.Lexeme != want.Name.Lexeme {
		want, got := want.Name.Lexeme, assign.Name.Lexeme
		t.Errorf("wrong value for assign.Name. want='%v' got='%v'", want, got)
		return false
	}

	return testExpression(assign.Value, want.Value, t)
}

func testLogical(stmt ast.Expression, want *ast.Logical, t *testing.T) bool {
	logical, isOk := stmt.(*ast.Logical)
	if !isOk {
		t.Errorf("stmt is not a *ast.Logical. got='%T'", stmt)
		return false
	}

	if logical.Operator.Lexeme != want.Operator.Lexeme {
		want, got := want.Operator, logical.Operator
		t.Errorf("wrong value for logical.Operator. want='%v' got='%v'", want, got)
		return false
	}

	return testExpression(logical.Left, want.Left, t) && testExpression(logical.Right, want.Right, t)
}

func testExpression(got ast.Expression, want ast.Expression, t *testing.T) bool {
	if want == nil {
		if got != nil {
			t.Errorf("wrong value for got. want='nil' got='%v'", got)
			return false
		}
		return true
	}

	switch want := want.(type) {
	case *ast.Unary:
		return testUnary(got, want, t)
	case *ast.Literal:
		return testLiteral(got, want, t)
	case *ast.Grouping:
		return testGrouping(got, want, t)
	case *ast.Binary:
		return testBinary(got, want, t)
	case *ast.Ternary:
		return testTernary(got, want, t)
	case *ast.Variable:
		return testVariable(got, want, t)
	case *ast.Assignment:
		return testAssignment(got, want, t)
	case *ast.Logical:
		return testLogical(got, want, t)
	default:
		t.Errorf("expression %T does not have a testing function. consider adding one", want)
		return false
	}
}

func testLetStmt(stmt ast.Statement, want *ast.LetStmt, t *testing.T) bool {
	let, isOk := stmt.(*ast.LetStmt)
	if !isOk {
		t.Errorf("stmt is not a *ast.LetStmt. got='%T'", stmt)
		return false
	}

	if let.Name.Lexeme != want.Name.Lexeme {
		want, got := want.Name.Lexeme, let.Name.Lexeme
		t.Errorf("wrong value for let.Name. want='%s' got='%s'", want, got)
		return false
	}

	if let.Value.String() != want.Value.String() {
		want, got := want.Value, let.Value
		t.Errorf("wrong value for let.Value. want='%v' got='%v'", want, got)
		return false
	}

	return true
}

func testPrintStmt(stmt ast.Statement, want *ast.PrintStmt, t *testing.T) bool {
	print, isOk := stmt.(*ast.PrintStmt)
	if !isOk {
		t.Errorf("stmt is not a *ast.PrintStmt. got='%T'", stmt)
		return false
	}

	if print.Exp.String() != want.Exp.String() {
		want, got := want.Exp, print.Exp
		t.Errorf("wrong value for print.Exp. want='%v' got='%v'", want, got)
		return false
	}

	return true
}

func testExprStmt(stmt ast.Statement, want *ast.ExpressionStmt, t *testing.T) bool {
	expr, isOk := stmt.(*ast.ExpressionStmt)
	if !isOk {
		t.Errorf("stmt is not a *ast.ExpressionStmt. got='%T'", stmt)
		return false
	}

	if expr.Exp.String() != want.Exp.String() {
		want, got := want.Exp, expr.Exp
		t.Errorf("wrong value for expr.Exp. want='%v' got='%v'", want, got)
		return false
	}

	return true
}

func testBlockStmt(stmt ast.Statement, want *ast.BlockStmt, t *testing.T) bool {
	block, isOk := stmt.(*ast.BlockStmt)
	if !isOk {
		t.Errorf("stmt is not a *ast.BlockStmt. got='%T'", stmt)
		return false
	}

	if len(block.Stmts) != len(want.Stmts) {
		want, got := len(want.Stmts), len(block.Stmts)
		t.Errorf("wrong number of statements in block. want='%d' got='%d'", want, got)
		return false
	}

	return true
}

func testIfStmt(stmt ast.Statement, want *ast.IfStmt, t *testing.T) bool {
	ifStmt, isOk := stmt.(*ast.IfStmt)
	if !isOk {
		t.Errorf("stmt is not a *ast.IfStmt. got='%T'", stmt)
		return false
	}

	if ifStmt.Condition.String() != want.Condition.String() {
		want, got := want.Condition, ifStmt.Condition
		t.Errorf("wrong value for ifStmt.Condition. want='%v' got='%v'", want, got)
		return false
	}

	return testStmt(ifStmt.Then, want.Then, t) && testStmt(ifStmt.OrElse, want.OrElse, t)

}

func testStmt(stmt ast.Statement, want ast.Statement, t *testing.T) bool {
	if want == nil {
		if stmt != nil {
			t.Errorf("wrong value for stmt. want='nil' got='%v'", stmt)
			return false
		}
		return true
	}
	switch want := want.(type) {
	case *ast.LetStmt:
		return testLetStmt(stmt, want, t)
	case *ast.PrintStmt:
		return testPrintStmt(stmt, want, t)
	case *ast.ExpressionStmt:
		return testExprStmt(stmt, want, t)
	case *ast.BlockStmt:
		return testBlockStmt(stmt, want, t)
	case *ast.IfStmt:
		return testIfStmt(stmt, want, t)
	case *ast.WhileStmt:
		return testWhile(stmt, want, t)
	default:
		t.Errorf("statement %T does not have a testing function. consider adding one", want)
		return false
	}

}

func testWhile(got ast.Statement, want *ast.WhileStmt, t *testing.T) bool {
	while, isOk := got.(*ast.WhileStmt)
	if !isOk {
		t.Errorf("got='%T' want a *ast.WhileStmt.", got)
		return false
	}
	if while.Condition.String() != want.Condition.String() {
		t.Errorf("got wrong conditional expression. want='%v' got='%v'", want.Condition.String(), while.Condition.String())
		return false
	}

	return testStmt(while.Body, want.Body, t)
}
