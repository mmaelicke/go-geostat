package estimator

import (
	"math"
	"testing"
)

func TestMatheronEstimator(t *testing.T) {
	diffs := []float64{2, -1, 3}
	// sum of squares: 4 + 1 + 9 = 14
	// n = 3
	// expected: 14 / (2*3) = 2.333...
	est := &Matheron{}
	got := est.Compute(diffs)
	if math.IsNaN(got) {
		t.Fatalf("unexpected error: %v", got)
	}
	expected := 14.0 / 6.0
	if math.Abs(got-expected) > 1e-9 {
		t.Errorf("Matheron: got %v, want %v", got, expected)
	}
}

func TestCressieEstimator(t *testing.T) {
	diffs := []float64{2, -1, 3}
	// mean(sqrt(abs(diffs))) = (sqrt(2) + 1 + sqrt(3)) / 3
	mean := (math.Sqrt(2) + 1 + math.Sqrt(3)) / 3
	n := float64(len(diffs))
	numerator := math.Pow(mean, 4)
	denominator := 0.457 + 0.494/n + 0.045/(n*n)
	expected := numerator / (2 * denominator)

	est := &Cressie{}
	got := est.Compute(diffs)
	if math.IsNaN(got) {
		t.Fatalf("unexpected error: %v", got)
	}
	if math.Abs(got-expected) > 1e-9 {
		t.Errorf("Cressie: got %v, want %v", got, expected)
	}
}
