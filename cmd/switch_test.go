package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestSwitchCommand(t *testing.T) {
	// Test command structure
	if switchCmd.Use != "switch" {
		t.Errorf("Expected switch command use to be 'switch', got '%s'", switchCmd.Use)
	}

	if switchCmd.Short == "" {
		t.Errorf("Switch command should have short description")
	}

	if switchCmd.Long == "" {
		t.Errorf("Switch command should have long description")
	}

	if switchCmd.RunE == nil {
		t.Errorf("Switch command should have RunE function")
	}

	// Test aliases
	expectedAliases := []string{"sw", "s", "cd"}
	if len(switchCmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(switchCmd.Aliases))
	}

	for i, expected := range expectedAliases {
		if i < len(switchCmd.Aliases) && switchCmd.Aliases[i] != expected {
			t.Errorf("Expected alias %d to be '%s', got '%s'", i, expected, switchCmd.Aliases[i])
		}
	}
}

func TestSwitchCommandStructure(t *testing.T) {
	// Test that switch command has proper configuration
	if switchCmd.Use != "switch" {
		t.Errorf("Expected switchCmd use to be 'switch', got '%s'", switchCmd.Use)
	}

	if switchCmd.Short == "" {
		t.Errorf("switchCmd should have short description")
	}

	if switchCmd.Long == "" {
		t.Errorf("switchCmd should have long description")
	}

	if switchCmd.RunE == nil {
		t.Errorf("switchCmd should have RunE function")
	}

	// Test command content
	expectedShort := "Switch to a different worktree"
	if switchCmd.Short != expectedShort {
		t.Errorf("Expected short description '%s', got '%s'", expectedShort, switchCmd.Short)
	}

	expectedLong := "Interactively select and switch to a different git worktree."
	if switchCmd.Long != expectedLong {
		t.Errorf("Expected long description '%s', got '%s'", expectedLong, switchCmd.Long)
	}
}

func TestSwitchCommandAliases(t *testing.T) {
	// Test that all expected aliases are present
	expectedAliases := map[string]bool{
		"sw": false,
		"s":  false,
		"cd": false,
	}

	for _, alias := range switchCmd.Aliases {
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

func TestSwitchCommandRegistered(t *testing.T) {
	// Check that switch command is registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "switch" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Switch command should be registered with root command")
	}
}

func TestSwitchCommandHelp(t *testing.T) {
	// Test switch command help
	testCmd := &cobra.Command{
		Use:     switchCmd.Use,
		Short:   switchCmd.Short,
		Long:    switchCmd.Long,
		Aliases: switchCmd.Aliases,
	}

	var buf bytes.Buffer
	testCmd.SetOut(&buf)
	testCmd.SetErr(&buf)
	testCmd.SetArgs([]string{"--help"})

	err := testCmd.Execute()
	if err != nil {
		t.Errorf("Switch help should not error: %v", err)
	}

	output := buf.String()

	// Basic check that help was generated
	if len(output) == 0 {
		t.Error("Help output should not be empty")
	}
}

func TestSwitchCommandByAlias(t *testing.T) {
	// Test that aliases work correctly by checking they're properly set
	aliases := switchCmd.Aliases

	expectedAliases := []string{"sw", "s", "cd"}
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
			t.Errorf("Expected alias '%s' not found in switch command", expected)
		}
	}
}

func TestSwitchCommandFlags(t *testing.T) {
	// Test that switch command doesn't have unexpected flags
	// Switch command should be simple with no additional flags
	flags := switchCmd.Flags()

	// Should have inherited flags from cobra (help, etc.) and our plain flag
	if flags.HasFlags() {
		// Check if there are any unexpected flags
		expectedFlags := map[string]bool{
			"help":        true,
			"plain":       true,
			"interactive": true,
		}
		flags.VisitAll(func(flag *pflag.Flag) {
			if !expectedFlags[flag.Name] {
				t.Errorf("Unexpected flag '%s' found in switch command", flag.Name)
			}
		})
	}
}

func TestSwitchCommandValidation(t *testing.T) {
	// Test command validation
	if switchCmd.Args != nil {
		// Test that args function works (if any)
		err := switchCmd.Args(switchCmd, []string{})
		if err != nil {
			t.Errorf("Args validation failed: %v", err)
		}
	}

	// Test that command is properly configured
	if switchCmd.Use == "" {
		t.Error("Switch command Use should not be empty")
	}

	if switchCmd.Short == "" {
		t.Error("Switch command Short description should not be empty")
	}

	// Test that RunE is set (required for execution)
	if switchCmd.RunE == nil {
		t.Error("Switch command should have RunE function set")
	}
}

