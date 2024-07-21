package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list [transit name]",
		Short: "List available Transits or commands in a specified Transit",
		Long:  `List all available transits stored in the transit/cmds directory, or list commands in a specified transit.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				listAllTransits()
			} else {
				listCommandsInTransit(args[0])
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func listAllTransits() {
	transits, selectedTransits, err := middleware.ListTransit("list")
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Error listing transits: %v", err)
	}
	if len(transits) == 0 {
		fmt.Println("No transits found.")
		return
	}
	withoutType := strings.Split(selectedTransits, ".yaml")[0]
	listCommands(withoutType)
}

func listCommandsInTransit(transitName string) {
	commands, err := middleware.ReadCommandsInTransit(transitName)
	if err != nil {
			log.SetFlags(0)
		log.Fatalf("Error reading transit: %v", err)
	}

	fmt.Printf("Commands in transit %s:\n", transitName)
	for _, cmd := range commands {
		fmt.Println(" -", cmd)
	}
}

func listCommands(transitName string) {
	commands, err := middleware.ReadCommandsInTransit(transitName)
	if err != nil {
				log.SetFlags(0)
		log.Fatalf("Error reading transit: %v", err)
	}

	fmt.Printf("Commands in transit %s:\n", transitName)
	for _, cmd := range commands {
		fmt.Println(" -", cmd)
	}
}
