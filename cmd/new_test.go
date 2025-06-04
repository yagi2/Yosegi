package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCommand(t *testing.T) {
	// Test command structure
	if newCmd.Use != "new [branch]" {
		t.Errorf("Expected new command use to be 'new [branch]', got '%s'", newCmd.Use)
	}

	if newCmd.Short == "" {
		t.Errorf("New command should have short description")
	}

	if newCmd.Long == "" {
		t.Errorf("New command should have long description")
	}

	if newCmd.RunE == nil {
		t.Errorf("New command should have RunE function")
	}

	// Test aliases
	expectedAliases := []string{"add", "create", "n"}
	if len(newCmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(newCmd.Aliases))
	}

	for i, expected := range expectedAliases {
		if i < len(newCmd.Aliases) && newCmd.Aliases[i] != expected {
			t.Errorf("Expected alias %d to be '%s', got '%s'", i, expected, newCmd.Aliases[i])
		}
	}
}

func TestNewCommandStructure(t *testing.T) {
	// Test that new command has proper configuration
	if newCmd.Use != "new [branch]" {
		t.Errorf("Expected newCmd use to be 'new [branch]', got '%s'", newCmd.Use)
	}

	if newCmd.Short == "" {
		t.Errorf("newCmd should have short description")
	}

	if newCmd.Long == "" {
		t.Errorf("newCmd should have long description")
	}

	if newCmd.RunE == nil {
		t.Errorf("newCmd should have RunE function")
	}

	// Test command content
	expectedShort := "Create a new git worktree"
	if newCmd.Short != expectedShort {
		t.Errorf("Expected short description '%s', got '%s'", expectedShort, newCmd.Short)
	}

	expectedLong := "Create a new git worktree interactively or with specified parameters."
	if newCmd.Long != expectedLong {
		t.Errorf("Expected long description '%s', got '%s'", expectedLong, newCmd.Long)
	}
}

func TestNewCommandAliases(t *testing.T) {
	// Test that all expected aliases are present
	expectedAliases := map[string]bool{
		"add":    false,
		"create": false,
		"n":      false,
	}

	for _, alias := range newCmd.Aliases {
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

func TestNewCommandFlags(t *testing.T) {
	// Reset flags to default values before testing
	createBranch = false
	worktreePath = ""

	// Test create-branch flag
	createBranchFlag := newCmd.Flags().Lookup("create-branch")
	if createBranchFlag == nil {
		t.Error("create-branch flag should exist")
	} else {
		if createBranchFlag.Shorthand != "b" {
			t.Errorf("Expected create-branch flag shorthand to be 'b', got '%s'", createBranchFlag.Shorthand)
		}
		if createBranchFlag.DefValue != "false" {
			t.Errorf("Expected create-branch flag default to be 'false', got '%s'", createBranchFlag.DefValue)
		}
		if createBranchFlag.Usage == "" {
			t.Error("create-branch flag should have usage text")
		}
	}

	// Test path flag
	pathFlag := newCmd.Flags().Lookup("path")
	if pathFlag == nil {
		t.Error("path flag should exist")
	} else {
		if pathFlag.Shorthand != "p" {
			t.Errorf("Expected path flag shorthand to be 'p', got '%s'", pathFlag.Shorthand)
		}
		if pathFlag.DefValue != "" {
			t.Errorf("Expected path flag default to be empty, got '%s'", pathFlag.DefValue)
		}
		if pathFlag.Usage == "" {
			t.Error("path flag should have usage text")
		}
	}
}

func TestNewCommandFlagValues(t *testing.T) {
	// Test that flag variables are properly connected

	// Reset to known state
	createBranch = false
	worktreePath = ""

	// Test setting create-branch flag
	err := newCmd.Flags().Set("create-branch", "true")
	if err != nil {
		t.Errorf("Failed to set create-branch flag: %v", err)
	}
	if !createBranch {
		t.Error("createBranch variable should be true after setting flag")
	}

	// Test setting path flag
	err = newCmd.Flags().Set("path", "/test/path")
	if err != nil {
		t.Errorf("Failed to set path flag: %v", err)
	}
	if worktreePath != "/test/path" {
		t.Errorf("Expected worktreePath to be '/test/path', got '%s'", worktreePath)
	}

	// Reset flags
	createBranch = false
	worktreePath = ""
}

func TestNewCommandRegistered(t *testing.T) {
	// Check that new command is registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "new" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("New command should be registered with root command")
	}
}

func TestNewCommandHelp(t *testing.T) {
	// Test new command help
	testCmd := &cobra.Command{
		Use:     newCmd.Use,
		Short:   newCmd.Short,
		Long:    newCmd.Long,
		Aliases: newCmd.Aliases,
	}

	// Add flags like the real command
	testCmd.Flags().BoolP("create-branch", "b", false, "Create a new branch")
	testCmd.Flags().StringP("path", "p", "", "Path for the new worktree")

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)
	testCmd.SetArgs([]string{"--help"})

	err := testCmd.Execute()
	if err != nil {
		t.Errorf("New help should not error: %v", err)
	}

	output := buf.String()

	// Basic check that help was generated
	if len(output) == 0 {
		t.Error("Help output should not be empty")
	}
}

