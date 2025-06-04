package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yagi2/yosegi/internal/git"
)

func TestNewSelector(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: true},
		{Path: "/path/to/feature", Branch: "feature", Commit: "def456", IsCurrent: false},
	}

	tests := []struct {
		name        string
		worktrees   []git.Worktree
		title       string
		action      string
		allowDelete bool
	}{
		{
			name:        "Basic selector",
			worktrees:   worktrees,
			title:       "Select Worktree",
			action:      "select",
			allowDelete: false,
		},
		{
			name:        "Selector with delete allowed",
			worktrees:   worktrees,
			title:       "Remove Worktree",
			action:      "remove",
			allowDelete: true,
		},
		{
			name:        "Empty worktree list",
			worktrees:   []git.Worktree{},
			title:       "Empty List",
			action:      "select",
			allowDelete: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewSelector(tt.worktrees, tt.title, tt.action, tt.allowDelete)

			if len(model.worktrees) != len(tt.worktrees) {
				t.Errorf("Expected %d worktrees, got %d", len(tt.worktrees), len(model.worktrees))
			}

			if model.title != tt.title {
				t.Errorf("Expected title '%s', got '%s'", tt.title, model.title)
			}

			if model.action != tt.action {
				t.Errorf("Expected action '%s', got '%s'", tt.action, model.action)
			}

			if model.allowDelete != tt.allowDelete {
				t.Errorf("Expected allowDelete %v, got %v", tt.allowDelete, model.allowDelete)
			}

			if model.cursor != 0 {
				t.Errorf("Expected cursor to be 0, got %d", model.cursor)
			}

			if model.selectedPath != "" {
				t.Errorf("Expected selectedPath to be empty, got '%s'", model.selectedPath)
			}

			if model.quitting {
				t.Errorf("Expected quitting to be false")
			}
		})
	}
}

func TestSelectorInit(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: true},
	}

	model := NewSelector(worktrees, "Test", "select", false)
	cmd := model.Init()

	if cmd != nil {
		t.Errorf("Expected Init() to return nil, got %v", cmd)
	}
}

func TestSelectorUpdate(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: true},
		{Path: "/path/to/feature1", Branch: "feature1", Commit: "def456", IsCurrent: false},
		{Path: "/path/to/feature2", Branch: "feature2", Commit: "ghi789", IsCurrent: false},
	}

	tests := []struct {
		name           string
		initialModel   SelectorModel
		msg            tea.Msg
		expectedCursor int
		expectedQuit   bool
		expectedAction string
		expectedPath   string
	}{
		{
			name:           "Move down",
			initialModel:   NewSelector(worktrees, "Test", "select", false),
			msg:            tea.KeyMsg{Type: tea.KeyDown},
			expectedCursor: 1,
			expectedQuit:   false,
		},
		{
			name: "Move down from bottom (should stay at bottom)",
			initialModel: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "select", false)
				m.cursor = 2 // Last item
				return m
			}(),
			msg:            tea.KeyMsg{Type: tea.KeyDown},
			expectedCursor: 2,
			expectedQuit:   false,
		},
		{
			name: "Move up",
			initialModel: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "select", false)
				m.cursor = 1
				return m
			}(),
			msg:            tea.KeyMsg{Type: tea.KeyUp},
			expectedCursor: 0,
			expectedQuit:   false,
		},
		{
			name:           "Move up from top (should stay at top)",
			initialModel:   NewSelector(worktrees, "Test", "select", false),
			msg:            tea.KeyMsg{Type: tea.KeyUp},
			expectedCursor: 0,
			expectedQuit:   false,
		},
		{
			name:           "Vim-style move down (j)",
			initialModel:   NewSelector(worktrees, "Test", "select", false),
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			expectedCursor: 1,
			expectedQuit:   false,
		},
		{
			name: "Vim-style move up (k)",
			initialModel: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "select", false)
				m.cursor = 1
				return m
			}(),
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			expectedCursor: 0,
			expectedQuit:   false,
		},
		{
			name:           "Select with Enter",
			initialModel:   NewSelector(worktrees, "Test", "select", false),
			msg:            tea.KeyMsg{Type: tea.KeyEnter},
			expectedCursor: 0,
			expectedQuit:   false, // quitting flag is not set, only tea.Quit command is returned
			expectedAction: "select",
			expectedPath:   "/path/to/main",
		},
		{
			name:         "Quit with q",
			initialModel: NewSelector(worktrees, "Test", "select", false),
			msg:          tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectedQuit: true,
		},
		{
			name:         "Quit with Ctrl+C",
			initialModel: NewSelector(worktrees, "Test", "select", false),
			msg:          tea.KeyMsg{Type: tea.KeyCtrlC},
			expectedQuit: true,
		},
		{
			name:           "Delete with d (when allowed)",
			initialModel:   NewSelector(worktrees, "Test", "remove", true),
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			expectedCursor: 0,
			expectedQuit:   false, // quitting flag is not set, only tea.Quit command is returned
			expectedAction: "delete",
			expectedPath:   "/path/to/main",
		},
		{
			name:           "Delete with d (when not allowed)",
			initialModel:   NewSelector(worktrees, "Test", "select", false),
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}},
			expectedCursor: 0,
			expectedQuit:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.initialModel
			newModel, cmd := model.Update(tt.msg)

			updatedModel := newModel.(SelectorModel)

			if updatedModel.cursor != tt.expectedCursor {
				t.Errorf("Expected cursor %d, got %d", tt.expectedCursor, updatedModel.cursor)
			}

			if updatedModel.quitting != tt.expectedQuit {
				t.Errorf("Expected quitting %v, got %v", tt.expectedQuit, updatedModel.quitting)
			}

			if tt.expectedPath != "" {
				if updatedModel.selectedPath != tt.expectedPath {
					t.Errorf("Expected selectedPath '%s', got '%s'", tt.expectedPath, updatedModel.selectedPath)
				}
			}

			if tt.expectedAction != "" {
				if updatedModel.action != tt.expectedAction {
					t.Errorf("Expected action '%s', got '%s'", tt.expectedAction, updatedModel.action)
				}
			}

			// Check if quit command was returned when expected
			if tt.expectedQuit && cmd == nil {
				t.Errorf("Expected quit command but got nil")
			}
		})
	}
}

