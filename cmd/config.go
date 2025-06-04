package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yagi2/cli-vibe-go/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage yosegi configuration",
	Long:  "Initialize, view, or modify yosegi configuration settings.",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration",
	Long:  "Create a default configuration file in ~/.config/yosegi/config.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		fmt.Println("✅ Default configuration file created successfully")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println("Current Configuration:")
		fmt.Printf("  Default Worktree Path: %s\n", cfg.DefaultWorktreePath)
		fmt.Printf("  Auto Create Branch: %t\n", cfg.Git.AutoCreateBranch)
		fmt.Printf("  Show Icons: %t\n", cfg.UI.ShowIcons)
		fmt.Printf("  Confirm Delete: %t\n", cfg.UI.ConfirmDelete)
		fmt.Printf("  Max Path Length: %d\n", cfg.UI.MaxPathLength)
		
		if len(cfg.Aliases) > 0 {
			fmt.Println("  Aliases:")
			for alias, command := range cfg.Aliases {
				fmt.Printf("    %s -> %s\n", alias, command)
			}
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}