package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var reorderCmd = &cobra.Command{
	Use:               "reorder [transit name]",
	Short:             "Reorder commands in a transit configuration",
	Long:              `Change the execution order of commands in a specific transit.`,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		var tr *transit.Transit

		if len(args) == 0 {
			var err error
			tr, err = ui.SelectTransit("reorder")
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

		if len(tr.Commands) <= 1 {
			fmt.Printf("Transit '%s' has %d command(s); need at least 2 to reorder.\n", tr.Name, len(tr.Commands))
			return
		}

		for {
			items := make([]string, len(tr.Commands))
			for i, c := range tr.Commands {
				items[i] = fmt.Sprintf("%d. %s", i+1, c)
			}

			prompt := promptui.Select{
				Label: "Select a command to move",
				Items: items,
			}

			fromIndex, _, err := prompt.Run()
			if err != nil {
				break
			}

			toPrompt := promptui.Prompt{
				Label: fmt.Sprintf("Enter new position (1-%d)", len(tr.Commands)),
				Validate: func(input string) error {
					pos, err := strconv.Atoi(input)
					if err != nil || pos < 1 || pos > len(tr.Commands) {
						return fmt.Errorf("must be a number between 1 and %d", len(tr.Commands))
					}
					return nil
				},
			}

			toPosStr, err := toPrompt.Run()
			if err != nil {
				break
			}

			toIndex, _ := strconv.Atoi(toPosStr)
			toIndex-- // 1-based to 0-based

			// Move element
			cmdToMove := tr.Commands[fromIndex]
			tr.Commands = append(tr.Commands[:fromIndex], tr.Commands[fromIndex+1:]...)
			tr.Commands = append(tr.Commands[:toIndex], append([]string{cmdToMove}, tr.Commands[toIndex:]...)...)

			if err := transit.SaveTransit(tr); err != nil {
				log.SetFlags(0)
				log.Fatalf("Failed to save reordered transit: %v", err)
			}

			fmt.Println("\nUpdated command order:")
			for i, c := range tr.Commands {
				fmt.Printf("  %d. %s\n", i+1, c)
			}

			continuePrompt := promptui.Prompt{
				Label:     "Reorder another command? [y/N]",
				IsConfirm: true,
			}
			res, _ := continuePrompt.Run()
			if strings.ToLower(res) != "y" {
				break
			}
		}
		fmt.Printf("✓ Finished reordering transit '%s'\n", tr.Name)
	},
}

func init() {
	rootCmd.AddCommand(reorderCmd)
}
