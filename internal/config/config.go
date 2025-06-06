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
	AutoCreateBranch             bool     `yaml:"auto_create_branch"`
	DeleteBranchOnWorktreeRemove bool     `yaml:"delete_branch_on_worktree_remove"`
	DefaultRemote                string   `yaml:"default_remote"`
	ExcludePatterns              []string `yaml:"exclude_patterns"`
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
	config, err := loadConfigFromFile()
	if err != nil {
		return defaultConfig(), nil
	}

	mergeWithDefaults(config)
	return config, nil
}

// loadConfigFromFile loads and parses the configuration file
func loadConfigFromFile() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// mergeWithDefaults merges the loaded config with default values for missing fields
func mergeWithDefaults(config *Config) {
	defaultCfg := defaultConfig()

	if config.DefaultWorktreePath == "" {
		config.DefaultWorktreePath = defaultCfg.DefaultWorktreePath
	}

	mergeThemeConfig(&config.Theme, &defaultCfg.Theme)
	mergeGitConfig(&config.Git, &defaultCfg.Git)
	mergeUIConfig(&config.UI, &defaultCfg.UI)

	if config.Aliases == nil {
		config.Aliases = defaultCfg.Aliases
	}
}

// mergeThemeConfig merges theme configuration with defaults
func mergeThemeConfig(config, defaultCfg *ThemeConfig) {
	if config.Primary == "" {
		config.Primary = defaultCfg.Primary
	}
	if config.Secondary == "" {
		config.Secondary = defaultCfg.Secondary
	}
	if config.Success == "" {
		config.Success = defaultCfg.Success
	}
	if config.Warning == "" {
		config.Warning = defaultCfg.Warning
	}
	if config.Error == "" {
		config.Error = defaultCfg.Error
	}
	if config.Muted == "" {
		config.Muted = defaultCfg.Muted
	}
	if config.Text == "" {
		config.Text = defaultCfg.Text
	}
}

// mergeGitConfig merges git configuration with defaults
func mergeGitConfig(config, defaultCfg *GitConfig) {
	if config.DefaultRemote == "" {
		config.DefaultRemote = defaultCfg.DefaultRemote
	}
	if config.ExcludePatterns == nil {
		config.ExcludePatterns = defaultCfg.ExcludePatterns
	}
}

// mergeUIConfig merges UI configuration with defaults
func mergeUIConfig(config, defaultCfg *UIConfig) {
	if config.MaxPathLength == 0 {
		config.MaxPathLength = defaultCfg.MaxPathLength
	}
	// Note: ShowIcons and ConfirmDelete will use Go's zero values (false)
	// if not explicitly set in config. This is expected behavior.
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
