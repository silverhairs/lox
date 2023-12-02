package interpreter

import "fmt"

type errReturn struct {
	value any
}

func (r *errReturn) String() string {
	return fmt.Sprintf("%v", r.value)
}
