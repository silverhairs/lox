package interpreter

import (
	"bytes"
	"fmt"
	"glox/lexer"
	"glox/parser"
	"glox/token"
	"math/rand"
	"strings"
	"testing"
)

func TestInterpret(t *testing.T) {
	stderr := bytes.NewBufferString("")
	stdout := bytes.NewBufferString("")
	x := rand.Float64()
	y := rand.Float64()
	fixtures := map[string]string{
		"1+1":                             "2",
		"!false":                          "true",
		"!true":                           "false",
		fmt.Sprintf("%v + %v", x, y):      fmt.Sprintf("%v", x+y),
		`"the number is"+1`:               "the number is1",
		`12+" is the number"`:             "12 is the number",
		`"hello world"`:                   "hello world",
		`(12+5*76/2)`:                     "202",
		"1==1.0000001":                    "false",
		`"yes" != "Yes"`:                  "true",
		`1>2?"1 is bigger":"2 is bigger"`: "2 is bigger",
	}

	for code, expected := range fixtures {

		lxr := lexer.New(fmt.Sprintf("print %s;", code))
		prsr := parser.New(lxr.Tokenize())
		intrprtr := New(stderr, stdout)

		if expr, err := prsr.Parse(); err != nil {
			t.Fatalf("failed to parse code %q. \ngot=%v \nexpected=%v", code, err.Error(), expected)
		} else {
			intrprtr.Interpret(expr)
			if stderr.String() != "" {
				t.Fatalf("failed to evaluate %q. expected=%v got=%v", code, expected, stderr.String())
			}
			actual := strings.TrimRight(stdout.String(), "\n")
			if actual != expected {
				t.Fatalf("failed to evaluate %q. expected=%q got=%q", code, expected, actual)
			}
		}
		if stderr.String() != "" && !strings.HasSuffix(stderr.String(), "\n") {
			t.Fatalf("stderr message must end with a new line")
		}
		if stdout.String() != "" && !strings.HasSuffix(stdout.String(), "\n") {
			t.Fatalf("stdout message must end with a new line")
		}
		stderr.Reset()
		stdout.Reset()
	}

	failures := map[string][]string{
		"1+false;": {
			"unsupported operands. This operation can only be performed with numbers and strings.",
			"+",
		},
		"1/0;":        {"division by zero", "/"},
		"1-false;":    {"Operator \"-\" only accepts number operands", "-"},
		"-\"yes\";":   {"Operator \"-\" only accepts number operands", "-"},
		"true*false;": {"Operator \"*\" only accepts number operands", "*"},
		"false>true;": {"Operator \">\" only accepts number operands", ">"},
		"false>=12;":  {"Operator \">=\" only accepts number operands", ">="},
		"true<=12;":   {"Operator \"<=\" only accepts number operands", "<="},
		"true<true;":  {"Operator \"<\" only accepts number operands", "<"},
	}

	for code, chunks := range failures {
		lxr := lexer.New(code)
		prsr := parser.New(lxr.Tokenize())

		if expr, err := prsr.Parse(); err != nil {
			t.Fatalf("failed to parse code %q", code)
		} else {
			intrprtr := New(stderr, stdout)
			intrprtr.Interpret(expr)

			if stderr.String() == "" {
				t.Fatalf("failed to capture exception for=%q. expected=%s, got=%s", code, chunks, stderr.String())
			} else {
				actual := stderr.String()
				errorMsgIsOk := true
				var notFound string
				for _, msg := range chunks {
					if !strings.Contains(actual, msg) {
						errorMsgIsOk = false
						notFound = msg
					}
				}
				if !errorMsgIsOk {
					t.Fatalf("%q => wrong error message. \nexpected contains=%q \ngot=%q", code, notFound, actual)
				}

			}
		}

		if stderr.String() != "" && !strings.HasSuffix(stderr.String(), "\n") {
			t.Fatalf("stderr message must end with a new line")
		}
		if stdout.String() != "" && !strings.HasSuffix(stdout.String(), "\n") {
			t.Fatalf("stdout message must end with a new line")
		}

		stdout.Reset()
		stderr.Reset()
	}

	vars := []struct {
		code  string
		name  string
		value any
	}{
		{code: `let number = 12;`, name: "number", value: 12},
		{code: `var seven = 7;`, name: "seven", value: 7},
		{code: `let is_boolean=true;`, name: "is_boolean", value: true},
		{code: `var name = "anya forger";`, name: "name", value: "anya forger"},
	}

	for _, variable := range vars {
		lxr := lexer.New(variable.code)
		prsr := parser.New(lxr.Tokenize())

		stmts, err := prsr.Parse()
		if err != nil {
			t.Fatalf("failed to parse code %q. \ngot=%v", variable.code, err.Error())
		}

		i := New(stderr, stdout)
		i.Interpret(stmts)

		if stderr.String() != "" {
			t.Fatalf("caught exception when evaluating code=%q. got=%v", variable.code, stderr.String())
		}

		tok := token.Token{Type: token.IDENTIFIER, Lexeme: variable.name, Literal: nil, Line: 1}
		got := i.Env.Get(tok)

		expected := variable.value
		if num, isOk := expected.(int); isOk {
			expected = float64(num)
		}

		if got != expected {
			t.Fatalf("failed to keep state of defined variable in code=%q. got='%v'\nexpected='%v'.", variable.code, got, expected)
		}

		stderr.Reset()
		stdout.Reset()
	}

}
