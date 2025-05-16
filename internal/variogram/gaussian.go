package variogram

import (
	"math"

	"github.com/mmaelicke/go-geostat/internal/types"
)

type Gaussian struct {
	types.BaseParams
	profile types.Profile
}

func (g *Gaussian) Name() string {
	return "gaussian"
}

func (g *Gaussian) Range() float64 {
	return g.BaseParams.Range
}

func (g *Gaussian) Sill() float64 {
	return g.BaseParams.Sill
}

func (g *Gaussian) Nugget() float64 {
	return g.BaseParams.Nugget
}

func (g *Gaussian) Evaluate(h float64) float64 {
	r := g.Range()
	sill := g.Sill()
	nugget := g.Nugget()
	a := r / 2.0
	return nugget + sill*(1.0-math.Exp(-(h*h)/(a*a)))
}

func (g *Gaussian) Map(h []float64) []float64 {
	variances := make([]float64, len(h))
	for i, h_i := range h {
		variances[i] = g.Evaluate(h_i)
	}
	return variances
}

func (g *Gaussian) Profile() types.Profile {
	return g.profile
}

func (g *Gaussian) SetProfile(p types.Profile) {
	g.profile = p
}
