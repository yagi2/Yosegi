package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRemoveCommand(t *testing.T) {
	// Test command structure
	if removeCmd.Use != "remove" {
		t.Errorf("Expected remove command use to be 'remove', got '%s'", removeCmd.Use)
	}

	if removeCmd.Short == "" {
		t.Errorf("Remove command should have short description")
	}

	if removeCmd.Long == "" {
		t.Errorf("Remove command should have long description")
	}

	if removeCmd.RunE == nil {
		t.Errorf("Remove command should have RunE function")
	}

	// Test aliases
	expectedAliases := []string{"rm", "delete", "del", "r"}
	if len(removeCmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(removeCmd.Aliases))
	}

	for i, expected := range expectedAliases {
		if i < len(removeCmd.Aliases) && removeCmd.Aliases[i] != expected {
			t.Errorf("Expected alias %d to be '%s', got '%s'", i, expected, removeCmd.Aliases[i])
		}
	}
}

func TestRemoveCommandStructure(t *testing.T) {
	// Test that remove command has proper configuration
	if removeCmd.Use != "remove" {
		t.Errorf("Expected removeCmd use to be 'remove', got '%s'", removeCmd.Use)
	}

	if removeCmd.Short == "" {
		t.Errorf("removeCmd should have short description")
	}

	if removeCmd.Long == "" {
		t.Errorf("removeCmd should have long description")
	}

	if removeCmd.RunE == nil {
		t.Errorf("removeCmd should have RunE function")
	}

	// Test command content
	expectedShort := "Remove a git worktree"
	if removeCmd.Short != expectedShort {
		t.Errorf("Expected short description '%s', got '%s'", expectedShort, removeCmd.Short)
	}

	expectedLong := "Interactively select and remove a git worktree."
	if removeCmd.Long != expectedLong {
		t.Errorf("Expected long description '%s', got '%s'", expectedLong, removeCmd.Long)
	}
}

func TestRemoveCommandAliases(t *testing.T) {
	// Test that all expected aliases are present
	expectedAliases := map[string]bool{
		"rm":     false,
		"delete": false,
		"del":    false,
		"r":      false,
	}

	for _, alias := range removeCmd.Aliases {
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

func TestRemoveCommandFlags(t *testing.T) {
	// Reset flags to default values before testing
	forceRemove = false

	// Test force flag
	forceFlag := removeCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("force flag should exist")
	} else {
		if forceFlag.Shorthand != "f" {
			t.Errorf("Expected force flag shorthand to be 'f', got '%s'", forceFlag.Shorthand)
		}
		if forceFlag.DefValue != "false" {
			t.Errorf("Expected force flag default to be 'false', got '%s'", forceFlag.DefValue)
		}
		if forceFlag.Usage == "" {
			t.Error("force flag should have usage text")
		}
	}
}

func TestRemoveCommandFlagValues(t *testing.T) {
	// Test that flag variables are properly connected
	
	// Reset to known state
	forceRemove = false

	// Test setting force flag
	err := removeCmd.Flags().Set("force", "true")
	if err != nil {
		t.Errorf("Failed to set force flag: %v", err)
	}
	if !forceRemove {
		t.Error("forceRemove variable should be true after setting flag")
	}

	// Reset flag
	forceRemove = false
}

func TestRemoveCommandRegistered(t *testing.T) {
	// Check that remove command is registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "remove" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Remove command should be registered with root command")
	}
}

func TestRemoveCommandHelp(t *testing.T) {
	// Test remove command help
	testCmd := &cobra.Command{
		Use:     removeCmd.Use,
		Short:   removeCmd.Short,
		Long:    removeCmd.Long,
		Aliases: removeCmd.Aliases,
	}

	// Add flags like the real command
	testCmd.Flags().BoolP("force", "f", false, "Force removal even if worktree is dirty")

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)
	testCmd.SetArgs([]string{"--help"})

	err := testCmd.Execute()
	if err != nil {
		t.Errorf("Remove help should not error: %v", err)
	}

	output := buf.String()
	
	// Basic check that help was generated
	if len(output) == 0 {
		t.Error("Help output should not be empty")
	}
}

func TestRemoveCommandByAlias(t *testing.T) {
	// Test that aliases work correctly by checking they're properly set
	aliases := removeCmd.Aliases
	
	expectedAliases := []string{"rm", "delete", "del", "r"}
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
			t.Errorf("Expected alias '%s' not found in remove command", expected)
		}
	}
}

