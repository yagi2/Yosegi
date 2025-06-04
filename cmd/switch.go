package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var (
	switchPlainOutput bool
)

var switchCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch to a different worktree",
	Long:    "Interactively select and switch to a different git worktree.",
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

		// Check for direct argument (worktree path or branch name)
		if len(args) > 0 {
			target := args[0]
			var selectedWorktree *git.Worktree

			// Try to find by path or branch name
			for _, wt := range worktrees {
				if wt.Path == target || wt.Branch == target || filepath.Base(wt.Path) == target {
					selectedWorktree = &wt
					break
				}
			}

			if selectedWorktree == nil {
				return fmt.Errorf("worktree '%s' not found", target)
			}

			// Output the path for shell integration
			absPath, err := filepath.Abs(selectedWorktree.Path)
			if err != nil {
				absPath = selectedWorktree.Path
			}
			fmt.Printf("CD:%s\n", absPath)

			// Show help message unless suppressed by environment variable
			if os.Getenv("YOSEGI_SHELL_INTEGRATION") == "" {
				fmt.Fprintf(os.Stderr, "\n# To enable automatic directory switching, set up shell integration:\n")
				fmt.Fprintf(os.Stderr, "# For bash: source /path/to/yosegi/scripts/shell_integration.bash\n")
				fmt.Fprintf(os.Stderr, "# For zsh: source /path/to/yosegi/scripts/shell_integration.zsh\n")
				fmt.Fprintf(os.Stderr, "# For fish: source /path/to/yosegi/scripts/shell_integration.fish\n")
				fmt.Fprintf(os.Stderr, "# Or manually run: cd %s\n", absPath)
				fmt.Fprintf(os.Stderr, "# Set YOSEGI_SHELL_INTEGRATION=1 to suppress this message\n")
			}

			return nil
		}

		// Check if TTY is available or plain output is requested
		if switchPlainOutput || !isatty() {
			// Plain output mode - just list worktrees for manual selection
			fmt.Println("Available worktrees:")
			for i, wt := range worktrees {
				current := ""
				if wt.IsCurrent {
					current = " (current)"
				}
				fmt.Printf("  %d. %s [%s]%s\n", i+1, wt.Path, wt.Branch, current)
			}
			fmt.Println("\nUsage: yosegi switch <path|branch>")
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

			// Show help message unless suppressed by environment variable
			if os.Getenv("YOSEGI_SHELL_INTEGRATION") == "" {
				fmt.Fprintf(os.Stderr, "\n# To enable automatic directory switching, set up shell integration:\n")
				fmt.Fprintf(os.Stderr, "# For bash: source /path/to/yosegi/scripts/shell_integration.bash\n")
				fmt.Fprintf(os.Stderr, "# For zsh: source /path/to/yosegi/scripts/shell_integration.zsh\n")
				fmt.Fprintf(os.Stderr, "# For fish: source /path/to/yosegi/scripts/shell_integration.fish\n")
				fmt.Fprintf(os.Stderr, "# Or manually run: cd %s\n", absPath)
				fmt.Fprintf(os.Stderr, "# Set YOSEGI_SHELL_INTEGRATION=1 to suppress this message\n")
			}
		}

		return nil
	},
}

func init() {
	switchCmd.Flags().BoolVarP(&switchPlainOutput, "plain", "", false, "Output in plain text format")
	rootCmd.AddCommand(switchCmd)
}
