package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of transit-cli, injected at build time via ldflags.
var Version = "1.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of transit-cli",
	Long:  `Print the version number of transit-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Transit version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
