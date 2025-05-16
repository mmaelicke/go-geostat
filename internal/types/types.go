package types

import (
	"time"
)

type Profile struct {
	PairwiseTime   time.Duration
	BinningTime    time.Duration
	HistogramTime  time.Duration
	SemivarTime    time.Duration
	EmpiricalTime  time.Duration
	InitialGuess   time.Duration
	FitTime        time.Duration
	TotalTime      time.Duration
	KFitTime       time.Duration
	KInitMeanTime  time.Duration
	KMatMeanTime   time.Duration
	KSolvMeanTime  time.Duration
	KTotalMeanTime time.Duration
	STotalMeanTime time.Duration
}

type SpatialFunction interface {
	Evaluate(float64) float64
	Map([]float64) []float64
	Name() string
	Range() float64
	Sill() float64
	Nugget() float64
	Profile() Profile
	SetProfile(Profile)
}

type BaseParams struct {
	Range  float64 `json:"range"`
	Sill   float64 `json:"sill"`
	Nugget float64 `json:"nugget"`
}

type Points struct {
	Points []Point
	Is3D   bool
}

type SpatialSample interface {
	Length() int
	Sample(size int) Points
	Read() Points
}

type Distance interface {
	Compute(p1, p2 *Point) float64
	Set3D(is3D bool)
}

type Estimator interface {
	Compute(differences []float64) float64
	Map(differences []float64, indices []int, numLags int) ([]float64, []bool)
}

type SampleVariogram interface {
	GetEdges() []float64
	GetHistogram() []int
	GetSemivariances() []float64
}

type SpatialInterpolator interface {
	Fit(condition Points)
	Interpolate(p Points) ([]Estimation, error)
	Profile() Profile
}

type EstimationError uint8

const (
	ErrNone EstimationError = iota
	ErrNoConditionPoints
	ErrSingularMatrix
)

type Estimation struct {
	Field    float64
	Variance float64
	ErrCode  EstimationError
}
