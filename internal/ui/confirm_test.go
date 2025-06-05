package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewConfirm(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		message  string
		expected ConfirmModel
	}{
		{
			name:    "creates confirm model with title and message",
			title:   "Test Title",
			message: "Test Message",
			expected: ConfirmModel{
				title:    "Test Title",
				message:  "Test Message",
				selected: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewConfirm(tt.title, tt.message)
			if model.title != tt.expected.title {
				t.Errorf("Expected title %s, got %s", tt.expected.title, model.title)
			}
			if model.message != tt.expected.message {
				t.Errorf("Expected message %s, got %s", tt.expected.message, model.message)
			}
			if model.selected != tt.expected.selected {
				t.Errorf("Expected selected %v, got %v", tt.expected.selected, model.selected)
			}
		})
	}
}

func TestConfirmUpdate(t *testing.T) {
	tests := []struct {
		name           string
		initial        ConfirmModel
		msg            tea.Msg
		expectedSelect bool
		shouldQuit     bool
	}{
		{
			name:           "up key selects yes",
			initial:        NewConfirm("Test", "Message"),
			msg:            tea.KeyMsg{Type: tea.KeyUp},
			expectedSelect: true,
			shouldQuit:     false,
		},
		{
			name:           "down key selects no",
			initial:        ConfirmModel{selected: true},
			msg:            tea.KeyMsg{Type: tea.KeyDown},
			expectedSelect: false,
			shouldQuit:     false,
		},
		{
			name:           "left key selects yes",
			initial:        NewConfirm("Test", "Message"),
			msg:            tea.KeyMsg{Type: tea.KeyLeft},
			expectedSelect: true,
			shouldQuit:     false,
		},
		{
			name:           "right key selects no",
			initial:        ConfirmModel{selected: true},
			msg:            tea.KeyMsg{Type: tea.KeyRight},
			expectedSelect: false,
			shouldQuit:     false,
		},
		{
			name:           "y key selects yes and quits",
			initial:        NewConfirm("Test", "Message"),
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			expectedSelect: true,
			shouldQuit:     true,
		},
		{
			name:           "n key selects no and quits",
			initial:        ConfirmModel{selected: true},
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			expectedSelect: false,
			shouldQuit:     true,
		},
		{
			name:           "enter key confirms current selection",
			initial:        ConfirmModel{selected: true},
			msg:            tea.KeyMsg{Type: tea.KeyEnter},
			expectedSelect: true,
			shouldQuit:     true,
		},
		{
			name:           "q key cancels",
			initial:        NewConfirm("Test", "Message"),
			msg:            tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectedSelect: false,
			shouldQuit:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, cmd := tt.initial.Update(tt.msg)
			model := updatedModel.(ConfirmModel)

			if model.selected != tt.expectedSelect {
				t.Errorf("Expected selected %v, got %v", tt.expectedSelect, model.selected)
			}

			if tt.shouldQuit && cmd == nil {
				t.Error("Expected quit command, got nil")
			}
		})
	}
}

func TestConfirmGetResult(t *testing.T) {
	tests := []struct {
		name     string
		model    ConfirmModel
		expected ConfirmResult
	}{
		{
			name: "confirmed yes",
			model: ConfirmModel{
				confirmed: true,
				selected:  true,
			},
			expected: ConfirmResult{
				Confirmed: true,
				Cancelled: false,
			},
		},
		{
			name: "confirmed no",
			model: ConfirmModel{
				confirmed: true,
				selected:  false,
			},
			expected: ConfirmResult{
				Confirmed: false,
				Cancelled: false,
			},
		},
		{
			name: "cancelled",
			model: ConfirmModel{
				cancelled: true,
			},
			expected: ConfirmResult{
				Confirmed: false,
				Cancelled: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.GetResult()
			if result.Confirmed != tt.expected.Confirmed {
				t.Errorf("Expected Confirmed %v, got %v", tt.expected.Confirmed, result.Confirmed)
			}
			if result.Cancelled != tt.expected.Cancelled {
				t.Errorf("Expected Cancelled %v, got %v", tt.expected.Cancelled, result.Cancelled)
			}
		})
	}
}

func TestConfirmView(t *testing.T) {
	t.Run("shows empty string when cancelled", func(t *testing.T) {
		model := ConfirmModel{cancelled: true}
		view := model.View()
		if view != "" {
			t.Errorf("Expected empty view when cancelled, got %s", view)
		}
	})

	t.Run("shows empty string when confirmed", func(t *testing.T) {
		model := ConfirmModel{confirmed: true}
		view := model.View()
		if view != "" {
			t.Errorf("Expected empty view when confirmed, got %s", view)
		}
	})

	t.Run("shows content when active", func(t *testing.T) {
		model := NewConfirm("Test Title", "Test Message")
		view := model.View()
		if view == "" {
			t.Error("Expected non-empty view when active")
		}
	})
}
