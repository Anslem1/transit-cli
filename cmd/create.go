package cmd

import (
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/spf13/cobra"
)

var (
	CreateCmd = &cobra.Command{
		Use:           "create <transit(s)>",
		Short:         "Creates one or more empty transits",
		Long:          `Use "transit create <transit(s)>" to create one or more empty transits.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				err := middleware2.NewTransitError(1, "Expected at least one transit name to be provided but got none")
				if transitErr, ok := err.(*middleware2.TransitError); ok {
					fmt.Printf("Error: %s\n", transitErr.Message)
				} else {
					fmt.Println(err)
				}
				return
			}

			success, err := middleware.CreateEmptyTransit(args)
			if err != nil {
				if transitErr, ok := err.(*middleware2.TransitError); ok {
					fmt.Printf("%s", transitErr.Message)
				} else {
					fmt.Println(err)
				}
			} else if success {
				fmt.Println("Transit(s) created successfully.")
			}
		},
	}
)
func init() {
	rootCmd.AddCommand(CreateCmd)
}
