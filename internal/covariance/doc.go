// Package covariance implements one-dimensional covariance functions for the turning bands method.
//
// # Mathematical Background
//
// The turning bands method (TBM) simulates random fields in multi-dimensional space by performing
// a series of one-dimensional simulations along lines. The one-dimensional covariance functions
// implemented in this package are mathematically derived from their multi-dimensional counterparts
// to ensure that the resulting simulation reproduces the desired spatial structure.
//
// For a given multi-dimensional covariance function C_d(r), where r is the distance in d-dimensional
// space, the corresponding one-dimensional covariance function C_1(h) is derived through a
// mathematical transformation. This transformation ensures that when the one-dimensional processes
// are combined in the higher dimension, they reproduce the target multi-dimensional covariance
// structure.
//
// # Implemented Functions
//
//  1. Spherical Covariance
//     The one-dimensional spherical covariance function is derived from the multi-dimensional
//     spherical model. For lag h:
//     C_1(h) = σ² * (1 - 1.5(h/a) + 0.5(h/a)³)  for h ≤ a
//     C_1(h) = 0                                for h > a
//     where:
//     - σ² is the sill (variance)
//     - a is the range parameter
//
//  2. Exponential Covariance
//     The one-dimensional exponential covariance function is derived from the multi-dimensional
//     exponential model. For lag h:
//     C_1(h) = σ² * exp(-3h/a)
//     where:
//     - σ² is the sill (variance)
//     - a is the range parameter
//     Note: The factor 3 ensures that the practical range (where C(h) ≈ 0.05σ²) matches the
//     range parameter a.
//
//  3. Gaussian Covariance
//     The one-dimensional Gaussian covariance function is derived from the multi-dimensional
//     Gaussian model. For lag h:
//     C_1(h) = σ² * exp(-3(h/a)²)
//     where:
//     - σ² is the sill (variance)
//     - a is the range parameter
//     Note: The factor 3 ensures that the practical range (where C(h) ≈ 0.05σ²) matches the
//     range parameter a.
//
//  4. Nugget Effect
//     The nugget effect represents a discontinuity at the origin, modeling measurement error
//     or micro-scale variation. For lag h:
//     C_1(h) = σ²  for h = 0
//     C_1(h) = 0   for h > 0
//     where σ² is the nugget variance.
//
// # Usage Example
//
// To create and use a spherical covariance function:
//
//	// Create a spherical covariance function with range 10 and sill 1
//	sph := Spherical{BaseParams{Range: 10.0, Sill: 1.0}}
//
//	// Evaluate the covariance at different lags
//	c0 := sph.Evaluate(0.0)   // At origin
//	c5 := sph.Evaluate(5.0)   // At half range
//	c10 := sph.Evaluate(10.0) // At range
//
// # References
//
// The mathematical derivation of these one-dimensional covariance functions can be found in:
//   - Mantoglou, A., & Wilson, J. L. (1982). The turning bands method for simulation of random
//     fields using line generation by a spectral method. Water Resources Research, 18(5), 1379-1394.
//   - Emery, X. (2008). A turning bands program for conditional co-simulation of cross-correlated
//     Gaussian random fields. Computers & Geosciences, 34(12), 1850-1862.
package covariance
