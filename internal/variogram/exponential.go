package variogram

import (
	"math"

	"github.com/mmaelicke/go-geostat/internal/types"
)

type Exponential struct {
	types.BaseParams
	profile types.Profile
}

func (e *Exponential) Name() string {
	return "exponential"
}

func (e *Exponential) Range() float64 {
	return e.BaseParams.Range
}

func (e *Exponential) Sill() float64 {
	return e.BaseParams.Sill
}

func (e *Exponential) Nugget() float64 {
	return e.BaseParams.Nugget
}

func (e *Exponential) Evaluate(h float64) float64 {
	r := e.Range()
	sill := e.Sill()
	nugget := e.Nugget()
	a := r / 3.0
	return nugget + sill*(1.0-math.Exp(-h/a))
}

func (e *Exponential) Map(h []float64) []float64 {
	variances := make([]float64, len(h))
	for i, h_i := range h {
		variances[i] = e.Evaluate(h_i)
	}
	return variances
}

func (e *Exponential) Profile() types.Profile {
	return e.profile
}

func (e *Exponential) SetProfile(p types.Profile) {
	e.profile = p
}
