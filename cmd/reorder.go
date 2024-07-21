package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var (
	reorderCmd = &cobra.Command{
		Use:   "reorder <transit name>",
		Short: "Reorder commands in a transit configuration",
		Long: `
Reorder commands in an existing transit configuration.

Use this command to change the order of commands in a specific transit.
`,
		Run: func(cmd *cobra.Command, args []string) {
			var transitName string
			if len(args) == 0 {
				var err error
				_, transitName, err = middleware.ListTransit("reorder")
				if err != nil {
					log.SetFlags(0)
					log.Fatalf("Error fetching list of transits: %v", err)
				}
			} else {
				transitName = strings.Split(args[0], ".yaml")[0]
			}
			if transitName == "" {
				fmt.Println("No transit selected. Exiting.")
				return
			}
			
			transitName = strings.Split(transitName, ".yaml")[0]
			err := middleware.ReorderCommandsInTransit(transitName)
			if err != nil {
				fmt.Printf("Failed to reorder commands in transit %s: %v\n", transitName, err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(reorderCmd)
}
