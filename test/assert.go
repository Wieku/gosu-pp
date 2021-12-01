package test

import (
	"math"
	"testing"
)

func Assert(tested string, a, b interface{}, t *testing.T) {
	if a != b {
		t.Fatalf("%s: Expected: %s, got: %s", tested, a, b)
	}
}

func AssertFloat(tested string, a, b, epsilon float64, t *testing.T) {
	if math.IsNaN(b) || math.Abs(a-b) > epsilon {
		t.Fatalf("%s: Expected: %.17f, got: %.17f", tested, a, b)
	}
}
