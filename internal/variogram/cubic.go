package variogram

import "github.com/mmaelicke/go-geostat/internal/types"

type Cubic struct {
	types.BaseParams
	profile types.Profile
}

func (c *Cubic) Name() string {
	return "cubic"
}

func (c *Cubic) Range() float64 {
	return c.BaseParams.Range
}

func (c *Cubic) Sill() float64 {
	return c.BaseParams.Sill
}

func (c *Cubic) Nugget() float64 {
	return c.BaseParams.Nugget
}

func (c *Cubic) Evaluate(h float64) float64 {
	r := c.Range()
	sill := c.Sill()
	nugget := c.Nugget()

	if h <= r {
		h_r := h / r
		// Cubic model formula: γ(h) = c₀ + c₁[7(h/a)² - 8.75(h/a)³ + 3.5(h/a)⁵ - 0.75(h/a)⁷]
		// where c₀ is nugget and c₁ is sill
		h_r2 := h_r * h_r
		h_r3 := h_r2 * h_r
		h_r5 := h_r3 * h_r2
		h_r7 := h_r5 * h_r2

		return nugget + sill*(7*h_r2-(35.0/4.0)*h_r3+(7.0/2.0)*h_r5-(3.0/4.0)*h_r7)
	}
	return nugget + sill
}

func (c *Cubic) Map(h []float64) []float64 {
	variances := make([]float64, len(h))
	for i, h_i := range h {
		variances[i] = c.Evaluate(h_i)
	}
	return variances
}

func (c *Cubic) Profile() types.Profile {
	return c.profile
}

func (c *Cubic) SetProfile(p types.Profile) {
	c.profile = p
}
