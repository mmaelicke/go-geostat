package covariance

import (
	"math"
	"testing"
)

func TestSpherical(t *testing.T) {
	s := Spherical{BaseParams{Range: 10.0, Sill: 1.0}}

	tests := []struct {
		h      float64
		expect float64
	}{
		{0.0, 1.0},    // At origin
		{5.0, 0.3125}, // At half range
		{10.0, 0.0},   // At range
		{15.0, 0.0},   // Beyond range
	}

	for _, tt := range tests {
		got := s.Evaluate(tt.h)
		if math.Abs(got-tt.expect) > 1e-6 {
			t.Errorf("Spherical.Evaluate(%v) = %v, want %v", tt.h, got, tt.expect)
		}
	}
}

func TestExponential(t *testing.T) {
	e := Exponential{BaseParams{Range: 10.0, Sill: 1.0}}

	tests := []struct {
		h      float64
		expect float64
	}{
		{0.0, 1.0},       // At origin
		{5.0, 0.223130},  // At half range
		{10.0, 0.049787}, // At range
		{20.0, 0.002479}, // Beyond range
	}

	for _, tt := range tests {
		got := e.Evaluate(tt.h)
		if math.Abs(got-tt.expect) > 1e-6 {
			t.Errorf("Exponential.Evaluate(%v) = %v, want %v", tt.h, got, tt.expect)
		}
	}
}

func TestGaussian(t *testing.T) {
	g := Gaussian{BaseParams{Range: 10.0, Sill: 1.0}}

	tests := []struct {
		h      float64
		expect float64
	}{
		{0.0, 1.0},       // At origin
		{5.0, 0.4723},    // At half range
		{10.0, 0.049787}, // At range
		{20.0, 0.000123}, // Beyond range
	}

	for _, tt := range tests {
		got := g.Evaluate(tt.h)
		if math.Abs(got-tt.expect) > 1e-3 {
			t.Errorf("Gaussian.Evaluate(%v) = %v, want %v", tt.h, got, tt.expect)
		}
	}
}

func TestNugget(t *testing.T) {
	n := Nugget{Value: 1.0}

	tests := []struct {
		h      float64
		expect float64
	}{
		{0.0, 1.0}, // At origin
		{0.1, 0.0}, // Near origin
		{1.0, 0.0}, // Away from origin
	}

	for _, tt := range tests {
		got := n.Evaluate(tt.h)
		if math.Abs(got-tt.expect) > 1e-6 {
			t.Errorf("Nugget.Evaluate(%v) = %v, want %v", tt.h, got, tt.expect)
		}
	}
}