func TestSelectorUpdateWithEmptyList(t *testing.T) {
	model := NewSelector([]git.Worktree{}, "Empty", "select", false)

	// Test Enter with empty list
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel := newModel.(SelectorModel)

	if updatedModel.selectedPath != "" {
		t.Errorf("Expected empty selectedPath with empty list, got '%s'", updatedModel.selectedPath)
	}

	if updatedModel.quitting {
		t.Errorf("Should not quit when pressing Enter on empty list")
	}

	if cmd != nil {
		t.Errorf("Should not return command when pressing Enter on empty list")
	}
}

func TestSelectorView(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: true},
		{Path: "/path/to/feature", Branch: "feature", Commit: "def456", IsCurrent: false},
	}

	tests := []struct {
		name                string
		model               SelectorModel
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name:  "Normal view",
			model: NewSelector(worktrees, "Git Worktrees", "select", false),
			expectedContains: []string{
				"🌲 Git Worktrees",
				"main",
				"feature",
				"/path/to/main",
				"/path/to/feature",
				"enter select",
				"q quit",
			},
			expectedNotContains: []string{
				"d delete",
			},
		},
		{
			name:  "View with delete allowed",
			model: NewSelector(worktrees, "Remove Worktree", "remove", true),
			expectedContains: []string{
				"🌲 Remove Worktree",
				"d delete",
				"enter remove",
			},
		},
		{
			name:  "Empty worktree list",
			model: NewSelector([]git.Worktree{}, "Empty", "select", false),
			expectedContains: []string{
				"🌲 Empty",
				"No worktrees found",
				"Press q to quit",
			},
			expectedNotContains: []string{
				"enter select",
			},
		},
		{
			name: "Quitting model",
			model: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "select", false)
				m.quitting = true
				return m
			}(),
			expectedContains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.model.View()

			for _, expected := range tt.expectedContains {
				if !strings.Contains(view, expected) {
					t.Errorf("Expected view to contain '%s', but it didn't.\nView: %s", expected, view)
				}
			}

			for _, notExpected := range tt.expectedNotContains {
				if strings.Contains(view, notExpected) {
					t.Errorf("Expected view to NOT contain '%s', but it did.\nView: %s", notExpected, view)
				}
			}
		})
	}
}

