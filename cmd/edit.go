package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:               "edit [transit name]",
	Short:             "Edit commands within a specified transit",
	Long:              `Allows editing or updating an existing command within a specified transit.`,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		var tr *transit.Transit

		if len(args) == 0 {
			var err error
			tr, err = ui.SelectTransit("edit")
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
			fmt.Printf("Transit '%s' has no commands to edit.\n", tr.Name)
			return
		}

		index, currentCmd, err := ui.SelectCommand("Select a command to edit", tr.Commands)
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Selection cancelled: %v", err)
		}

		prompt := promptui.Prompt{
			Label:   "Edit command",
			Default: currentCmd,
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("command cannot be empty")
				}
				return nil
			},
		}

		newCmd, err := prompt.Run()
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Edit cancelled: %v", err)
		}

		tr.Commands[index] = strings.TrimSpace(newCmd)
		if err := transit.SaveTransit(tr); err != nil {
			log.SetFlags(0)
			log.Fatalf("Failed to save transit '%s': %v", tr.Name, err)
		}

		fmt.Printf("✓ Updated command in transit '%s':\n  Before: %s\n  After:  %s\n", tr.Name, currentCmd, newCmd)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
