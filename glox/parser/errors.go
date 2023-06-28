package parser

import "glox/token"

type parseError struct {
	tok     token.Token
	message string
}

func error(tok token.Token, message string) *parseError {
	return &parseError{tok, message}
}