func TestSelectorGetResult(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: true},
		{Path: "/path/to/feature", Branch: "feature", Commit: "def456", IsCurrent: false},
	}

	tests := []struct {
		name           string
		setupModel     func() SelectorModel
		expectedAction string
		expectedPath   string
		expectedBranch string
	}{
		{
			name: "Normal selection",
			setupModel: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "select", false)
				m.selectedPath = "/path/to/main"
				return m
			},
			expectedAction: "select",
			expectedPath:   "/path/to/main",
			expectedBranch: "main",
		},
		{
			name: "Delete action",
			setupModel: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "delete", true)
				m.selectedPath = "/path/to/feature"
				m.action = "delete"
				return m
			},
			expectedAction: "delete",
			expectedPath:   "/path/to/feature",
			expectedBranch: "feature",
		},
		{
			name: "Quit without selection",
			setupModel: func() SelectorModel {
				m := NewSelector(worktrees, "Test", "select", false)
				m.quitting = true
				return m
			},
			expectedAction: "quit",
			expectedPath:   "",
			expectedBranch: "",
		},
		{
			name: "No selection made",
			setupModel: func() SelectorModel {
				return NewSelector(worktrees, "Test", "select", false)
			},
			expectedAction: "quit",
			expectedPath:   "",
			expectedBranch: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			result := model.GetResult()

			if result.Action != tt.expectedAction {
				t.Errorf("Expected action '%s', got '%s'", tt.expectedAction, result.Action)
			}

			if result.Worktree.Path != tt.expectedPath {
				t.Errorf("Expected path '%s', got '%s'", tt.expectedPath, result.Worktree.Path)
			}

			if result.Worktree.Branch != tt.expectedBranch {
				t.Errorf("Expected branch '%s', got '%s'", tt.expectedBranch, result.Worktree.Branch)
			}
		})
	}
}

func TestShortenPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Short path",
			input:    "/short/path",
			expected: "/short/path",
		},
		{
			name:     "Exact length path",
			input:    strings.Repeat("a", 50),
			expected: strings.Repeat("a", 50),
		},
		{
			name:     "Long path",
			input:    "/very/long/path/that/exceeds/the/maximum/length/limit/and/should/be/shortened",
			expected: "..." + "/very/long/path/that/exceeds/the/maximum/length/limit/and/should/be/shortened"[len("/very/long/path/that/exceeds/the/maximum/length/limit/and/should/be/shortened")-47:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shortenPath(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}

			// Ensure result is never longer than 50 characters
			if len(result) > 50 {
				t.Errorf("Result length %d exceeds maximum of 50", len(result))
			}
		})
	}
}

func TestKeyMappings(t *testing.T) {
	// Test that key mappings are properly defined
	if keys.Up.Keys() == nil {
		t.Error("Up key mapping should be defined")
	}

	if keys.Down.Keys() == nil {
		t.Error("Down key mapping should be defined")
	}

	if keys.Enter.Keys() == nil {
		t.Error("Enter key mapping should be defined")
	}

	if keys.Quit.Keys() == nil {
		t.Error("Quit key mapping should be defined")
	}

	if keys.Delete.Keys() == nil {
		t.Error("Delete key mapping should be defined")
	}

	// Test that help text is not empty
	if keys.Up.Help().Key == "" {
		t.Error("Up key should have help text")
	}

	if keys.Down.Help().Key == "" {
		t.Error("Down key should have help text")
	}
}

func TestSelectorWithCurrentWorktree(t *testing.T) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: false},
		{Path: "/path/to/feature", Branch: "feature", Commit: "def456", IsCurrent: true},
	}

	model := NewSelector(worktrees, "Test", "select", false)
	view := model.View()

	// Should show different styling for current worktree
	if !strings.Contains(view, "feature") {
		t.Error("View should contain the current worktree branch")
	}

	if !strings.Contains(view, "/path/to/feature") {
		t.Error("View should contain the current worktree path")
	}
}

// Benchmark tests
func BenchmarkSelectorView(b *testing.B) {
	worktrees := make([]git.Worktree, 100)
	for i := 0; i < 100; i++ {
		worktrees[i] = git.Worktree{
			Path:      "/path/to/worktree" + string(rune(i)),
			Branch:    "branch" + string(rune(i)),
			Commit:    "commit" + string(rune(i)),
			IsCurrent: i == 0,
		}
	}

	model := NewSelector(worktrees, "Benchmark", "select", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkSelectorUpdate(b *testing.B) {
	worktrees := []git.Worktree{
		{Path: "/path/to/main", Branch: "main", Commit: "abc123", IsCurrent: true},
		{Path: "/path/to/feature", Branch: "feature", Commit: "def456", IsCurrent: false},
	}

	model := NewSelector(worktrees, "Benchmark", "select", false)
	msg := tea.KeyMsg{Type: tea.KeyDown}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel, _ := model.Update(msg)
		model = newModel.(SelectorModel)
	}
}
