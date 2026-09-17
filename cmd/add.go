package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:               "add [transit name] [commands...]",
	Short:             "Adds new commands to a transit configuration",
	Long:              `Adds one or more commands to an existing transit. You can specify commands on the command line or enter them interactively.`,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		var tr *transit.Transit

		if len(args) == 0 {
			var err error
			tr, err = ui.SelectTransit("add commands to")
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error selecting transit: %v", err)
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

		existingMap := make(map[string]bool)
		for _, c := range tr.Commands {
			existingMap[c] = true
		}

		var addedCommands []string

		// Add commands passed via CLI args
		if len(args) > 1 {
			for _, c := range args[1:] {
				c = strings.TrimSpace(c)
				if c != "" && !existingMap[c] {
					tr.Commands = append(tr.Commands, c)
					existingMap[c] = true
					addedCommands = append(addedCommands, c)
				}
			}
		} else {
			// Interactive input loop
			fmt.Println("Enter commands to add (leave blank to finish):")
			for {
				command, err := ui.PromptInput("New command")
				if err != nil || command == "" {
					break
				}
				if existingMap[command] {
					fmt.Printf("⚠️  Command '%s' already in transit '%s', skipped\n", command, tr.Name)
					continue
				}
				tr.Commands = append(tr.Commands, command)
				existingMap[command] = true
				addedCommands = append(addedCommands, command)
			}
		}

		if len(addedCommands) > 0 {
			if err := transit.SaveTransit(tr); err != nil {
				log.SetFlags(0)
				log.Fatalf("Failed to save transit: %v", err)
			}
			for _, c := range addedCommands {
				fmt.Printf("✓ Added command: %s\n", c)
			}
			fmt.Printf("Successfully added %d command(s) to '%s'\n", len(addedCommands), tr.Name)
		} else {
			fmt.Println("No new commands were added.")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
