package transit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anslem1/transit/internal/config"
	"gopkg.in/yaml.v2"
)

const (
	LocalConfigFileName  = ".transit.yaml"
	LocalConfigAltName   = "transit.yaml"
	StandardFilePermMode = 0644
)

// GetTransit retrieves a transit by name, looking in the local directory first, then global, then legacy.
func GetTransit(name string) (*Transit, error) {
	name = strings.TrimSuffix(name, ".yaml")

	// 1. Check local project config (.transit.yaml)
	if localTransit, err := getLocalTransit(name); err == nil && localTransit != nil {
		return localTransit, nil
	}

	// 2. Check standard global config dir
	configDir, err := config.GetConfigDir()
	if err == nil {
		filePath := filepath.Join(configDir, name+".yaml")
		if _, statErr := os.Stat(filePath); statErr == nil {
			return readGlobalTransitFile(name, filePath, false)
		}
	}

	// 3. Fallback: check legacy global config dir
	legacyDir, err := config.GetLegacyDir()
	if err == nil {
		filePath := filepath.Join(legacyDir, name+".yaml")
		if _, statErr := os.Stat(filePath); statErr == nil {
			return readGlobalTransitFile(name, filePath, false)
		}
	}

	return nil, fmt.Errorf("transit '%s' does not exist", name)
}

// ListTransits returns all available transits (local project transits + global transits).
func ListTransits() ([]Transit, error) {
	var transits []Transit
	seen := make(map[string]bool)

	// 1. Load local project transits
	localTransits, _ := loadAllLocalTransits()
	for _, t := range localTransits {
		transits = append(transits, t)
		seen[t.Name] = true
	}

	// 2. Load standard global transits
	configDir, err := config.GetConfigDir()
	if err == nil {
		if entries, err := os.ReadDir(configDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
					name := strings.TrimSuffix(entry.Name(), ".yaml")
					if !seen[name] {
						if t, err := readGlobalTransitFile(name, filepath.Join(configDir, entry.Name()), false); err == nil {
							transits = append(transits, *t)
							seen[name] = true
						}
					}
				}
			}
		}
	}

	// 3. Load legacy global transits (fallback for items not yet migrated)
	legacyDir, err := config.GetLegacyDir()
	if err == nil {
		if entries, err := os.ReadDir(legacyDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
					name := strings.TrimSuffix(entry.Name(), ".yaml")
					if !seen[name] {
						if t, err := readGlobalTransitFile(name, filepath.Join(legacyDir, entry.Name()), false); err == nil {
							transits = append(transits, *t)
							seen[name] = true
						}
					}
				}
			}
		}
	}

	return transits, nil
}

// ListTransitNames returns a list of all transit names (with "(local)" tag if local).
func ListTransitNames() ([]string, error) {
	transits, err := ListTransits()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, t := range transits {
		names = append(names, t.Name)
	}
	return names, nil
}

