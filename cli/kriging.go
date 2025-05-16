package cli

import "github.com/spf13/cobra"

func init() {
	krigingCmd := &cobra.Command{
		Use:   "krig",
		Short: "Kriging implementation",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	rootCmd.AddCommand(krigingCmd)
}
