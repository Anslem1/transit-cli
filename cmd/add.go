package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use:   "add <transit name>",
		Short: "Adds new commands to a transit configuration",
		Long: `
			adds one or more commands to an existing transit.

	Use this command to add commands to a specific transit for later execution. 
	You will be prompted to enter commands interactively.
`,
		Run: func(cmd *cobra.Command, args []string) {
			var transitName string
			if len(args) == 0 {
				var err error
				_, transitName, err = middleware.ListTransit("add")
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
			_, err := middleware.ReadCommandsInTransit(strings.Split(transitName, ".yaml")[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			// Prompt for initial command(s) to add
			var initialCommands []string
			if len(args) > 1 {
				initialCommands = args[1:]
			}

			transitName = strings.Split(transitName, ".yaml")[0]
			err = middleware.AddCommandsToTransit(transitName, initialCommands)
			if err != nil {
				fmt.Printf("Failed to add commands to transit %s: %v\n", transitName, err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
