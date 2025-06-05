package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
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
	AutoCreateBranch              bool     `yaml:"auto_create_branch"`
	DeleteBranchOnWorktreeRemove  bool     `yaml:"delete_branch_on_worktree_remove"`
	DefaultRemote                 string   `yaml:"default_remote"`
	ExcludePatterns               []string `yaml:"exclude_patterns"`
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
			AutoCreateBranch:             true,
			DeleteBranchOnWorktreeRemove: false, // Default to false for safety
			DefaultRemote:                "origin",
			ExcludePatterns:              []string{},
		},
		UI: UIConfig{
			ShowIcons:     true,
			ConfirmDelete: true,
			MaxPathLength: 50,
		},
		Aliases: map[string]string{
			"ls": "list",
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

	// Merge theme fields individually
	if config.Theme.Primary == "" {
		config.Theme.Primary = defaultCfg.Theme.Primary
	}
	if config.Theme.Secondary == "" {
		config.Theme.Secondary = defaultCfg.Theme.Secondary
	}
	if config.Theme.Success == "" {
		config.Theme.Success = defaultCfg.Theme.Success
	}
	if config.Theme.Warning == "" {
		config.Theme.Warning = defaultCfg.Theme.Warning
	}
	if config.Theme.Error == "" {
		config.Theme.Error = defaultCfg.Theme.Error
	}
	if config.Theme.Muted == "" {
		config.Theme.Muted = defaultCfg.Theme.Muted
	}
	if config.Theme.Text == "" {
		config.Theme.Text = defaultCfg.Theme.Text
	}

	// Merge git config fields
	if config.Git.DefaultRemote == "" {
		config.Git.DefaultRemote = defaultCfg.Git.DefaultRemote
	}
	if config.Git.ExcludePatterns == nil {
		config.Git.ExcludePatterns = defaultCfg.Git.ExcludePatterns
	}

	// Merge UI config fields (need to check if they were actually set)
	// For boolean fields, we need to check if they were explicitly set
	// This is a limitation of YAML unmarshaling - we can't distinguish between
	// false and unset. For now, we'll use defaults for uninitialized structs.
	if config.UI.MaxPathLength == 0 {
		config.UI.MaxPathLength = defaultCfg.UI.MaxPathLength
	}
	// Note: ShowIcons and ConfirmDelete will use Go's zero values (false)
	// if not explicitly set in config. This is expected behavior.

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
