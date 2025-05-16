// Package kriging implements ordinary kriging interpolation methods.
package kriging

import "fmt"

// ErrInvalidPoints is returned when the input points for kriging are invalid.
// This can occur when points are nil, empty, or contain invalid coordinates.
type ErrInvalidPoints struct {
	// Reason contains a detailed description of why the points are invalid
	Reason string
}

// Error implements the error interface for ErrInvalidPoints.
func (e ErrInvalidPoints) Error() string {
	return fmt.Sprintf("invalid points: %s", e.Reason)
}

// ErrSingularMatrix is returned when the kriging matrix becomes singular
// and cannot be solved. This typically occurs when points are too close together
// or when there are linear dependencies in the spatial configuration.
type ErrSingularMatrix struct {
	// Size is the dimension of the singular matrix
	Size int
	// Reason contains a detailed description of why the matrix is singular
	Reason string
}

// Error implements the error interface for ErrSingularMatrix.
func (e ErrSingularMatrix) Error() string {
	return fmt.Sprintf("singular matrix of size %d: %s", e.Size, e.Reason)
}

// ErrInvalidModel is returned when the spatial function model provided
// for kriging is invalid. This can occur when the model parameters are
// invalid or when the model type is not supported.
type ErrInvalidModel struct {
	// Reason contains a detailed description of why the model is invalid
	Reason string
}

// Error implements the error interface for ErrInvalidModel.
func (e ErrInvalidModel) Error() string {
	return fmt.Sprintf("invalid model: %s", e.Reason)
}

// ErrGridCreation is returned when there is an error creating the
// interpolation grid. This can occur when grid parameters are invalid
// or when there are memory allocation issues.
type ErrGridCreation struct {
	// Reason contains a detailed description of why grid creation failed
	Reason string
}

// Error implements the error interface for ErrGridCreation.
func (e ErrGridCreation) Error() string {
	return fmt.Sprintf("grid creation failed: %s", e.Reason)
}

// ErrInterpolation is returned when kriging interpolation fails at a specific point.
// This can occur due to numerical instability, singular matrices, or other
// computational issues at the given location.
type ErrInterpolation struct {
	// Point is a string representation of the location where interpolation failed
	Point string
	// Reason contains a detailed description of why interpolation failed
	Reason string
}

// Error implements the error interface for ErrInterpolation.
func (e ErrInterpolation) Error() string {
	return fmt.Sprintf("interpolation failed at point %s: %s", e.Point, e.Reason)
}
