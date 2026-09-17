package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:               "delete [transit name(s)]",
	Short:             "Removes a transit or group of transits",
	Long:              `Use "transit delete <transit-name>" to delete one or more transits.`,
	SilenceErrors:     true,
	SilenceUsage:      true,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		toDelete := args
		if len(toDelete) == 0 {
			tr, err := ui.SelectTransit("delete")
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error: %v", err)
			}
			toDelete = []string{tr.Name}
		}

		for _, name := range toDelete {
			name = strings.TrimSuffix(name, ".yaml")
			err := transit.DeleteTransit(name)
			if err != nil {
				log.SetFlags(0)
				log.Printf("Error deleting transit '%s': %v\n", name, err)
			} else {
				fmt.Printf("✓ Deleted transit '%s'\n", name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
