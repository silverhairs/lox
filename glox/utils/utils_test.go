package utils

import (
	"testing"
)

type example struct {
	index  int
	String string
}

func TestMap(t *testing.T) {

	fixtures := []example{
		{index: 0, String: "The value is 0"},
		{index: 1, String: "The value is 1"},
		{index: 2, String: "The value is 2"},
		{index: 3, String: "The value is 3"},
		{index: 4, String: "The value is 4"},
		{index: 5, String: "The value is 5"},
		{index: 6, String: "The value is 6"},
	}

	result := Map(fixtures, func(f example) string { return f.String })

	for i, got := range result {
		expected := fixtures[i].String

		if expected != got {
			t.Fatalf("%d failed -> expected=%q got=%q", i, expected, got)
		}
	}
}
