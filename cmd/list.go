package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var (
	plainOutput bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all git worktrees",
	Long:    "Display an interactive list of all git worktrees in the repository.",
	Aliases: []string{"ls", "l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		manager, err := git.NewManager()
		if err != nil {
			return fmt.Errorf("failed to initialize git manager: %w", err)
		}

		worktrees, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list worktrees: %w", err)
		}

		// Check if plain output is explicitly requested
		if plainOutput {
			// Plain output mode
			if len(worktrees) == 0 {
				fmt.Println("No worktree found")
				return nil
			}

			fmt.Println("Git Worktrees:")
			for _, wt := range worktrees {
				current := ""
				if wt.IsCurrent {
					current = " (current)"
				}
				fmt.Printf("  %s [%s]%s\n", wt.Path, wt.Branch, current)
			}
			return nil
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

func init() {
	listCmd.Flags().BoolVarP(&plainOutput, "plain", "", false, "Output in plain text format")
	rootCmd.AddCommand(listCmd)
}
