package middleware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

// RunSearch handles the search logic based on user input
func SearchTransit(args []string) {
	if len(args) == 0 {
		prompt := promptui.Prompt{
			Label: "Enter command to search for",
		}
		searchTerm, err := prompt.Run()
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Prompt failed %v\n", err)
			return
		}
		searchAcrossAllTransits(searchTerm)
	} else {
		transitName := args[0]

		// Check if the specified transit exists before prompting for the command
		if !transitExists(transitName) {
			log.SetFlags(0)
			log.Fatalf("The transit '%s' does not exist\n", transitName)
			return
		}

		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Enter command to search for in transit '%s'", transitName),
		}
		searchTerm, err := prompt.Run()
		if err != nil {
			log.SetFlags(0)
			log.Fatalf("Prompt failed %v\n", err)
			return
		}
		searchWithinTransit(transitName, searchTerm)
	}
}

// transitExists checks if a specified transit exists
func transitExists(transitName string) bool {
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Failed to get library path: %v\n", err)
		return false
	}

	transitPath := filepath.Join(libraryPath, "transit", "cmds")
	filePath := filepath.Join(transitPath, transitName+".yaml")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// searchAcrossAllTransits searches for a command across all transits
func searchAcrossAllTransits(searchTerm string) {
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Failed to get library path: %v\n", err)
		return
	}

	transitPath := filepath.Join(libraryPath, "transit", "cmds")

	// Ensure the transitPath directory exists
	err = ensureDirExists(transitPath)
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Failed to ensure transit directory exists: %v\n", err)
		return
	}

	files, err := os.ReadDir(transitPath)
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Failed to read transit directory: %v\n", err)
		return
	}

	found := false
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") {
			transitName := strings.TrimSuffix(file.Name(), ".yaml")
			commands, err := readCommandsFromTransit(transitName)
			if err != nil {
				log.Printf("Failed to read transit %s: %v\n", transitName, err)
				continue
			}
			for _, cmd := range commands {
				if strings.Contains(cmd, searchTerm) {
					fmt.Printf("The command you are looking for is in the transit '%s'\n", transitName)
					found = true
				}
			}
		}
	}

	if !found {
		fmt.Printf("The command '%s' was not found in any transits\n", searchTerm)
	}
}

// searchWithinTransit searches for a command within a specific transit
func searchWithinTransit(transitName, searchTerm string) {
	commands, err := readCommandsFromTransit(transitName)
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Failed to read transit %s: %v\n", transitName, err)
		return
	}

	found := false
	for i, cmd := range commands {
		if strings.Contains(cmd, searchTerm) {
			fmt.Printf("The command '%s' is at position %d in transit '%s'\n", searchTerm, i+1, transitName)
			found = true
		}
	}

	if !found {
		fmt.Printf("The command '%s' does not exist in transit '%s'\n", searchTerm, transitName)
	}
}

// readCommandsFromTransit reads the commands from a transit YAML file
func readCommandsFromTransit(transitName string) ([]string, error) {
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return nil, err
	}

	transitPath := filepath.Join(libraryPath, "transit", "cmds")
	filePath := filepath.Join(transitPath, transitName+".yaml")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("the transit %s does not exist", transitName)
	}

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	var cmdFile CommandFile
	err = yaml.Unmarshal(yamlFile, &cmdFile)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML file: %v", err)
	}

	return cmdFile.Commands, nil
}

// ensureDirExists ensures that a directory exists, creating it if necessary
func ensureDirExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}
	return nil
}
