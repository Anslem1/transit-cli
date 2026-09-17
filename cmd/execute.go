package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Anslem1/transit/internal/runner"
	"github.com/Anslem1/transit/internal/transit"
	"github.com/Anslem1/transit/internal/ui"
	"github.com/spf13/cobra"
)

var (
	skipPrompts     bool
	parallel        bool
	continueOnError bool
	dryRun          bool
)

var executeCmd = &cobra.Command{
	Use:               "execute [transit name] [args...]",
	Short:             "Execute the Transit",
	Long:              `Run a list of commands from a transit file with optional interactive confirmation or in parallel. You can pass arguments to substitute placeholders ($1, $2, $@).`,
	ValidArgsFunction: TransitNameCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		var tr *transit.Transit
		var runArgs []string

		if len(args) == 0 {
			var err error
			tr, err = ui.SelectTransit("execute")
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error selecting transit: %v", err)
			}
		} else {
			transitName := strings.TrimSuffix(args[0], ".yaml")
			var err error
			tr, err = transit.GetTransit(transitName)
			if err != nil {
				log.SetFlags(0)
				log.Fatalf("Error reading transit '%s': %v", transitName, err)
			}
			if len(args) > 1 {
				runArgs = args[1:]
			}
		}

		if tr == nil {
			fmt.Println("No transit selected. Exiting.")
			return
		}

		opts := runner.ExecuteOptions{
			Args:            runArgs,
			SkipPrompt:      skipPrompts || parallel,
			Parallel:        parallel,
			ContinueOnError: continueOnError,
			DryRun:          dryRun,
		}

		if err := runner.ExecuteCommands(tr.Commands, opts); err != nil {
			log.SetFlags(0)
			log.Fatalf("Execution failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
	executeCmd.Flags().BoolVarP(&skipPrompts, "skip", "s", false, "Skip confirmation prompts and execute all commands at once")
	executeCmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "Execute all commands in parallel/concurrently")
	executeCmd.Flags().BoolVarP(&continueOnError, "continue-on-error", "c", false, "Continue running remaining commands even if a command fails")
	executeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview commands that would be executed without running them")
}
