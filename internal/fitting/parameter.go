package fitting

import (
	"fmt"
	"math"

	"github.com/mmaelicke/go-geostat/internal/types"
	"github.com/mmaelicke/go-geostat/internal/variogram"
	"gonum.org/v1/gonum/optimize"
)

func EstimateParameterFromSampleVariogram(v types.SampleVariogram) (types.BaseParams, error) {
	edges := v.GetEdges()
	semvar := v.GetSemivariances()

	sill := 0.0
	sill_idx := 0
	nugget := semvar[0]

	for i, s := range semvar {
		if s > sill {
			sill = s
			sill_idx = i
		}
	}

	sill -= nugget
	r := edges[sill_idx]

	return types.BaseParams{
		Range:  r,
		Sill:   sill,
		Nugget: nugget,
	}, nil
}

// objectiveFunction implements the least-squares objective function for variogram fitting
type objectiveFunction struct {
	edges       []float64
	semivars    []float64
	model       types.SpatialFunction
	modelParams types.BaseParams
}

func (f *objectiveFunction) Func(x []float64) float64 {
	// Ensure parameters are positive
	if x[0] <= 0 || x[1] <= 0 || x[2] < 0 {
		return math.Inf(1)
	}

	// Update model parameters
	f.modelParams.Range = x[0]
	f.modelParams.Sill = x[1]
	f.modelParams.Nugget = x[2]

	// Calculate sum of squared differences
	sum := 0.0
	for i, h := range f.edges {
		modelVar := f.model.Evaluate(h)
		diff := modelVar - f.semivars[i]
		sum += diff * diff
	}
	return sum
}

func FitVariogram(v types.SampleVariogram, initial types.BaseParams, modelName string) (types.SpatialFunction, error) {
	// Create the variogram model
	model, err := variogram.NewVariogram(modelName, initial)
	if err != nil {
		return nil, fmt.Errorf("failed to create variogram model: %w", err)
	}

	// Get experimental variogram data
	edges := v.GetEdges()
	semivars := v.GetSemivariances()

	// Create objective function
	obj := &objectiveFunction{
		edges:       edges,
		semivars:    semivars,
		model:       model,
		modelParams: initial,
	}

	// Set up optimization problem
	problem := optimize.Problem{
		Func: obj.Func,
	}

	// Set up optimization settings
	settings := &optimize.Settings{
		MajorIterations: 100,
		FuncEvaluations: 1000,
		Converger: &optimize.FunctionConverge{
			Absolute:   1e-6,
			Iterations: 10,
		},
	}

	// Initial guess
	x0 := []float64{initial.Range, initial.Sill, initial.Nugget}

	// Run optimization with Nelder-Mead method
	method := &optimize.NelderMead{}

	result, err := optimize.Minimize(problem, x0, settings, method)
	if err != nil {
		return nil, fmt.Errorf("optimization failed: %w", err)
	}

	// Create final model with optimized parameters
	finalParams := types.BaseParams{
		Range:  result.X[0],
		Sill:   result.X[1],
		Nugget: result.X[2],
	}

	finalModel, err := variogram.NewVariogram(modelName, finalParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create final variogram model: %w", err)
	}

	return finalModel, nil
}
