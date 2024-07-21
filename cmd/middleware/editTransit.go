package middleware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Anslem1/transit/cmd/middleware/middleware2"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

func UpdateCommandsInTransit(transitName string, commands []string) error {
	libraryPath, err := middleware2.GetLibraryPath()
	if err != nil {
		return fmt.Errorf("failed to get library path: %v", err)
	}

	transitPath := filepath.Join(libraryPath, "transit", "cmds", transitName+".yaml")

	file, err := os.OpenFile(transitPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open transit file: %v", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	// Wrap commands in CommandList struct
	commandList := CommandList{
		Commands: commands,
	}

	return encoder.Encode(commandList)
}

func EditTransit(transitName string) {
	commands, err := ReadCommandsInTransit(transitName)
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Error reading transit: %v", err)
	}

	if len(commands) == 0 {
		log.SetFlags(0)
		log.Printf("There are no commands in transit %s.\n", transitName)
		return
	}

	prompt := promptui.Select{
		Label: "Select a command to edit",
		Items: commands,
	}

	index, result, err := prompt.Run()
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Prompt failed: %v", err)
	}

	validate := func(input string) error {
		if input == "" {
			return fmt.Errorf("command cannot be empty")
		}
		return nil
	}

	editPrompt := promptui.Prompt{
		Label:    "Edit command",
		Default:  result,
		Validate: validate,
	}

	newCommand, err := editPrompt.Run()
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Prompt failed: %v", err)
	}

	commands[index] = newCommand

	err = UpdateCommandsInTransit(transitName, commands)
	if err != nil {
		log.SetFlags(0)
		log.Fatalf("Error updating transit: %v", err)
	}

	fmt.Printf("Updated command in transit %s: %s\n", transitName, newCommand)
}
