package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewInput(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		prompts  []string
		defaults []string
	}{
		{
			name:     "Single input",
			title:    "Enter Value",
			prompts:  []string{"Name"},
			defaults: []string{"default-name"},
		},
		{
			name:     "Multiple inputs",
			title:    "Enter Values",
			prompts:  []string{"Name", "Email", "Age"},
			defaults: []string{"John", "john@example.com", "30"},
		},
		{
			name:     "No defaults",
			title:    "Enter Info",
			prompts:  []string{"Field1", "Field2"},
			defaults: []string{},
		},
		{
			name:     "Partial defaults",
			title:    "Enter Info",
			prompts:  []string{"Field1", "Field2", "Field3"},
			defaults: []string{"default1", "default2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewInput(tt.title, tt.prompts, tt.defaults)

			if model.title != tt.title {
				t.Errorf("Expected title '%s', got '%s'", tt.title, model.title)
			}

			if len(model.inputs) != len(tt.prompts) {
				t.Errorf("Expected %d inputs, got %d", len(tt.prompts), len(model.inputs))
			}

			if len(model.values) != len(tt.prompts) {
				t.Errorf("Expected %d values slots, got %d", len(tt.prompts), len(model.values))
			}

			if model.focused != 0 {
				t.Errorf("Expected focused to be 0, got %d", model.focused)
			}

			if model.submitted {
				t.Errorf("Expected submitted to be false")
			}

			if model.cancelled {
				t.Errorf("Expected cancelled to be false")
			}

			// Check that first input is focused
			if !model.inputs[0].Focused() {
				t.Errorf("Expected first input to be focused")
			}

			// Check defaults
			for i, defaultVal := range tt.defaults {
				if i < len(model.inputs) && defaultVal != "" {
					if model.inputs[i].Value() != defaultVal {
						t.Errorf("Expected input %d to have default value '%s', got '%s'", i, defaultVal, model.inputs[i].Value())
					}
				}
			}

			// Check placeholders
			for i, prompt := range tt.prompts {
				if model.inputs[i].Placeholder != prompt {
					t.Errorf("Expected input %d to have placeholder '%s', got '%s'", i, prompt, model.inputs[i].Placeholder)
				}
			}
		})
	}
}

func TestInputInit(t *testing.T) {
	model := NewInput("Test", []string{"Field1"}, []string{})
	cmd := model.Init()

	// Should return textinput.Blink command
	if cmd == nil {
		t.Errorf("Expected Init() to return textinput.Blink command, got nil")
	}
}

func TestInputUpdate(t *testing.T) {
	tests := []struct {
		name           string
		setupModel     func() InputModel
		msg            tea.Msg
		expectedFocus  int
		expectedSubmit bool
		expectedCancel bool
	}{
		{
			name: "Tab to next input",
			setupModel: func() InputModel {
				return NewInput("Test", []string{"Field1", "Field2"}, []string{})
			},
			msg:           tea.KeyMsg{Type: tea.KeyTab},
			expectedFocus: 1,
		},
		{
			name: "Tab from last input (wraps to first)",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.focused = 1
				m.inputs[0].Blur()
				m.inputs[1].Focus()
				return m
			},
			msg:           tea.KeyMsg{Type: tea.KeyTab},
			expectedFocus: 0,
		},
		{
			name: "Shift+Tab to previous input",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.focused = 1
				m.inputs[0].Blur()
				m.inputs[1].Focus()
				return m
			},
			msg:           tea.KeyMsg{Type: tea.KeyShiftTab},
			expectedFocus: 0,
		},
		{
			name: "Shift+Tab from first input (wraps to last)",
			setupModel: func() InputModel {
				return NewInput("Test", []string{"Field1", "Field2"}, []string{})
			},
			msg:           tea.KeyMsg{Type: tea.KeyShiftTab},
			expectedFocus: 1,
		},
		{
			name: "Enter submits when on last input",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.focused = 1
				m.inputs[0].Blur()
				m.inputs[1].Focus()
				m.inputs[0].SetValue("value1")
				m.inputs[1].SetValue("value2")
				return m
			},
			msg:            tea.KeyMsg{Type: tea.KeyEnter},
			expectedFocus:  1, // Focus should remain on last input before submit
			expectedSubmit: true,
		},
		{
			name: "Enter moves to next input when not on last",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.inputs[0].SetValue("value1")
				return m
			},
			msg:           tea.KeyMsg{Type: tea.KeyEnter},
			expectedFocus: 1,
		},
		{
			name: "Enter submits when all inputs filled",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.inputs[0].SetValue("value1")
				m.inputs[1].SetValue("value2")
				return m
			},
			msg:            tea.KeyMsg{Type: tea.KeyEnter},
			expectedSubmit: true,
		},
		{
			name: "Ctrl+C cancels",
			setupModel: func() InputModel {
				return NewInput("Test", []string{"Field1"}, []string{})
			},
			msg:            tea.KeyMsg{Type: tea.KeyCtrlC},
			expectedCancel: true,
		},
		{
			name: "Esc cancels",
			setupModel: func() InputModel {
				return NewInput("Test", []string{"Field1"}, []string{})
			},
			msg:            tea.KeyMsg{Type: tea.KeyEsc},
			expectedCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			newModel, cmd := model.Update(tt.msg)
			updatedModel := newModel.(InputModel)

			if tt.expectedFocus >= 0 && updatedModel.focused != tt.expectedFocus {
				t.Errorf("Expected focused to be %d, got %d", tt.expectedFocus, updatedModel.focused)
			}

			if updatedModel.submitted != tt.expectedSubmit {
				t.Errorf("Expected submitted to be %v, got %v", tt.expectedSubmit, updatedModel.submitted)
			}

			if updatedModel.cancelled != tt.expectedCancel {
				t.Errorf("Expected cancelled to be %v, got %v", tt.expectedCancel, updatedModel.cancelled)
			}

			// Check that appropriate input is focused
			if tt.expectedFocus >= 0 && len(updatedModel.inputs) > tt.expectedFocus {
				if !updatedModel.inputs[tt.expectedFocus].Focused() {
					t.Errorf("Expected input %d to be focused", tt.expectedFocus)
				}
			}

			// Check quit command when submitted or cancelled
			if (tt.expectedSubmit || tt.expectedCancel) && cmd == nil {
				t.Errorf("Expected quit command when submitted/cancelled")
			}
		})
	}
}

