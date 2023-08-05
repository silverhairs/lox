package token

import "testing"

func TestLookupIdentifier(t *testing.T) {
	tests := []struct {
		word     string
		expected TokenType
	}{
		{word: "nil", expected: NIL},
		{word: "class", expected: CLASS},
		{word: "if", expected: IF},
		{word: "maybe12", expected: IDENTIFIER},
		{word: "else", expected: ELSE},
		{word: "and", expected: AND},
		{word: "or", expected: OR},
		{word: "return", expected: RETURN},
		{word: "true", expected: TRUE},
		{word: "false", expected: FALSE},
		{word: "print", expected: PRINT},
		{word: "fn", expected: FUNCTION},
		{word: "while", expected: WHILE},
		{word: "for", expected: FOR},
		{word: "this", expected: THIS},
		{word: "super", expected: SUPER},
		{word: "let", expected: LET},
		{word: "var", expected: LET},
		{word: "func", expected: IDENTIFIER},
		{word: "struct", expected: IDENTIFIER},
		{word: "interface", expected: IDENTIFIER},
		{word: "type", expected: IDENTIFIER},
	}

	for _, test := range tests {
		got := LookupIdentifier(test.word)
		expected := test.expected

		if got != expected {
			t.Fatalf("failed to get token type for '%s'. got='%v' expected='%v'", test.word, got, expected)
		}
	}
}
