package middleware2

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// getLibraryPath returns the path to the Library (or equivalent) directory based on OS.
func GetLibraryPath() (string, error) {
	var libraryPath string
	switch runtime.GOOS {
	case "darwin":
		libraryPath = filepath.Join(os.Getenv("HOME"), "Library") // macOS
	case "windows":
		libraryPath = filepath.Join(os.Getenv("SystemDrive"), "Users", os.Getenv("USERNAME"), "Documents") // Windows
	case "linux":
		libraryPath = filepath.Join(os.Getenv("HOME"), "Documents") // Linux (using Documents folder for example)
	default:
		return "", fmt.Errorf("unsupported OS")
	}
	return libraryPath, nil
}