package env

import (
	"glox/exception"
	"glox/token"
)

type Environment struct {
	values map[string]any
}

func New() *Environment {
	return &Environment{values: make(map[string]any)}
}

func (env *Environment) Define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) Get(name token.Token) any {
	if val, isOk := env.values[name.Lexeme]; isOk {
		return val
	}
	return exception.Runtime(name, "undefined variable '"+name.Lexeme+"'.")
}
