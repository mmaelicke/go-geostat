/*
Package geostat provides a comprehensive geostatistical library in Go, implementing
various spatial analysis and simulation methods.

# Overview

The library is organized into several packages:

  - kriging: Implementation of various kriging methods

  - Ordinary kriging for 2D and 3D data

  - Neighbor selection and optimization

  - Variance estimation

  - sgs: Sequential Gaussian Simulation

  - Multiple realizations generation

  - Parallel processing support

  - Progress tracking

  - Neighbor optimization

  - variogram: Variogram modeling and fitting

  - Theoretical variogram models (spherical, exponential, gaussian)

  - Parameter estimation and fitting

  - Model validation

  - empirical: Empirical variogram calculation

  - Flexible lag definition

  - Multiple estimator types

  - Robust calculation methods

  - distance: Distance metrics

  - Euclidean distance in 2D and 3D

  - Custom metric support

  - Optimized calculations

  - types: Common data structures

  - Point and Points types

  - Spatial function interfaces

  - Error types and handling

  - io: Input/Output operations

  - CSV file handling (io/csv)

  - JSON support (io/json)

  - ASCII grid files (io/asc)

# Getting Started

Basic kriging example:

	import (
		"github.com/mmaelicke/go-geostat/internal/kriging"
		"github.com/mmaelicke/go-geostat/internal/types"
		"github.com/mmaelicke/go-geostat/internal/variogram"
	)

	// Create a variogram model
	model, _ := variogram.NewVariogram("spherical", types.BaseParams{
		Range: 100,
		Sill:  1.0,
	})

	// Create kriging interpolator
	kr := kriging.New(model, 10, nil, false)

	// Fit and interpolate
	kr.Fit(points)
	estimations, _ := kr.Interpolate(newPoints)

Basic SGS example:

	import (
		"github.com/mmaelicke/go-geostat/internal/sgs"
		"github.com/mmaelicke/go-geostat/internal/types"
		"github.com/mmaelicke/go-geostat/internal/variogram"
	)

	// Create SGS simulator
	simulator := sgs.New(model, 10, nil, true)

	// Generate realizations
	simulator.Fit(points)
	simulations, _ := simulator.Simulate(targetPoints, 5)

# Features

The library provides:

  - Efficient implementations of geostatistical algorithms
  - Support for both 2D and 3D spatial data
  - Parallel processing capabilities
  - Progress tracking for long-running operations
  - Comprehensive error handling
  - Example datasets
  - Extensive documentation and examples

# Installation

To install the library:

	go get github.com/mmaelicke/go-geostat

# Contributing

Contributions are welcome! Please see the repository README for guidelines.

# References

The implementations are based on:

  - Deutsch, C.V. and Journel, A.G. (1998) "GSLIB: Geostatistical Software Library and User's Guide"
  - Goovaerts, P. (1997) "Geostatistics for Natural Resources Evaluation"
  - Chilès, J.P. and Delfiner, P. (2012) "Geostatistics: Modeling Spatial Uncertainty"
*/
package main
