package cmd

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/git"
)

func TestListCommand(t *testing.T) {
	// Test command structure
	if listCmd.Use != "list" {
		t.Errorf("Expected list command use to be 'list', got '%s'", listCmd.Use)
	}

	if listCmd.Short == "" {
		t.Errorf("List command should have short description")
	}

	if listCmd.Long == "" {
		t.Errorf("List command should have long description")
	}

	if listCmd.RunE == nil {
		t.Errorf("List command should have RunE function")
	}

	// Test aliases
	expectedAliases := []string{"ls", "l"}
	if len(listCmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(listCmd.Aliases))
	}

	for i, expected := range expectedAliases {
		if i < len(listCmd.Aliases) && listCmd.Aliases[i] != expected {
			t.Errorf("Expected alias %d to be '%s', got '%s'", i, expected, listCmd.Aliases[i])
		}
	}
}

func TestListCommandStructure(t *testing.T) {
	// Test that list command has proper configuration
	if listCmd.Use != "list" {
		t.Errorf("Expected listCmd use to be 'list', got '%s'", listCmd.Use)
	}

	if listCmd.Short == "" {
		t.Errorf("listCmd should have short description")
	}

	if listCmd.Long == "" {
		t.Errorf("listCmd should have long description")
	}

	if listCmd.RunE == nil {
		t.Errorf("listCmd should have RunE function")
	}

	// Test command content
	expectedShort := "List all git worktrees"
	if listCmd.Short != expectedShort {
		t.Errorf("Expected short description '%s', got '%s'", expectedShort, listCmd.Short)
	}

	expectedLong := "Display an interactive list of all git worktrees in the repository."
	if listCmd.Long != expectedLong {
		t.Errorf("Expected long description '%s', got '%s'", expectedLong, listCmd.Long)
	}
}

func TestListCommandAliases(t *testing.T) {
	// Test that all expected aliases are present
	expectedAliases := map[string]bool{
		"ls": false,
		"l":  false,
	}

	for _, alias := range listCmd.Aliases {
		if _, exists := expectedAliases[alias]; exists {
			expectedAliases[alias] = true
		} else {
			t.Errorf("Unexpected alias '%s'", alias)
		}
	}

	for alias, found := range expectedAliases {
		if !found {
			t.Errorf("Expected alias '%s' not found", alias)
		}
	}
}

func TestListCommandRegistered(t *testing.T) {
	// Check that list command is registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "list" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("List command should be registered with root command")
	}
}

func TestListCommandHelp(t *testing.T) {
	// Test list command help
	testCmd := &cobra.Command{
		Use:     listCmd.Use,
		Short:   listCmd.Short,
		Long:    listCmd.Long,
		Aliases: listCmd.Aliases,
	}

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)
	testCmd.SetArgs([]string{"--help"})

	err := testCmd.Execute()
	if err != nil {
		t.Errorf("List help should not error: %v", err)
	}

	output := buf.String()

	// Basic check that help was generated
	if len(output) == 0 {
		t.Error("Help output should not be empty")
	}
}

func TestListCommandByAlias(t *testing.T) {
	// Test that aliases work correctly by checking they're properly set
	aliases := listCmd.Aliases

	expectedAliases := []string{"ls", "l"}
	if len(aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(aliases))
		return
	}

	aliasMap := make(map[string]bool)
	for _, alias := range aliases {
		aliasMap[alias] = true
	}

	for _, expected := range expectedAliases {
		if !aliasMap[expected] {
			t.Errorf("Expected alias '%s' not found in list command", expected)
		}
	}
}

func TestListCommandFlags(t *testing.T) {
	// Test that list command has the expected flags
	flags := listCmd.Flags()

	// Test --print flag exists
	printFlag := flags.Lookup("print")
	if printFlag == nil {
		t.Error("List command should have --print flag")
	} else {
		// Test flag properties
		if printFlag.Shorthand != "p" {
			t.Errorf("Expected --print shorthand to be 'p', got '%s'", printFlag.Shorthand)
		}

		if printFlag.DefValue != "false" {
			t.Errorf("Expected --print default value to be 'false', got '%s'", printFlag.DefValue)
		}

		if printFlag.Usage == "" {
			t.Error("--print flag should have usage description")
		}
	}

	// Test that print flag is a bool
	if printFlag != nil {
		if printFlag.Value.Type() != "bool" {
			t.Errorf("Expected --print flag to be bool type, got '%s'", printFlag.Value.Type())
		}
	}
}

func TestListCommandValidation(t *testing.T) {
	// Test command validation
	if listCmd.Args != nil {
		// Test that args function works (if any)
		err := listCmd.Args(listCmd, []string{})
		if err != nil {
			t.Errorf("Args validation failed: %v", err)
		}
	}

	// Test that command is properly configured
	if listCmd.Use == "" {
		t.Error("List command Use should not be empty")
	}

	if listCmd.Short == "" {
		t.Error("List command Short description should not be empty")
	}

	// Test that RunE is set (required for execution)
	if listCmd.RunE == nil {
		t.Error("List command should have RunE function set")
	}
}

func TestListCommandInheritance(t *testing.T) {
	// Test that list command inherits from root command correctly
	parent := listCmd.Parent()
	if parent == nil {
		t.Error("List command should have a parent command")
		return
	}

	if parent.Name() != "yosegi" {
		t.Errorf("Expected parent command to be 'yosegi', got '%s'", parent.Name())
	}
}

func TestListCommandPath(t *testing.T) {
	// Test command path construction
	expectedPath := "yosegi list"
	actualPath := listCmd.CommandPath()

	if actualPath != expectedPath {
		t.Errorf("Expected command path '%s', got '%s'", expectedPath, actualPath)
	}
}

