package middleware

import (
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"os"
	"path/filepath"
	"strings"
)

// DeleteTransit deletes the specified YAML files from the transit/cmds directory.
func DeleteTransit(filenames []string) error {
	// Determine the Library (or equivalent) directory path based on OS
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return middleware2.NewTransitError(1, fmt.Sprintf("failed to get library path: %v", err))
	}

	// Define the path to transit/cmds directory
	transitPath := filepath.Join(libraryPath, "transit", "cmds")

	var deletedFiles []string
	var notFoundFiles []string
	var deletionErrors []string

	for _, filename := range filenames {
		// Check if the file exists
		filePath := filepath.Join(transitPath, filename+".yaml")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			notFoundFiles = append(notFoundFiles, filename)
			continue
		}

		// Delete the file
		if err := os.Remove(filePath); err != nil {
			deletionErrors = append(deletionErrors, fmt.Sprintf("%s: %v", filename, err))
			continue
		}
		deletedFiles = append(deletedFiles, filename)
	}

	// Create a detailed error message if there are any errors
	var errorMessages []string
	if len(notFoundFiles) > 0 {
		for _, file := range notFoundFiles {
			errorMessages = append(errorMessages, fmt.Sprintf("%s transit could not be found", file))
		}
	}
	if len(deletionErrors) > 0 {
		errorMessages = append(errorMessages, deletionErrors...)
	}

	// Return a combined error if there are any error messages
	if len(errorMessages) > 0 {
		return middleware2.NewTransitError(2, strings.Join(errorMessages, "\n"))
	}

	if len(deletedFiles) == 0 {
		return middleware2.NewTransitError(3, "None of the specified transits exist")
	}

	// Print deleted files
	for _, file := range deletedFiles {
		fmt.Printf("%s transit deleted\n", file)
	}

	return nil
}