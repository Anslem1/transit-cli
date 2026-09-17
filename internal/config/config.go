package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the standard directory where global transit files are stored.
// On macOS: ~/Library/Application Support/transit/cmds (or ~/.config/transit/cmds)
// On Linux: ~/.config/transit/cmds (XDG compliant)
// On Windows: %APPDATA%\transit\cmds
func GetConfigDir() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, ".config")
	}

	transitDir := filepath.Join(baseDir, "transit", "cmds")
	return transitDir, nil
}

// GetLegacyDir returns the legacy directory used in previous versions of Transit
// (e.g., ~/Library/transit/cmds on macOS, ~/Documents/transit/cmds on Linux).
func GetLegacyDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "transit", "cmds"), nil
	}
	return filepath.Join(home, "Documents", "transit", "cmds"), nil
}

// EnsureDir ensures that the directory exists with standard directory permissions.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
