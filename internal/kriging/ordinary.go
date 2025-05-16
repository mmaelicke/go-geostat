package kriging

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/mmaelicke/go-geostat/internal/distance"
	"github.com/mmaelicke/go-geostat/internal/types"
	"gonum.org/v1/gonum/mat"
)

type Params struct {
	MaxDistance float64
	MaxPoints   int
	InRange     bool
	dist        types.Distance
}

type StepProfile struct {
	InitTime  time.Duration
	MatTime   time.Duration
	SolvTime  time.Duration
	TotalTime time.Duration
}

type neighbor struct {
	p   *types.Point
	idx int
	d   float64
}

type krigResult struct {
	Index      int
	Estimation types.Estimation
	Profile    StepProfile
	Err        error
}

type OrdinaryKriging struct {
	sf        types.SpatialFunction
	condition types.Points
	params    Params
	profile   types.Profile
	dm        *mat.Dense
	//kd        *kdtree.Tree
	isFitted bool
}

func New(sf types.SpatialFunction, maxPoints int, dist types.Distance, inRange bool) *OrdinaryKriging {
	if dist == nil {
		dist = &distance.EuclideanDistance{}
	}
	return &OrdinaryKriging{
		sf: sf,
		params: Params{
			MaxDistance: math.Inf(1),
			MaxPoints:   maxPoints,
			InRange:     inRange,
			dist:        dist,
		},
		isFitted: false,
	}
}

func (k *OrdinaryKriging) SetDM(dm *mat.Dense) {
	k.dm = dm
	k.isFitted = true
}

func (k *OrdinaryKriging) Fit(condition types.Points) {
	// Filter out NaN values from condition points
	validPoints := make([]types.Point, 0, len(condition.Points))
	for _, p := range condition.Points {
		if !math.IsNaN(p.Value) {
			validPoints = append(validPoints, p)
		}
	}

	k.condition = types.Points{
		Points: validPoints,
		Is3D:   condition.Is3D,
	}
	n := len(validPoints)
	dist := k.params.dist
	nugget := k.sf.Nugget()
	prof := k.sf.Profile()
	k.params.dist.Set3D(condition.Is3D)
	if nugget == 0.0 {
		nugget = 1e-10
	}

	start := time.Now()
	// For ease of reading, I create a actual square distance matrix
	dm := mat.NewDense(n, n, nil)
	for i := range validPoints {
		for j := range validPoints {
			if i == j {
				dm.Set(i, j, nugget)
			} else {
				v := k.sf.Evaluate(dist.Compute(&validPoints[i], &validPoints[j]))
				dm.Set(i, j, v)
			}
		}
	}
	k.dm = dm
	prof.FitTime = time.Since(start)
	k.profile = prof
	k.isFitted = true
}

func (k *OrdinaryKriging) Profile() types.Profile {
	return k.profile
}

