package middleware

import (
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

func ReadCommandsInTransit(filename string) ([]string, error) {
	// Determine the Library (or equivalent) directory path based on OS
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return nil, err
	}

	// Define the path to transit/cmds directory
	transitPath := filepath.Join(libraryPath, "transit", "cmds")

	// Check if the file exists
	filePath := filepath.Join(transitPath, filename+".yaml")
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
	
		return nil, fmt.Errorf("transit named %s does not exist", filename)
	}

	// Read file contents
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	// Unmarshal YAML into CommandList struct
	var commands CommandList
	if err := yaml.Unmarshal(fileContents, &commands); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}
	return commands.Commands, nil
}
