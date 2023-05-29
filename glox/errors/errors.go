package errors

import "fmt"

type GloxError struct {
	Line    int32
	Message string
}

func Report(line int32, msg string) *GloxError {
	return &GloxError{
		Line:    line,
		Message: msg,
	}
}

func (e *GloxError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}
