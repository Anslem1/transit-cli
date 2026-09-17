package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for commands across all transits",
	Long:  `Search for commands containing a search term across all transits.`,
	Run: func(cmd *cobra.Command, args []string) {
		var query string
		if len(args) == 0 {
			var err error
			query, err = ui.PromptInput("Search query")
			if err != nil || strings.TrimSpace(query) == "" {
				fmt.Println("No search query entered.")
				return
			}
		} else {
			query = strings.Join(args, " ")
		}

		results, err := transit.SearchCommands(query)
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Error searching transits: %v", err)
		}

		if len(results) == 0 {
			fmt.Printf("No commands found matching '%s'\n", query)
			return
		}

		fmt.Printf("Matches for '%s':\n", query)
		for transitName, matches := range results {
			fmt.Printf("\n📁 Transit: %s\n", transitName)
			for _, m := range matches {
				fmt.Printf("   • %s\n", m)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
