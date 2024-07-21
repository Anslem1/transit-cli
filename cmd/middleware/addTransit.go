package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

func AddCommandsToTransit(filename string, initialCommands []string) error {
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return err
	}

	transitPath := filepath.Join(libraryPath, "transit", "cmds")
	filePath := filepath.Join(transitPath, filename+".yaml")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("the transit %s you are trying to update does not exist", filename)
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

	// Create a map for quick lookup of existing commands
	existingCommandsMap := make(map[string]bool)
	for _, cmd := range cmdFile.Commands {
		existingCommandsMap[cmd] = true
	}

	// Add initial commands if provided
	var addedCommands []string
	for _, command := range initialCommands {
		parsedCommands := parseCommands(command)
		for _, cmd := range parsedCommands {
			if existingCommandsMap[cmd] {
				fmt.Printf("\"%s\" already exists in transit '%s', skipped\n", cmd, filename)
			} else {
				cmdFile.Commands = append(cmdFile.Commands, cmd)
				existingCommandsMap[cmd] = true
				addedCommands = append(addedCommands, cmd)
			}
		}
	}

	// Interactive prompt to add new commands
	for {
		prompt := promptui.Prompt{
			Label: "Enter new command (or leave blank to finish)",
		}

		command, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			break
		}

		if command == "" {
			break
		}

		parsedCommands := parseCommands(command)
		for _, cmd := range parsedCommands {
			if existingCommandsMap[cmd] {
				fmt.Printf("\"%s\" already exists in transit '%s', skipped\n", cmd, filename)
			} else {
				cmdFile.Commands = append(cmdFile.Commands, cmd)
				existingCommandsMap[cmd] = true
				addedCommands = append(addedCommands, cmd)
			}
		}
	}

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

	if len(addedCommands) > 0 {
		for _, v := range addedCommands {
			fmt.Printf("Command %v added successfully to transit '%s'\n", v, filename)
		}
		fmt.Printf("Added %v new commands to the %v transit\n", len(addedCommands), filename)
	} else {
		fmt.Println("No new commands were added.")
	}

	return nil
}

// parseCommands splits commands while respecting quotes
func parseCommands(command string) []string {
	var commands []string
	var sb strings.Builder
	inQuotes := false

	for _, r := range command {
		switch r {
		case ',':
			if inQuotes {
				sb.WriteRune(r)
			} else {
				commands = append(commands, strings.TrimSpace(sb.String()))
				sb.Reset()
			}
		case '"', '\'':
			inQuotes = !inQuotes
			sb.WriteRune(r)
		default:
			sb.WriteRune(r)
		}
	}
	if sb.Len() > 0 {
		commands = append(commands, strings.TrimSpace(sb.String()))
	}

	return commands
}