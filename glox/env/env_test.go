package env

import (
	"glox/token"
	"strings"
	"testing"
)

func TestDefine(t *testing.T) {
	e := New()
	vars := map[string]any{
		"age":     25,
		"name":    `"boris"`,
		"isHuman": true,
	}

	for name, value := range vars {
		e.Define(name, value)
		if e.values[name] != value {
			t.Fatalf("failed to define a variable %q with a value. expected=%v got=%v", name, value, e.values[name])
		}
	}
}

func TestGet(t *testing.T) {
	e := New()
	vars := map[string]any{
		"age":     25,
		"name":    `"boris"`,
		"isHuman": true,
	}
	for name, value := range vars {
		e.Define(name, value)
		tok := token.Token{Type: token.IDENTIFIER, Lexeme: name, Literal: nil, Line: 1}
		got := e.Get(tok)

		if got != value {
			t.Fatalf("failed to get variable %q. expected=%v got=%v", name, value, got)
		}
	}
	undefined := token.Token{Type: token.IDENTIFIER, Lexeme: "undefined", Literal: nil, Line: 1}
	got := e.Get(undefined)

	if err, isErr := got.(error); isErr {
		if !strings.Contains(err.Error(), "undefined variable") || !strings.Contains(err.Error(), "RuntimeException") {
			t.Fatalf("wrong error message. expected a RuntimeException with message 'undefined variable'. got=%q", err.Error())
		}
	} else {
		t.Fatalf("failed to capture 'undefined variable'exception. got=%v", got)
	}
}

func TestAssign(t *testing.T) {
	e := New()
	e.Define("name", `"boris"`)

	err := e.Assign(token.Token{Type: token.IDENTIFIER, Lexeme: "name", Literal: nil, Line: 1}, "anya")

	if err != nil {
		t.Fatalf("failed to assign new value to variable=%v. got=%s", "name", err.Error())
	}

	if e.values["name"] != "anya" {
		t.Fatalf("failed to assign new value to variable. expected=%v got=%v", "anya", e.values["name"])
	}

	undefined := token.Token{Type: token.IDENTIFIER, Lexeme: "nothing", Literal: nil, Line: 1}

	err = e.Assign(undefined, 12)
	if err == nil {
		t.Fatal("assigning an undefined variable should results on an error.")
	}

	msg := err.Error()
	if !strings.Contains(msg, "undefined variable") && !strings.Contains(msg, undefined.Lexeme) {
		t.Fatalf("wrong error message for undefined variable. got='%s' expected contains=['%s', '%s']", msg, "undefined variable", undefined.Lexeme)
	}

}
