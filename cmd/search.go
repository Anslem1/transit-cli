package cmd

import (
	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [shuttle name]",
	Short: "Search for commands in shuttles",
	Long:  "Search for commands across all shuttles or within a specific shuttle.",
	Run: func(cmd *cobra.Command, args []string) {
		middleware.SearchTransit(args)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	// rootCmd.fl
}
