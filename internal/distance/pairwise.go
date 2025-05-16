package distance

import "github.com/mmaelicke/go-geostat/internal/types"

func PairwiseDistances(points []types.Point, dist types.Distance, withDifferences bool) ([]float64, []float64) {
	n := len(points)
	result := make([]float64, 0, n*(n-1)/2)
	differences := make([]float64, 0, n*(n-1)/2)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result = append(result, dist.Compute(&points[i], &points[j]))
			if withDifferences {
				differences = append(differences, points[i].Value-points[j].Value)
			}
		}
	}
	return result, differences
}