func TestInputView(t *testing.T) {
	tests := []struct {
		name                string
		setupModel          func() InputModel
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name: "Normal input view",
			setupModel: func() InputModel {
				return NewInput("Enter Information", []string{"Name", "Email"}, []string{"John", ""})
			},
			expectedContains: []string{
				"📝 Enter Information",
				"tab/shift+tab navigate",
				"enter submit",
				"esc cancel",
			},
		},
		{
			name: "Submitted model (empty view)",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1"}, []string{})
				m.submitted = true
				return m
			},
			expectedContains: []string{},
		},
		{
			name: "Cancelled model (empty view)",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1"}, []string{})
				m.cancelled = true
				return m
			},
			expectedContains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			view := model.View()

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

func TestInputGetResult(t *testing.T) {
	tests := []struct {
		name              string
		setupModel        func() InputModel
		expectedSubmitted bool
		expectedValues    []string
	}{
		{
			name: "Submitted with values",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.inputs[0].SetValue("value1")
				m.inputs[1].SetValue("value2")
				m.values = []string{"value1", "value2"}
				m.submitted = true
				return m
			},
			expectedSubmitted: true,
			expectedValues:    []string{"value1", "value2"},
		},
		{
			name: "Cancelled",
			setupModel: func() InputModel {
				m := NewInput("Test", []string{"Field1"}, []string{})
				m.cancelled = true
				return m
			},
			expectedSubmitted: false,
			expectedValues:    []string{""},
		},
		{
			name: "Not submitted or cancelled",
			setupModel: func() InputModel {
				return NewInput("Test", []string{"Field1"}, []string{})
			},
			expectedSubmitted: false,
			expectedValues:    []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.setupModel()
			result := model.GetResult()

			if result.Submitted != tt.expectedSubmitted {
				t.Errorf("Expected submitted to be %v, got %v", tt.expectedSubmitted, result.Submitted)
			}

			if len(result.Values) != len(tt.expectedValues) {
				t.Errorf("Expected %d values, got %d", len(tt.expectedValues), len(result.Values))
				return
			}

			for i, expected := range tt.expectedValues {
				if i < len(result.Values) && result.Values[i] != expected {
					t.Errorf("Expected value %d to be '%s', got '%s'", i, expected, result.Values[i])
				}
			}
		})
	}
}

func TestInputNextPrevInput(t *testing.T) {
	model := NewInput("Test", []string{"Field1", "Field2", "Field3"}, []string{})

	// Test nextInput
	originalFocused := model.focused
	model.nextInput()
	if model.focused != (originalFocused+1)%len(model.inputs) {
		t.Errorf("nextInput didn't move to next input correctly")
	}

	// Test wrapping
	model.focused = len(model.inputs) - 1 // Last input
	model.nextInput()
	if model.focused != 0 {
		t.Errorf("nextInput didn't wrap to first input")
	}

	// Test prevInput
	model.focused = 1
	model.prevInput()
	if model.focused != 0 {
		t.Errorf("prevInput didn't move to previous input correctly")
	}

	// Test wrapping backwards
	model.focused = 0
	model.prevInput()
	if model.focused != len(model.inputs)-1 {
		t.Errorf("prevInput didn't wrap to last input")
	}
}

