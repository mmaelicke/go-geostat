package sgs

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mmaelicke/go-geostat/internal/kriging"
	"github.com/mmaelicke/go-geostat/internal/types"
	"gonum.org/v1/gonum/stat/distuv"
)

type SGS struct {
	sf              types.SpatialFunction
	newInterpolator func() types.SpatialInterpolator
	profile         types.Profile
	condition       types.Points
	dist            types.Distance
	isFitted        bool
	maxPoints       int
	useNeighbors    bool
	progress        *progressTracker
}

func New(sf types.SpatialFunction, maxPoints int, dist types.Distance, showProgress bool) *SGS {
	newInterpolator := func() types.SpatialInterpolator {
		return kriging.New(sf, maxPoints, dist, false)
	}

	var pt *progressTracker
	if showProgress {
		pt = newProgressTracker()
	}

	return &SGS{
		sf:              sf,
		newInterpolator: newInterpolator,
		isFitted:        false,
		profile:         types.Profile{},
		dist:            dist,
		maxPoints:       maxPoints,
		useNeighbors:    true,
		progress:        pt,
	}
}

func (s *SGS) Fit(p types.Points) {
	s.condition = p
	s.isFitted = true
}

func (s *SGS) Interpolate(p types.Points) ([]types.Estimation, error) {
	// we want the spatial interpolator interface, thus Interpolate is a simulation with n=1
	estimations, err := s.Simulate(p, 1)
	if err != nil {
		return nil, err
	}
	return estimations[0], nil
}

func (s *SGS) Simulate(p types.Points, n int) ([][]types.Estimation, error) {
	if !s.isFitted {
		return nil, fmt.Errorf("the SGS needs first condition points")
	}
	if s.progress != nil {
		defer s.progress.close()
	}

	simulations := make([][]types.Estimation, n)
	errors := make([]error, n)

	wg := sync.WaitGroup{}

	for sim_idx := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// currently only sequential is implemented
			// maybe if we have more, we might want to use a wrapper here
			simulations[i], errors[i] = s.sequential(p, i)
		}(sim_idx)
	}

	wg.Wait()
	nerrs := 0
	for _, err := range errors {
		if err != nil {
			nerrs++
		}
	}
	if nerrs > 0 {
		fmt.Printf("Warning: %d simulations failed\n", nerrs)
	}

	return simulations, nil
}

type neighbor struct {
	point types.Point
	dist  float64
	idx   int
}

// findClosestNeighbors pre-calculates the closest neighbors for each simulation point
func (s *SGS) findClosestNeighbors(simPoints types.Points) [][]neighbor {
	closestNeighbors := make([][]neighbor, len(simPoints.Points))

	for i, point := range simPoints.Points {
		// Calculate distances to all condition points
		neighbors := make([]neighbor, len(s.condition.Points))
		for j, cond := range s.condition.Points {
			d := s.dist.Compute(&point, &cond)
			neighbors[j] = neighbor{
				point: cond,
				dist:  d,
				idx:   j,
			}
		}
		// Sort by distance
		sort.Slice(neighbors, func(i, j int) bool {
			return neighbors[i].dist < neighbors[j].dist
		})
		// Keep only maxPoints closest
		if len(neighbors) > s.maxPoints {
			neighbors = neighbors[:s.maxPoints]
		}
		closestNeighbors[i] = neighbors
	}
	return closestNeighbors
}

// getConditionPoints returns the condition points to use for kriging at a specific location
func (s *SGS) getConditionPoints(point types.Point, neighbors []neighbor, points []types.Point, pMask []bool) []types.Point {
	newCondition := make([]types.Point, 0, s.maxPoints)

	// If using neighbor optimization
	if s.useNeighbors && neighbors != nil {
		// Add original closest condition points
		for _, n := range neighbors {
			newCondition = append(newCondition, n.point)
		}

		// Add previously simulated points if they're closer than our current neighbors
		maxDist := neighbors[len(neighbors)-1].dist
		for i, m := range pMask {
			if m && !math.IsNaN(points[i].Value) {
				d := s.dist.Compute(&point, &points[i])
				if d < maxDist {
					newCondition = append(newCondition, points[i])
				}
			}
		}
	} else {
		// Traditional approach: use all condition points
		newCondition = make([]types.Point, len(s.condition.Points))
		copy(newCondition, s.condition.Points)

		// Add all valid simulated points
		for i, m := range pMask {
			if m && !math.IsNaN(points[i].Value) {
				newCondition = append(newCondition, points[i])
			}
		}
	}

	return newCondition
}

func (s *SGS) sequential(p types.Points, routineId int) ([]types.Estimation, error) {
	n := len(p.Points)
	pMask := make([]bool, n)

	// create a copy of the points
	points := make([]types.Point, n)
	copy(points, p.Points)

	estimations := make([]types.Estimation, n)

	// create a random index order for the points
	rand_idx := rand.Perm(n)

	var totalFitTime time.Duration
	var totalInterpTime time.Duration
	totalStart := time.Now()

	// Pre-calculate neighbors if optimization is enabled
	var closestNeighbors [][]neighbor
	if s.useNeighbors {
		closestNeighbors = s.findClosestNeighbors(p)
	}

	interpolator := s.newInterpolator()

	for it, idx := range rand_idx {
		// Get condition points for current location
		var neighbors []neighbor
		if s.useNeighbors {
			neighbors = closestNeighbors[idx]
		}
		newCondition := s.getConditionPoints(points[idx], neighbors, points, pMask)

		// Time the Fit operation
		fitStart := time.Now()

		// Fit using the selected condition points
		interpolator.Fit(types.Points{Points: newCondition, Is3D: p.Is3D})

		fitTime := time.Since(fitStart)
		totalFitTime += fitTime

		loc := types.Points{
			Points: []types.Point{points[idx]},
			Is3D:   p.Is3D,
		}

		// Time the Interpolate operation
		interpStart := time.Now()
		est, err := interpolator.Interpolate(loc)
		interpTime := time.Since(interpStart)
		totalInterpTime += interpTime
		if s.progress != nil {
			s.progress.send(Progress{
				id:         routineId,
				current:    it,
				total:      len(rand_idx),
				fitTime:    fitTime,
				interpTime: interpTime,
			})
		}

		if err != nil {
			// If kriging fails, set NaN and continue
			points[idx].Value = math.NaN()
			estimations[idx] = types.Estimation{
				Field: math.NaN(),
			}
			continue
		}

		// Check if the estimation was successful
		if est[0].ErrCode != types.ErrNone || math.IsNaN(est[0].Field) {
			points[idx].Value = math.NaN()
			estimations[idx] = types.Estimation{
				Field: math.NaN(),
			}
			continue
		}

		nrm := distuv.Normal{Mu: est[0].Field, Sigma: math.Sqrt(est[0].Variance)}
		simulation := nrm.Rand()

		points[idx].Value = simulation
		estimations[idx] = types.Estimation{
			Field: simulation,
		}
		pMask[idx] = true
	}

	s.profile.KFitTime = totalFitTime / time.Duration(len(rand_idx))
	s.profile.KTotalMeanTime = totalInterpTime / time.Duration(len(rand_idx))
	s.profile.STotalMeanTime = time.Since(totalStart)

	return estimations, nil
}

func (s *SGS) Profile() types.Profile {
	return s.profile
}