func TestRemoveCommandValidation(t *testing.T) {
	// Test command validation
	if removeCmd.Args != nil {
		// Test that args function works (if any)
		err := removeCmd.Args(removeCmd, []string{})
		if err != nil {
			t.Errorf("Args validation failed: %v", err)
		}
	}

	// Test that command is properly configured
	if removeCmd.Use == "" {
		t.Error("Remove command Use should not be empty")
	}

	if removeCmd.Short == "" {
		t.Error("Remove command Short description should not be empty")
	}

	// Test that RunE is set (required for execution)
	if removeCmd.RunE == nil {
		t.Error("Remove command should have RunE function set")
	}
}

func TestRemoveCommandInheritance(t *testing.T) {
	// Test that remove command inherits from root command correctly
	parent := removeCmd.Parent()
	if parent == nil {
		t.Error("Remove command should have a parent command")
		return
	}

	if parent.Name() != "yosegi" {
		t.Errorf("Expected parent command to be 'yosegi', got '%s'", parent.Name())
	}
}

func TestRemoveCommandPath(t *testing.T) {
	// Test command path construction
	expectedPath := "yosegi remove"
	actualPath := removeCmd.CommandPath()
	
	if actualPath != expectedPath {
		t.Errorf("Expected command path '%s', got '%s'", expectedPath, actualPath)
	}
}

func TestRemoveCommandUsage(t *testing.T) {
	// Test usage string generation
	usage := removeCmd.UsageString()
	
	// Basic check that usage string is generated
	if len(usage) == 0 {
		t.Error("Usage string should not be empty")
	}
	
	// Should contain the command name
	if !strings.Contains(usage, "remove") {
		t.Error("Usage should contain command name 'remove'")
	}
}

func TestRemoveCommandFlagShorthands(t *testing.T) {
	// Test flag shorthand functionality
	testCmd := &cobra.Command{
		Use: "remove",
	}
	
	var testForceRemove bool
	testCmd.Flags().BoolVarP(&testForceRemove, "force", "f", false, "Force removal even if worktree is dirty")

	// Test short flag parsing
	testCmd.SetArgs([]string{"-f"})
	testCmd.ParseFlags([]string{"-f"})

	if !testForceRemove {
		t.Error("Short flag -f should set force to true")
	}
}

func TestRemoveCommandBehavior(t *testing.T) {
	// Test command configuration behavior
	if removeCmd.DisableFlagsInUseLine {
		t.Error("Remove command should not disable flags in usage line")
	}

	if removeCmd.Hidden {
		t.Error("Remove command should not be hidden")
	}

	if removeCmd.Deprecated != "" {
		t.Error("Remove command should not be deprecated")
	}
}

func TestRemoveCommandCompletion(t *testing.T) {
	// Test command completion setup
	if removeCmd.ValidArgsFunction != nil {
		// If completion is set up, test it
		completions, directive := removeCmd.ValidArgsFunction(removeCmd, []string{}, "")
		
		// Remove command shouldn't complete arguments since it's interactive
		if len(completions) != 0 {
			t.Errorf("Remove command should not provide argument completions, got %d", len(completions))
		}
		
		// Directive should indicate no file completion
		if directive != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("Expected ShellCompDirectiveNoFileComp, got %d", directive)
		}
	}
}

func TestRemoveCommandGroups(t *testing.T) {
	// Test command grouping (if any)
	if removeCmd.GroupID != "" {
		// If grouped, verify it's in the correct group
		t.Logf("Remove command is in group: %s", removeCmd.GroupID)
	}
}

func TestRemoveCommandAnnotations(t *testing.T) {
	// Test command annotations (if any)
	if removeCmd.Annotations != nil && len(removeCmd.Annotations) > 0 {
		t.Logf("Remove command has %d annotations", len(removeCmd.Annotations))
		for key, value := range removeCmd.Annotations {
			t.Logf("Annotation %s: %s", key, value)
		}
	}
}

// Test integration with parent command
func TestRemoveCommandIntegration(t *testing.T) {
	// Verify the command is properly integrated
	found := false
	var foundCmd *cobra.Command
	
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "remove" {
			found = true
			foundCmd = cmd
			break
		}
	}
	
	if !found {
		t.Fatal("Remove command not found in root command")
	}
	
	// Test that it's the same command
	if foundCmd != removeCmd {
		t.Error("Found command is not the same as removeCmd")
	}
	
	// Test aliases are preserved
	if len(foundCmd.Aliases) != 4 {
		t.Errorf("Expected 4 aliases, got %d", len(foundCmd.Aliases))
	}
	
	// Test flags are preserved
	if !foundCmd.Flags().HasFlags() {
		t.Error("Remove command should have flags")
	}
}

