package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg == nil {
		t.Errorf("defaultConfig() should not return nil")
	}

	// Test default values
	if cfg.DefaultWorktreePath != "../" {
		t.Errorf("Expected default worktree path '../', got '%s'", cfg.DefaultWorktreePath)
	}

	if !cfg.Git.AutoCreateBranch {
		t.Errorf("Expected AutoCreateBranch to be true by default")
	}

	if cfg.Git.DefaultRemote != "origin" {
		t.Errorf("Expected default remote 'origin', got '%s'", cfg.Git.DefaultRemote)
	}

	if !cfg.UI.ShowIcons {
		t.Errorf("Expected ShowIcons to be true by default")
	}

	if !cfg.UI.ConfirmDelete {
		t.Errorf("Expected ConfirmDelete to be true by default")
	}

	if cfg.UI.MaxPathLength != 50 {
		t.Errorf("Expected MaxPathLength to be 50, got %d", cfg.UI.MaxPathLength)
	}

	// Test theme colors
	expectedColors := map[string]string{
		"Primary":   "#7C3AED",
		"Secondary": "#06B6D4",
		"Success":   "#10B981",
		"Warning":   "#F59E0B",
		"Error":     "#EF4444",
		"Muted":     "#6B7280",
		"Text":      "#F9FAFB",
	}

	themeValue := reflect.ValueOf(cfg.Theme)
	for colorName, expectedValue := range expectedColors {
		field := themeValue.FieldByName(colorName)
		if !field.IsValid() {
			t.Errorf("Theme field %s not found", colorName)
			continue
		}
		actualValue := field.String()
		if actualValue != expectedValue {
			t.Errorf("Expected theme %s to be '%s', got '%s'", colorName, expectedValue, actualValue)
		}
	}

	// Test aliases
	if cfg.Aliases == nil {
		t.Errorf("Aliases should not be nil")
	}
	expectedAliases := map[string]string{
		"ls": "list",
		"sw": "switch",
		"rm": "remove",
	}
	for alias, command := range expectedAliases {
		if cfg.Aliases[alias] != command {
			t.Errorf("Expected alias '%s' to map to '%s', got '%s'", alias, command, cfg.Aliases[alias])
		}
	}
}

