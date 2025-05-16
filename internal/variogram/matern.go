package variogram

import (
	"math"

	"github.com/mmaelicke/go-geostat/internal/types"
)

type Matern struct {
	types.BaseParams
	profile types.Profile
	// Smoothness parameter (nu) controls the differentiability of the process
	// Common values are 0.5 (exponential), 1.5, 2.5
	Nu float64
}

func (m *Matern) Name() string {
	return "matern"
}

func (m *Matern) Range() float64 {
	return m.BaseParams.Range
}

func (m *Matern) Sill() float64 {
	return m.BaseParams.Sill
}

func (m *Matern) Nugget() float64 {
	return m.BaseParams.Nugget
}

func (m *Matern) Evaluate(h float64) float64 {
	r := m.Range()
	sill := m.Sill()
	nugget := m.Nugget()
	nu := m.Nu

	if h <= 0 {
		return nugget
	}

	// Matern model formula:
	// γ(h) = c₀ + c₁[1 - (2^(1-ν)/Γ(ν)) * (h/a)^ν * K_ν(h/a)]
	// where:
	// c₀ is nugget
	// c₁ is sill
	// ν (nu) is the smoothness parameter
	// a is the range
	// K_ν is the modified Bessel function of the second kind of order ν
	// Γ is the gamma function

	// For numerical stability, we use a simplified form when h is very small
	if h < 1e-10 {
		return nugget
	}

	h_r := h / r
	term := math.Pow(2, 1-nu) / math.Gamma(nu)
	term *= math.Pow(h_r, nu)
	term *= BesselK(nu, h_r)

	return nugget + sill*(1-term)
}

func (m *Matern) Map(h []float64) []float64 {
	variances := make([]float64, len(h))
	for i, h_i := range h {
		variances[i] = m.Evaluate(h_i)
	}
	return variances
}

func (m *Matern) Profile() types.Profile {
	return m.profile
}

func (m *Matern) SetProfile(p types.Profile) {
	m.profile = p
}
