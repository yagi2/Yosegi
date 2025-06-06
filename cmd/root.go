package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/config"
	"github.com/yagi2/yosegi/internal/ui"
)

// Version information (set by main.go)
var (
	version = "dev"
	commit  = "unknown" 
	date    = "unknown"
	builtBy = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "yosegi",
	Short: "Interactive git worktree management tool",
	Long: `Yosegi is a CLI tool for managing git worktrees with an interactive interface.
It provides visual and intuitive commands to create, list, and manage git worktrees.`,
	Version: getVersionString(),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Use the same functionality as list command
		return listCmd.RunE(cmd, args)
	},
}

// SetVersionInfo sets the version information from build-time variables
func SetVersionInfo(v, c, d, b string) {
	version = v
	commit = c
	date = d
	builtBy = b
	rootCmd.Version = getVersionString()
}

// getVersionString returns a formatted version string
func getVersionString() string {
	if version == "dev" {
		return fmt.Sprintf("dev (commit: %s, built: %s, go: %s)", 
			commit, date, runtime.Version())
	}
	return fmt.Sprintf("%s (commit: %s, built: %s, by: %s, go: %s)", 
		version, commit, date, builtBy, runtime.Version())
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Load configuration and initialize theme
	cfg, err := config.Load()
	if err == nil {
		ui.InitializeTheme(cfg)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
