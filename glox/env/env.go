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

func (env *Environment) Assign(name token.Token, value any) error {
	if _, isOk := env.values[name.Lexeme]; !isOk {
		return exception.Runtime(name, "undefined variable '"+name.Lexeme+"'.")
	}
	env.values[name.Lexeme] = value
	return nil
}
