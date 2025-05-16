package estimator

import (
	"math"
	"sync"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func mapByIndices(differences []float64, indices []int, numLags int, estimator types.Estimator) ([]float64, []bool) {
	var wg sync.WaitGroup

	results := make([]float64, numLags)
	mask := make([]bool, numLags)

	for lag := 0; lag < numLags; lag++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			diffs := make([]float64, 0, 100)
			for j, idx := range indices {
				if idx == i {
					diffs = append(diffs, differences[j])
				}
			}

			results[i] = estimator.Compute(diffs)
			mask[i] = math.IsNaN(results[i])

		}(lag)
	}
	wg.Wait()

	return results, mask

}
