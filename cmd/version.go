package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)


var verbose bool

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of transit-cli",
    Long:  `Print the version number of transit-cli.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Transit version 1.0.0")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
