package variogram

import "github.com/mmaelicke/go-geostat/internal/types"

type Spherical struct {
	types.BaseParams
	profile types.Profile
}

func (s *Spherical) Name() string {
	return "spherical"
}

func (s *Spherical) Range() float64 {
	return s.BaseParams.Range
}

func (s *Spherical) Sill() float64 {
	return s.BaseParams.Sill
}

func (s *Spherical) Nugget() float64 {
	return s.BaseParams.Nugget
}

func (s *Spherical) Evaluate(h float64) float64 {
	r := s.Range()
	sill := s.Sill()
	nugget := s.Nugget()
	if h <= r {
		h_r := h / r
		return nugget + sill*(1.5*h_r-0.5*h_r*h_r*h_r)
	}
	return nugget + sill
}

func (s *Spherical) Map(h []float64) []float64 {
	variances := make([]float64, len(h))
	for i, h_i := range h {
		variances[i] = s.Evaluate(h_i)
	}
	return variances
}

func (s *Spherical) Profile() types.Profile {
	return s.profile
}

func (s *Spherical) SetProfile(p types.Profile) {
	s.profile = p
}
