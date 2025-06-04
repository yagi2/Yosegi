package cmd

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to a different worktree",
	Long:  "Interactively select and switch to a different git worktree.",
	Aliases: []string{"sw", "s", "cd"},
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := git.NewManager()
		if err != nil {
			return fmt.Errorf("failed to initialize git manager: %w", err)
		}

		worktrees, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list worktrees: %w", err)
		}

		if len(worktrees) == 0 {
			fmt.Println("No worktrees found")
			return nil
		}

		// Interactive mode
		model := ui.NewSelector(worktrees, "Switch Worktree", "switch", false)
		program := tea.NewProgram(model)
		
		finalModel, err := program.Run()
		if err != nil {
			return fmt.Errorf("failed to run interactive interface: %w", err)
		}

		result := finalModel.(ui.SelectorModel).GetResult()
		if result.Action == "quit" {
			return nil
		}

		if result.Action == "switch" {
			// Output the path for shell integration
			absPath, err := filepath.Abs(result.Worktree.Path)
			if err != nil {
				absPath = result.Worktree.Path
			}
			
			fmt.Printf("CD:%s\n", absPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}