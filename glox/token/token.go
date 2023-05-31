package token

import "fmt"

type TokenType string

const (
	// Single-character tokens
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"
	COMMA       = ","
	DOT         = "."
	MINUS       = "-"
	PLUS        = "+"
	SEMICOLON   = ";"
	SLASH       = "/"
	STAR        = "*"

	// One or two characters tokens
	BANG          = "!"
	BANG_EQUAL    = "!="
	EQUAL         = "="
	EQUAL_EQUAL   = "=="
	GREATER       = ">"
	GREATER_EQUAL = ">="
	LESS          = "<"
	LESS_EQUAL    = "<="

	// Literals
	IDENTIFIER = "IDENT"
	STRING     = "STR"
	NUMBER     = "NUM"

	// Keywords
	AND      = "and"
	CLASS    = "class"
	ELSE     = "else"
	IF       = "if"
	FALSE    = "false"
	FUNCTION = "fn"
	FOR      = "for"
	OR       = "or"
	NIL      = "nil"
	PRINT    = "print"
	RETURN   = "return"
	SUPER    = "super"
	THIS     = "this"
	TRUE     = "true"
	LET      = "let"
	WHILE    = "while"

	EOF = ""
)

type Token struct {
	Type    TokenType
	Lexeme  any
	Literal string
	Line    int
}

func New(Type TokenType, Literal string, Lexeme any, Line int) *Token {
	return &Token{
		Literal: Literal,
		Type:    Type,
		Lexeme:  Lexeme,
		Line:    Line,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %+v %s", t.Type, t.Lexeme, t.Literal)
}
