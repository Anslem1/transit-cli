package middleware

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
)


func ExecuteCommandsInTransit(commands []string, skipPrompt bool) {
	for _, cmd := range commands {
		if !skipPrompt {
			prompt := promptui.Prompt{
				Label: fmt.Sprintf("Execute this command: %s [y/n]", cmd),
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				continue
			}

			if result == "y" || result == "" {
				err := runCommand(cmd)
				if err != nil {
					fmt.Printf("Error executing command: %s\n", err)
				}
			} else {
				fmt.Printf("Skipped command: %s\n", cmd)
			}
		} else {
			err := runCommand(cmd)
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
			}
		}
	}
}

// runCommand runs a single command using exec.Command
func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command) // Use "bash" on Unix-like systems, "cmd", "/C", command on Windows

	// Connect the command's input, output, and error streams to the current process's streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// Run the command and wait for it to complete
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}