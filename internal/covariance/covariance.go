// Package covariance implements one-dimensional covariance functions for the turning bands method.
package covariance

import "math"

// CovarianceFunction defines the interface for all 1D covariance functions
type CovarianceFunction interface {
	// Evaluate returns the covariance value at lag h
	Evaluate(h float64) float64
	// Range returns the range parameter of the covariance function
	Range() float64
	// Sill returns the sill parameter of the covariance function
	Sill() float64
}

// BaseParams contains common parameters for covariance functions
type BaseParams struct {
	Range float64 // Range parameter
	Sill  float64 // Sill parameter
}

// Spherical implements the 1D spherical covariance function
type Spherical struct {
	BaseParams
}

// Exponential implements the 1D exponential covariance function
type Exponential struct {
	BaseParams
}

// Gaussian implements the 1D Gaussian covariance function
type Gaussian struct {
	BaseParams
}

// Nugget implements the nugget effect
type Nugget struct {
	Value float64 // Nugget value
}

// Evaluate implements the spherical covariance function
func (s Spherical) Evaluate(h float64) float64 {
	if h >= s.BaseParams.Range {
		return 0
	}
	return s.BaseParams.Sill * (1 - 1.5*(h/s.BaseParams.Range) + 0.5*math.Pow(h/s.BaseParams.Range, 3))
}

// Range returns the range parameter
func (s Spherical) Range() float64 {
	return s.BaseParams.Range
}

// Sill returns the sill parameter
func (s Spherical) Sill() float64 {
	return s.BaseParams.Sill
}

// Evaluate implements the exponential covariance function
func (e Exponential) Evaluate(h float64) float64 {
	return e.BaseParams.Sill * math.Exp(-3*h/e.BaseParams.Range)
}

// Range returns the range parameter
func (e Exponential) Range() float64 {
	return e.BaseParams.Range
}

// Sill returns the sill parameter
func (e Exponential) Sill() float64 {
	return e.BaseParams.Sill
}

// Evaluate implements the Gaussian covariance function
func (g Gaussian) Evaluate(h float64) float64 {
	return g.BaseParams.Sill * math.Exp(-3*math.Pow(h/g.BaseParams.Range, 2))
}

// Range returns the range parameter
func (g Gaussian) Range() float64 {
	return g.BaseParams.Range
}

// Sill returns the sill parameter
func (g Gaussian) Sill() float64 {
	return g.BaseParams.Sill
}

// Evaluate implements the nugget effect
func (n Nugget) Evaluate(h float64) float64 {
	if h == 0 {
		return n.Value
	}
	return 0
}

// Range returns 0 for nugget effect
func (n Nugget) Range() float64 {
	return 0
}

// Sill returns the nugget value
func (n Nugget) Sill() float64 {
	return n.Value
}
