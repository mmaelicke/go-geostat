package csv

import (
	"testing"
)

func TestReadCopperData(t *testing.T) {
	// Read the copper data from meuse.txt
	data, err := ReadCSV("meuse.txt", "x", "y", "", "", "copper", "", false)
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	// The meuse dataset has 155 samples
	expectedLength := 155
	if len(data.Points) != expectedLength {
		t.Errorf("Expected %d samples, got %d", expectedLength, len(data.Points))
	}

	// Check if we have valid data
	if len(data.Points) > 0 {
		// Check first point
		firstPoint := data.Points[0]
		if firstPoint.X == 0 || firstPoint.Y == 0 || firstPoint.Value == 0 {
			t.Error("First point has zero values")
		}

		// Check last point
		lastPoint := data.Points[len(data.Points)-1]
		if lastPoint.X == 0 || lastPoint.Y == 0 || lastPoint.Value == 0 {
			t.Error("Last point has zero values")
		}

		// Check that all points have valid coordinates and values
		for i, p := range data.Points {
			if p.X == 0 || p.Y == 0 || p.Value == 0 {
				t.Errorf("Point %d has zero values: X=%f, Y=%f, Value=%f", i, p.X, p.Y, p.Value)
			}
		}
	}

	// Verify that we're not in 3D mode
	if data.Is3D {
		t.Error("Data should not be in 3D mode")
	}
}
