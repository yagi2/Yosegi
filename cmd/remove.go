package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/config"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var (
	forceRemove bool
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a git worktree",
	Long:    "Interactively select and remove a git worktree.",
	Aliases: []string{"rm", "delete", "del", "r"},
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

		// Filter out current worktree (can't remove current worktree)
		var removableWorktrees []git.Worktree
		for _, wt := range worktrees {
			if !wt.IsCurrent {
				removableWorktrees = append(removableWorktrees, wt)
			}
		}

		if len(removableWorktrees) == 0 {
			fmt.Println("No removable worktrees found (cannot remove current worktree)")
			return nil
		}

		// Interactive mode
		model := ui.NewSelector(removableWorktrees, "Remove Worktree", "remove", true)
		program := tea.NewProgram(model)

		finalModel, err := program.Run()
		if err != nil {
			return fmt.Errorf("failed to run interactive interface: %w", err)
		}

		result := finalModel.(ui.SelectorModel).GetResult()
		if result.Action == "quit" {
			return nil
		}

		if result.Action == "remove" || result.Action == "delete" || result.Action == "select" {
			// Confirm removal
			confirmModel := ui.NewConfirm(
				"Confirm Removal",
				fmt.Sprintf("Remove worktree at %s?", result.Worktree.Path),
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
			fmt.Printf("Removing worktree at '%s'...\n", result.Worktree.Path)
			err = manager.Remove(result.Worktree.Path, forceRemove)
			if err != nil {
				return fmt.Errorf("failed to remove worktree: %w", err)
			}

			fmt.Printf("✅ Successfully removed worktree at '%s'\n", result.Worktree.Path)

			// Check if we should also delete the branch
			cfg, err := config.Load()
			if err != nil {
				cfg = &config.Config{}
			}

			// Skip branch deletion for detached HEAD or bare repository
			if result.Worktree.Branch == "(detached)" || result.Worktree.Branch == "(bare)" {
				return nil
			}

			// Determine if we should delete the branch
			deleteBranch := cfg.Git.DeleteBranchOnWorktreeRemove

			// Check for unpushed commits
			hasUnpushed, unpushedCount, err := manager.HasUnpushedCommits(result.Worktree.Branch)
			if err == nil && hasUnpushed {
				// Show warning and ask for confirmation
				warningModel := ui.NewConfirm(
					"Branch Deletion Warning",
					fmt.Sprintf("Branch '%s' has %d unpushed commits. Delete branch anyway?", result.Worktree.Branch, unpushedCount),
				)
				program := tea.NewProgram(warningModel)

				finalWarningModel, err := program.Run()
				if err != nil {
					return fmt.Errorf("failed to run warning dialog: %w", err)
				}

				warningResult := finalWarningModel.(ui.ConfirmModel).GetResult()
				if warningResult.Cancelled || !warningResult.Confirmed {
					deleteBranch = false
				} else {
					deleteBranch = true
				}
			} else if !deleteBranch {
				// Ask if user wants to delete the branch (if not configured to auto-delete)
				confirmBranchModel := ui.NewConfirm(
					"Delete Branch",
					fmt.Sprintf("Also delete the local branch '%s'?", result.Worktree.Branch),
				)
				program := tea.NewProgram(confirmBranchModel)

				finalBranchModel, err := program.Run()
				if err != nil {
					return fmt.Errorf("failed to run branch deletion dialog: %w", err)
				}

				branchResult := finalBranchModel.(ui.ConfirmModel).GetResult()
				deleteBranch = !branchResult.Cancelled && branchResult.Confirmed
			}

			// Delete the branch if confirmed
			if deleteBranch {
				fmt.Printf("Deleting branch '%s'...\n", result.Worktree.Branch)
				err = manager.DeleteBranch(result.Worktree.Branch, forceRemove || hasUnpushed)
				if err != nil {
					// Don't fail the whole operation if branch deletion fails
					fmt.Printf("⚠️  Warning: Failed to delete branch: %v\n", err)
				} else {
					fmt.Printf("✅ Successfully deleted branch '%s'\n", result.Worktree.Branch)
				}
			}
		}

		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal even if worktree is dirty")
	rootCmd.AddCommand(removeCmd)
}
