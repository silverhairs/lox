package interpreter

import (
	"fmt"
	"glox/ast"
	"glox/env"
	"glox/object"
	"glox/token"
)

type LoxFunction struct {
	object.Callable[*Interpreter]
	declaration *ast.Function
}

func NewFunction(declaration *ast.Function) *LoxFunction {
	return &LoxFunction{declaration: declaration}
}

func (fn *LoxFunction) Call(i *Interpreter, args []token.Token) any {
	env := env.New(i.Env)
	for i, param := range fn.declaration.Params {
		arg := args[i]
		env.Define(param.Lexeme, arg)
	}
	i.executeBlock(fn.declaration.Body, env)
	return nil
}

func (fn *LoxFunction) Arity() int {
	return len(fn.declaration.Params)
}

func (fn *LoxFunction) String() string {
	return fmt.Sprintf("<fn '%s'>", fn.declaration.Name.Lexeme)
}
