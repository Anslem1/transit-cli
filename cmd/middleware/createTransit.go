package middleware

import (
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

// CreateEmptyTransit creates multiple empty transits without any commands.
func CreateEmptyTransit(filenames []string) (bool, error) {
	// Determine the Library (or equivalent) directory path based on OS
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return false, middleware2.NewTransitError(1, fmt.Sprintf("failed to get library path: %v", err))
	}

	// Define the path to transit/cmds directory
	transitPath := filepath.Join(libraryPath, "transit", "cmds")

	// Create transit/cmds directory if it doesn't exist
	if err := os.MkdirAll(transitPath, 0755); err != nil {
		return false, middleware2.NewTransitError(2, fmt.Sprintf("failed to create directory: %v", err))
	}

	var errors []string

	for _, filename := range filenames {
		// Data to be written to YAML with no commands
		data := CommandList{
			Commands: []string{},
		}

		// Check if the file already exists
		filePath := filepath.Join(transitPath, filename+".yaml")
		if _, err := os.Stat(filePath); err == nil {
			errors = append(errors, fmt.Sprintf("Error: transit %s already exists, skipped", filename))
			continue
		}

		// Create the file inside transit/cmds directory
		file, err := os.Create(filePath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to create transit %s: %v", filename, err))
			continue
		}
		defer file.Close()

		// Write data to the YAML file
		encoder := yaml.NewEncoder(file)
		defer encoder.Close()
		if err := encoder.Encode(data); err != nil {
			errors = append(errors, fmt.Sprintf("failed to encode YAML for transit %s: %v", filename, err))
			continue
		}
		fmt.Printf("Empty transit %s created. You can add commands later with the `transit add %s` command\n", filename, filename)
	}

	if len(errors) > 0 {

		return false, middleware2.NewTransitError(5, fmt.Sprintf("%v\n", strings.Join(errors, "\n")))
	}
	return true, nil
}
