package variogram

import "math"

// BesselK returns the modified Bessel function of the second kind of order nu
// This is an approximation that works well for nu = 0.5, 1.5, 2.5
// For other values, a more sophisticated implementation would be needed
func BesselK(nu, x float64) float64 {
	// For small x, use asymptotic expansion
	if x < 1e-10 {
		return math.Inf(1)
	}

	// For large x, use asymptotic expansion
	if x > 100 {
		return math.Sqrt(math.Pi/(2*x)) * math.Exp(-x)
	}

	// For nu = 0.5, we have a closed form
	if math.Abs(nu-0.5) < 1e-10 {
		return math.Sqrt(math.Pi/(2*x)) * math.Exp(-x)
	}

	// For nu = 1.5, we have a closed form
	if math.Abs(nu-1.5) < 1e-10 {
		return math.Sqrt(math.Pi/(2*x)) * (1 + 1/x) * math.Exp(-x)
	}

	// For nu = 2.5, we have a closed form
	if math.Abs(nu-2.5) < 1e-10 {
		return math.Sqrt(math.Pi/(2*x)) * (1 + 3/x + 3/(x*x)) * math.Exp(-x)
	}

	// For other values, use a simple approximation
	// This is not as accurate but should work for most practical purposes
	return math.Sqrt(math.Pi/(2*x)) * math.Exp(-x) * (1 + nu/x)
}
