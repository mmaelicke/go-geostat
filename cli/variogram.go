package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mmaelicke/go-geostat/internal/distance"
	"github.com/mmaelicke/go-geostat/internal/empirical"
	"github.com/mmaelicke/go-geostat/internal/estimator"
	"github.com/mmaelicke/go-geostat/internal/kriging"
	"github.com/mmaelicke/go-geostat/internal/sgs"
	"github.com/mmaelicke/go-geostat/internal/types"
	"github.com/mmaelicke/go-geostat/io/asc"
	"github.com/mmaelicke/go-geostat/io/csv"
	"github.com/spf13/cobra"
)

// Config holds all configuration options for the variogram command
type Config struct {
	// Input/Output options
	CSVPath      string
	OutputPath   string
	OutputFormat string

	// Column specifications
	XCol     string
	YCol     string
	ZCol     string
	TCol     string
	ValueCol string

	// Variogram parameters
	NLags  int
	MaxLag float64

	// Model parameters
	ModelName     string
	DistType      string
	EstimatorName string
	TimeFormat    string

	// Processing options
	MaxPoints int
	DX        float64
	DY        float64
	DZ        float64

	// Flags
	Performance bool
	Fit         bool
	UseKriging  bool
	KrigingOnly bool
	UseSGS      bool
	SGSOnly     bool
	SGSSimCount int
}

// newDefaultConfig returns a Config with default values
func newDefaultConfig() *Config {
	return &Config{
		NLags:         10,
		MaxLag:        0,
		MaxPoints:     100,
		DX:            1.0,
		DY:            1.0,
		DZ:            1.0,
		SGSSimCount:   1,
		ModelName:     "spherical",
		DistType:      "euclidean",
		EstimatorName: "matheron",
	}
}

func init() {
	config := newDefaultConfig()

	varioCmd := &cobra.Command{
		Use:   "vario",
		Short: "Compute empirical variogram and fit model",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runVariogram(config); err != nil {
				log.Fatalf("Error running variogram: %v", err)
			}
		},
	}

	// Input/Output flags
	varioCmd.Flags().StringVar(&config.CSVPath, "csv", "", "Path to input CSV file")
	varioCmd.Flags().StringVar(&config.OutputPath, "output", "", "Path to output file")
	varioCmd.Flags().StringVar(&config.OutputFormat, "format", "json", "Output format (json, csv)")

	// Column specification flags
	varioCmd.Flags().StringVar(&config.XCol, "x", "x", "X coordinate column name")
	varioCmd.Flags().StringVar(&config.YCol, "y", "y", "Y coordinate column name")
	varioCmd.Flags().StringVar(&config.ZCol, "z", "", "Z coordinate column name")
	varioCmd.Flags().StringVar(&config.TCol, "t", "", "Time column name")
	varioCmd.Flags().StringVar(&config.ValueCol, "value", "value", "Value column name")

	// Variogram parameter flags
	varioCmd.Flags().IntVar(&config.NLags, "nlags", 10, "Number of lags")
	varioCmd.Flags().Float64Var(&config.MaxLag, "maxlag", 0, "Maximum lag distance")

	// Model parameter flags
	varioCmd.Flags().StringVar(&config.ModelName, "model", "spherical", "Variogram model type")
	varioCmd.Flags().StringVar(&config.DistType, "dist", "euclidean", "Distance metric")
	varioCmd.Flags().StringVar(&config.EstimatorName, "estimator", "matheron", "Variogram estimator")
	varioCmd.Flags().StringVar(&config.TimeFormat, "timeformat", "", "Time format string")

	// Processing option flags
	varioCmd.Flags().IntVar(&config.MaxPoints, "maxpoints", 100, "Maximum number of points to use")
	varioCmd.Flags().Float64Var(&config.DX, "dx", 1.0, "X grid spacing")
	varioCmd.Flags().Float64Var(&config.DY, "dy", 1.0, "Y grid spacing")
	varioCmd.Flags().Float64Var(&config.DZ, "dz", 1.0, "Z grid spacing")

	// Feature flags
	varioCmd.Flags().BoolVar(&config.Performance, "perf", false, "Enable performance profiling")
	varioCmd.Flags().BoolVar(&config.Fit, "fit", false, "Fit variogram model")
	varioCmd.Flags().BoolVar(&config.UseKriging, "krig", false, "Perform kriging")
	varioCmd.Flags().BoolVar(&config.KrigingOnly, "krigonly", false, "Only perform kriging")
	varioCmd.Flags().BoolVar(&config.UseSGS, "sgs", false, "Perform sequential Gaussian simulation")
	varioCmd.Flags().BoolVar(&config.SGSOnly, "sgsonly", false, "Only perform sequential Gaussian simulation")
	varioCmd.Flags().IntVar(&config.SGSSimCount, "nsim", 1, "Number of SGS simulations")

	rootCmd.AddCommand(varioCmd)
}

