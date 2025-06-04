package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yagi2/cli-vibe-go/internal/config"
	"github.com/yagi2/cli-vibe-go/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "yosegi",
	Short: "Interactive git worktree management tool",
	Long: `Yosegi is a CLI tool for managing git worktrees with an interactive interface.
It provides visual and intuitive commands to create, switch, and manage git worktrees.

Shell Integration:
To enable directory switching functionality, add the following to your shell config:

Bash (~/.bashrc):
  source /path/to/yosegi/scripts/shell_integration.bash

Zsh (~/.zshrc):
  source /path/to/yosegi/scripts/shell_integration.zsh

Fish (~/.config/fish/config.fish):
  source /path/to/yosegi/scripts/shell_integration.fish`,
	Version: "0.1.0",
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