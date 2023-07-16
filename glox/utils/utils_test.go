package utils

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		Value  int
		String string
	}{
		{Value: 1, String: "The value is 1"},
		{Value: 2, String: "The value is 2"},
		{Value: 3, String: "The value is 3"},
		{Value: 4, String: "The value is 4"},
		{Value: 5, String: "The value is 5"},
		{Value: 6, String: "The value is 6"},
		{Value: 7, String: "The value is 7"},
	}

	for i, test := range tests {
		expected := fmt.Sprintf("The value is %d", test.Value)
		got := test.String

		if expected != got {
			t.Fatalf("%d failed -> expected=%q got=%q", i, expected, got)
		}
	}
}
