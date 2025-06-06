package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yagi2/yosegi/internal/git"
)

// KeyboardSelector provides arrow key navigation without full TUI
type KeyboardSelector struct {
	worktrees []git.Worktree
	cursor    int
	input     FileInterface
	output    FileInterface
}

// NewKeyboardSelector creates a new keyboard-based selector
func NewKeyboardSelector(worktrees []git.Worktree, input, output *os.File) *KeyboardSelector {
	return newKeyboardSelectorWithFiles(worktrees, input, output)
}

// newKeyboardSelectorWithFiles is the testable version that accepts interfaces
func newKeyboardSelectorWithFiles(worktrees []git.Worktree, input, output FileInterface) *KeyboardSelector {
	return &KeyboardSelector{
		worktrees: worktrees,
		cursor:    0,
		input:     input,
		output:    output,
	}
}

// Run starts the keyboard selector
func (k *KeyboardSelector) Run() (*git.Worktree, error) {
	if len(k.worktrees) == 0 {
		return nil, fmt.Errorf("no worktrees found")
	}

	// Set terminal to raw mode for key capture
	if err := k.setRawMode(); err != nil {
		return nil, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer k.restoreMode()

	// Initial render
	k.render()

	// Main input loop
	for {
		key, err := k.readKey()
		if err != nil {
			return nil, fmt.Errorf("failed to read key: %w", err)
		}

		switch key {
		case "up", "k":
			if k.cursor > 0 {
				k.cursor--
				k.render()
			}
		case "down", "j":
			if k.cursor < len(k.worktrees)-1 {
				k.cursor++
				k.render()
			}
		case "enter":
			k.clearScreen()
			return &k.worktrees[k.cursor], nil
		case "q", "ctrl+c":
			k.clearScreen()
			return nil, fmt.Errorf("selection cancelled")
		}
	}
}

// setRawMode puts the terminal in raw mode to capture individual keystrokes
func (k *KeyboardSelector) setRawMode() error {
	// Use stty to set raw mode
	cmd := exec.Command("stty", "-echo", "-icanon", "min", "1", "time", "0")
	
	// Convert to *os.File if possible, otherwise use default stdin/stdout
	if osInput, ok := k.input.(*os.File); ok {
		cmd.Stdin = osInput
	} else {
		cmd.Stdin = os.Stdin
	}
	
	if osOutput, ok := k.output.(*os.File); ok {
		cmd.Stdout = osOutput
		cmd.Stderr = osOutput
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	return cmd.Run()
}

// restoreMode restores the terminal to its original state
func (k *KeyboardSelector) restoreMode() {
	// Restore terminal settings
	cmd := exec.Command("stty", "echo", "icanon")
	
	// Convert to *os.File if possible, otherwise use default stdin/stdout
	if osInput, ok := k.input.(*os.File); ok {
		cmd.Stdin = osInput
	} else {
		cmd.Stdin = os.Stdin
	}
	
	if osOutput, ok := k.output.(*os.File); ok {
		cmd.Stdout = osOutput
		cmd.Stderr = osOutput
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	cmd.Run() // Ignore errors during cleanup
}

// readKey reads a single key from input and returns a normalized key name
func (k *KeyboardSelector) readKey() (string, error) {
	buf := make([]byte, 4)
	n, err := k.input.Read(buf)
	if err != nil {
		return "", err
	}

	// Parse key sequences
	switch {
	case n == 1:
		switch buf[0] {
		case 13, 10: // Enter
			return "enter", nil
		case 3: // Ctrl+C
			return "ctrl+c", nil
		case 106: // j
			return "j", nil
		case 107: // k
			return "k", nil
		case 113: // q
			return "q", nil
		}
	case n == 3 && buf[0] == 27 && buf[1] == 91: // ESC [ sequence
		switch buf[2] {
		case 65: // Up arrow
			return "up", nil
		case 66: // Down arrow
			return "down", nil
		}
	}

	return "", nil // Unknown key, ignore
}

// render draws the current state of the selector
func (k *KeyboardSelector) render() {
	// Clear screen and move cursor to top
	fmt.Fprint(k.output, "\033[2J\033[H")

	// Title
	fmt.Fprintf(k.output, "\033[1m🌲 Git Worktrees\033[0m\n")
	fmt.Fprintf(k.output, "%s\n", strings.Repeat("-", 60))

	// Worktree list
	for i, wt := range k.worktrees {
		status := "  "
		if wt.IsCurrent {
			status = "* "
		}

		// Highlight current selection
		if i == k.cursor {
			fmt.Fprintf(k.output, "\033[7m") // Reverse video
		}

		fmt.Fprintf(k.output, "%s%s (%s)\033[0m\n", status, wt.Path, wt.Branch)
	}

	// Help text
	fmt.Fprintf(k.output, "%s\n", strings.Repeat("-", 60))
	fmt.Fprintf(k.output, "\033[2m↑/k up  ↓/j down  Enter select  q quit\033[0m\n")
}

// clearScreen clears the screen
func (k *KeyboardSelector) clearScreen() {
	fmt.Fprint(k.output, "\033[2J\033[H")
}