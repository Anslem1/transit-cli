package middleware

import (
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
	"strings"
)

// ListTransits lists all available transit YAML files in the transit/cmds directory.
func ListTransit(listFrom string) ([]string, string, error) {
	// Determine the Library (or equivalent) directory path based on OS
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get library path: %v", err)
	}

	// Define the path to transit/cmds directory
	transitPath := filepath.Join(libraryPath, "transit", "cmds")

	// Open the transit/cmds directory
	dirEntries, err := os.ReadDir(transitPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read transit directory: %v", err)
	}

	// Collect the names and filenames of all YAML files (transits)
	var transits []string
	var transitFilenames []string
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			transitName := strings.TrimSuffix(entry.Name(), ".yaml")
			transits = append(transits, transitName)
			transitFilenames = append(transitFilenames, entry.Name())
		}
	}

	// Prompt the user to select a transit
	if len(transits) == 0 {
		return nil, "", fmt.Errorf("no transits found")
	}

	var prompt promptui.Select

	if listFrom == "delete" {
		return transits, "", nil
	} else if listFrom == "execute" {
		prompt = promptui.Select{
			Label: "Select a transit to execute",
			Items: transits,
		}
	} else if listFrom == "list" {
		prompt = promptui.Select{
			Label: "Select a transit to list",
			Items: transits,
		}
	} else if listFrom == "edit" {
		prompt = promptui.Select{
			Label: "Select a transit to edit a command",
			Items: transits,
		}
	} else if listFrom == "add" {
		prompt = promptui.Select{
			Label: "Select a transit to add a command",
			Items: transits,
		}
	} else if listFrom == "reorder" {
		prompt = promptui.Select{
			Label: "Select a transit to reoder a command",
			Items: transits,
		}
	} else if listFrom == "remove" {
		prompt = promptui.Select{
			Label: "Select a transit you want to remove a command",
			Items: transits,
		}
	}
	// Run the prompt and get the selected transit name
	_, result, err := prompt.Run()
	if err != nil {
		return nil, "", fmt.Errorf("prompt failed: %v", err)
	}

	// Find the corresponding filename for the selected transit
	var selectedFilename string
	for i, transit := range transits {
		if transit == result {
			selectedFilename = transitFilenames[i]
			break
		}
	}

	if selectedFilename == "" {
		return nil, "", fmt.Errorf("selected transit filename not found")
	}
	return transits, selectedFilename, nil
}
