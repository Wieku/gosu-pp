package test

import (
	"math"
	"testing"
)

func Assert(a, b interface{}, t *testing.T) {
	if a != b {
		t.Fatalf("Expected: %s, got: %s", a, b)
	}
}

func AssertFloat(a, b, epsilon float64, t *testing.T) {
	if math.IsNaN(b) || math.Abs(a-b) > epsilon {
		t.Fatalf("Expected: %f, got: %f", a, b)
	}
}
