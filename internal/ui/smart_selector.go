package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagi2/yosegi/internal/git"
)

// SmartSelectWorktree selects the best UI method based on TTY capabilities
func SmartSelectWorktree(worktrees []git.Worktree) (*git.Worktree, error) {
	if len(worktrees) == 0 {
		return nil, fmt.Errorf("no worktrees found")
	}

	// Detect TTY capability
	capability := DetectTTYCapability()

	switch capability {
	case FullTTYControl:
		return bubbleTeaSelector(worktrees)
	case BasicTTYControl:
		return keyboardSelector(worktrees)
	case NoTTYControl:
		return fallbackSelector(worktrees)
	default:
		return fallbackSelector(worktrees)
	}
}

// bubbleTeaSelector uses the full Bubble Tea TUI
func bubbleTeaSelector(worktrees []git.Worktree) (*git.Worktree, error) {
	model := NewSelector(worktrees, "Git Worktrees", "select", false)
	program := tea.NewProgram(model)

	finalModel, err := program.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run interactive interface: %w", err)
	}

	result := finalModel.(SelectorModel).GetResult()
	if result.Action == "select" {
		return &result.Worktree, nil
	}
	return nil, fmt.Errorf("selection cancelled")
}

// keyboardSelector uses the custom keyboard interface
func keyboardSelector(worktrees []git.Worktree) (*git.Worktree, error) {
	input, output, cleanup, err := GetTTYFiles(BasicTTYControl)
	if err != nil {
		// Fallback to simple selector if TTY setup fails
		return SimpleSelectWorktree(worktrees, os.Stderr, os.Stdin)
	}
	defer cleanup()

	selector := NewKeyboardSelector(worktrees, input, output)
	return selector.Run()
}

// fallbackSelector uses the simple numbered interface
func fallbackSelector(worktrees []git.Worktree) (*git.Worktree, error) {
	// Return first non-current worktree automatically
	for _, wt := range worktrees {
		if !wt.IsCurrent {
			return &wt, nil
		}
	}

	// If all worktrees are current, return the first one
	if len(worktrees) > 0 {
		return &worktrees[0], nil
	}

	return nil, fmt.Errorf("no worktrees found")
}
