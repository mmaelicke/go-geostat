package distance

import (
	"math"
	"testing"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func TestPairwiseDistances_Euclidean2D(t *testing.T) {
	points := []types.Point{
		{X: 0, Y: 0},
		{X: 3, Y: 0},
		{X: 0, Y: 4},
	}
	dist := &EuclideanDistance{Is3D: false}
	got, _ := PairwiseDistances(points, dist, false)

	expected := []float64{3, 4, 5} // (0,0)-(3,0):3, (0,0)-(0,4):4, (3,0)-(0,4):5
	if len(got) != len(expected) {
		t.Fatalf("expected %d distances, got %d", len(expected), len(got))
	}
	for i, want := range expected {
		if math.Abs(got[i]-want) > 1e-9 {
			t.Errorf("distance %d: got %v, want %v", i, got[i], want)
		}
	}
}

func TestPairwiseDistances_Euclidean3D(t *testing.T) {
	points := []types.Point{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: -2, Z: 2},
		{X: -2, Y: 2, Z: -1},
	}
	dist := &EuclideanDistance{Is3D: true}
	got, _ := PairwiseDistances(points, dist, false)

	// Distances:
	// (0,0,0)-(1,-2,2): sqrt(1^2 + (-2)^2 + 2^2) = sqrt(1+4+4) = 3
	// (0,0,0)-(-2,2,-1): sqrt((-2)^2 + 2^2 + (-1)^2) = sqrt(4+4+1) = 3
	// (1,-2,2)-(-2,2,-1): sqrt((-3)^2 + 4^2 + (-3)^2) = sqrt(9+16+9) = sqrt(34)
	expected := []float64{3, 3, math.Sqrt(34)}
	if len(got) != len(expected) {
		t.Fatalf("expected %d distances, got %d", len(expected), len(got))
	}
	for i, want := range expected {
		if math.Abs(got[i]-want) > 1e-9 {
			t.Errorf("distance %d: got %v, want %v", i, got[i], want)
		}
	}
}

func TestEuclideanDistance_2D(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   types.Point
		expected float64
	}{
		{
			name:     "origin to (1,1)",
			p1:       types.Point{X: 0, Y: 0},
			p2:       types.Point{X: 1, Y: 1},
			expected: math.Sqrt(2),
		},
		{
			name:     "negative coordinates",
			p1:       types.Point{X: -1, Y: -1},
			p2:       types.Point{X: 1, Y: 1},
			expected: 2 * math.Sqrt(2),
		},
		{
			name:     "same point",
			p1:       types.Point{X: 5, Y: 5},
			p2:       types.Point{X: 5, Y: 5},
			expected: 0,
		},
	}

	d := &EuclideanDistance{Is3D: false}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Compute(&tt.p1, &tt.p2)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Compute() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEuclideanDistance_3D(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   types.Point
		expected float64
	}{
		{
			name:     "origin to (1,1,1)",
			p1:       types.Point{X: 0, Y: 0, Z: 0},
			p2:       types.Point{X: 1, Y: 1, Z: 1},
			expected: math.Sqrt(3),
		},
		{
			name:     "negative coordinates",
			p1:       types.Point{X: -1, Y: -1, Z: -1},
			p2:       types.Point{X: 1, Y: 1, Z: 1},
			expected: 2 * math.Sqrt(3),
		},
		{
			name:     "same point",
			p1:       types.Point{X: 5, Y: 5, Z: 5},
			p2:       types.Point{X: 5, Y: 5, Z: 5},
			expected: 0,
		},
	}

	d := &EuclideanDistance{Is3D: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Compute(&tt.p1, &tt.p2)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Compute() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestManhattanDistance_2D(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   types.Point
		expected float64
	}{
		{
			name:     "origin to (1,1)",
			p1:       types.Point{X: 0, Y: 0},
			p2:       types.Point{X: 1, Y: 1},
			expected: 2,
		},
		{
			name:     "negative coordinates",
			p1:       types.Point{X: -1, Y: -1},
			p2:       types.Point{X: 1, Y: 1},
			expected: 4,
		},
		{
			name:     "same point",
			p1:       types.Point{X: 5, Y: 5},
			p2:       types.Point{X: 5, Y: 5},
			expected: 0,
		},
	}

	d := &ManhattanDistance{Is3D: false}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Compute(&tt.p1, &tt.p2)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Compute() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestChebyshevDistance_2D(t *testing.T) {
	tests := []struct {
		name     string
		p1, p2   types.Point
		expected float64
	}{
		{
			name:     "origin to (1,1)",
			p1:       types.Point{X: 0, Y: 0},
			p2:       types.Point{X: 1, Y: 1},
			expected: 1,
		},
		{
			name:     "negative coordinates",
			p1:       types.Point{X: -1, Y: -1},
			p2:       types.Point{X: 1, Y: 1},
			expected: 2,
		},
		{
			name:     "same point",
			p1:       types.Point{X: 5, Y: 5},
			p2:       types.Point{X: 5, Y: 5},
			expected: 0,
		},
	}

	d := &ChebyshevDistance{Is3D: false}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Compute(&tt.p1, &tt.p2)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Compute() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func BenchmarkEuclideanDistance_2D(b *testing.B) {
	d := &EuclideanDistance{Is3D: false}
	p1 := &types.Point{X: 1, Y: 1}
	p2 := &types.Point{X: 2, Y: 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Compute(p1, p2)
	}
}

func BenchmarkEuclideanDistance_3D(b *testing.B) {
	d := &EuclideanDistance{Is3D: true}
	p1 := &types.Point{X: 1, Y: 1, Z: 1}
	p2 := &types.Point{X: 2, Y: 2, Z: 2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Compute(p1, p2)
	}
}
