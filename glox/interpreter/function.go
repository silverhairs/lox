package interpreter

import (
	"fmt"
	"glox/ast"
	"glox/env"
)

type LoxFunction struct {
	declaration *ast.Function
}

func NewFunction(declaration *ast.Function) *LoxFunction {
	return &LoxFunction{declaration: declaration}
}

func (fn *LoxFunction) Call(i *Interpreter, args []any) (value any) {
	defer func() {
		if r := recover(); r != nil {
			if rtrn, isReturn := r.(errReturn); !isReturn {
				panic(r)
			} else {
				value = rtrn.value
				return
			}
		}
	}()
	env := env.New(i.Env)
	for i, param := range fn.declaration.Params {
		arg := args[i]
		env.Define(param.Lexeme, arg)
	}
	i.executeBlock(fn.declaration.Body, env)
	return value
}

func (fn *LoxFunction) Arity() int {
	return len(fn.declaration.Params)
}

func (fn *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", fn.declaration.Name.Lexeme)
}
