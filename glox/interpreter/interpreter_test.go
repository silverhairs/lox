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
		"print 1+1;":                        "2",
		"print !false;":                     "true",
		"print !true;":                      "false",
		fmt.Sprintf("print %v + %v;", x, y): fmt.Sprintf("%v", x+y),
		`print "the number is"+1;`:          "the number is1",
		`print 12+" is the number";`:        "12 is the number",
		`print "hello world";`:              "hello world",
		`print (12+5*76/2);`:                "202",
		"print 1==1.0000001;":               "false",
		`print "yes" != "Yes";`:             "true",
		`let x = 1>2? "1 is bigger":"2 is bigger"; print x;`: "2 is bigger",
		`true and false;`:  "false",
		`true or false;`:   "true",
		`12 and false;`:    "false",
		`12 or false;`:     "12",
		`!true and false;`: "false",
		`!false and true;`: "true",
		`!false or true;`:  "true",
		`false or 12;`:     "12",
		`false and 12;`:    "false",
		`!true and 12;`:    "false",
		`!false and 12;`:   "12",
		`true and 12;`:     "12",
		`let age = 12; if(age>0 and age<18){ print "minor"; }`:                                        "minor",
		`let age = 19; if(age>0 and age<18){ print "minor"; } else { print "adult"; }`:                "adult",
		`let age = 21; if(age>18 or age > 21){ print "can drink"; }`:                                  "can drink",
		`let age = 17; if(age>18 or age > 21){ print "can drink"; } else { print "can't drink"; }`:    "can't drink",
		`let count=0; while(count<1){count=count+1;}`:                                                 "1",
		`let count=0; while(count<5){count=count+1;}`:                                                 "1\n2\n3\n4\n5",
		`fun greets(name){print "Hello "+name+"!";}greets("John");`:                                   "Hello John!\n<nil>",
		`fun count(n) {if(n > 1) count(n-1); print n;} count(5);`:                                     "1\n<nil>\n2\n<nil>\n3\n<nil>\n4\n<nil>\n5\n<nil>",
		`fun add(x,y){ return x+y; } let five = add(2,3); print five;`:                                "5",
		`fun concat(base,suffix){return base+suffix;} let word = concat("humor", "ist"); print word;`: "humorist",
		`fun test(max){let x=0; while(true){ if(x==max){ return;} print x; x=x+1; }} test(5);`:        "0\n1\n1\n2\n2\n3\n3\n4\n4\n5\n<nil>",
	}

	for code, expected := range fixtures {

		lxr := lexer.New(code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := parser.New(tokens)
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
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := parser.New(tokens)

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
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("Scanning failed with exception='%v'", err.Error())
		}
		prsr := parser.New(tokens)

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

	errors := []struct {
		code     string
		patterns []string
	}{
		{
			code:     `while(num==10) print num;`,
			patterns: []string{"RuntimeException", "undefined variable 'num'"},
		},
	}

	for _, failure := range errors {
		lxr := lexer.New(failure.code)
		tokens, err := lxr.Tokenize()
		if err != nil {
			t.Fatalf("failed to tokenize code '%v'", failure.code)
		}
		stmts, err := parser.New(tokens).Parse()
		if err != nil {
			t.Fatalf("failed to parse code `%v`. got='%s'", failure.code, err.Error())
		}
		i := New(stderr, stdout)
		i.Interpret(stmts)
		got := stderr.String()
		for _, pattern := range failure.patterns {
			if !strings.Contains(got, pattern) {
				t.Fatalf("%v ->failed to catch error. expected='%v' got='%v'", failure.code, pattern, got)
			}
		}

		stderr.Reset()
		stdout.Reset()
	}

}