func TestNewCommandByAlias(t *testing.T) {
	// Test that aliases work correctly by checking they're properly set
	aliases := newCmd.Aliases

	expectedAliases := []string{"add", "create", "n"}
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
			t.Errorf("Expected alias '%s' not found in new command", expected)
		}
	}
}

func TestNewCommandValidation(t *testing.T) {
	// Test command validation
	if newCmd.Args != nil {
		// Test that args function works (if any)
		err := newCmd.Args(newCmd, []string{"test-branch"})
		if err != nil {
			t.Errorf("Args validation failed: %v", err)
		}
	}

	// Test that command is properly configured
	if newCmd.Use == "" {
		t.Error("New command Use should not be empty")
	}

	if newCmd.Short == "" {
		t.Error("New command Short description should not be empty")
	}

	// Test that RunE is set (required for execution)
	if newCmd.RunE == nil {
		t.Error("New command should have RunE function set")
	}
}

func TestNewCommandInheritance(t *testing.T) {
	// Test that new command inherits from root command correctly
	parent := newCmd.Parent()
	if parent == nil {
		t.Error("New command should have a parent command")
		return
	}

	if parent.Name() != "yosegi" {
		t.Errorf("Expected parent command to be 'yosegi', got '%s'", parent.Name())
	}
}

func TestNewCommandPath(t *testing.T) {
	// Test command path construction
	expectedPath := "yosegi new"
	actualPath := newCmd.CommandPath()

	if actualPath != expectedPath {
		t.Errorf("Expected command path '%s', got '%s'", expectedPath, actualPath)
	}
}

func TestNewCommandUsage(t *testing.T) {
	// Test usage string generation
	usage := newCmd.UsageString()

	// Basic check that usage string is generated
	if len(usage) == 0 {
		t.Error("Usage string should not be empty")
	}

	// Should contain the command name
	if !strings.Contains(usage, "new") {
		t.Error("Usage should contain command name 'new'")
	}
}

func TestNewCommandFlagShorthands(t *testing.T) {
	// Test flag shorthand functionality
	testCmd := &cobra.Command{
		Use: "new [branch]",
	}

	var testCreateBranch bool
	var testWorktreePath string

	testCmd.Flags().BoolVarP(&testCreateBranch, "create-branch", "b", false, "Create a new branch")
	testCmd.Flags().StringVarP(&testWorktreePath, "path", "p", "", "Path for the new worktree")

	// Test short flag parsing
	testCmd.SetArgs([]string{"-b", "-p", "/test/path", "test-branch"})
	if err := testCmd.ParseFlags([]string{"-b", "-p", "/test/path"}); err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	if !testCreateBranch {
		t.Error("Short flag -b should set create-branch to true")
	}

	if testWorktreePath != "/test/path" {
		t.Errorf("Short flag -p should set path to '/test/path', got '%s'", testWorktreePath)
	}
}

func TestNewCommandBehavior(t *testing.T) {
	// Test command configuration behavior
	if newCmd.DisableFlagsInUseLine {
		t.Error("New command should not disable flags in usage line")
	}

	if newCmd.Hidden {
		t.Error("New command should not be hidden")
	}

	if newCmd.Deprecated != "" {
		t.Error("New command should not be deprecated")
	}
}

func TestNewCommandCompletion(t *testing.T) {
	// Test command completion setup
	if newCmd.ValidArgsFunction != nil {
		// If completion is set up, test it
		completions, directive := newCmd.ValidArgsFunction(newCmd, []string{}, "")

		// We don't expect specific completions but should not error
		_ = completions
		_ = directive
	}
}

func TestNewCommandGroups(t *testing.T) {
	// Test command grouping (if any)
	if newCmd.GroupID != "" {
		// If grouped, verify it's in the correct group
		t.Logf("New command is in group: %s", newCmd.GroupID)
	}
}

