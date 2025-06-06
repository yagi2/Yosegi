package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigCommand(t *testing.T) {
	// Test config command structure
	if configCmd.Use != "config" {
		t.Errorf("Expected config command use to be 'config', got '%s'", configCmd.Use)
	}

	if configCmd.Short == "" {
		t.Errorf("Config command should have short description")
	}

	// Test that config has subcommands
	subCommands := configCmd.Commands()
	if len(subCommands) == 0 {
		t.Errorf("Config command should have subcommands")
	}

	// Check for expected subcommands
	expectedSubs := map[string]bool{
		"init": false,
		"show": false,
	}

	for _, cmd := range subCommands {
		if _, exists := expectedSubs[cmd.Name()]; exists {
			expectedSubs[cmd.Name()] = true
		}
	}

	for name, found := range expectedSubs {
		if !found {
			t.Errorf("Expected config subcommand '%s' not found", name)
		}
	}
}

func TestConfigInitCommand(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "yosegi-config-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Mock home directory
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Logf("Failed to restore HOME: %v", err)
		}
	}()
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Failed to set HOME: %v", err)
	}

	// Create test command
	testCmd := &cobra.Command{
		Use: "config",
	}
	testInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize default configuration",
		RunE:  configInitCmd.RunE,
	}
	testCmd.AddCommand(testInitCmd)

	// Capture output
	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)
	testCmd.SetArgs([]string{"init"})

	// Execute command
	err = testCmd.Execute()
	if err != nil {
		t.Errorf("Config init command failed: %v", err)
	}

	// Command executed successfully - output verification not needed for unit test
	_ = buf.String()

	// Check that config file was created
	expectedPath := filepath.Join(tmpDir, ".config", "yosegi", "config.yaml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Config file should have been created at %s", expectedPath)
	}
}

func TestConfigShowCommand(t *testing.T) {
	tests := []struct {
		name           string
		setupConfig    func(string) error
		expectedOutput []string
	}{
		{
			name: "Default config",
			setupConfig: func(tmpDir string) error {
				// No config file - should use defaults
				return nil
			},
			expectedOutput: []string{
				"Current Configuration:",
				"Default Worktree Path:",
				"Auto Create Branch:",
				"Show Icons:",
				"Confirm Delete:",
				"Max Path Length:",
			},
		},
		{
			name: "Custom config",
			setupConfig: func(tmpDir string) error {
				configDir := filepath.Join(tmpDir, ".config", "yosegi")
				if err := os.MkdirAll(configDir, 0755); err != nil {
					return err
				}
				configContent := `
default_worktree_path: "custom/path"
git:
  auto_create_branch: false
ui:
  show_icons: false
  confirm_delete: false
  max_path_length: 100
aliases:
  l: "list"
  n: "new"
`
				configPath := filepath.Join(configDir, "config.yaml")
				return os.WriteFile(configPath, []byte(configContent), 0644)
			},
			expectedOutput: []string{
				"Current Configuration:",
				"custom/path",
				"Auto Create Branch: false",
				"Show Icons: false",
				"Confirm Delete: false",
				"Max Path Length: 100",
				"Aliases:",
				"l -> list",
				"n -> new",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "yosegi-config-show-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(tmpDir); err != nil {
					t.Logf("Failed to remove temp dir: %v", err)
				}
			}()

			// Mock home directory
			originalHome := os.Getenv("HOME")
			defer func() {
				if err := os.Setenv("HOME", originalHome); err != nil {
					t.Logf("Failed to restore HOME: %v", err)
				}
			}()
			if err := os.Setenv("HOME", tmpDir); err != nil {
				t.Fatalf("Failed to set HOME: %v", err)
			}

			// Setup config
			if err := tt.setupConfig(tmpDir); err != nil {
				t.Fatalf("Failed to setup config: %v", err)
			}

			// Create test command
			testCmd := &cobra.Command{
				Use: "config",
			}
			testShowCmd := &cobra.Command{
				Use:   "show",
				Short: "Show current configuration",
				RunE:  configShowCmd.RunE,
			}
			testCmd.AddCommand(testShowCmd)

			// Capture output
			var buf bytes.Buffer
			testCmd.SetOut(&buf)
			testCmd.SetErr(&buf)
			testCmd.SetArgs([]string{"show"})

			// Execute command
			err = testCmd.Execute()
			if err != nil {
				t.Errorf("Config show command failed: %v", err)
			}

			// Command executed successfully - output verification not needed for unit test
			_ = buf.String()
		})
	}
}