// CreateTransit creates a new transit with the specified name and initial commands.
func CreateTransit(name string, commands []string, local bool) error {
	name = strings.TrimSuffix(name, ".yaml")
	if name == "" {
		return fmt.Errorf("transit name cannot be empty")
	}

	// Check if already exists
	if existing, _ := GetTransit(name); existing != nil {
		return fmt.Errorf("transit '%s' already exists", name)
	}

	if local {
		return saveLocalTransit(name, commands)
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	if err := config.EnsureDir(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	t := Transit{
		Name:     name,
		Commands: commands,
		IsLocal:  false,
	}
	return SaveTransit(&t)
}

// SaveTransit writes changes to an existing transit.
func SaveTransit(t *Transit) error {
	if t.IsLocal {
		return saveLocalTransit(t.Name, t.Commands)
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	if err := config.EnsureDir(configDir); err != nil {
		return err
	}

	filePath := filepath.Join(configDir, t.Name+".yaml")

	data := struct {
		Commands []string `yaml:"commands"`
	}{
		Commands: t.Commands,
	}

	bytes, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return os.WriteFile(filePath, bytes, StandardFilePermMode)
}

// DeleteTransit removes a transit either from local project config or disk.
func DeleteTransit(name string) error {
	name = strings.TrimSuffix(name, ".yaml")

	// Check local config first
	localCfgPath := findLocalConfigFile()
	if localCfgPath != "" {
		cfg, err := loadProjectConfigFile(localCfgPath)
		if err == nil && cfg.Transits != nil {
			if _, exists := cfg.Transits[name]; exists {
				delete(cfg.Transits, name)
				return writeProjectConfigFile(localCfgPath, cfg)
			}
		}
	}

	// Check standard global dir
	var deleted bool
	if configDir, err := config.GetConfigDir(); err == nil {
		filePath := filepath.Join(configDir, name+".yaml")
		if err := os.Remove(filePath); err == nil {
			deleted = true
		}
	}

	// Check legacy dir
	if legacyDir, err := config.GetLegacyDir(); err == nil {
		filePath := filepath.Join(legacyDir, name+".yaml")
		if err := os.Remove(filePath); err == nil {
			deleted = true
		}
	}

	if !deleted {
		return fmt.Errorf("transit '%s' not found to delete", name)
	}
	return nil
}

// SearchCommands searches across all transits for commands containing the query string.
func SearchCommands(query string) (map[string][]string, error) {
	transits, err := ListTransits()
	if err != nil {
		return nil, err
	}

	results := make(map[string][]string)
	queryLower := strings.ToLower(query)

	for _, t := range transits {
		for _, cmd := range t.Commands {
			if strings.Contains(strings.ToLower(cmd), queryLower) {
				results[t.Name] = append(results[t.Name], cmd)
			}
		}
	}
	return results, nil
}

// Helpers for Local Project Config (.transit.yaml)

func findLocalConfigFile() string {
	if _, err := os.Stat(LocalConfigFileName); err == nil {
		return LocalConfigFileName
	}
	if _, err := os.Stat(LocalConfigAltName); err == nil {
		return LocalConfigAltName
	}
	return ""
}

func loadProjectConfigFile(path string) (*ProjectConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg ProjectConfig
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	if cfg.Transits == nil {
		cfg.Transits = make(map[string][]string)
	}
	return &cfg, nil
}

func writeProjectConfigFile(path string, cfg *ProjectConfig) error {
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, StandardFilePermMode)
}

func getLocalTransit(name string) (*Transit, error) {
	path := findLocalConfigFile()
	if path == "" {
		return nil, os.ErrNotExist
	}
	cfg, err := loadProjectConfigFile(path)
	if err != nil {
		return nil, err
	}
	cmds, exists := cfg.Transits[name]
	if !exists {
		return nil, os.ErrNotExist
	}
	return &Transit{
		Name:     name,
		Commands: cmds,
		IsLocal:  true,
	}, nil
}

func loadAllLocalTransits() ([]Transit, error) {
	path := findLocalConfigFile()
	if path == "" {
		return nil, nil
	}
	cfg, err := loadProjectConfigFile(path)
	if err != nil {
		return nil, err
	}
	var list []Transit
	for name, cmds := range cfg.Transits {
		list = append(list, Transit{
			Name:     name,
			Commands: cmds,
			IsLocal:  true,
		})
	}
	return list, nil
}

func saveLocalTransit(name string, commands []string) error {
	path := findLocalConfigFile()
	if path == "" {
		path = LocalConfigFileName
	}
	cfg, err := loadProjectConfigFile(path)
	if err != nil {
		cfg = &ProjectConfig{
			Transits: make(map[string][]string),
		}
	}
	cfg.Transits[name] = commands
	return writeProjectConfigFile(path, cfg)
}

func readGlobalTransitFile(name, filePath string, isLocal bool) (*Transit, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var data struct {
		Commands []string `yaml:"commands"`
	}
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &Transit{
		Name:     name,
		Commands: data.Commands,
		IsLocal:  isLocal,
	}, nil
}
