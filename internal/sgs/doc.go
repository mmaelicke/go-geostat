/*
Package sgs implements Sequential Gaussian Simulation (SGS) for spatial data.

SGS is a geostatistical simulation method that generates multiple equally probable
realizations of a spatial field while honoring:
  - The spatial correlation structure (variogram)
  - The observed data points (conditioning)
  - The global histogram of the data

# Algorithm Overview

The SGS algorithm works as follows:

1. Define a random path through all unsampled locations
2. At each location:
  - Use kriging to estimate the local mean and variance
  - Draw a random value from the local conditional distribution
  - Add the simulated value to the conditioning data

3. Repeat for the next location until all points are simulated

# Key Features

  - Parallel simulation of multiple realizations
  - Neighbor search optimization for large datasets
  - Progress tracking for long-running simulations
  - Support for both 2D and 3D spatial fields
  - Customizable kriging parameters

# Basic Usage

	import "github.com/mmaelicke/go-geostat/internal/sgs"

	// Create a new SGS simulator
	simulator := sgs.New(model, maxPoints, dist, showProgress)

	// Fit using conditioning data
	simulator.Fit(conditionPoints)

	// Generate multiple realizations
	simulations, err := simulator.Simulate(targetPoints, numSimulations)

# Performance Optimization

The package includes several optimizations:

  - Neighbor selection: Only the closest points are used for kriging
  - Parallel processing: Multiple realizations are simulated concurrently
  - Memory efficiency: Points are processed sequentially

# Implementation Details

The SGS type implements both the SpatialInterpolator and SpatialSimulator interfaces:

  - Interpolate(): Returns a single realization (n=1 simulation)
  - Simulate(): Returns multiple realizations
  - Fit(): Sets the conditioning points
  - Profile(): Returns performance metrics

The simulation process uses kriging as the local estimator, with the following steps:

1. Random Path Generation:
  - Creates a random permutation of target points
  - Ensures different paths for each realization

2. Local Estimation:
  - Uses ordinary kriging for mean and variance
  - Selects nearest neighbors for conditioning
  - Handles failed estimations gracefully

3. Value Simulation:
  - Draws from local Gaussian distribution
  - Uses estimated mean and variance
  - Updates conditioning data immediately

# Error Handling

The package handles various error conditions:

  - Invalid input points
  - Kriging failures at specific locations
  - Memory allocation issues
  - Numerical instabilities

Failed estimations are marked with NaN values and the simulation continues.

# Progress Tracking

For long-running simulations, the package provides progress tracking:

  - Percentage completion
  - Time per point
  - Kriging performance metrics
  - Overall simulation statistics

# References

For more information about Sequential Gaussian Simulation:

  - Deutsch, C.V. and Journel, A.G. (1998) "GSLIB: Geostatistical Software Library and User's Guide"
  - Gómez-Hernández, J.J. and Journel, A.G. (1993) "Joint Sequential Simulation of Multi-Gaussian Fields"
*/
package sgs
