package native

import (
	"math"
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	clock := Clock[any]()
	got := clock.call(nil, []any{})
	want := (float64(time.Now().UnixNano()) / float64(time.Millisecond)) / 1000

	val, isOk := got.(float64)
	if !isOk {
		t.Fatalf("clock time is not a float64. got='%T'", got)
	}

	margin := math.Round((want/val)*100) / 100
	if margin > 1 || margin < 0.9 {
		t.Fatalf("wrong value for time. want='~%v' got='%v'", want, val)
	}

	if clock.String() != NATIVE_FN_STR {
		t.Errorf("clock.String() as a wrong value. want='%s' got='%s'", NATIVE_FN_STR, clock.String())
		t.Fail()
	}

	if clock.Arity() != 0 {
		t.Errorf("clock.Arity() has wrong value. want='0' got='%d'", clock.Arity())
	}
}
