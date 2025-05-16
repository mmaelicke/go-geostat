package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func WriteVarioCSVToWriter(w io.Writer, v types.SampleVariogram, m types.SpatialFunction) error {
	csvw := csv.NewWriter(w)

	if m != nil {
		metadata := fmt.Sprintf("# model: %s, range: %f, sill: %f, nugget: %f\n", m.Name(), m.Range(), m.Sill(), m.Nugget())
		w.Write([]byte(metadata))
	}

	header := []string{"lag", "count", "upper_edge", "semivariance"}
	if m != nil {
		header = append(header, "model")
	}
	csvw.Write(header)

	edges := v.GetEdges()
	semivariances := v.GetSemivariances()
	histogram := v.GetHistogram()

	if len(edges) != len(semivariances) || len(edges) != len(histogram) {
		return fmt.Errorf("edges, semivariances, and histogram must have the same length")
	}

	for e := range edges {
		row := []string{
			fmt.Sprintf("%d", e),
			fmt.Sprintf("%d", histogram[e]),
			fmt.Sprintf("%f", edges[e]),
			fmt.Sprintf("%f", semivariances[e]),
		}
		if m != nil {
			row = append(row, fmt.Sprintf("%f", m.Evaluate(edges[e])))
		}
		csvw.Write(row)
	}
	csvw.Flush()
	return nil
}

func WriteVarioCSV(path string, v types.SampleVariogram, m types.SpatialFunction) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	return WriteVarioCSVToWriter(f, v, m)
}