func (k *OrdinaryKriging) Interpolate(p types.Points) ([]types.Estimation, error) {
	if !k.isFitted {
		return []types.Estimation{}, fmt.Errorf("kriging model not fitted")
	}

	results := make(chan krigResult, len(p.Points))

	wg := sync.WaitGroup{}
	for i, c := range p.Points {
		wg.Add(1)
		go func(i int, c types.Point) {
			defer wg.Done()

			est, prof, err := k.krige(c)
			results <- krigResult{
				Index:      i,
				Estimation: est,
				Profile:    prof,
				Err:        err,
			}
		}(i, c)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	estimations := make([]types.Estimation, len(p.Points))
	initSum := time.Duration(0)
	matSum := time.Duration(0)
	solvSum := time.Duration(0)
	totalSum := time.Duration(0)
	n := 0

	for r := range results {
		if r.Err != nil || r.Estimation.ErrCode != types.ErrNone {
			estimations[r.Index] = types.Estimation{
				Field:    math.NaN(),
				Variance: math.NaN(),
			}
		} else {
			estimations[r.Index] = r.Estimation
			initSum += r.Profile.InitTime
			matSum += r.Profile.MatTime
			solvSum += r.Profile.SolvTime
			totalSum += r.Profile.TotalTime
			n++
		}
	}
	if n == 0 {
		n = 1
	}
	k.profile.KInitMeanTime = initSum / time.Duration(n)
	k.profile.KMatMeanTime = matSum / time.Duration(n)
	k.profile.KSolvMeanTime = solvSum / time.Duration(n)
	k.profile.KTotalMeanTime = totalSum / time.Duration(n)
	return estimations, nil
}

func (k *OrdinaryKriging) krige(p types.Point) (types.Estimation, StepProfile, error) {
	if !k.isFitted {
		return types.Estimation{}, StepProfile{}, fmt.Errorf("kriging model not fitted")
	}
	prof := StepProfile{}
	maxp := k.params.MaxPoints

	start := time.Now()
	startTotal := start

	allNeighbors := make([]neighbor, len(k.condition.Points))
	for i, c := range k.condition.Points {
		d := k.params.dist.Compute(&c, &p)
		allNeighbors[i] = neighbor{
			p:   &k.condition.Points[i],
			d:   d,
			idx: i,
		}
	}

	sort.Slice(allNeighbors, func(i, j int) bool {
		return allNeighbors[i].d < allNeighbors[j].d
	})

	// Check if we have enough neighbors
	if len(allNeighbors) == 0 {
		return types.Estimation{ErrCode: types.ErrNoConditionPoints}, StepProfile{}, nil
	}

	// Adjust maxp if we don't have enough neighbors
	if maxp > len(allNeighbors) {
		maxp = len(allNeighbors)
	}

	neighbors := allNeighbors[:maxp]

	// Only check range if we have more points than maxp
	if k.params.InRange && len(allNeighbors) > maxp && allNeighbors[maxp].d < k.sf.Range() {
		return types.Estimation{ErrCode: types.ErrNoConditionPoints}, StepProfile{}, nil
	}

	prof.InitTime = time.Since(start)

	start = time.Now()
	aData := make([]float64, (maxp+1)*(maxp+1))
	bData := make([]float64, maxp+1)

	for i := 0; i < maxp+1; i++ {
		for j := 0; j < maxp+1; j++ {
			if i == maxp || j == maxp {
				if i == maxp && j == maxp {
					aData[i*(maxp+1)+j] = 0
				} else {
					aData[i*(maxp+1)+j] = 1
				}
			} else {
				aData[i*(maxp+1)+j] = k.dm.At(neighbors[i].idx, neighbors[j].idx)
			}
		}
	}

	for i := range neighbors {
		v := k.sf.Evaluate(k.params.dist.Compute(&p, neighbors[i].p))
		bData[i] = v
	}
	bData[maxp] = 1
	prof.MatTime = time.Since(start)

	start = time.Now()
	A := mat.NewDense(maxp+1, maxp+1, aData)
	b := mat.NewVecDense(maxp+1, bData)

	var L mat.VecDense
	err := L.SolveVec(A, b)
	if err != nil {
		return types.Estimation{ErrCode: types.ErrSingularMatrix}, StepProfile{}, fmt.Errorf("error solving linear system: %v", err)
	}
	prof.SolvTime = time.Since(start)

	field := 0.0
	for i := range neighbors {
		field += L.AtVec(i) * neighbors[i].p.Value
	}

	// Calculate error variance for Ordinary Kriging
	variance := 0.0
	for i := range neighbors {
		variance += L.AtVec(i) * bData[i]
	}
	variance += L.AtVec(maxp)

	estimation := types.Estimation{
		Field:    field,
		Variance: variance,
		ErrCode:  types.ErrNone,
	}
	prof.TotalTime = time.Since(startTotal)
	return estimation, prof, nil
}
