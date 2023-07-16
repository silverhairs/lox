package exception

import (
	"fmt"
	"glox/token"
)

func Generic(line int, where string, msg string) error {
	out := fmt.Sprintf("%s at %s", msg, where)
	return fmt.Errorf("GloxGenericException(%s)\n[line %d]", out, line)
}

// Calls `e.New` with an empty string for the `where` argument.
func Short(line int, msg string) error {
	return Generic(line, "", msg)
}

func Runtime(token token.Token, message string) error {
	return fmt.Errorf("RuntimeException(%q, %s)\n[line: %d]", token.Lexeme, message, token.Line)
}
