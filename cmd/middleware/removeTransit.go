package middleware

import (
	"bufio"
	"fmt"
	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// CommandFile defines the structure of your YAML file
type CommandFile struct {
	Commands []string `yaml:"commands"`
}

// RemoveCommandsFromTransit removes specified commands from a YAML file and saves the updated file
func RemoveCommandsFromTransit(filename string, commandsToRemove []string) error {
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return err
	}

	transitPath := filepath.Join(libraryPath, "transit", "cmds")
	filePath := filepath.Join(transitPath, filename+".yaml")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("the transit %s does not exist", filename)
	}

	// Read the YAML file
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %v", err)
	}

	// Unmarshal the YAML file into the CommandFile struct
	var cmdFile CommandFile
	err = yaml.Unmarshal(yamlFile, &cmdFile)
	if err != nil {
		return fmt.Errorf("error unmarshaling YAML file: %v", err)
	}

	// Filter out the commands to be removed
	var updatedCommands []string
	commandsToRemoveMap := make(map[string]bool)
	for _, cmd := range commandsToRemove {
		commandsToRemoveMap[cmd] = true
	}
	for _, cmd := range cmdFile.Commands {
		if !commandsToRemoveMap[cmd] {
			updatedCommands = append(updatedCommands, cmd)
		}
	}

	// Update the Commands field with the filtered list
	cmdFile.Commands = updatedCommands

	// Marshal the updated CommandFile back into YAML format
	updatedYamlFile, err := yaml.Marshal(&cmdFile)
	if err != nil {
		return fmt.Errorf("error marshaling updated YAML file: %v", err)
	}

	// Write the updated YAML back to the file
	err = os.WriteFile(filePath, updatedYamlFile, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing updated YAML file: %v", err)
	}

	return nil
}

// GetUserSelectedCommands prompts the user to select commands to remove
func GetUserSelectedCommands(commands []string) ([]string, error) {
	selectedCommands := make(map[string]bool)

	for {
		// Sort commands alphabetically for better user experience
		sort.Strings(commands)

		// Print selected commands if any
		if len(selectedCommands) > 0 {
			fmt.Printf("Selected commands: ")
			for cmd := range selectedCommands {
				fmt.Printf("%s, ", cmd)
			}
			fmt.Println()
		}

		// Print commands with indices
		fmt.Println("Commands in the transit:")
		for i, cmd := range commands {
			fmt.Printf("[%d] %s\n", i+1, cmd)
		}
		fmt.Println()

		// Prompt user to enter command number(s) or press Enter to finish
		fmt.Print("Enter command number(s) to select/deselect (space-separated) or press Enter to finish: ")

		// Read user input
		input := getUserInput()

		// If input is empty, exit the loop
		if input == "" {
			fmt.Println("Exit")
			break
		}

		// Split input by spaces
		indexStrs := strings.Fields(input)

		var invalidIndices []string

		// Convert input strings to indices
		for _, indexStr := range indexStrs {
			// Trim spaces
			indexStr = strings.TrimSpace(indexStr)

			// Check if indexStr contains only numeric characters
			if _, err := strconv.Atoi(indexStr); err != nil {
				invalidIndices = append(invalidIndices, indexStr)
				continue
			}

			// Convert to integer index
			index, err := strconv.Atoi(indexStr)
			if err != nil || index < 1 || index > len(commands) {
				invalidIndices = append(invalidIndices, indexStr)
				continue
			}

			// Adjust index to zero-based
			selectedIndex := index - 1
			command := commands[selectedIndex]

			// Toggle selection
			if selectedCommands[command] {
				delete(selectedCommands, command)
			} else {
				selectedCommands[command] = true
			}
		}

		// Check for any invalid indices
		if len(invalidIndices) > 0 {
			fmt.Printf("The command(s) '%s' number does not exist in this transit. Please enter valid number(s).\n", strings.Join(invalidIndices, ", "))
		}
	}

	// Convert selectedCommands map to a slice
	var selectedCommandsSlice []string
	for cmd := range selectedCommands {
		selectedCommandsSlice = append(selectedCommandsSlice, cmd)
	}

	return selectedCommandsSlice, nil
}

// getUserInput reads user input from the console
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}