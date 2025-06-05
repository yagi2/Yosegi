package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
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

		// Interactive mode
		model := ui.NewSelector(worktrees, "Git Worktrees", "print path", true)
		program := tea.NewProgram(model)

		finalModel, err := program.Run()
		if err != nil {
			return fmt.Errorf("failed to run interactive interface: %w", err)
		}

		result := finalModel.(ui.SelectorModel).GetResult()
		switch result.Action {
		case "select":
			// Print the selected worktree path to stdout
			fmt.Println(result.Worktree.Path)

		case "create":
			// Call existing new command
			return newCmd.RunE(cmd, []string{})

		case "delete":
			// Call existing remove command with pre-selected worktree
			return runRemoveWithSelectedWorktree(result.Worktree)

		case "quit":
			// Do nothing, just exit
		}

		return nil
	},
}

// runRemoveWithSelectedWorktree runs remove command with a pre-selected worktree
func runRemoveWithSelectedWorktree(selectedWorktree git.Worktree) error {
	// Filter out current worktree (can't remove current worktree)
	if selectedWorktree.IsCurrent {
		return fmt.Errorf("cannot remove current worktree")
	}

	manager, err := git.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize git manager: %w", err)
	}

	// Confirm removal
	confirmModel := ui.NewConfirm(
		"Confirm Removal",
		fmt.Sprintf("Remove worktree at %s?", selectedWorktree.Path),
	)
	program := tea.NewProgram(confirmModel)

	finalConfirmModel, err := program.Run()
	if err != nil {
		return fmt.Errorf("failed to run confirmation interface: %w", err)
	}

	confirmResult := finalConfirmModel.(ui.ConfirmModel).GetResult()
	if confirmResult.Cancelled || !confirmResult.Confirmed {
		fmt.Println("Removal cancelled")
		return nil
	}

	// Remove the worktree
	fmt.Printf("Removing worktree at '%s'...\n", selectedWorktree.Path)
	err = manager.Remove(selectedWorktree.Path, false)
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	fmt.Printf("✅ Successfully removed worktree at '%s'\n", selectedWorktree.Path)
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
