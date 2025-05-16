package csv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReadCopperData(t *testing.T) {
	// Get the workspace root directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}
	// Move up one level to reach workspace root
	rootDir := filepath.Dir(filepath.Dir(wd))
	// Read the copper data from meuse.txt
	data, err := ReadCSV(rootDir+"/data/meuse.txt", "x", "y", "", "", "copper", "", false)
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