func TestRemoveCommandExample(t *testing.T) {
	// Test that command can be created and configured properly
	testCmd := &cobra.Command{
		Use:     removeCmd.Use,
		Short:   removeCmd.Short,
		Long:    removeCmd.Long,
		Aliases: removeCmd.Aliases,
	}

	// Test basic properties
	if testCmd.Name() != "remove" {
		t.Errorf("Expected command name 'remove', got '%s'", testCmd.Name())
	}

	if !testCmd.HasAlias("rm") {
		t.Error("Command should have alias 'rm'")
	}

	if !testCmd.HasAlias("delete") {
		t.Error("Command should have alias 'delete'")
	}

	if !testCmd.HasAlias("del") {
		t.Error("Command should have alias 'del'")
	}

	if !testCmd.HasAlias("r") {
		t.Error("Command should have alias 'r'")
	}
}

func TestRemoveCommandGlobalVariables(t *testing.T) {
	// Test that global variables are properly initialized
	originalForceRemove := forceRemove

	// Test default values
	if forceRemove != false {
		t.Error("forceRemove should default to false")
	}

	// Restore original values
	forceRemove = originalForceRemove
}

func TestRemoveCommandArgsAcceptance(t *testing.T) {
	// Test that command accepts the expected arguments
	// Remove command should accept no arguments (interactive mode)
	
	// Test with no args (should be valid)
	if removeCmd.Args != nil {
		err := removeCmd.Args(removeCmd, []string{})
		if err != nil {
			t.Errorf("Remove command should accept no arguments, got error: %v", err)
		}
	}
}

func TestRemoveCommandFlagUsage(t *testing.T) {
	// Test that force flag has proper usage description
	forceFlag := removeCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Fatal("force flag should exist")
	}

	expectedUsage := "Force removal even if worktree is dirty"
	if forceFlag.Usage != expectedUsage {
		t.Errorf("Expected force flag usage '%s', got '%s'", expectedUsage, forceFlag.Usage)
	}
}

func TestRemoveCommandFlagDefaults(t *testing.T) {
	// Test flag default values
	forceFlag := removeCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Fatal("force flag should exist")
	}

	if forceFlag.DefValue != "false" {
		t.Errorf("Expected force flag default value 'false', got '%s'", forceFlag.DefValue)
	}

	// Test flag type
	if forceFlag.Value.Type() != "bool" {
		t.Errorf("Expected force flag type 'bool', got '%s'", forceFlag.Value.Type())
	}
}

func TestRemoveCommandFlagBinding(t *testing.T) {
	// Test that flags are properly bound to variables
	originalForceRemove := forceRemove
	
	// Reset to known state
	forceRemove = false
	
	// Set flag and check variable
	err := removeCmd.Flags().Set("force", "true")
	if err != nil {
		t.Fatalf("Failed to set force flag: %v", err)
	}
	
	if !forceRemove {
		t.Error("forceRemove variable should be updated when flag is set")
	}
	
	// Reset flag and check variable
	err = removeCmd.Flags().Set("force", "false")
	if err != nil {
		t.Fatalf("Failed to reset force flag: %v", err)
	}
	
	if forceRemove {
		t.Error("forceRemove variable should be updated when flag is reset")
	}
	
	// Restore original value
	forceRemove = originalForceRemove
}

// Benchmark tests
func BenchmarkRemoveCommandCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := &cobra.Command{
			Use:     "remove",
			Short:   "Remove a git worktree",
			Long:    "Interactively select and remove a git worktree.",
			Aliases: []string{"rm", "delete", "del", "r"},
		}
		
		var testForceRemove bool
		cmd.Flags().BoolVarP(&testForceRemove, "force", "f", false, "Force removal even if worktree is dirty")
		
		_ = cmd
	}
}

func BenchmarkRemoveCommandHelp(b *testing.B) {
	testCmd := &cobra.Command{
		Use:     removeCmd.Use,
		Short:   removeCmd.Short,
		Long:    removeCmd.Long,
		Aliases: removeCmd.Aliases,
	}
	testCmd.Flags().BoolP("force", "f", false, "Force removal even if worktree is dirty")

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		testCmd.SetArgs([]string{"--help"})
		testCmd.Execute()
	}
}

func BenchmarkRemoveCommandFlagParsing(b *testing.B) {
	testCmd := &cobra.Command{
		Use: "remove",
	}
	
	var testForceRemove bool
	testCmd.Flags().BoolVarP(&testForceRemove, "force", "f", false, "Force removal even if worktree is dirty")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testCmd.ParseFlags([]string{"-f"})
	}
}