package config

import (
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	DefaultWorktreePath string            `yaml:"default_worktree_path"`
	Theme               ThemeConfig       `yaml:"theme"`
	Git                 GitConfig         `yaml:"git"`
	UI                  UIConfig          `yaml:"ui"`
	Aliases             map[string]string `yaml:"aliases"`
}

// ThemeConfig represents theme configuration
type ThemeConfig struct {
	Primary   string `yaml:"primary"`
	Secondary string `yaml:"secondary"`
	Success   string `yaml:"success"`
	Warning   string `yaml:"warning"`
	Error     string `yaml:"error"`
	Muted     string `yaml:"muted"`
	Text      string `yaml:"text"`
}

// GitConfig represents git-specific configuration
type GitConfig struct {
	AutoCreateBranch bool     `yaml:"auto_create_branch"`
	DefaultRemote    string   `yaml:"default_remote"`
	ExcludePatterns  []string `yaml:"exclude_patterns"`
}

// UIConfig represents UI-specific configuration
type UIConfig struct {
	ShowIcons     bool `yaml:"show_icons"`
	ConfirmDelete bool `yaml:"confirm_delete"`
	MaxPathLength int  `yaml:"max_path_length"`
}

// defaultConfig returns the default configuration
func defaultConfig() *Config {
	return &Config{
		DefaultWorktreePath: "../",
		Theme: ThemeConfig{
			Primary:   "#7C3AED",
			Secondary: "#06B6D4",
			Success:   "#10B981",
			Warning:   "#F59E0B",
			Error:     "#EF4444",
			Muted:     "#6B7280",
			Text:      "#F9FAFB",
		},
		Git: GitConfig{
			AutoCreateBranch: false,
			DefaultRemote:    "origin",
			ExcludePatterns:  []string{},
		},
		UI: UIConfig{
			ShowIcons:     true,
			ConfirmDelete: true,
			MaxPathLength: 50,
		},
		Aliases: map[string]string{
			"ls": "list",
			"sw": "switch",
			"rm": "remove",
		},
	}
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	// Check for local config first
	if _, err := os.Stat(".yosegi.yaml"); err == nil {
		return ".yosegi.yaml", nil
	}

	// Check for global config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "yosegi")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// Load loads the configuration from file or returns default config
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return defaultConfig(), nil
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return defaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig(), nil
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return defaultConfig(), nil
	}

	// Merge with defaults for missing fields
	defaultCfg := defaultConfig()
	if config.DefaultWorktreePath == "" {
		config.DefaultWorktreePath = defaultCfg.DefaultWorktreePath
	}
	if config.Theme.Primary == "" {
		config.Theme = defaultCfg.Theme
	}
	if config.Aliases == nil {
		config.Aliases = defaultCfg.Aliases
	}

	return &config, nil
}

// Save saves the configuration to file
func Save(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// InitConfig creates a default configuration file
func InitConfig() error {
	config := defaultConfig()
	return Save(config)
}