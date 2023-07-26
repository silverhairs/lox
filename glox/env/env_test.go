package env

import (
	"glox/token"
	"math/rand"
	"strings"
	"testing"
)

func TestDefine(t *testing.T) {
	global := Global()
	vars := map[string]any{
		"age":     25,
		"name":    `"boris"`,
		"isHuman": true,
	}

	for name, value := range vars {
		global.Define(name, value)
		if global.values[name] != value {
			t.Fatalf("failed to define a variable %q with a value. expected=%v got=%v", name, value, global.values[name])
		}
	}

	local := New(global)
	if local.enclosing != global {
		t.Fatalf("local env must have a reference to its enclosing env. expected=%v got=%v", global, local.enclosing)
	}

	expected := "bob"
	local.Define("name", expected)
	got := local.values["name"]
	if got != expected {
		t.Fatalf("failed to shadow variable in local scope. got='%v'. expected='%s'", got, expected)
	}
}

func TestGet(t *testing.T) {
	global := Global()
	vars := map[string]any{
		"age":     25,
		"name":    `"boris"`,
		"isHuman": true,
	}
	for name, value := range vars {
		global.Define(name, value)
		tok := token.Token{Type: token.IDENTIFIER, Lexeme: name, Literal: nil, Line: 1}
		got := global.Get(tok)

		if got != value {
			t.Fatalf("failed to get variable %q. expected=%v got=%v", name, value, got)
		}
	}
	undefined := token.Token{Type: token.IDENTIFIER, Lexeme: "undefined", Literal: nil, Line: 1}
	got := global.Get(undefined)

	if err, isErr := got.(error); isErr {
		if !(strings.Contains(err.Error(), "undefined variable") && strings.Contains(err.Error(), "RuntimeException")) {
			t.Fatalf("wrong error message. expected a RuntimeException with message 'undefined variable'. got=%q", err.Error())
		}
	} else {
		t.Fatalf("failed to capture 'undefined variable'exception. got=%v", got)
	}

	local := New(global)

	if local.enclosing != global {
		t.Fatalf("non-global scope must have a reference to their enclosing scope. got='%v' expected='%v'", local.enclosing, global)
	}

	for name := range vars {
		tok := token.Token{Type: token.IDENTIFIER, Lexeme: name, Literal: nil, Line: 1}
		got = local.Get(tok)
		expected := global.Get(tok)

		if expected != got {
			t.Fatalf("failed to get variable from enclosing environment in local scope. expected='%v'. got='%v'", expected, got)
		}
	}

	expected := rand.Int()
	local.Define(undefined.Lexeme, expected)
	got = local.Get(undefined)

	if expected != got {
		t.Fatalf("failed to get variable defined in local environment. got='%v' expected='%v'", got, expected)
	}

	gotGlobal := global.Get(undefined)
	if _, isErr := gotGlobal.(error); !isErr {
		t.Fatalf("variable defined in local scope should not be accessible in enclosing scope. got='%v'. expected='RuntimeException'", gotGlobal)
	}

}

func TestAssign(t *testing.T) {
	global := Global()
	global.Define("name", `"boris"`)
	tok := token.Token{Type: token.IDENTIFIER, Lexeme: "name", Literal: nil, Line: 1}

	err := global.Assign(tok, "anya")

	if err != nil {
		t.Fatalf("failed to assign new value to variable=%v. got=%s", "name", err.Error())
	}

	if global.values["name"] != "anya" {
		t.Fatalf("failed to assign new value to variable. expected=%v got=%v", "anya", global.values["name"])
	}
	undefined := token.Token{Type: token.IDENTIFIER, Lexeme: "nothing", Literal: nil, Line: 1}

	err = global.Assign(undefined, 12)
	if err == nil {
		t.Fatal("assigning an undefined variable should results on an error.")
	}

	msg := err.Error()
	if !(strings.Contains(msg, "undefined variable") && strings.Contains(msg, undefined.Lexeme)) {
		t.Fatalf("wrong error message for undefined variable. got='%s' expected contains=['%s', '%s']", msg, "undefined variable", undefined.Lexeme)
	}

	local := New(global)

	if local.enclosing != global {
		t.Fatalf("non-global scope must have a reference to their enclosing scope. got='%v' expected='%v'", local.enclosing, global)
	}

	expected := "bob"
	local.Assign(tok, expected)
	got := local.Get(tok)

	if expected != got {
		t.Fatalf("failed to assign value to global variable inside local scope. got='%v' expected='%v'", got, expected)
	}
}
