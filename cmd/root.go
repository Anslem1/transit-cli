package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/spf13/cobra"
)

var longDesc = `Welcome to Transit

Transit is a CLI tool designed to help developers save time by running a set of usual terminal commands all at once.
Here are some common Transit commands:

Creating Transits:
   create   Initializes a new, empty transit.
            Usage: transit create [transit_name]
            Example: transit create my-transit

Managing Commands in Transits:
   add      Adds a new command to an existing transit.
            Usage: transit add [transit_name]
            Example: transit add my-transit

   edit     Edits a particular command in a transit.
            Usage: transit edit [transit_name]
            Example: transit edit my-transit

   remove   Removes a command from a transit.
            Usage: transit remove [transit_name]
            Example: transit remove my-transit

   reorder  Changes the order in which commands are executed in a transit.
            Usage: transit reorder [transit_name]
            Example: transit reorder my-transit

   list     Displays all available transits or the commands within a specified transit.
            Usage: transit list [transit_name]
            Example: transit list my-transit

Executing Transits:
   execute  Runs the commands stored in a specified transit.
            Usage: transit execute [transit_name]
            Example: transit execute my-transit

Deleting Transits:
   delete   Deletes specified transits from your list.
            Usage: transit delete [transit_name]
            Example: transit delete my-transit

  Transit Version:
  version: Prints the version of the current transit.

You can also use 'transit [command]' without specifying the exact transit, and it will bring up a prompt to select the transit interactively.

For more information about a command, use "transit [command] --help".`

var (
	// cfgFile string

	rootCmd = &cobra.Command{
		Use:   "transit",
		Short: "Transit is a CLI tool for managing and executing commands",
		Long:  longDesc,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.SetFlags(0)
				log.Fatal(longDesc)
			}
		},
	}
)

// TransitNameCompletion provides shell autocompletion for existing transit names.
func TransitNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	transits, err := transit.ListTransitNames()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return transits, cobra.ShellCompDirectiveNoFileComp
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

