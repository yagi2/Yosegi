package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/yagi2/yosegi/internal/git"
)

// SimpleSelectWorktree displays a numbered list of worktrees and prompts for selection
// This is used when full TUI is not available (e.g., in command substitution)
func SimpleSelectWorktree(worktrees []git.Worktree, output *os.File, input *os.File) (*git.Worktree, error) {
	if len(worktrees) == 0 {
		return nil, fmt.Errorf("no worktrees found")
	}

	// Display worktree list
	fmt.Fprintln(output, "\n🌲 Git Worktrees:")
	fmt.Fprintln(output, strings.Repeat("-", 60))
	
	for i, wt := range worktrees {
		status := "  "
		if wt.IsCurrent {
			status = "* "
		}
		fmt.Fprintf(output, "%s%d) %s (%s)\n", status, i+1, wt.Path, wt.Branch)
	}
	
	fmt.Fprintln(output, strings.Repeat("-", 60))
	fmt.Fprint(output, "Select worktree (1-", len(worktrees), ") or 'q' to quit: ")

	// Read user input
	reader := bufio.NewReader(input)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		
		input = strings.TrimSpace(input)
		
		// Check for quit
		if input == "q" || input == "Q" {
			return nil, fmt.Errorf("selection cancelled")
		}
		
		// Try to parse as number
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Fprintf(output, "Invalid input. Please enter a number (1-%d) or 'q' to quit: ", len(worktrees))
			continue
		}
		
		// Check range
		if num < 1 || num > len(worktrees) {
			fmt.Fprintf(output, "Invalid selection. Please enter a number (1-%d) or 'q' to quit: ", len(worktrees))
			continue
		}
		
		return &worktrees[num-1], nil
	}
}