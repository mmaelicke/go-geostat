package estimator

import (
	"math"
)

type Matheron struct{}

func (m *Matheron) Compute(differences []float64) float64 {
	n := len(differences)
	if n == 0 {
		return math.NaN()
	}
	sum := 0.0
	for _, diff := range differences {
		sum += math.Pow(diff, 2)
	}
	sum /= (2 * float64(n))

	return sum
}

func (m *Matheron) Map(differences []float64, indices []int, numLags int) ([]float64, []bool) {
	return mapByIndices(differences, indices, numLags, m)
}

type Cressie struct{}

func (c *Cressie) Compute(differences []float64) float64 {
	n := len(differences)
	if n == 0 {
		return math.NaN()
	}
	sum := 0.0
	for _, diff := range differences {
		sum += math.Pow(math.Abs(diff), 0.5)
	}
	mean := sum / float64(n)
	numerator := math.Pow(mean, 4)
	denominator := 0.457 + 0.494/float64(n) + 0.045/(float64(n)*float64(n))
	return numerator / (2 * denominator)
}

func (c *Cressie) Map(differences []float64, indices []int, numLags int) ([]float64, []bool) {
	return mapByIndices(differences, indices, numLags, c)
}