func TestSwitchCommandInheritance(t *testing.T) {
	// Test that switch command inherits from root command correctly
	parent := switchCmd.Parent()
	if parent == nil {
		t.Error("Switch command should have a parent command")
		return
	}

	if parent.Name() != "yosegi" {
		t.Errorf("Expected parent command to be 'yosegi', got '%s'", parent.Name())
	}
}

func TestSwitchCommandPath(t *testing.T) {
	// Test command path construction
	expectedPath := "yosegi switch"
	actualPath := switchCmd.CommandPath()

	if actualPath != expectedPath {
		t.Errorf("Expected command path '%s', got '%s'", expectedPath, actualPath)
	}
}

func TestSwitchCommandUsage(t *testing.T) {
	// Test usage string generation
	usage := switchCmd.UsageString()

	// Basic check that usage string is generated
	if len(usage) == 0 {
		t.Error("Usage string should not be empty")
	}

	// Should contain the command name
	if !strings.Contains(usage, "switch") {
		t.Error("Usage should contain command name 'switch'")
	}
}

func TestSwitchCommandExample(t *testing.T) {
	// Test that command can be created and configured properly
	testCmd := &cobra.Command{
		Use:     switchCmd.Use,
		Short:   switchCmd.Short,
		Long:    switchCmd.Long,
		Aliases: switchCmd.Aliases,
	}

	// Test basic properties
	if testCmd.Name() != "switch" {
		t.Errorf("Expected command name 'switch', got '%s'", testCmd.Name())
	}

	if !testCmd.HasAlias("sw") {
		t.Error("Command should have alias 'sw'")
	}

	if !testCmd.HasAlias("s") {
		t.Error("Command should have alias 's'")
	}

	if !testCmd.HasAlias("cd") {
		t.Error("Command should have alias 'cd'")
	}
}

func TestSwitchCommandBehavior(t *testing.T) {
	// Test command configuration behavior
	if switchCmd.DisableFlagsInUseLine {
		t.Error("Switch command should not disable flags in usage line")
	}

	if switchCmd.Hidden {
		t.Error("Switch command should not be hidden")
	}

	if switchCmd.Deprecated != "" {
		t.Error("Switch command should not be deprecated")
	}
}

func TestSwitchCommandCompletion(t *testing.T) {
	// Test command completion setup
	if switchCmd.ValidArgsFunction != nil {
		// If completion is set up, test it
		completions, directive := switchCmd.ValidArgsFunction(switchCmd, []string{}, "")

		// Switch command shouldn't complete arguments since it's interactive
		if len(completions) != 0 {
			t.Errorf("Switch command should not provide argument completions, got %d", len(completions))
		}

		// Directive should indicate no file completion
		if directive != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("Expected ShellCompDirectiveNoFileComp, got %d", directive)
		}
	}
}

func TestSwitchCommandGroups(t *testing.T) {
	// Test command grouping (if any)
	if switchCmd.GroupID != "" {
		// If grouped, verify it's in the correct group
		t.Logf("Switch command is in group: %s", switchCmd.GroupID)
	}
}

func TestSwitchCommandAnnotations(t *testing.T) {
	// Test command annotations (if any)
	if len(switchCmd.Annotations) > 0 {
		t.Logf("Switch command has %d annotations", len(switchCmd.Annotations))
		for key, value := range switchCmd.Annotations {
			t.Logf("Annotation %s: %s", key, value)
		}
	}
}

// Test integration with parent command
func TestSwitchCommandIntegration(t *testing.T) {
	// Verify the command is properly integrated
	found := false
	var foundCmd *cobra.Command

	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "switch" {
			found = true
			foundCmd = cmd
			break
		}
	}

	if !found {
		t.Fatal("Switch command not found in root command")
	}

	// Test that it's the same command
	if foundCmd != switchCmd {
		t.Error("Found command is not the same as switchCmd")
	}

	// Test aliases are preserved
	if len(foundCmd.Aliases) != 3 {
		t.Errorf("Expected 3 aliases, got %d", len(foundCmd.Aliases))
	}
}

func TestSwitchCommandArgsAcceptance(t *testing.T) {
	// Test that command accepts the expected arguments
	// Switch command should accept no arguments (interactive mode)

	// Test with no args (should be valid)
	if switchCmd.Args != nil {
		err := switchCmd.Args(switchCmd, []string{})
		if err != nil {
			t.Errorf("Switch command should accept no arguments, got error: %v", err)
		}
	}
}

func TestSwitchCommandNoFlags(t *testing.T) {
	// Test that switch command has no custom flags (only inherited ones)
	localFlags := switchCmd.LocalFlags()

	// Should have our plain flag as local flag
	expectedLocalFlags := map[string]bool{
		"plain":       true,
		"interactive": true,
	}

	if localFlags.HasFlags() {
		localFlags.VisitAll(func(flag *pflag.Flag) {
			if !expectedLocalFlags[flag.Name] {
				t.Errorf("Unexpected local flag '%s' found in switch command", flag.Name)
			}
		})
	}
}

