package cmd

import (
	"fmt"
	"log"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var localCreate bool

var CreateCmd = &cobra.Command{
	Use:           "create [transit name(s)]",
	Short:         "Creates one or more empty transits",
	Long:          `Use "transit create <transit-name>" to create one or more transits. Use --local to create in project-local .transit.yaml.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		names := args
		if len(names) == 0 {
			name, err := ui.PromptInput("Enter transit name")
			if err != nil || name == "" {
				fmt.Println("No transit name provided. Exiting.")
				return
			}
			names = []string{name}
		}

		for _, name := range names {
			err := transit.CreateTransit(name, []string{}, localCreate)
			if err != nil {
				log.SetFlags(0)
				log.Printf("Error creating transit '%s': %v\n", name, err)
				continue
			}
			scope := "global"
			if localCreate {
				scope = "local (.transit.yaml)"
			}
			fmt.Printf("✓ Created %s transit '%s'. Add commands using 'transit add %s'.\n", scope, name, name)
		}
	},
}

func init() {
	rootCmd.AddCommand(CreateCmd)
	CreateCmd.Flags().BoolVarP(&localCreate, "local", "l", false, "Create transit in project-local .transit.yaml")
}