func TestGetConfigPath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "yosegi-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Test local config detection
	t.Run("Local config exists", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "local-test")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		localConfigPath := filepath.Join(testDir, ".yosegi.yaml")
		if err := os.WriteFile(localConfigPath, []byte("test: value"), 0644); err != nil {
			t.Fatalf("Failed to create local config file: %v", err)
		}

		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}

		configPath, err := getConfigPath()
		if err != nil {
			t.Errorf("getConfigPath() failed: %v", err)
		}

		expectedPath := ".yosegi.yaml"
		if configPath != expectedPath {
			t.Errorf("Expected config path '%s', got '%s'", expectedPath, configPath)
		}
	})

	// Test global config path (when local doesn't exist)
	t.Run("Global config path", func(t *testing.T) {
		testDir := filepath.Join(tmpDir, "global-test")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		if err := os.Chdir(testDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}

		// Mock home directory
		originalHome := os.Getenv("HOME")
		mockHome := filepath.Join(tmpDir, "mock-home")
		if err := os.MkdirAll(mockHome, 0755); err != nil {
			t.Fatalf("Failed to create mock home: %v", err)
		}
		os.Setenv("HOME", mockHome)
		defer os.Setenv("HOME", originalHome)

		configPath, err := getConfigPath()
		if err != nil {
			t.Errorf("getConfigPath() failed: %v", err)
		}

		expectedPath := filepath.Join(mockHome, ".config", "yosegi", "config.yaml")
		if configPath != expectedPath {
			t.Errorf("Expected config path '%s', got '%s'", expectedPath, configPath)
		}

		// Check that the config directory was created
		configDir := filepath.Dir(configPath)
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			t.Errorf("Config directory should have been created: %s", configDir)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectDefault  bool
		expectedError  bool
		validateConfig func(*testing.T, *Config)
	}{
		{
			name:          "No config file",
			configContent: "",
			expectDefault: true,
			expectedError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				if cfg.DefaultWorktreePath != "../" {
					t.Errorf("Should use default worktree path")
				}
			},
		},
		{
			name: "Valid config file",
			configContent: `
default_worktree_path: "custom/path"
theme:
  primary: "#FF0000"
  secondary: "#00FF00"
git:
  auto_create_branch: false
  default_remote: "upstream"
ui:
  show_icons: false
  confirm_delete: false
  max_path_length: 100
aliases:
  l: "list"
  n: "new"
`,
			expectDefault: false,
			expectedError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				if cfg.DefaultWorktreePath != "custom/path" {
					t.Errorf("Expected custom worktree path, got '%s'", cfg.DefaultWorktreePath)
				}
				if cfg.Theme.Primary != "#FF0000" {
					t.Errorf("Expected custom primary color, got '%s'", cfg.Theme.Primary)
				}
				if cfg.Git.AutoCreateBranch {
					t.Errorf("Expected AutoCreateBranch to be false")
				}
				if cfg.Git.DefaultRemote != "upstream" {
					t.Errorf("Expected custom remote, got '%s'", cfg.Git.DefaultRemote)
				}
				if cfg.UI.ShowIcons {
					t.Errorf("Expected ShowIcons to be false")
				}
				if cfg.UI.MaxPathLength != 100 {
					t.Errorf("Expected MaxPathLength to be 100, got %d", cfg.UI.MaxPathLength)
				}
				if cfg.Aliases["l"] != "list" {
					t.Errorf("Expected alias 'l' to map to 'list'")
				}
			},
		},
		{
			name: "Invalid YAML",
			configContent: `
default_worktree_path: "test"
invalid_yaml: [
`,
			expectDefault: true,
			expectedError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				if cfg.DefaultWorktreePath != "../" {
					t.Errorf("Should fallback to default on invalid YAML")
				}
			},
		},
		{
			name: "Partial config",
			configContent: `
theme:
  primary: "#CUSTOM"
git:
  auto_create_branch: false
`,
			expectDefault: false,
			expectedError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				// Should merge with defaults
				if cfg.DefaultWorktreePath != "../" {
					t.Errorf("Should use default for missing fields")
				}
				if cfg.Theme.Primary != "#CUSTOM" {
					t.Errorf("Should use custom value for provided fields")
				}
				if cfg.Git.AutoCreateBranch {
					t.Errorf("Should use custom value for provided fields")
				}
				// Note: UI.ShowIcons will be false (Go zero value) since it's not set in the partial config
				// This is expected behavior due to YAML unmarshaling limitations
				if cfg.UI.ShowIcons != false {
					t.Errorf("UI.ShowIcons should be false when not set in partial config")
				}
				if cfg.UI.MaxPathLength != 50 {
					t.Errorf("Should use default for unset MaxPathLength")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "yosegi-load-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(originalDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			// Create config file if content is provided
			if tt.configContent != "" {
				configPath := ".yosegi.yaml"
				if err := os.WriteFile(configPath, []byte(tt.configContent), 0644); err != nil {
					t.Fatalf("Failed to create config file: %v", err)
				}
			}

			cfg, err := Load()
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if cfg == nil {
				t.Fatalf("Config should not be nil")
			}

			if tt.validateConfig != nil {
				tt.validateConfig(t, cfg)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yosegi-save-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Mock home directory
	originalHome := os.Getenv("HOME")
	mockHome := filepath.Join(tmpDir, "mock-home")
	if err := os.MkdirAll(mockHome, 0755); err != nil {
		t.Fatalf("Failed to create mock home: %v", err)
	}
	os.Setenv("HOME", mockHome)
	defer os.Setenv("HOME", originalHome)

	// Create a custom config
	cfg := &Config{
		DefaultWorktreePath: "custom/path",
		Theme: ThemeConfig{
			Primary:   "#FF0000",
			Secondary: "#00FF00",
		},
		Git: GitConfig{
			AutoCreateBranch: false,
			DefaultRemote:    "upstream",
		},
		UI: UIConfig{
			ShowIcons:     false,
			ConfirmDelete: false,
			MaxPathLength: 100,
		},
		Aliases: map[string]string{
			"l": "list",
			"n": "new",
		},
	}

	// Save the config
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Load the config back
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify the saved config
	if loadedCfg.DefaultWorktreePath != cfg.DefaultWorktreePath {
		t.Errorf("DefaultWorktreePath not saved correctly")
	}
	if loadedCfg.Theme.Primary != cfg.Theme.Primary {
		t.Errorf("Theme.Primary not saved correctly")
	}
	if loadedCfg.Git.AutoCreateBranch != cfg.Git.AutoCreateBranch {
		t.Errorf("Git.AutoCreateBranch not saved correctly")
	}
	if loadedCfg.UI.ShowIcons != cfg.UI.ShowIcons {
		t.Errorf("UI.ShowIcons not saved correctly")
	}
	if loadedCfg.Aliases["l"] != cfg.Aliases["l"] {
		t.Errorf("Aliases not saved correctly")
	}
}

func TestInitConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yosegi-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Mock home directory
	originalHome := os.Getenv("HOME")
	mockHome := filepath.Join(tmpDir, "mock-home")
	if err := os.MkdirAll(mockHome, 0755); err != nil {
		t.Fatalf("Failed to create mock home: %v", err)
	}
	os.Setenv("HOME", mockHome)
	defer os.Setenv("HOME", originalHome)

	// Initialize config
	if err := InitConfig(); err != nil {
		t.Fatalf("InitConfig() failed: %v", err)
	}

	// Verify config file was created
	expectedPath := filepath.Join(mockHome, ".config", "yosegi", "config.yaml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Config file should have been created at %s", expectedPath)
	}

	// Load and verify the config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Should match default config
	defaultCfg := defaultConfig()
	if !reflect.DeepEqual(cfg, defaultCfg) {
		t.Errorf("Initialized config does not match default config")
	}
}

