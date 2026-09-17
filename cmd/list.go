package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var listAllFlag bool

var listCmd = &cobra.Command{
	Use:               "list [transit name]",
	Short:             "List available Transits or commands in a specified Transit",
	Long:              `List all available transits, or display commands inside a specific transit. Use --all to list all transits without prompting.`,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		if listAllFlag {
			printAllTransitsSummary()
			return
		}

		if len(args) == 0 {
			tr, err := ui.SelectTransit("list commands for")
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error: %v", err)
			}
			printTransitCommands(tr)
			return
		}

		transitName := strings.TrimSuffix(args[0], ".yaml")
		tr, err := transit.GetTransit(transitName)
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Error reading transit '%s': %v", transitName, err)
		}
		printTransitCommands(tr)
	},
}

func printAllTransitsSummary() {
	transits, err := transit.ListTransits()
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Error listing transits: %v", err)
	}
	if len(transits) == 0 {
		fmt.Println("No transits found.")
		return
	}
	fmt.Println("Available Transits:")
	for _, t := range transits {
		scope := "global"
		if t.IsLocal {
			scope = "local (.transit.yaml)"
		}
		fmt.Printf(" • %-20s (%d commands) [%s]\n", t.Name, len(t.Commands), scope)
	}
}

func printTransitCommands(tr *transit.Transit) {
	scope := "global"
	if tr.IsLocal {
		scope = "local (.transit.yaml)"
	}
	fmt.Printf("\nCommands in transit '%s' [%s]:\n", tr.Name, scope)
	if len(tr.Commands) == 0 {
		fmt.Println("  (no commands)")
		return
	}
	for i, cmd := range tr.Commands {
		fmt.Printf("  %d. %s\n", i+1, cmd)
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listAllFlag, "all", "a", false, "List all transits with a summary of their commands")
}
