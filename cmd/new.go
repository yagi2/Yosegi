package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagi2/yosegi/internal/config"
	"github.com/yagi2/yosegi/internal/git"
	"github.com/yagi2/yosegi/internal/ui"
)

var (
	createBranch    bool
	createBranchSet bool // Track if the flag was explicitly set
	worktreePath    string
)

var newCmd = &cobra.Command{
	Use:     "new [branch]",
	Short:   "Create a new git worktree",
	Long:    "Create a new git worktree interactively or with specified parameters.",
	Aliases: []string{"add", "create", "n"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{
				Git: config.GitConfig{
					AutoCreateBranch: false,
				},
			}
		}

		manager, err := git.NewManager()
		if err != nil {
			return fmt.Errorf("failed to initialize git manager: %w", err)
		}

		var branch string
		var path string

		// If branch is provided as argument, use it
		if len(args) > 0 {
			branch = args[0]
		}

		// Interactive mode for missing parameters
		if branch == "" || worktreePath == "" {
			var model ui.InputModel

			// If both branch and path are missing, use the enhanced worktree input with auto-generation
			if branch == "" && worktreePath == "" {
				model = ui.NewWorktreeInput("Create New Worktree", cfg.DefaultWorktreePath)
			} else {
				// Fallback to regular input for partial inputs
				prompts := []string{}
				defaults := []string{}

				if branch == "" {
					prompts = append(prompts, "Branch name (e.g., feature/new-feature)")
					defaults = append(defaults, "")
				}

				if worktreePath == "" {
					prompts = append(prompts, "Worktree directory path (e.g., ../feature-branch)")
					if branch != "" {
						defaults = append(defaults, filepath.Join(cfg.DefaultWorktreePath, branch))
					} else {
						defaults = append(defaults, "")
					}
				}

				model = ui.NewInput("Create New Worktree", prompts, defaults)
			}

			program := tea.NewProgram(model)

			finalModel, err := program.Run()
			if err != nil {
				return fmt.Errorf("failed to run interactive interface: %w", err)
			}

			result := finalModel.(ui.InputModel).GetResult()
			if !result.Submitted {
				fmt.Println("Cancelled")
				return nil
			}

			values := result.Values
			idx := 0

			if branch == "" {
				branch = strings.TrimSpace(values[idx])
				idx++
			}

			if worktreePath == "" {
				path = strings.TrimSpace(values[idx])
			}
		} else {
			path = worktreePath
		}

		// Validate inputs
		if branch == "" {
			return fmt.Errorf("branch name is required")
		}
		if path == "" {
			return fmt.Errorf("worktree path is required")
		}

		// Create the worktree
		// Use config auto_create_branch if createBranch flag is not explicitly set
		if !createBranchSet && cfg.Git.AutoCreateBranch {
			createBranch = true
		}

		fmt.Printf("Creating worktree '%s' at '%s'...\n", branch, path)
		err = manager.Add(path, branch, createBranch)
		if err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}

		fmt.Printf("✅ Successfully created worktree '%s' at '%s'\n", branch, path)

		return nil
	},
}

func init() {
	flags := newCmd.Flags()
	flags.BoolVarP(&createBranch, "create-branch", "b", false, "Create a new branch")
	flags.StringVarP(&worktreePath, "path", "p", "", "Path for the new worktree")

	// Mark that create-branch flag was explicitly set
	newCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		createBranchSet = flags.Changed("create-branch")
		return nil
	}

	rootCmd.AddCommand(newCmd)
}
