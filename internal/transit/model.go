package transit

// Transit represents a named workflow containing a list of commands.
type Transit struct {
	Name     string   `yaml:"name,omitempty"`
	Commands []string `yaml:"commands"`
	IsLocal  bool     `yaml:"-"`
}

// ProjectConfig represents the schema for a project-local .transit.yaml file.
type ProjectConfig struct {
	Transits map[string][]string `yaml:"transits"`
}