func TestAllInputsFilled(t *testing.T) {
	tests := []struct {
		name     string
		model    InputModel
		expected bool
	}{
		{
			name: "All inputs filled",
			model: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.inputs[0].SetValue("value1")
				m.inputs[1].SetValue("value2")
				return m
			}(),
			expected: true,
		},
		{
			name: "Some inputs empty",
			model: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.inputs[0].SetValue("value1")
				m.inputs[1].SetValue("")
				return m
			}(),
			expected: false,
		},
		{
			name: "All inputs empty",
			model: func() InputModel {
				return NewInput("Test", []string{"Field1", "Field2"}, []string{})
			}(),
			expected: false,
		},
		{
			name: "Inputs with only whitespace",
			model: func() InputModel {
				m := NewInput("Test", []string{"Field1", "Field2"}, []string{})
				m.inputs[0].SetValue("  ")
				m.inputs[1].SetValue("\t")
				return m
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.allInputsFilled()
			if result != tt.expected {
				t.Errorf("Expected allInputsFilled() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInputWithCharacterInput(t *testing.T) {
	model := NewInput("Test", []string{"Name"}, []string{})

	// Simulate typing characters
	chars := []rune{'h', 'e', 'l', 'l', 'o'}
	for _, char := range chars {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{char},
		}
		newModel, _ := model.Update(msg)
		model = newModel.(InputModel)
	}

	// Check that the input contains the typed text
	if model.inputs[0].Value() != "hello" {
		t.Errorf("Expected input value to be 'hello', got '%s'", model.inputs[0].Value())
	}
}

func TestInputCharLimit(t *testing.T) {
	model := NewInput("Test", []string{"Name"}, []string{})

	// Check that inputs have char limit set
	for i, input := range model.inputs {
		if input.CharLimit != 200 {
			t.Errorf("Expected input %d to have char limit 200, got %d", i, input.CharLimit)
		}
	}
}

func TestInputEdgeCases(t *testing.T) {
	t.Run("Single input immediate submit", func(t *testing.T) {
		model := NewInput("Test", []string{"Field1"}, []string{})
		model.inputs[0].SetValue("value")

		// Press Enter on single input with value
		newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updatedModel := newModel.(InputModel)

		if !updatedModel.submitted {
			t.Errorf("Single filled input should submit on Enter")
		}

		if cmd == nil {
			t.Errorf("Should return quit command on submit")
		}
	})

	t.Run("Empty inputs submit attempt", func(t *testing.T) {
		model := NewInput("Test", []string{"Field1", "Field2"}, []string{})

		// Try to submit with empty inputs
		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updatedModel := newModel.(InputModel)

		if updatedModel.submitted {
			t.Errorf("Should not submit with empty inputs")
		}

		// Should move to next input instead
		if updatedModel.focused != 1 {
			t.Errorf("Should move to next input when current is empty")
		}
	})
}

// Test input model with textinput integration
func TestInputTextInputIntegration(t *testing.T) {
	model := NewInput("Test", []string{"Name", "Email"}, []string{"John", "john@example.com"})

	// Verify textinput models are properly configured
	for i, input := range model.inputs {
		if input.Placeholder == "" {
			t.Errorf("Input %d should have placeholder set", i)
		}

		if input.CharLimit != 200 {
			t.Errorf("Input %d should have char limit 200", i)
		}

		// First input should be focused
		if i == 0 && !input.Focused() {
			t.Errorf("First input should be focused")
		} else if i != 0 && input.Focused() {
			t.Errorf("Only first input should be focused initially")
		}
	}
}

// Benchmark tests
func BenchmarkInputView(b *testing.B) {
	model := NewInput("Benchmark Test", []string{"Field1", "Field2", "Field3"}, []string{"default1", "default2", "default3"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkInputUpdate(b *testing.B) {
	model := NewInput("Benchmark Test", []string{"Field1", "Field2"}, []string{})
	msg := tea.KeyMsg{Type: tea.KeyTab}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newModel, _ := model.Update(msg)
		model = newModel.(InputModel)
	}
}

func BenchmarkNewInput(b *testing.B) {
	prompts := []string{"Name", "Email", "Phone", "Address", "City", "State", "Zip"}
	defaults := []string{"John", "john@example.com", "555-1234", "123 Main St", "Anytown", "ST", "12345"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewInput("Benchmark", prompts, defaults)
	}
}
