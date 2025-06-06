package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/config"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var (
	printMode bool
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

		// Check if --print flag is used
		if printMode {
			if len(worktrees) == 0 {
				return fmt.Errorf("no worktrees found")
			}

			// Use smart selector that adapts to TTY capabilities
			selectedWorktree, err := ui.SmartSelectWorktree(worktrees)
			if err != nil {
				return err
			}

			// Print the selected worktree path to stdout
			fmt.Println(selectedWorktree.Path)
			return nil
		}

		// Non-TTY mode (e.g., command substitution without --print)
		if !isatty.IsTerminal(os.Stdout.Fd()) {
			if len(worktrees) == 0 {
				return fmt.Errorf("no worktrees found")
			}

			// Non-interactive fallback: return first non-current worktree
			for _, wt := range worktrees {
				if !wt.IsCurrent {
					fmt.Println(wt.Path)
					return nil
				}
			}

			// If all worktrees are current (unlikely), output the first one
			fmt.Fprintln(os.Stderr, "Warning: No non-current worktree found")
			return fmt.Errorf("no suitable worktree found")
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
	if selectedWorktree.IsCurrent {
		return fmt.Errorf("cannot remove current worktree")
	}

	manager, err := git.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize git manager: %w", err)
	}

	// Confirm worktree removal
	if !confirmWorktreeRemoval(selectedWorktree.Path) {
		fmt.Println("Removal cancelled")
		return nil
	}

	// Remove the worktree
	if err := removeWorktree(manager, selectedWorktree.Path); err != nil {
		return err
	}

	// Handle branch deletion if applicable
	return handleBranchDeletion(manager, selectedWorktree.Branch)
}

// confirmWorktreeRemoval shows confirmation dialog for worktree removal
func confirmWorktreeRemoval(path string) bool {
	confirmModel := ui.NewConfirm(
		"Confirm Removal",
		fmt.Sprintf("Remove worktree at %s?", path),
	)
	program := tea.NewProgram(confirmModel)

	finalModel, err := program.Run()
	if err != nil {
		return false
	}

	result := finalModel.(ui.ConfirmModel).GetResult()
	return !result.Cancelled && result.Confirmed
}

// removeWorktree removes the specified worktree
func removeWorktree(manager git.Manager, path string) error {
	fmt.Printf("Removing worktree at '%s'...\n", path)
	if err := manager.Remove(path, false); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}
	fmt.Printf("✅ Successfully removed worktree at '%s'\n", path)
	return nil
}

// handleBranchDeletion handles the branch deletion logic
func handleBranchDeletion(manager git.Manager, branch string) error {
	// Skip branch deletion for special branches
	if branch == "(detached)" || branch == "(bare)" {
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}

	deleteBranch, err := shouldDeleteBranch(manager, branch, cfg.Git.DeleteBranchOnWorktreeRemove)
	if err != nil {
		return err
	}

	if deleteBranch {
		return deleteBranchWithConfirmation(manager, branch)
	}
	return nil
}

// shouldDeleteBranch determines if the branch should be deleted
func shouldDeleteBranch(manager git.Manager, branch string, autoDelete bool) (bool, error) {
	hasUnpushed, unpushedCount, err := manager.HasUnpushedCommits(branch)
	if err == nil && hasUnpushed {
		return confirmUnpushedBranchDeletion(branch, unpushedCount), nil
	}

	if !autoDelete {
		return confirmBranchDeletion(branch), nil
	}

	return autoDelete, nil
}

// confirmUnpushedBranchDeletion shows warning for unpushed commits
func confirmUnpushedBranchDeletion(branch string, unpushedCount int) bool {
	warningModel := ui.NewConfirm(
		"Branch Deletion Warning",
		fmt.Sprintf("Branch '%s' has %d unpushed commits. Delete branch anyway?", branch, unpushedCount),
	)
	program := tea.NewProgram(warningModel)

	finalModel, err := program.Run()
	if err != nil {
		return false
	}

	result := finalModel.(ui.ConfirmModel).GetResult()
	return !result.Cancelled && result.Confirmed
}

// confirmBranchDeletion asks user if they want to delete the branch
func confirmBranchDeletion(branch string) bool {
	confirmModel := ui.NewConfirm(
		"Delete Branch",
		fmt.Sprintf("Also delete the local branch '%s'?", branch),
	)
	program := tea.NewProgram(confirmModel)

	finalModel, err := program.Run()
	if err != nil {
		return false
	}

	result := finalModel.(ui.ConfirmModel).GetResult()
	return !result.Cancelled && result.Confirmed
}

// deleteBranchWithConfirmation deletes the branch and shows result
func deleteBranchWithConfirmation(manager git.Manager, branch string) error {
	fmt.Printf("Deleting branch '%s'...\n", branch)
	hasUnpushed, _, _ := manager.HasUnpushedCommits(branch)

	if err := manager.DeleteBranch(branch, hasUnpushed); err != nil {
		fmt.Printf("⚠️  Warning: Failed to delete branch: %v\n", err)
		return nil // Don't fail the whole operation
	}

	fmt.Printf("✅ Successfully deleted branch '%s'\n", branch)
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolVarP(&printMode, "print", "p", false, "Show interactive selector on stderr and print selected path to stdout (for use in scripts)")
}
