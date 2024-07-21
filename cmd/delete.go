package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:           "delete [transit name(s)]",
		Short:         "Removes a transit or group of transits at the same time",
		Long:          `Use "transit delete <transit-name>" to delete one or more transits.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			var transits []string

			if len(args) == 0 {
				// List transits and prompt user to select
				var err error
				transits, _, err = middleware.ListTransit("delete")
				if err != nil {
					log.SetFlags(0)
					log.Fatalf("Error fetching list of transits: %v", err)
				}

				// Display the list and prompt the user to select
				fmt.Println("Available transits:")
				for i, transit := range transits {
					fmt.Printf("%d: %s\n", i+1, transit)
				}

				// Read user input for transit numbers to delete
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Select transit number(s) to delete (comma-separated or space-separated for multiple): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					log.Fatalf("Error reading input: %v", err)
				}

				// Check if input is empty
				if input == "\n" || input == "" {
					fmt.Println("No transits selected for deletion")
					return
				}

				// Split input by both commas and spaces
				input = strings.TrimSpace(input)
				inputs := strings.FieldsFunc(input, func(r rune) bool {
					return r == ',' || r == ' '
				})

				var selectedTransits []string
				for _, input := range inputs {
					index, err := strconv.Atoi(input)
					if err != nil || index < 1 || index > len(transits) {
						fmt.Printf("Invalid selection: %s\n", input)
						continue
					}
					selectedTransits = append(selectedTransits, transits[index-1])
				}

				// Pass all selected transits to middleware.DeleteTransit
				err = middleware.DeleteTransit(selectedTransits)
				if err != nil {
					if transitErr, ok := err.(*middleware2.TransitError); ok {
						fmt.Printf("Error: %s\n", transitErr.Message)
					} else {
						fmt.Println(err)
					}
				}

			} else {
				// Directly delete the transits specified as command line arguments
				transits = args
				err := middleware.DeleteTransit(transits)
				if err != nil {
					if transitErr, ok := err.(*middleware2.TransitError); ok {
						fmt.Printf("Error: %s\n", transitErr.Message)
					} else {
						fmt.Println(err)
					}
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}
