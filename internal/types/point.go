package types

import (
	"math"
	"time"

	"gonum.org/v1/gonum/spatial/kdtree"
)

type Point struct {
	X, Y, Z float64
	Value   float64
	Time    time.Time
	HasTime bool
	Is3D    bool
}

// implement the Comparable interface from gonum
func (p *Point) Dims() int {
	if p.Is3D {
		return 3
	}
	return 2
}

func (p *Point) Compare(c kdtree.Comparable, d kdtree.Dim) float64 {
	q := c.(*Point)
	switch d {
	case 0:
		return p.X - q.X
	case 1:
		return p.Y - q.Y
	case 2:
		return p.Z - q.Z
	default:
		panic("dimension out of bounds")
	}
}

func (p *Point) Distance(c kdtree.Comparable) float64 {
	sum := 0.0
	for i := 0; i < p.Dims(); i++ {
		sum += math.Pow(p.Compare(c, kdtree.Dim(i)), 2)
	}
	return sum
}
