package cmd

import (
	"fmt"
	"os"
	"runtime"

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

			// For --print mode, try to get terminal access for interaction
			var inputFile *os.File
			var outputFile *os.File = os.Stderr
			var interactive bool = false
			
			// Try to open controlling terminal
			if runtime.GOOS != "windows" {
				if ctty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
					inputFile = ctty
					outputFile = ctty
					interactive = true
					defer ctty.Close()
				}
			}
			
			// If we couldn't get terminal access, fallback to non-interactive
			if !interactive {
				// Non-interactive fallback: return first non-current worktree
				for _, wt := range worktrees {
					if !wt.IsCurrent {
						fmt.Println(wt.Path)
						return nil
					}
				}
				// If all worktrees are current, just return the first one
				if len(worktrees) > 0 {
					fmt.Println(worktrees[0].Path)
					return nil
				}
				return fmt.Errorf("no worktrees found")
			}
			
			// Use simple numbered selector for --print mode
			selectedWorktree, err := ui.SimpleSelectWorktree(worktrees, outputFile, inputFile)
			if err != nil {
				// If selection fails, fallback to first non-current
				for _, wt := range worktrees {
					if !wt.IsCurrent {
						fmt.Println(wt.Path)
						return nil
					}
				}
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

	// Check if we should also delete the branch
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}

	// Skip branch deletion for detached HEAD or bare repository
	if selectedWorktree.Branch == "(detached)" || selectedWorktree.Branch == "(bare)" {
		return nil
	}

	// Determine if we should delete the branch
	deleteBranch := cfg.Git.DeleteBranchOnWorktreeRemove

	// Check for unpushed commits
	hasUnpushed, unpushedCount, err := manager.HasUnpushedCommits(selectedWorktree.Branch)
	if err == nil && hasUnpushed {
		// Show warning and ask for confirmation
		warningModel := ui.NewConfirm(
			"Branch Deletion Warning",
			fmt.Sprintf("Branch '%s' has %d unpushed commits. Delete branch anyway?", selectedWorktree.Branch, unpushedCount),
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
			fmt.Sprintf("Also delete the local branch '%s'?", selectedWorktree.Branch),
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
		fmt.Printf("Deleting branch '%s'...\n", selectedWorktree.Branch)
		err = manager.DeleteBranch(selectedWorktree.Branch, hasUnpushed)
		if err != nil {
			// Don't fail the whole operation if branch deletion fails
			fmt.Printf("⚠️  Warning: Failed to delete branch: %v\n", err)
		} else {
			fmt.Printf("✅ Successfully deleted branch '%s'\n", selectedWorktree.Branch)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolVarP(&printMode, "print", "p", false, "Show interactive selector on stderr and print selected path to stdout (for use in scripts)")
}
