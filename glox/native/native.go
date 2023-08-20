package native

import "time"

const NATIVE_FN_STR = "<native fn>"

type native[T any] struct {
	call     func(i T, argumets []any) any
	arity    int
	toString string
}

func (n *native[T]) Call(i T, args []any) any {
	return n.call(i, args)
}

func (n *native[T]) Arity() int {
	return n.arity
}

func (n *native[T]) String() string {
	if n.toString == "" {
		return NATIVE_FN_STR
	}
	return n.toString
}

func Clock[T any]() native[T] {
	return native[T]{
		arity: 0,
		call: func(i T, argumets []any) any {
			return float64(time.Now().UnixNano()) / float64(time.Second)
		},
	}
}