func TestNewCommandAnnotations(t *testing.T) {
	// Test command annotations (if any)
	if len(newCmd.Annotations) > 0 {
		t.Logf("New command has %d annotations", len(newCmd.Annotations))
		for key, value := range newCmd.Annotations {
			t.Logf("Annotation %s: %s", key, value)
		}
	}
}

// Test integration with parent command
func TestNewCommandIntegration(t *testing.T) {
	// Verify the command is properly integrated
	found := false
	var foundCmd *cobra.Command

	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "new [branch]" {
			found = true
			foundCmd = cmd
			break
		}
	}

	if !found {
		t.Fatal("New command not found in root command")
	}

	// Test that it's the same command
	if foundCmd != newCmd {
		t.Error("Found command is not the same as newCmd")
	}

	// Test aliases are preserved
	if len(foundCmd.Aliases) != 3 {
		t.Errorf("Expected 3 aliases, got %d", len(foundCmd.Aliases))
	}

	// Test flags are preserved
	if !foundCmd.Flags().HasFlags() {
		t.Error("New command should have flags")
	}
}

func TestNewCommandExample(t *testing.T) {
	// Test that command can be created and configured properly
	testCmd := &cobra.Command{
		Use:     newCmd.Use,
		Short:   newCmd.Short,
		Long:    newCmd.Long,
		Aliases: newCmd.Aliases,
	}

	// Test basic properties
	if testCmd.Name() != "new" {
		t.Errorf("Expected command name 'new', got '%s'", testCmd.Name())
	}

	if !testCmd.HasAlias("add") {
		t.Error("Command should have alias 'add'")
	}

	if !testCmd.HasAlias("create") {
		t.Error("Command should have alias 'create'")
	}

	if !testCmd.HasAlias("n") {
		t.Error("Command should have alias 'n'")
	}
}

func TestNewCommandGlobalVariables(t *testing.T) {
	// Test that global variables are properly initialized
	originalCreateBranch := createBranch
	originalWorktreePath := worktreePath

	// Test default values
	if createBranch != false {
		t.Error("createBranch should default to false")
	}

	if worktreePath != "" {
		t.Error("worktreePath should default to empty string")
	}

	// Restore original values
	createBranch = originalCreateBranch
	worktreePath = originalWorktreePath
}

func TestNewCommandArgsAcceptance(t *testing.T) {
	// Test that command accepts the expected arguments
	// New command should accept 0 or 1 argument (branch name)

	// Test with no args (should be valid)
	if newCmd.Args != nil {
		err := newCmd.Args(newCmd, []string{})
		if err != nil {
			t.Errorf("New command should accept no arguments, got error: %v", err)
		}
	}

	// Test with one arg (should be valid)
	if newCmd.Args != nil {
		err := newCmd.Args(newCmd, []string{"feature-branch"})
		if err != nil {
			t.Errorf("New command should accept one argument, got error: %v", err)
		}
	}
}

// Benchmark tests
func BenchmarkNewCommandCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := &cobra.Command{
			Use:     "new [branch]",
			Short:   "Create a new git worktree",
			Long:    "Create a new git worktree interactively or with specified parameters.",
			Aliases: []string{"add", "create", "n"},
		}

		var testCreateBranch bool
		var testWorktreePath string
		cmd.Flags().BoolVarP(&testCreateBranch, "create-branch", "b", false, "Create a new branch")
		cmd.Flags().StringVarP(&testWorktreePath, "path", "p", "", "Path for the new worktree")

		_ = cmd
	}
}

func BenchmarkNewCommandHelp(b *testing.B) {
	testCmd := &cobra.Command{
		Use:     newCmd.Use,
		Short:   newCmd.Short,
		Long:    newCmd.Long,
		Aliases: newCmd.Aliases,
	}
	testCmd.Flags().BoolP("create-branch", "b", false, "Create a new branch")
	testCmd.Flags().StringP("path", "p", "", "Path for the new worktree")

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		testCmd.SetArgs([]string{"--help"})
		if err := testCmd.Execute(); err != nil {
			// Log error but continue benchmark
			b.Logf("Command execution error: %v", err)
		}
	}
}

func BenchmarkNewCommandFlagParsing(b *testing.B) {
	testCmd := &cobra.Command{
		Use: "new [branch]",
	}

	var testCreateBranch bool
	var testWorktreePath string
	testCmd.Flags().BoolVarP(&testCreateBranch, "create-branch", "b", false, "Create a new branch")
	testCmd.Flags().StringVarP(&testWorktreePath, "path", "p", "", "Path for the new worktree")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := testCmd.ParseFlags([]string{"-b", "-p", "/test/path"}); err != nil {
			// Log error but continue benchmark
			b.Logf("Flag parsing error: %v", err)
		}
	}
}
