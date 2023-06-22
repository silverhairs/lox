package lexer

import (
	"glox/token"
	"testing"
)

func TestTokenize(t *testing.T) {
	input := `
	let age = 12;
	5 + 10
	1 - 2
	5 > 0
	1 < 12
	2 >= 1
	1 <= 1
	class Example {}
	fn call_me()
	true != false
	// this is a comment
	this.name
	super.person
	`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "age"},
		{token.EQUAL, "="},
		{token.NUMBER, "12"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "5"},
		{token.PLUS, "+"},
		{token.NUMBER, "10"},
		{token.NUMBER, "1"},
		{token.MINUS, "-"},
		{token.NUMBER, "2"},
		{token.NUMBER, "5"},
		{token.GREATER, ">"},
		{token.NUMBER, "0"},
		{token.NUMBER, "1"},
		{token.LESS, "<"},
		{token.NUMBER, "12"},
		{token.NUMBER, "2"},
		{token.GREATER_EQ, ">="},
		{token.NUMBER, "1"},
		{token.NUMBER, "1"},
		{token.LESS_EQ, "<="},
		{token.NUMBER, "1"},
		{token.CLASS, "class"},
		{token.IDENTIFIER, "Example"},
		{token.L_BRACE, "{"},
		{token.R_BRACE, "}"},
		{token.FUNCTION, "fn"},
		{token.IDENTIFIER, "call_me"},
		{token.L_PAREN, "("},
		{token.R_PAREN, ")"},
		{token.TRUE, "true"},
		{token.BANG_EQ, "!="},
		{token.FALSE, "false"},
		{token.COMMENT_L, "// this is a comment"},
		{token.THIS, "this"},
		{token.DOT, "."},
		{token.IDENTIFIER, "name"},
		{token.SUPER, "super"},
		{token.DOT, "."},
		{token.IDENTIFIER, "person"},
		{token.EOF, ""},
	}

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) != len(tests) {
		t.Fatalf("Wrong number of tokens. expected: %d got %d", len(tests), len(tokens))
	}

	for i, tok := range tokens {
		testTok := tests[i]

		if tok.Type != testTok.expectedType {
			t.Fatalf("wrong token type at test %d. expected %q got %q", i, testTok.expectedType, tok.Type)
		}

		if tok.Lexeme != testTok.expectedLexeme {
			t.Fatalf("wrong lexeme at test %d. expected %q got %q", i, testTok.expectedLexeme, tok.Lexeme)
		}
	}

}
