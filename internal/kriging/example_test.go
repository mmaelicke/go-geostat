package kriging_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmaelicke/go-geostat/internal/distance"
	"github.com/mmaelicke/go-geostat/internal/empirical"
	"github.com/mmaelicke/go-geostat/internal/estimator"
	"github.com/mmaelicke/go-geostat/internal/kriging"
	"github.com/mmaelicke/go-geostat/internal/types"
	"github.com/mmaelicke/go-geostat/internal/variogram"
	"github.com/mmaelicke/go-geostat/io/csv"
)

func ExampleNew_withPancakeData() {
	// Get the workspace root directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}
	// Move up one level to reach workspace root
	rootDir := filepath.Dir(filepath.Dir(wd))

	// Load the pancake dataset
	data, err := csv.ReadCSV(filepath.Join(rootDir, "data", "pancake.csv"), "", "", "", "", "", "", false)
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}
	points := data.Read()

	// Create distance metric and estimator
	dist := &distance.EuclideanDistance{Is3D: false}
	est := &estimator.Matheron{}

	// Calculate empirical variogram with 25 lags
	vg := empirical.NewEmpiricalVariogram(points, 25, 500, dist, est)
	err = vg.Compute()
	if err != nil {
		fmt.Printf("Error computing variogram: %v\n", err)
		return
	}

	// Fit a spherical model
	model, err := vg.Fit("spherical")
	if err != nil {
		fmt.Printf("Error fitting model: %v\n", err)
		return
	}

	// Create kriging interpolator with 10 neighbors
	kr := kriging.New(model, 10, dist, false)
	kr.Fit(points)

	// Create a grid for interpolation
	grid, err := kriging.DenseGrid(points, 10.0, 10.0, 0)
	if err != nil {
		fmt.Printf("Error creating grid: %v\n", err)
		return
	}

	// Perform interpolation
	estimations, err := kr.Interpolate(grid)
	if err != nil {
		fmt.Printf("Error interpolating: %v\n", err)
		return
	}

	// Print some statistics about the interpolation
	var sumField, sumVar float64
	var minField, maxField = estimations[0].Field, estimations[0].Field
	var minVar, maxVar = estimations[0].Variance, estimations[0].Variance

	for _, est := range estimations {
		sumField += est.Field
		sumVar += est.Variance
		if est.Field < minField {
			minField = est.Field
		}
		if est.Field > maxField {
			maxField = est.Field
		}
		if est.Variance < minVar {
			minVar = est.Variance
		}
		if est.Variance > maxVar {
			maxVar = est.Variance
		}
	}

	meanField := sumField / float64(len(estimations))
	meanVar := sumVar / float64(len(estimations))

	fmt.Printf("Interpolation Results:\n")
	fmt.Printf("Number of points interpolated: %d\n", len(estimations))
	fmt.Printf("Mean estimated value: %.2f\n", meanField)
	fmt.Printf("Value range: %.2f - %.2f\n", minField, maxField)
	fmt.Printf("Mean estimation variance: %.2f\n", meanVar)
	fmt.Printf("Variance range: %.2f - %.2f\n", minVar, maxVar)

	// Output:
	// Interpolation Results:
	// Number of points interpolated: 2500
	// Mean estimated value: 187.05
	// Value range: 91.27 - 243.37
	// Mean estimation variance: 251.95
	// Variance range: 171.12 - 549.08
}

func ExampleNew() {
	// Create a simple spherical model
	model, _ := variogram.NewVariogram("spherical", types.BaseParams{
		Range: 100,
		Sill:  1.0,
	})

	// Create distance metric
	dist := &distance.EuclideanDistance{Is3D: false}

	// Create kriging interpolator
	kr := kriging.New(model, 10, dist, false)

	// Create some sample points
	points := types.Points{
		Points: []types.Point{
			{X: 0, Y: 0, Value: 1.0},
			{X: 1, Y: 1, Value: 2.0},
			{X: 2, Y: 2, Value: 3.0},
		},
		Is3D: false,
	}

	// Fit the model
	kr.Fit(points)

	// Create a point to interpolate
	newPoint := types.Points{
		Points: []types.Point{{X: 0.5, Y: 0.5}},
		Is3D:   false,
	}

	// Perform interpolation
	estimations, _ := kr.Interpolate(newPoint)
	fmt.Printf("Estimated value at (0.5, 0.5): %.2f\n", estimations[0].Field)
	fmt.Printf("Estimation variance: %.2f\n", estimations[0].Variance)

	// Output:
	// Estimated value at (0.5, 0.5): 1.50
	// Estimation variance: 0.01
}

func ExampleNew_with3D() {
	// Create a simple exponential model
	model, _ := variogram.NewVariogram("exponential", types.BaseParams{
		Range: 100,
		Sill:  1.0,
	})

	// Create 3D distance metric
	dist := &distance.EuclideanDistance{Is3D: true}

	// Create kriging interpolator
	kr := kriging.New(model, 10, dist, false)

	// Create some 3D sample points
	points := types.Points{
		Points: []types.Point{
			{X: 0, Y: 0, Z: 0, Value: 1.0},
			{X: 1, Y: 1, Z: 1, Value: 2.0},
			{X: 2, Y: 2, Z: 2, Value: 3.0},
		},
		Is3D: true,
	}

	// Fit the model
	kr.Fit(points)

	// Create a 3D point to interpolate
	newPoint := types.Points{
		Points: []types.Point{{X: 0.5, Y: 0.5, Z: 0.5}},
		Is3D:   true,
	}

	// Perform interpolation
	estimations, _ := kr.Interpolate(newPoint)
	fmt.Printf("Estimated value at (0.5, 0.5, 0.5): %.2f\n", estimations[0].Field)
	fmt.Printf("Estimation variance: %.2f\n", estimations[0].Variance)

	// Output:
	// Estimated value at (0.5, 0.5, 0.5): 1.50
	// Estimation variance: 0.03
}
