package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/spf13/cobra"
)

var (
	editCmd = &cobra.Command{
		Use:   "edit [transit name]",
		Short: "Edit commands within a specified transit",
		Long:  `Allows editing, updating, or reordering commands within a specified transit.`,
		Run: func(cmd *cobra.Command, args []string) {
			var transitName string
			if len(args) == 0 {
				var err error
				_, transitName, err = middleware.ListTransit("edit")

				if err != nil {
					log.SetFlags(0)
					log.Fatalf("Error fetching list of transits: %v", err)
				}

			} else {
				transitName = strings.Split(args[0], ".yaml")[0]
			}
			if transitName == "" {
				fmt.Println("No transit selected. Exiting.")
				return
			}
			_, err := middleware.ReadCommandsInTransit(strings.Split(transitName, ".yaml")[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			transitName = strings.Split(transitName, ".yaml")[0]
			middleware.EditTransit(transitName)
		},
	}
)

func init() {
	rootCmd.AddCommand(editCmd)
}
