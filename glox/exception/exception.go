package exception

import (
	"fmt"
)

func Generic(line int, where string, msg string) error {
	message := fmt.Sprintf("%s @ %s", msg, where)
	return fmt.Errorf("GloxGenericError([line %d] Error: %s)", line, message)
}

// Calls `e.New` with an empty string for the `where` argument.
func Short(line int, msg string) error {
	return Generic(line, "", msg)
}

func RuntimeError(body any, message string) error {
	return fmt.Errorf("RuntimeError(%+v, %s)", body, message)
}
