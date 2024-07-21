package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

// CommandFile defines the structure of your YAML file


func ReorderCommandsInTransit(filename string) error {
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

	// Check if there are any commands to reorder
	if len(cmdFile.Commands) == 0 {
		fmt.Printf("No commands to reorder in transit '%s'.\n", filename)
		return nil
	}

	// Interactive prompt to reorder commands
	for {
		// Select the command to move
		items := make([]string, len(cmdFile.Commands))
		for i, cmd := range cmdFile.Commands {
			items[i] = fmt.Sprintf("%d. %s", i+1, cmd)
		}

		prompt := promptui.Select{
			Label: "Select a command to move",
			Items: items,
		}

		index, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			break
		}

		// Display the list of positions with commands for the new position
		fmt.Println("Available positions:")
		for i, cmd := range cmdFile.Commands {
			fmt.Printf("%d. %s\n", i+1, cmd)
		}

		// Prompt for the new position
		positionPrompt := promptui.Prompt{
			Label: "Enter new position (number) or select from list",
			Validate: func(input string) error {
				pos, err := strconv.Atoi(input)
				if err != nil || pos < 1 || pos > len(cmdFile.Commands) {
					return fmt.Errorf("invalid position")
				}
				return nil
			},
		}

		positionStr, err := positionPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			break
		}

		newPosition, _ := strconv.Atoi(positionStr)

		// Move the command to the new position
		cmd := cmdFile.Commands[index]
		cmdFile.Commands = append(cmdFile.Commands[:index], cmdFile.Commands[index+1:]...)
		cmdFile.Commands = append(cmdFile.Commands[:newPosition-1], append([]string{cmd}, cmdFile.Commands[newPosition-1:]...)...)

		// Confirm the reordering
		fmt.Println("Commands reordered:")
		for i, cmd := range cmdFile.Commands {
			fmt.Printf("%d. %s\n", i+1, cmd)
		}

		// Ask if the user wants to continue reordering
		confirmPrompt := promptui.Prompt{
			Label:     "Do you want to reorder another command? (y/n)",
			AllowEdit: true,
		}

		confirm, err := confirmPrompt.Run()
		if err != nil || confirm != "y" {
			break
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

	fmt.Printf("Successfully reordered commands in transit '%s'.\n", filename)

	return nil
}
