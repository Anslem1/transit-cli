package middleware

import (
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// CommandList holds a list of commands to be written to YAML.
type CommandList struct {
	Commands []string `yaml:"commands"`
}

// WriteCommandsToFile writes a list of commands to a YAML file.
func WriteCommandsToTransit(filename string, commands []string) (bool, error) {
	// Data to be written to YAML
	data := CommandList{
		Commands: commands,
	}

	// Determine the Library (or equivalent) directory path based on OS
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return false, err
	}

	// Define the path to transit/cmds directory
	transitPath := filepath.Join(libraryPath, "transit", "cmds")

	// Create transit/cmds directory if it doesn't exist
	if err := os.MkdirAll(transitPath, 0755); err != nil {
		return false, fmt.Errorf("failed to create directory: %v", err)
	}

	// Check if the file already exists
	filePath := filepath.Join(transitPath, filename+".yaml")
	if _, err := os.Stat(filePath); err == nil {
		return false, fmt.Errorf("transit %s already exists, please use a different name for your transit", filename)
	}

	// Create the file inside transit/cmds directory
	file, err := os.Create(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to create transit: %v", err)
	}
	defer file.Close()

	// Write data to the YAML file
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	if err := encoder.Encode(data); err != nil {
		return false, fmt.Errorf("failed to encode YAML: %v", err)
	}
	fmt.Printf("Commands written to transit %s Use 'transit execute %s' to run your transit command(s).\n", filename, filename)
	return true, nil
}