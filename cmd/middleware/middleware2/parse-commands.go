package middleware2

import "strings"

// parseCommands parses the commands string into individual commands based on commas
func ParseCommands(commandsString string) []string {
	// Remove leading and trailing spaces
	commandsString = strings.TrimSpace(commandsString)

	// Split commandsString into individual commands based on commas
	commands := strings.Split(commandsString, ",")

	// Trim spaces from each command
	for i := range commands {
		commands[i] = strings.TrimSpace(commands[i])
	}

	// Remove empty strings
	var cleanedCommands []string
	for _, cmd := range commands {
		if cmd != "" {
			cleanedCommands = append(cleanedCommands, cmd)
		}
	}

	return cleanedCommands
}