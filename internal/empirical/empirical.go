package empirical

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mmaelicke/go-geostat/internal/distance"
	"github.com/mmaelicke/go-geostat/internal/estimator"
	"github.com/mmaelicke/go-geostat/internal/fitting"
	"github.com/mmaelicke/go-geostat/internal/lagging"
	"github.com/mmaelicke/go-geostat/internal/types"
)

type intermediate struct {
	distances   []float64
	groups      []int
	histogram   []int
	differences []float64
}

type Properties struct {
	numLags   int
	maxLag    float64
	dist      types.Distance
	estimator types.Estimator
}

type EmpiricalVariogram struct {
	sample       types.Points
	semivariance []float64
	binEdges     []float64
	isCalulated  bool
	Properties
	intermediate
	profile types.Profile
}

var logger = slog.Default()

func NewEmpiricalVariogram(sample types.Points, numLags int, maxLag float64, dist types.Distance, e types.Estimator) *EmpiricalVariogram {
	if dist == nil {
		dist = &distance.EuclideanDistance{}
	}
	dist.Set3D(sample.Is3D)

	if e == nil {
		e = &estimator.Matheron{}
	}

	return &EmpiricalVariogram{
		sample:       sample,
		semivariance: make([]float64, 0),
		binEdges:     make([]float64, numLags),
		Properties: Properties{
			numLags:   numLags,
			maxLag:    maxLag,
			dist:      dist,
			estimator: e,
		},
		intermediate: intermediate{
			distances: make([]float64, 0),
			groups:    make([]int, 0),
			histogram: make([]int, 0),
		},
	}
}

func (v *EmpiricalVariogram) Compute() error {
	startTotal := time.Now()

	start := time.Now()
	v.intermediate.distances, v.intermediate.differences = distance.PairwiseDistances(v.sample.Points, v.dist, true)
	v.profile.PairwiseTime = time.Since(start)

	start = time.Now()
	var err error
	v.binEdges, err = lagging.CalculateEdges(v.intermediate.distances, v.numLags, v.maxLag)
	if err != nil {
		return err
	}
	v.intermediate.groups = make([]int, len(v.intermediate.distances))
	v.intermediate.groups = lagging.GetEdgeIndex(v.intermediate.distances, v.binEdges)
	v.profile.BinningTime = time.Since(start)

	start = time.Now()
	v.intermediate.histogram = make([]int, v.numLags)
	for _, g := range v.intermediate.groups {
		if g >= 0 && g < v.numLags {
			v.intermediate.histogram[g]++
		}
	}
	v.profile.HistogramTime = time.Since(start)

	start = time.Now()
	var mask []bool
	v.semivariance, mask = v.estimator.Map(v.intermediate.differences, v.intermediate.groups, v.numLags)
	v.profile.SemivarTime = time.Since(start)

	v.profile.EmpiricalTime = time.Since(startTotal)

	for i, m := range mask {
		if m {
			logger.Warn(fmt.Sprintf("For lag %d, no semi-variance was computed", i))
		}
	}
	v.isCalulated = true

	return nil
}

func (v *EmpiricalVariogram) GetSemivariances() []float64 {
	return v.semivariance
}

func (v *EmpiricalVariogram) GetEdges() []float64 {
	return v.binEdges
}

func (v *EmpiricalVariogram) GetHistogram() []int {
	return v.intermediate.histogram
}

func (v *EmpiricalVariogram) GetProperties() Properties {
	return v.Properties
}

func (v *EmpiricalVariogram) GetProfile() types.Profile {
	return v.profile
}

func (v *EmpiricalVariogram) Fit(modelName string) (types.SpatialFunction, error) {
	if !v.isCalulated {
		return nil, fmt.Errorf("empirical variogram is not calculated")
	}
	start := time.Now()
	startTotal := start
	profile := v.profile
	params, err := fitting.EstimateParameterFromSampleVariogram(v)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate parameters: %w", err)
	}
	profile.InitialGuess = time.Since(start)

	start = time.Now()
	model, err := fitting.FitVariogram(v, params, modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to fit variogram: %w", err)
	}
	profile.FitTime = time.Since(start)
	profile.TotalTime = profile.EmpiricalTime + time.Since(startTotal)
	model.SetProfile(profile)

	return model, nil
}
