# go-geostat

A comprehensive geostatistical library in Go, implementing various spatial analysis and simulation methods.

## Overview

The library is organized into several packages:

### Kriging (`internal/kriging`)
- Ordinary kriging for 2D and 3D data
- Neighbor selection and optimization
- Variance estimation

### Sequential Gaussian Simulation (`internal/sgs`)
- Multiple realizations generation
- Parallel processing support
- Progress tracking
- Neighbor optimization

### Variogram Modeling (`internal/variogram`)
- Theoretical variogram models (spherical, exponential, gaussian)
- Parameter estimation and fitting
- Model validation

### Empirical Variogram (`internal/empirical`)
- Flexible lag definition
- Multiple estimator types
- Robust calculation methods

### Distance Metrics (`internal/distance`)
- Euclidean distance in 2D and 3D
- Custom metric support
- Optimized calculations

### Common Types (`internal/types`)
- Point and Points types
- Spatial function interfaces
- Error types and handling

### Input/Output (`io`)
- CSV file handling (`io/csv`)
- JSON support (`io/json`)
- ASCII grid files (`io/asc`)

## Installation

To install the library:

```bash
go get github.com/mmaelicke/go-geostat
```

## Usage Examples

### Basic Kriging

```go
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
```

### Sequential Gaussian Simulation

```go
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
```

## Features

- Efficient implementations of geostatistical algorithms
- Support for both 2D and 3D spatial data
- Parallel processing capabilities
- Progress tracking for long-running operations
- Comprehensive error handling
- Example datasets
- Extensive documentation and examples

## Documentation

Each package includes detailed documentation and examples. View the documentation locally using:

```bash
godoc -http=:6060
```

Then visit:
- http://localhost:6060/pkg/github.com/mmaelicke/go-geostat/ for the main documentation
- http://localhost:6060/pkg/github.com/mmaelicke/go-geostat/internal/kriging/ for kriging
- http://localhost:6060/pkg/github.com/mmaelicke/go-geostat/internal/sgs/ for SGS
- etc.

## Command Line Interface

The library includes a CLI for common geostatistical operations:

```bash
go install github.com/mmaelicke/go-geostat@latest
go-geostat --help
```

## References

The implementations are based on:

- Deutsch, C.V. and Journel, A.G. (1998) "GSLIB: Geostatistical Software Library and User's Guide"
- Goovaerts, P. (1997) "Geostatistics for Natural Resources Evaluation"
- Chilès, J.P. and Delfiner, P. (2012) "Geostatistics: Modeling Spatial Uncertainty"
