package ui

import (
	"fmt"
	"strings"

	"github.com/Anslem1/transit/internal/transit"
	"github.com/manifoldco/promptui"
)

// SelectTransit prompts the user to select an available transit from local and global sources.
func SelectTransit(actionLabel string) (*transit.Transit, error) {
	transits, err := transit.ListTransits()
	if err != nil {
		return nil, fmt.Errorf("failed to list transits: %w", err)
	}
	if len(transits) == 0 {
		return nil, fmt.Errorf("no transits found. Create one with 'transit create <name>'")
	}

	displayItems := make([]string, len(transits))
	for i, t := range transits {
		tag := ""
		if t.IsLocal {
			tag = " (local)"
		}
		displayItems[i] = fmt.Sprintf("%s%s", t.Name, tag)
	}

	prompt := promptui.Select{
		Label: fmt.Sprintf("Select a transit to %s", actionLabel),
		Items: displayItems,
	}

	selectedIndex, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("selection cancelled: %w", err)
	}

	return &transits[selectedIndex], nil
}

// PromptInput prompts the user for text input with a given label.
func PromptInput(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	res, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}

// SelectCommand prompts the user to choose one command from a list.
func SelectCommand(label string, commands []string) (int, string, error) {
	if len(commands) == 0 {
		return -1, "", fmt.Errorf("no commands to select from")
	}
	prompt := promptui.Select{
		Label: label,
		Items: commands,
	}
	idx, res, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}
	return idx, res, nil
}
