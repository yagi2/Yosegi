package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/config"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "yosegi",
	Short: "Interactive git worktree management tool",
	Long: `Yosegi is a CLI tool for managing git worktrees with an interactive interface.
It provides visual and intuitive commands to create, list, and manage git worktrees.`,
	Version: "0.1.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := git.NewManager()
		if err != nil {
			return fmt.Errorf("failed to initialize git manager: %w", err)
		}

		worktrees, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list worktrees: %w", err)
		}

		// Interactive mode
		model := ui.NewSelector(worktrees, "Git Worktrees", "view", false)
		program := tea.NewProgram(model)

		finalModel, err := program.Run()
		if err != nil {
			return fmt.Errorf("failed to run interactive interface: %w", err)
		}

		result := finalModel.(ui.SelectorModel).GetResult()
		if result.Action == "select" {
			fmt.Printf("Selected worktree: %s (%s)\n", result.Worktree.Path, result.Worktree.Branch)
		}

		return nil
	},
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
