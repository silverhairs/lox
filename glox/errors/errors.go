package errors

import "fmt"

type GloxError struct {
	Line    int32
	Message string
}

func Report(line int32, where string, msg string) *GloxError {
	return &GloxError{
		Line:    line,
		Message: fmt.Sprintf("%s @ %s", msg, where),
	}
}

func (e *GloxError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}
