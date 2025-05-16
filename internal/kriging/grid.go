package kriging

import (
	"fmt"
	"math"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func DenseGrid(p types.Points, dx, dy, dz float64) (types.Points, error) {
	if p.Is3D {
		return dense3d(p, dx, dy, dz)
	}
	return dense2d(p, dx, dy)
}

func dense2d(p types.Points, dx, dy float64) (types.Points, error) {
	if p.Is3D {
		return types.Points{}, fmt.Errorf("points are 3D")
	}
	minx := math.Inf(1)
	maxx := math.Inf(-1)
	miny := math.Inf(1)
	maxy := math.Inf(-1)

	for _, c := range p.Points {
		if c.X < minx {
			minx = c.X
		}
		if c.X > maxx {
			maxx = c.X
		}
		if c.Y < miny {
			miny = c.Y
		}
		if c.Y > maxy {
			maxy = c.Y
		}
	}

	nx := int(math.Ceil((maxx - minx) / dx))
	ny := int(math.Ceil((maxy - miny) / dy))

	grid := make([]types.Point, nx*ny)

	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			grid[i*ny+j] = types.Point{
				X: minx + float64(i)*dx,
				Y: miny + float64(j)*dy,
			}
		}
	}

	return types.Points{
		Points: grid,
		Is3D:   false,
	}, nil
}

func dense3d(p types.Points, dx, dy, dz float64) (types.Points, error) {
	if !p.Is3D {
		return types.Points{}, fmt.Errorf("points are not 3D")
	}
	minx := math.Inf(1)
	maxx := math.Inf(-1)
	miny := math.Inf(1)
	maxy := math.Inf(-1)
	minz := math.Inf(1)
	maxz := math.Inf(-1)

	for _, c := range p.Points {
		if c.X < minx {
			minx = c.X
		}
		if c.X > maxx {
			maxx = c.X
		}
		if c.Y < miny {
			miny = c.Y
		}
		if c.Y > maxy {
			maxy = c.Y
		}
		if c.Z < minz {
			minz = c.Z
		}
		if c.Z > maxz {
			maxz = c.Z
		}
	}

	nx := int(math.Ceil((maxx - minx) / dx))
	ny := int(math.Ceil((maxy - miny) / dy))
	nz := int(math.Ceil((maxz - minz) / dz))

	grid := make([]types.Point, nx*ny*nz)

	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			for k := 0; k < nz; k++ {
				grid[i*ny*nz+j*nz+k] = types.Point{
					X: minx + float64(i)*dx,
					Y: miny + float64(j)*dy,
					Z: minz + float64(k)*dz,
				}
			}
		}
	}

	return types.Points{
		Points: grid,
		Is3D:   true,
	}, nil
}