func TestListCommandUsage(t *testing.T) {
	// Test usage string generation
	usage := listCmd.UsageString()

	// Basic check that usage string is generated
	if len(usage) == 0 {
		t.Error("Usage string should not be empty")
	}

	// Should contain the command name
	if !strings.Contains(usage, "list") {
		t.Error("Usage should contain command name 'list'")
	}
}

func TestListCommandExample(t *testing.T) {
	// Test that command can be created and configured properly
	testCmd := &cobra.Command{
		Use:     listCmd.Use,
		Short:   listCmd.Short,
		Long:    listCmd.Long,
		Aliases: listCmd.Aliases,
	}

	// Test basic properties
	if testCmd.Name() != "list" {
		t.Errorf("Expected command name 'list', got '%s'", testCmd.Name())
	}

	if !testCmd.HasAlias("ls") {
		t.Error("Command should have alias 'ls'")
	}

	if !testCmd.HasAlias("l") {
		t.Error("Command should have alias 'l'")
	}
}

func TestListCommandBehavior(t *testing.T) {
	// Test command configuration behavior
	if listCmd.DisableFlagsInUseLine {
		t.Error("List command should not disable flags in usage line")
	}

	if listCmd.Hidden {
		t.Error("List command should not be hidden")
	}

	if listCmd.Deprecated != "" {
		t.Error("List command should not be deprecated")
	}
}

func TestListCommandCompletion(t *testing.T) {
	// Test command completion setup
	if listCmd.ValidArgsFunction != nil {
		// If completion is set up, test it
		completions, directive := listCmd.ValidArgsFunction(listCmd, []string{}, "")

		// List command shouldn't complete arguments since it takes none
		if len(completions) != 0 {
			t.Errorf("List command should not provide argument completions, got %d", len(completions))
		}

		// Directive should indicate no file completion
		if directive != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("Expected ShellCompDirectiveNoFileComp, got %d", directive)
		}
	}
}

func TestListCommandGroups(t *testing.T) {
	// Test command grouping (if any)
	if listCmd.GroupID != "" {
		// If grouped, verify it's in the correct group
		t.Logf("List command is in group: %s", listCmd.GroupID)
	}
}

func TestListCommandAnnotations(t *testing.T) {
	// Test command annotations (if any)
	if len(listCmd.Annotations) > 0 {
		t.Logf("List command has %d annotations", len(listCmd.Annotations))
		for key, value := range listCmd.Annotations {
			t.Logf("Annotation %s: %s", key, value)
		}
	}
}

// Test integration with parent command
func TestListCommandIntegration(t *testing.T) {
	// Verify the command is properly integrated
	found := false
	var foundCmd *cobra.Command

	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "list" {
			found = true
			foundCmd = cmd
			break
		}
	}

	if !found {
		t.Fatal("List command not found in root command")
	}

	// Test that it's the same command
	if foundCmd != listCmd {
		t.Error("Found command is not the same as listCmd")
	}

	// Test aliases are preserved
	if len(foundCmd.Aliases) != 2 {
		t.Errorf("Expected 2 aliases, got %d", len(foundCmd.Aliases))
	}
}

// Benchmark test
func BenchmarkListCommandCreation(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = &cobra.Command{
			Use:     "list",
			Short:   "List all git worktrees",
			Long:    "Display an interactive list of all git worktrees in the repository.",
			Aliases: []string{"ls", "l"},
		}
	}
}

func BenchmarkListCommandHelp(b *testing.B) {
	testCmd := &cobra.Command{
		Use:     listCmd.Use,
		Short:   listCmd.Short,
		Long:    listCmd.Long,
		Aliases: listCmd.Aliases,
	}

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)

	b.ResetTimer()
	for b.Loop() {
		buf.Reset()
		testCmd.SetArgs([]string{"--help"})
		if err := testCmd.Execute(); err != nil {
			// Log error but continue benchmark
			b.Logf("Command execution error: %v", err)
		}
	}
}

func TestRunRemoveWithSelectedWorktree(t *testing.T) {
	// Skip TUI tests on Windows CI
	if runtime.GOOS == "windows" && os.Getenv("CI") != "" {
		t.Skip("Skipping TUI test on Windows CI")
	}

	tests := []struct {
		name        string
		worktree    git.Worktree
		expectError bool
		errorMsg    string
	}{
		{
			name: "Cannot remove current worktree",
			worktree: git.Worktree{
				Path:      "/current/path",
				Branch:    "main",
				IsCurrent: true,
			},
			expectError: true,
			errorMsg:    "cannot remove current worktree",
		},
		{
			name: "Non-current worktree removal test structure",
			worktree: git.Worktree{
				Path:      "/test/path",
				Branch:    "feature",
				IsCurrent: false,
			},
			expectError: true, // Will fail in test environment
			errorMsg:    "failed to", // Generic error prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runRemoveWithSelectedWorktree(tt.worktree)
			
			if (err != nil) != tt.expectError {
				t.Errorf("runRemoveWithSelectedWorktree() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if err != nil && tt.errorMsg != "" {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestRunRemoveWithSelectedWorktreeCurrentCheck(t *testing.T) {
	// Test specifically for current worktree check
	currentWorktree := git.Worktree{
		Path:      "/current",
		Branch:    "main", 
		IsCurrent: true,
	}
	
	err := runRemoveWithSelectedWorktree(currentWorktree)
	
	if err == nil {
		t.Error("Expected error when trying to remove current worktree")
		return
	}
	
	expectedMsg := "cannot remove current worktree"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
