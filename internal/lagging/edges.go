package lagging

import (
	"fmt"
	"math"
)

func CalculateEdges(distances []float64, numLags int, maxLag float64) ([]float64, error) {
	edges := make([]float64, numLags)

	max := math.Inf(-1)

	for _, distance := range distances {
		if distance > max {
			max = distance
		}
	}
	max = math.Min(max, maxLag)
	if max == math.Inf(-1) {
		return nil, fmt.Errorf("no edges could be calculated: max distance is infinite")
	}

	step := max / float64(numLags)

	for i := range edges {
		edges[i] = float64(i+1) * step
	}

	return edges, nil
}

func GetEdgeIndex(distances []float64, edges []float64) []int {
	indices := make([]int, len(distances))
	for i, distance := range distances {
		indices[i] = -1 // default to -1 (not assigned)
		lower := 0.0
		for j, upper := range edges {
			if distance >= lower && distance < upper {
				indices[i] = j
				break
			}
			lower = upper
		}
	}
	return indices
}