func TestSwitchCommandSimplicity(t *testing.T) {
	// Test that switch command is properly configured as a simple command

	// Should not have pre/post run functions
	if switchCmd.PreRun != nil {
		t.Error("Switch command should not have PreRun function")
	}

	if switchCmd.PostRun != nil {
		t.Error("Switch command should not have PostRun function")
	}

	if switchCmd.PreRunE != nil {
		t.Error("Switch command should not have PreRunE function")
	}

	if switchCmd.PostRunE != nil {
		t.Error("Switch command should not have PostRunE function")
	}

	// Should have RunE (main execution)
	if switchCmd.RunE == nil {
		t.Error("Switch command should have RunE function")
	}

	// Should not have Run (we use RunE for error handling)
	if switchCmd.Run != nil {
		t.Error("Switch command should use RunE instead of Run")
	}
}

func TestSwitchCommandParentRelationship(t *testing.T) {
	// Test the parent-child relationship
	parent := switchCmd.Parent()
	if parent == nil {
		t.Fatal("Switch command should have a parent")
	}

	// Parent should be root command
	if parent.Name() != "yosegi" {
		t.Errorf("Expected parent name 'yosegi', got '%s'", parent.Name())
	}

	// Check that parent contains this command
	found := false
	for _, child := range parent.Commands() {
		if child == switchCmd {
			found = true
			break
		}
	}

	if !found {
		t.Error("Parent command should contain switch command in its children")
	}
}

func TestSwitchCommandAliasUniqueness(t *testing.T) {
	// Test that aliases don't conflict with other commands
	allCommands := rootCmd.Commands()
	usedNames := make(map[string]string)

	// Collect all command names and aliases
	for _, cmd := range allCommands {
		usedNames[cmd.Name()] = cmd.Name()
		for _, alias := range cmd.Aliases {
			if existing, exists := usedNames[alias]; exists && existing != cmd.Name() {
				t.Errorf("Alias '%s' conflicts between commands '%s' and '%s'", alias, existing, cmd.Name())
			}
			usedNames[alias] = cmd.Name()
		}
	}

	// Specifically check switch command aliases
	for _, alias := range switchCmd.Aliases {
		if cmdName, exists := usedNames[alias]; !exists || cmdName != "switch" {
			if exists && cmdName != "switch" {
				t.Errorf("Switch alias '%s' conflicts with command '%s'", alias, cmdName)
			}
		}
	}
}

func TestSwitchCommandDocumentation(t *testing.T) {
	// Test that documentation strings are properly set
	if switchCmd.Short == "" {
		t.Error("Switch command should have short description")
	}

	if switchCmd.Long == "" {
		t.Error("Switch command should have long description")
	}

	// Test content quality (basic checks)
	if !strings.Contains(strings.ToLower(switchCmd.Short), "switch") {
		t.Error("Short description should mention 'switch'")
	}

	if !strings.Contains(strings.ToLower(switchCmd.Long), "worktree") {
		t.Error("Long description should mention 'worktree'")
	}
}

func TestSwitchCommandConsistency(t *testing.T) {
	// Test consistency with other commands in the suite

	// Should follow same naming pattern
	if !strings.HasSuffix(switchCmd.Use, "switch") && switchCmd.Use != "switch" {
		t.Error("Switch command Use should be 'switch'")
	}

	// Should have RunE like other commands
	if switchCmd.RunE == nil {
		t.Error("Switch command should have RunE function like other commands")
	}

	// Should be registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd == switchCmd {
			found = true
			break
		}
	}

	if !found {
		t.Error("Switch command should be registered with root command")
	}
}

// Benchmark tests
func BenchmarkSwitchCommandCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := &cobra.Command{
			Use:     "switch",
			Short:   "Switch to a different worktree",
			Long:    "Interactively select and switch to a different git worktree.",
			Aliases: []string{"sw", "s", "cd"},
		}
		_ = cmd
	}
}

func BenchmarkSwitchCommandHelp(b *testing.B) {
	testCmd := &cobra.Command{
		Use:     switchCmd.Use,
		Short:   switchCmd.Short,
		Long:    switchCmd.Long,
		Aliases: switchCmd.Aliases,
	}

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

func BenchmarkSwitchCommandPathResolution(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = switchCmd.CommandPath()
	}
}

func BenchmarkSwitchCommandAliasCheck(b *testing.B) {
	testCmd := &cobra.Command{
		Use:     switchCmd.Use,
		Aliases: switchCmd.Aliases,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = testCmd.HasAlias("sw")
		_ = testCmd.HasAlias("s")
		_ = testCmd.HasAlias("cd")
	}
}
