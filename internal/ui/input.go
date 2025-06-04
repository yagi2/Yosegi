package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	title       string
	inputs      []textinput.Model
	focused     int
	submitted   bool
	cancelled   bool
	values      []string
}

type InputResult struct {
	Values    []string
	Submitted bool
}

func NewInput(title string, prompts []string, defaults []string) InputModel {
	inputs := make([]textinput.Model, len(prompts))
	
	for i, prompt := range prompts {
		input := textinput.New()
		input.Placeholder = prompt
		input.CharLimit = 200
		
		if i < len(defaults) && defaults[i] != "" {
			input.SetValue(defaults[i])
		}
		
		if i == 0 {
			input.Focus()
		}
		
		inputs[i] = input
	}

	return InputModel{
		title:  title,
		inputs: inputs,
		values: make([]string, len(prompts)),
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit

		case tea.KeyEnter:
			// If on last input or all inputs filled, submit
			if m.focused == len(m.inputs)-1 || m.allInputsFilled() {
				for i, input := range m.inputs {
					m.values[i] = input.Value()
				}
				m.submitted = true
				return m, tea.Quit
			}
			// Otherwise, move to next input
			m.nextInput()

		case tea.KeyTab, tea.KeyShiftTab:
			if msg.Type == tea.KeyTab {
				m.nextInput()
			} else {
				m.prevInput()
			}
		}
	}

	// Update focused input
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m InputModel) View() string {
	if m.submitted || m.cancelled {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render(fmt.Sprintf("📝 %s", m.title)))
	b.WriteString("\n\n")

	// Input fields
	for i, input := range m.inputs {
		b.WriteString(input.View())
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	// Help text
	b.WriteString("\n\n")
	helpText := "tab/shift+tab navigate • enter submit • esc cancel"
	b.WriteString(HelpStyle.Render(helpText))

	return BorderStyle.Render(b.String())
}

func (m InputModel) GetResult() InputResult {
	return InputResult{
		Values:    m.values,
		Submitted: m.submitted,
	}
}

func (m *InputModel) nextInput() {
	m.inputs[m.focused].Blur()
	m.focused = (m.focused + 1) % len(m.inputs)
	m.inputs[m.focused].Focus()
}

func (m *InputModel) prevInput() {
	m.inputs[m.focused].Blur()
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
	m.inputs[m.focused].Focus()
}

func (m InputModel) allInputsFilled() bool {
	for _, input := range m.inputs {
		if strings.TrimSpace(input.Value()) == "" {
			return false
		}
	}
	return true
}