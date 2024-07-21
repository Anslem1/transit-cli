package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var skipPrompts bool

var executeCmd = &cobra.Command{
	Use:   "execute [transit name]",
	Short: "Execute the Transit",
	Long:  `Run a list of commands from a transit file with optional interactive confirmation.`,
	Run: func(cmd *cobra.Command, args []string) {
		var transitName string

		if len(args) == 0 {
			// List transits and prompt for selection only if no transit name provided
			var err error
			_, transitName, err = middleware.ListTransit("execute")
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
		commands, err := middleware.ReadCommandsInTransit(strings.Split(transitName, ".yaml")[0])
		if err != nil {
		log.SetFlags(0)
			log.Fatalf("Error reading transit '%s': %v", transitName, err)
		}

		middleware.ExecuteCommandsInTransit(commands, skipPrompts)
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolVarP(&skipPrompts, "skip", "s", false, "Skip confirmation prompts and execute all commands at once")
}
