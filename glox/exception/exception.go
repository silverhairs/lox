package exception

import (
	"fmt"
	"glox/token"
)

const (
	PARSE_EXCEPTION   = "ParseException"
	RUNTIME_EXCEPTION = "RuntimeException"
	GENERIC_EXCEPTION = "GenericException"
)

func Generic(line int, where string, msg string) error {
	out := fmt.Sprintf("%s(%s at %s)", GENERIC_EXCEPTION, msg, where)
	return fmt.Errorf("unhandled exception: %s\n[line %d]", out, line)
}

// Calls `Generic` with an empty string for the `where` argument.
func Short(line int, msg string) error {
	return Generic(line, "", msg)
}

func Runtime(token token.Token, message string) error {
	return fmt.Errorf("unhandled exception: %s(%q, %s)\n[line: %d]", RUNTIME_EXCEPTION, token.Lexeme, message, token.Line)
}

func Parse(tok token.Token) error {
	return fmt.Errorf("unhandled exception: %s(%q, illegal token)\n[line: %d]", PARSE_EXCEPTION, tok.Lexeme, tok.Line)
}
