package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:               "remove [transit name]",
	Short:             "Removes command(s) from a transit",
	Long:              `Remove a command from an existing transit interactively.`,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		var tr *transit.Transit

		if len(args) == 0 {
			var err error
			tr, err = ui.SelectTransit("remove command from")
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error: %v", err)
			}
		} else {
			transitName := strings.TrimSuffix(args[0], ".yaml")
			var err error
			tr, err = transit.GetTransit(transitName)
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error reading transit '%s': %v", transitName, err)
			}
		}

		if tr == nil {
			fmt.Println("No transit selected. Exiting.")
			return
		}

		if len(tr.Commands) == 0 {
			fmt.Printf("Transit '%s' has no commands to remove.\n", tr.Name)
			return
		}

		index, selectedCmd, err := ui.SelectCommand("Select command to remove", tr.Commands)
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Removal cancelled: %v", err)
		}

		// Remove from slice
		tr.Commands = append(tr.Commands[:index], tr.Commands[index+1:]...)

		if err := transit.SaveTransit(tr); err != nil {
			log.SetFlags(0)
			log.Fatalf("Failed to save transit '%s': %v", tr.Name, err)
		}

		fmt.Printf("✓ Removed '%s' from transit '%s'\n", selectedCmd, tr.Name)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