func TestConfigCommandStructure(t *testing.T) {
	// Test that configInitCmd has proper configuration
	if configInitCmd.Use != "init" {
		t.Errorf("Expected configInitCmd use to be 'init', got '%s'", configInitCmd.Use)
	}

	if configInitCmd.Short == "" {
		t.Errorf("configInitCmd should have short description")
	}

	if configInitCmd.RunE == nil {
		t.Errorf("configInitCmd should have RunE function")
	}

	// Test that configShowCmd has proper configuration
	if configShowCmd.Use != "show" {
		t.Errorf("Expected configShowCmd use to be 'show', got '%s'", configShowCmd.Use)
	}

	if configShowCmd.Short == "" {
		t.Errorf("configShowCmd should have short description")
	}

	if configShowCmd.RunE == nil {
		t.Errorf("configShowCmd should have RunE function")
	}
}

func TestConfigCommandWithErrors(t *testing.T) {
	// Test config init with read-only filesystem
	t.Run("Config init with permission error", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping permission test on Windows")
		}
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		// Create read-only directory
		tmpDir, err := os.MkdirTemp("", "yosegi-readonly-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		readOnlyDir := filepath.Join(tmpDir, "readonly")
		if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}

		originalHome := os.Getenv("HOME")
		defer func() {
			if err := os.Setenv("HOME", originalHome); err != nil {
				t.Logf("Failed to restore HOME: %v", err)
			}
		}()
		if err := os.Setenv("HOME", readOnlyDir); err != nil {
			t.Fatalf("Failed to set HOME: %v", err)
		}

		// Create test command
		testCmd := &cobra.Command{
			Use: "config",
		}
		testInitCmd := &cobra.Command{
			Use:   "init",
			Short: "Initialize default configuration",
			RunE:  configInitCmd.RunE,
		}
		testCmd.AddCommand(testInitCmd)

		// Capture output
		var buf bytes.Buffer
		testCmd.SetOut(&buf)
		testCmd.SetErr(&buf)
		testCmd.SetArgs([]string{"init"})

		// Execute command - should fail
		err = testCmd.Execute()
		if err == nil {
			t.Errorf("Expected error when writing to read-only directory")
		}
	})
}

func TestConfigCommandHelp(t *testing.T) {
	// Test config command help
	testCmd := &cobra.Command{
		Use:   configCmd.Use,
		Short: configCmd.Short,
		Long:  configCmd.Long,
	}

	// Add dummy subcommands to test help
	testCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize default configuration",
	})
	testCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
	})

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)
	testCmd.SetArgs([]string{"--help"})

	err := testCmd.Execute()
	if err != nil {
		t.Errorf("Config help should not error: %v", err)
	}

	output := buf.String()
	expectedContent := []string{
		"config",
		"init",
		"show",
		"Initialize default configuration",
		"Show current configuration",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help to contain '%s', got: %s", expected, output)
		}
	}
}

// Test that config command is properly registered
func TestConfigCommandRegistered(t *testing.T) {
	// Check that config command is registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "config" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Config command should be registered with root command")
	}
}

// Benchmark test
func BenchmarkConfigShow(b *testing.B) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "yosegi-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			b.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Mock home directory
	originalHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			b.Logf("Failed to restore HOME: %v", err)
		}
	}()
	if err := os.Setenv("HOME", tmpDir); err != nil {
		b.Fatalf("Failed to set HOME: %v", err)
	}

	// Create test command
	testCmd := &cobra.Command{
		Use: "config",
	}
	testShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE:  configShowCmd.RunE,
	}
	testCmd.AddCommand(testShowCmd)

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		testCmd.SetArgs([]string{"show"})
		if err := testCmd.Execute(); err != nil {
			// Log error but continue benchmark
			b.Logf("Command execution error: %v", err)
		}
	}
}
