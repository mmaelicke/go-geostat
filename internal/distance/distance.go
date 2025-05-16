package distance

import (
	"math"

	"github.com/mmaelicke/go-geostat/internal/types"
)

type EuclideanDistance struct {
	Is3D bool
}

func (d *EuclideanDistance) Compute(p1, p2 *types.Point) float64 {
	if d.Is3D {
		return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2) + math.Pow(p1.Z-p2.Z, 2))
	}
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

func (d *EuclideanDistance) Set3D(is3D bool) {
	d.Is3D = is3D
}

type ChebyshevDistance struct {
	Is3D bool
}

func (d *ChebyshevDistance) Compute(p1, p2 *types.Point) float64 {
	if d.Is3D {
		return math.Max(math.Abs(p1.X-p2.X), math.Max(math.Abs(p1.Y-p2.Y), math.Abs(p1.Z-p2.Z)))
	}
	return math.Max(math.Abs(p1.X-p2.X), math.Abs(p1.Y-p2.Y))
}

func (d *ChebyshevDistance) Set3D(is3D bool) {
	d.Is3D = is3D
}

type ManhattanDistance struct {
	Is3D bool
}

func (d *ManhattanDistance) Compute(p1, p2 *types.Point) float64 {
	if d.Is3D {
		return math.Abs(p1.X-p2.X) + math.Abs(p1.Y-p2.Y) + math.Abs(p1.Z-p2.Z)
	}
	return math.Abs(p1.X-p2.X) + math.Abs(p1.Y-p2.Y)
}

func (d *ManhattanDistance) Set3D(is3D bool) {
	d.Is3D = is3D
}