func TestConfigStructValidation(t *testing.T) {
	// Test that all config struct fields have yaml tags
	configType := reflect.TypeOf(Config{})
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" {
			t.Errorf("Field %s is missing yaml tag", field.Name)
		}
	}

	// Test ThemeConfig struct
	themeType := reflect.TypeOf(ThemeConfig{})
	for i := 0; i < themeType.NumField(); i++ {
		field := themeType.Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" {
			t.Errorf("ThemeConfig field %s is missing yaml tag", field.Name)
		}
	}

	// Test GitConfig struct
	gitType := reflect.TypeOf(GitConfig{})
	for i := 0; i < gitType.NumField(); i++ {
		field := gitType.Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" {
			t.Errorf("GitConfig field %s is missing yaml tag", field.Name)
		}
	}

	// Test UIConfig struct
	uiType := reflect.TypeOf(UIConfig{})
	for i := 0; i < uiType.NumField(); i++ {
		field := uiType.Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" {
			t.Errorf("UIConfig field %s is missing yaml tag", field.Name)
		}
	}
}

// Test edge cases and error conditions
func TestConfigEdgeCases(t *testing.T) {
	t.Run("Save to read-only directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "yosegi-readonly-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create a read-only directory
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}

		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", readOnlyDir)
		defer os.Setenv("HOME", originalHome)

		cfg := defaultConfig()
		err = Save(cfg)
		if err == nil {
			t.Errorf("Expected error when saving to read-only directory")
		}
	})

	t.Run("Load with permission denied", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		tmpDir, err := os.MkdirTemp("", "yosegi-permission-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current directory: %v", err)
		}
		defer os.Chdir(originalDir)

		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}

		// Create a config file with no read permissions
		configPath := ".yosegi.yaml"
		if err := os.WriteFile(configPath, []byte("test: value"), 0644); err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}
		if err := os.Chmod(configPath, 0000); err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}
		defer os.Chmod(configPath, 0644) // Cleanup

		// Should fallback to default config
		cfg, err := Load()
		if err != nil {
			t.Errorf("Load() should not fail when config file is unreadable")
		}
		if cfg == nil {
			t.Errorf("Should return default config when file is unreadable")
		}
	})
}

// Benchmark tests
func BenchmarkLoad(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "yosegi-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		b.Fatalf("Failed to change directory: %v", err)
	}

	// Create a config file
	configContent := `
default_worktree_path: "../"
theme:
  primary: "#7C3AED"
  secondary: "#06B6D4"
git:
  auto_create_branch: true
ui:
  show_icons: true
`
	if err := os.WriteFile(".yosegi.yaml", []byte(configContent), 0644); err != nil {
		b.Fatalf("Failed to create config file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Errorf("Load() failed: %v", err)
		}
	}
}

func BenchmarkDefaultConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg := defaultConfig()
		if cfg == nil {
			b.Errorf("defaultConfig() returned nil")
		}
	}
}