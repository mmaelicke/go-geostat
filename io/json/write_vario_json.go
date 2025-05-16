package json

import (
	"encoding/json"
	"io"

	"github.com/mmaelicke/go-geostat/internal/types"
	"github.com/mmaelicke/go-geostat/internal/variogram"
)

type param struct {
	Range  float64 `json:"range"`
	Sill   float64 `json:"sill"`
	Nugget float64 `json:"nugget"`
	Nu     float64 `json:"nu,omitempty"`
	Name   string  `json:"name,omitempty"`
}

type varioJson struct {
	Edges         []float64 `json:"edges"`
	Histogram     []int     `json:"histogram"`
	Semivariances []float64 `json:"semivariances"`
	Params        param     `json:"params,omitempty"`
}

func WriteVarioJsonToWriter(w io.Writer, v types.SampleVariogram, m types.SpatialFunction) error {
	vario := varioJson{
		Edges:         v.GetEdges(),
		Histogram:     v.GetHistogram(),
		Semivariances: v.GetSemivariances(),
	}
	if m != nil {
		vario.Params = param{
			Range:  m.Range(),
			Sill:   m.Sill(),
			Nugget: m.Nugget(),
			Name:   m.Name(),
		}
		if m.Name() == "matern" {
			vario.Params.Nu = m.(*variogram.Matern).Nu
		}
	}

	err := json.NewEncoder(w).Encode(vario)
	if err != nil {
		return err
	}

	return nil
}
