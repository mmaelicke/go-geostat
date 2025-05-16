/*
Package csv provides functionality for reading and writing spatial data in CSV format.

The package supports reading CSV files with spatial coordinates (2D/3D), values,
and optionally time information. It provides flexible column mapping and robust
error handling for parsing various CSV formats.

# Example Datasets

The package includes example datasets for testing and demonstration:

Pancake Dataset (data/pancake.csv):
A 2D synthetic dataset with 301 points representing a smooth surface with values
ranging from 85 to 244. The dataset is suitable for demonstrating kriging and
variogram analysis. The data follows a CSV format with three columns:

	x: X-coordinate (range: 0-500)
	y: Y-coordinate (range: 0-500)
	value: Measured value (range: 85-244)

Basic Usage:

	import "github.com/mmaelicke/go-geostat/io/csv"

	// Read the pancake dataset with default parameters
	data, err := csv.ReadCSV("data/pancake.csv", "", "", "", "", "", "", false)
	if err != nil {
		log.Fatal(err)
	}

	// Access the points
	points := data.Read()

	// Sample random subset of points
	sample := data.Sample(100)

The PointData type implements the types.SpatialSample interface, providing methods
for accessing and sampling the data:

  - Length(): Returns the number of points
  - Sample(size int): Returns a random subset of points
  - Read(): Returns all points as types.Points

For writing results:

  - WriteVarioCSV: Writes variogram results
  - WriteKrigCSV: Writes kriging results
*/
package csv
