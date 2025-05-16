package asc

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func WriteKrigAscToWriter(w io.Writer, gridList types.Points, values []float64) error {
	if gridList.Is3D {
		return fmt.Errorf("3D grids are not supported")
	}
	// Create copies of the input data
	grid := make([]types.Point, len(gridList.Points))
	copy(grid, gridList.Points)
	valuesCopy := make([]float64, len(values))
	copy(valuesCopy, values)

	type gridValue struct {
		Point types.Point
		Value float64
	}
	combined := make([]gridValue, len(grid))
	for i := range grid {
		combined[i] = gridValue{Point: grid[i], Value: valuesCopy[i]}
	}

	sort.Slice(combined, func(i, j int) bool {
		if combined[i].Point.Y != combined[j].Point.Y {
			return combined[i].Point.Y > combined[j].Point.Y
		}
		return combined[i].Point.X < combined[j].Point.X
	})

	// After sorting, extract sorted grid and values
	for i := range grid {
		grid[i] = combined[i].Point
		valuesCopy[i] = combined[i].Value
	}

	llx := grid[0].X
	lly := grid[len(grid)-1].Y

	xVals := make(map[float64]struct{})
	yVals := make(map[float64]struct{})
	for _, p := range grid {
		xVals[p.X] = struct{}{}
		yVals[p.Y] = struct{}{}
	}
	nx := len(xVals)
	ny := len(yVals)
	dx := (grid[len(grid)-1].X - grid[0].X) / float64(nx-1)

	w.Write([]byte(fmt.Sprintf("NCOLS %d\n", nx)))
	w.Write([]byte(fmt.Sprintf("NROWS %d\n", ny)))
	w.Write([]byte(fmt.Sprintf("XLLCORNER %f\n", llx)))
	w.Write([]byte(fmt.Sprintf("YLLCORNER %f\n", lly)))
	w.Write([]byte(fmt.Sprintf("CELLSIZE %f\n", dx)))
	w.Write([]byte("NODATA_VALUE -9999\n"))
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			idx := j*nx + i
			if math.IsNaN(valuesCopy[idx]) {
				if i == nx-1 {
					w.Write([]byte("-9999\n"))
				} else {
					w.Write([]byte("-9999 "))
				}
			} else {
				if i == nx-1 {
					w.Write([]byte(fmt.Sprintf("%f\n", valuesCopy[idx])))
				} else {
					w.Write([]byte(fmt.Sprintf("%f ", valuesCopy[idx])))
				}
			}
		}
	}

	return nil
}

func WriteKrigAsc(path string, gridList types.Points, values []float64) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	return WriteKrigAscToWriter(f, gridList, values)
}
