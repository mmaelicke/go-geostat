package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mmaelicke/go-geostat/internal/types"
)

// PointData holds a collection of spatial points
type PointData struct {
	Points []types.Point
	Is3D   bool
}

func (p PointData) Length() int {
	return len(p.Points)
}

func (p PointData) Read() types.Points {
	l := len(p.Points)
	refs := make([]types.Point, l)
	for i, p := range p.Points {
		refs[i] = p
	}
	return types.Points{
		Points: refs,
		Is3D:   p.Is3D,
	}
}

func (p PointData) Sample(size int) types.Points {
	l := len(p.Points)
	if size > l {
		return p.Read()
	}

	indices := make([]int, l)
	for i := range indices {
		indices[i] = i
	}

	// Shuffle the indices
	rand.Shuffle(l, func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// Take the first 'size' indices
	sample := make([]types.Point, size)
	for i := 0; i < size; i++ {
		sample[i] = p.Points[indices[i]]
	}

	return types.Points{
		Points: sample,
		Is3D:   p.Is3D,
	}
}

func ReadCSVFromReader(reader io.Reader, xCol, yCol, zCol, tCol, valueCol, timeFormat string, errorOnParse bool) (PointData, error) {
	csvReader := csv.NewReader(reader)

	header, err := csvReader.Read()
	if err != nil {
		return PointData{}, fmt.Errorf("failed to read header: %w", err)
	}

	if xCol == "" {
		xCol = "x"
	}
	if yCol == "" {
		yCol = "y"
	}
	if zCol == "" {
		zCol = "z"
	}
	if tCol == "" {
		tCol = "time"
	}
	if valueCol == "" {
		valueCol = "value"
	}
	if timeFormat == "" {
		timeFormat = "2006-01-02 15:04:05"
	}

	xIdx := -1
	yIdx := -1
	zIdx := -1
	tIdx := -1
	valueIdx := -1

	for i, col := range header {
		switch col {
		case xCol:
			xIdx = i
		case yCol:
			yIdx = i
		case zCol:
			zIdx = i
		case tCol:
			tIdx = i
		case valueCol:
			valueIdx = i
		}
	}

	if xIdx == -1 || yIdx == -1 || valueIdx == -1 {
		return PointData{}, fmt.Errorf("missing required columns. You need to specify at least x, y and value columns")
	}

	data := PointData{
		Points: make([]types.Point, 0),
		Is3D:   zIdx != -1,
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return PointData{}, fmt.Errorf("failed to read record: %w", err)
		}

		if strings.HasPrefix(record[0], "#") {
			continue
		}

		point := types.Point{}

		// Parse coordinates
		if x, err := strconv.ParseFloat(record[xIdx], 64); err == nil {
			point.X = x
		} else {
			if errorOnParse {
				return PointData{}, fmt.Errorf("failed to parse x: %w", err)
			}
			continue
		}

		if y, err := strconv.ParseFloat(record[yIdx], 64); err == nil {
			point.Y = y
		} else {
			if errorOnParse {
				return PointData{}, fmt.Errorf("failed to parse y: %w", err)
			}
			continue
		}

		if value, err := strconv.ParseFloat(record[valueIdx], 64); err == nil {
			point.Value = value
		} else {
			if errorOnParse {
				return PointData{}, fmt.Errorf("failed to parse value: %w", err)
			}
			continue
		}

		if data.Is3D {
			if z, err := strconv.ParseFloat(record[zIdx], 64); err == nil {
				point.Z = z
				point.Is3D = true
			} else {
				if errorOnParse {
					return PointData{}, fmt.Errorf("failed to parse z: %w", err)
				}
				continue
			}
		}

		if tIdx != -1 {
			if t, err := time.Parse(timeFormat, record[tIdx]); err == nil {
				point.Time = t
				point.HasTime = true
			} else {
				if errorOnParse {
					return PointData{}, fmt.Errorf("failed to parse time: %w", err)
				}
				continue
			}
		}

		data.Points = append(data.Points, point)
	}

	return data, nil
}

func ReadCSV(path, xCol, yCol, zCol, tCol, valueCol, timeFormat string, errorOnParse bool) (PointData, error) {
	file, err := os.Open(path)
	if err != nil {
		return PointData{}, err
	}
	defer file.Close()
	return ReadCSVFromReader(file, xCol, yCol, zCol, tCol, valueCol, timeFormat, errorOnParse)
}
