package sgs_test

import (
	"fmt"

	"github.com/mmaelicke/go-geostat/internal/distance"
	"github.com/mmaelicke/go-geostat/internal/empirical"
	"github.com/mmaelicke/go-geostat/internal/estimator"
	"github.com/mmaelicke/go-geostat/internal/kriging"
	"github.com/mmaelicke/go-geostat/internal/sgs"
	"github.com/mmaelicke/go-geostat/io/csv"
)

func Example_withPancakeData() {
	// Load the pancake dataset from the data directory
	data, err := csv.ReadCSV("../../data/pancake.csv", "", "", "", "", "", "", false)
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

	// Create SGS simulator with 10 neighbors
	simulator := sgs.New(model, 10, dist, false) // disable progress for example output
	simulator.Fit(points)

	// Create a grid for simulation
	grid, err := kriging.DenseGrid(points, 10.0, 10.0, 0)
	if err != nil {
		fmt.Printf("Error creating grid: %v\n", err)
		return
	}

	// Generate 5 realizations
	simulations, err := simulator.Simulate(grid, 5)
	if err != nil {
		fmt.Printf("Error simulating: %v\n", err)
		return
	}

	// Print only the dimensions and counts since values are random
	fmt.Printf("Simulation results:\n")
	fmt.Printf("Number of realizations: %d\n", len(simulations))
	fmt.Printf("Points per realization: %d\n", len(simulations[0]))
	fmt.Printf("Grid spacing: 10.0 x 10.0\n")

	// Output:
	// Simulation results:
	// Number of realizations: 5
	// Points per realization: 2500
	// Grid spacing: 10.0 x 10.0
}
