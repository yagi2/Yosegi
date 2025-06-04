package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/cli-vibe-go/internal/git"
	"github.com/yagi2/cli-vibe-go/internal/ui"
)

var (
	forceRemove bool
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a git worktree",
	Long:  "Interactively select and remove a git worktree.",
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

		if result.Action == "remove" || result.Action == "delete" {
			// Confirm removal
			confirmModel := ui.NewInput(
				"Confirm Removal",
				[]string{fmt.Sprintf("Type 'yes' to remove worktree at %s", result.Worktree.Path)},
				[]string{""},
			)
			program := tea.NewProgram(confirmModel)
			
			finalConfirmModel, err := program.Run()
			if err != nil {
				return fmt.Errorf("failed to run confirmation interface: %w", err)
			}

			confirmResult := finalConfirmModel.(ui.InputModel).GetResult()
			if !confirmResult.Submitted || len(confirmResult.Values) == 0 || confirmResult.Values[0] != "yes" {
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
		}

		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal even if worktree is dirty")
	rootCmd.AddCommand(removeCmd)
}