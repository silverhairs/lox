package interpreter

type Callable interface {
	Call(i *Interpreter, arguments []any) any
	Arity() int
	String() string
}
