package ast

import (
	"glox/token"
	"testing"
)

func TestPrint(t *testing.T) {

	input := []Expression{
		&Binary{
			Left:     &Literal{Value: "5"},
			Operator: token.Token{Type: token.ASTERISK, Lexeme: "*", Line: 1},
			Right:    &Unary{Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1}, Right: &Literal{Value: 10}},
		},
		&Unary{
			Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
			Right:    &Literal{Value: 123},
		},
		&Grouping{
			Exp: &Binary{
				Left:     &Literal{Value: 1},
				Operator: token.Token{Type: token.PLUS, Lexeme: "+", Line: 1},
				Right:    &Literal{Value: 2},
			},
		},
		&Ternary{
			Condition:      &Literal{Value: true},
			ThenOperator:   token.Token{Type: token.QUESTION_MARK, Lexeme: "?", Line: 1},
			Then:           &Literal{Value: 1},
			OrElseOperator: token.Token{Type: token.COLON, Lexeme: ":", Line: 1},
			OrElse:         &Literal{Value: 2},
		},
		&Literal{Value: "yes"},
		&Variable{Name: token.Token{Type: token.IDENTIFIER, Lexeme: "number", Line: 1}},
		&Assignment{
			Name:  token.Token{Type: token.IDENTIFIER, Lexeme: "number", Line: 1},
			Value: NewLiteralExpression(12),
		},
	}
	printer := NewPrinter()

	for _, test := range input {

		got := printer.Print(test)
		if got != test.String() {
			t.Fatalf("got %q, want %q", got, test.String())
		}

	}

	exp := &fakeExp{}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("accepted expressions that does not return a string")
		}
	}()
	printer.Print(exp)
}

type fakeExp struct{}

func (exp *fakeExp) String() string {
	return "(fake expression)"
}

func (exp *fakeExp) Type() ExpType {
	var res ExpType = "fake"
	return res
}

func (exp *fakeExp) Accept(Visitor) any {
	return 12
}
