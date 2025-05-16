/*
Package kriging implements ordinary kriging interpolation methods for spatial data.

Kriging is a method of interpolation that predicts values at unsampled locations using
a weighted average of observed values at nearby locations. The weights are determined
by the spatial correlation structure of the data, which is modeled using a variogram.

Key Features:

  - Ordinary kriging implementation
  - Support for various distance metrics
  - Efficient neighbor selection
  - Variance estimation
  - Support for 2D and 3D datasets

Basic Usage:

	import "github.com/mmaelicke/go-geostat/internal/kriging"

	// Create a new kriging interpolator
	kr := kriging.New(model, maxPoints, dist, false)

	// Fit the model to your data
	kr.Fit(points)

	// Perform interpolation
	estimations, err := kr.Interpolate(newPoints)

The package provides several error types for handling common kriging issues:

  - ErrInvalidPoints: When input points are invalid
  - ErrSingularMatrix: When the kriging matrix is singular
  - ErrInvalidModel: When the spatial function model is invalid
  - ErrGridCreation: When grid creation fails
  - ErrInterpolation: When interpolation fails at a specific point

For performance optimization, the package includes:

  - Efficient matrix operations
  - Neighbor selection algorithms
  - Parallel processing capabilities
*/
package kriging
