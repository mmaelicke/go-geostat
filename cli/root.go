package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-geostat",
	Short: "go-geostat is a tool for geostatistical analysis",
	Long: `go-geostat is a tool for geostatistical analysis.
	
	go-geostat is a tool for geostatistical analysis. It provides a command line interface to perform variogram analysis, kriging, and other geostatistical operations.
	
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
