package object

type Callable[T any] interface {
	Call(i T, arguments []any) any
	Arity() int
}