func runVariogram(config *Config) error {
	var data csv.PointData
	var err error

	if config.UseKriging && config.UseSGS {
		return fmt.Errorf("kriging and SGS cannot be performed at the same time")
	}

	if config.CSVPath != "" {
		data, err = csv.ReadCSV(config.CSVPath, config.XCol, config.YCol, config.ZCol,
			config.TCol, config.ValueCol, config.TimeFormat, false)
	} else {
		data, err = csv.ReadCSVFromReader(os.Stdin, config.XCol, config.YCol, config.ZCol,
			config.TCol, config.ValueCol, config.TimeFormat, false)
	}
	if err != nil {
		return fmt.Errorf("error reading CSV: %v", err)
	}

	points := data.Read()
	if config.MaxLag == 0 {
		config.MaxLag = 1e6
	}

	var dist types.Distance
	var est types.Estimator

	switch strings.ToLower(config.DistType) {
	case "chebyshev":
		dist = &distance.ChebyshevDistance{}
	case "manhattan":
		dist = &distance.ManhattanDistance{}
	case "euclidean":
		dist = &distance.EuclideanDistance{}
	default:
		log.Fatalf("Unsupported distance type: %s", config.DistType)
	}

	switch strings.ToLower(config.EstimatorName) {
	case "matheron":
		est = &estimator.Matheron{}
	case "cressie":
		est = &estimator.Cressie{}
	default:
		log.Fatalf("Unsupported estimator: %s", config.EstimatorName)
	}

	vg := empirical.NewEmpiricalVariogram(points, config.NLags, config.MaxLag, dist, est)
	err = vg.Compute()
	if err != nil {
		log.Fatalf("Error computing empirical variogram: %v", err)
	}

	var model types.SpatialFunction
	if config.Fit || config.UseKriging {
		model, err = vg.Fit(config.ModelName)
		if err != nil {
			log.Fatalf("Error fitting model: %v", err)
		}
	}

	if config.Performance {
		profile := vg.GetProfile()
		fmt.Println("# Variogram estimation runtime:")
		fmt.Printf("# Pairwise time:      %v\n", profile.PairwiseTime)
		fmt.Printf("# Binning time:       %v\n", profile.BinningTime)
		fmt.Printf("# Histogram time:     %v\n", profile.HistogramTime)
		fmt.Printf("# Semivariogram time: %v\n", profile.SemivarTime)
		fmt.Printf("# Sample time:        %v\n", profile.EmpiricalTime)
		if config.Fit {
			prof := model.Profile()
			fmt.Printf("# Initial guess:     %v\n", prof.InitialGuess)
			fmt.Printf("# Fit time:          %v\n", prof.FitTime)
		}
		fmt.Printf("# Total time:         %v\n", profile.TotalTime)
	}

	var estimation []types.Estimation
	var grid types.Points
	if config.UseKriging {
		grid, err = kriging.DenseGrid(points, config.DX, config.DY, config.DZ)
		if err != nil {
			log.Fatalf("Error creating dense grid: %v", err)
		}

		kr := kriging.New(model, config.MaxPoints, dist, false)
		kr.Fit(points)
		estimation, err = kr.Interpolate(grid)
		if err != nil {
			log.Fatalf("Error interpolating: %v", err)
		}

		if config.Performance {
			prof := kr.Profile()
			fmt.Println("# Kriging runtime:")
			fmt.Printf("# Fit time:          %v\n", prof.FitTime)
			fmt.Printf("# K-Init time:       %v\n", prof.KInitMeanTime)
			fmt.Printf("# K-Matrix time:     %v\n", prof.KMatMeanTime)
			fmt.Printf("# K-Solving time:    %v\n", prof.KSolvMeanTime)
			fmt.Printf("# K-Total time:      %v\n", prof.KTotalMeanTime)
		}
	}

	var sims [][]types.Estimation
	if config.UseSGS {
		grid, err = kriging.DenseGrid(points, config.DX, config.DY, config.DZ)
		if err != nil {
			log.Fatalf("Error creating dense grid: %v", err)
		}

		s := sgs.New(model, config.MaxPoints, dist, true)
		s.Fit(points)

		sims, err = s.Simulate(grid, config.SGSSimCount)
		if err != nil {
			log.Fatalf("Error simulating: %v", err)
		}
	}

	if !config.KrigingOnly && !config.SGSOnly {
		if config.OutputPath != "" {
			err = csv.WriteVarioCSV(config.OutputPath+"_variogram.csv", vg, model)
			if err != nil {
				log.Fatalf("Error writing output: %v", err)
			}
		} else {
			csv.WriteVarioCSVToWriter(os.Stdout, vg, model)
		}
	}
	if config.UseKriging {
		field := make([]float64, len(estimation))
		variance := make([]float64, len(estimation))
		for i, e := range estimation {
			field[i] = e.Field
			variance[i] = e.Variance
		}

		if config.OutputPath != "" {
			if config.OutputFormat == "asc" {
				asc.WriteKrigAsc(config.OutputPath+"_krig_field.asc", grid, field)
				asc.WriteKrigAsc(config.OutputPath+"_krig_variance.asc", grid, variance)
			}
			if config.OutputFormat == "csv" {
				csv.WriteKrigCSV(config.OutputPath+"_krig.csv", grid, estimation)
			}
		} else {
			if config.OutputFormat == "asc" {
				asc.WriteKrigAscToWriter(os.Stdout, grid, field)
			} else if config.OutputFormat == "csv" {
				csv.WriteKrigCSVToWriter(os.Stdout, grid, estimation)
			}
		}
	}
	if config.UseSGS {
		for sim_idx, sim := range sims {
			field := make([]float64, len(sim))
			for i, e := range sim {
				field[i] = e.Field
			}

			if config.OutputPath != "" {
				if config.OutputFormat == "asc" {
					asc.WriteKrigAsc(config.OutputPath+"_sgs_sim_"+strconv.Itoa(sim_idx)+".asc", grid, field)
				}
				if config.OutputFormat == "csv" {
					csv.WriteKrigCSV(config.OutputPath+"_sgs_sim_"+strconv.Itoa(sim_idx)+".csv", grid, sim)
				}
			} else {
				if config.OutputFormat == "asc" {
					fmt.Printf("--- Simulation %d ---\n", sim_idx)
					asc.WriteKrigAscToWriter(os.Stdout, grid, field)
				} else if config.OutputFormat == "csv" {
					fmt.Printf("--- Simulation %d ---\n", sim_idx)
					csv.WriteKrigCSVToWriter(os.Stdout, grid, sim)
				}
			}
		}
	}
	return nil
}
