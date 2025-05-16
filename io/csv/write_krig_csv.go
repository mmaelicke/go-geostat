package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func WriteKrigCSVToWriter(w io.Writer, gridList types.Points, estimation []types.Estimation) error {
	csvw := csv.NewWriter(w)

	header := []string{"x", "y"}
	if gridList.Is3D {
		header = append(header, "z")
	}
	header = append(header, "value", "variance")
	csvw.Write(header)

	for i, p := range gridList.Points {
		row := []string{fmt.Sprintf("%f", p.X), fmt.Sprintf("%f", p.Y)}
		if gridList.Is3D {
			row = append(row, fmt.Sprintf("%f", p.Z))
		}
		row = append(row, fmt.Sprintf("%f", estimation[i].Field), fmt.Sprintf("%f", estimation[i].Variance))
		csvw.Write(row)
	}
	csvw.Flush()
	return csvw.Error()
}

func WriteKrigCSV(path string, gridList types.Points, estimation []types.Estimation) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	return WriteKrigCSVToWriter(f, gridList, estimation)
}
