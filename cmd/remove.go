package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove <transit name>",
		Short: "Removes command(s) from a transit",
		Long: `
Remove one or more commands from an existing transit.

Use this command to remove specific commands from a transit configuration.
You will be prompted to select commands by their number to remove interactively.`,
		Run: func(cmd *cobra.Command, args []string) {
			var transitName string
			if len(args) == 0 {
				var err error
				_, transitName, err = middleware.ListTransit("remove")
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
			// Fetch existing commands from the transit
			transitName = strings.Split(transitName, ".yaml")[0]
			commands, err := middleware.ReadCommandsInTransit(transitName)
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Failed to fetch commands from transit %s: %v", transitName, err)
				return
			}

			if len(commands) == 0 {
				fmt.Printf("Transit '%s' does not contain any commands.\n", transitName)
				return
			}

			// Prompt user to select commands to remove
			selectedCommands, err := middleware.GetUserSelectedCommands(commands)
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Failed to select commands for removal: %v", err)
				return
			}

			if len(selectedCommands) == 0 {
				fmt.Println("No commands selected for removal. Exiting.")
				return
			}

			// Display the commands that are set to be deleted
			fmt.Printf("Commands set to be removed from transit '%s':\n", transitName)
			for _, cmd := range selectedCommands {
				fmt.Printf("- %s\n", cmd)
			}

			// Remove selected commands from the transit
			transitName = strings.Split(transitName, ".yaml")[0]
			err = middleware.RemoveCommandsFromTransit(transitName, selectedCommands)
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Failed to remove commands from transit %s: %v", transitName, err)
				return
			}

			fmt.Printf("Removed %d command(s) from transit '%s'\n", len(selectedCommands), transitName)
		},
	}
)

func init() {
	rootCmd.AddCommand(removeCmd)
}
