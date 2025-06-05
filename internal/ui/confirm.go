package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type confirmKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Yes   key.Binding
	No    key.Binding
	Enter key.Binding
	Quit  key.Binding
}

var confirmKeys = confirmKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "select yes"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "select no"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "select yes"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "select no"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yes"),
	),
	No: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "no"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/esc", "cancel"),
	),
}

type ConfirmModel struct {
	title     string
	message   string
	selected  bool // true = yes, false = no
	confirmed bool
	cancelled bool
}

type ConfirmResult struct {
	Confirmed bool // true if user selected yes, false if no or cancelled
	Cancelled bool // true if user cancelled (q/esc)
}

func NewConfirm(title, message string) ConfirmModel {
	return ConfirmModel{
		title:    title,
		message:  message,
		selected: false, // default to "no" for safety
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, confirmKeys.Quit):
			m.cancelled = true
			return m, tea.Quit

		case key.Matches(msg, confirmKeys.Up), key.Matches(msg, confirmKeys.Left):
			m.selected = true

		case key.Matches(msg, confirmKeys.Down), key.Matches(msg, confirmKeys.Right):
			m.selected = false

		case key.Matches(msg, confirmKeys.Yes):
			m.selected = true
			m.confirmed = true
			return m, tea.Quit

		case key.Matches(msg, confirmKeys.No):
			m.selected = false
			m.confirmed = true
			return m, tea.Quit

		case key.Matches(msg, confirmKeys.Enter):
			m.confirmed = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	if m.cancelled || m.confirmed {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render(fmt.Sprintf("⚠️  %s", m.title)))
	b.WriteString("\n\n")

	// Message
	b.WriteString(NormalStyle.Render(m.message))
	b.WriteString("\n\n")

	// Options
	yesStyle := NormalStyle
	noStyle := NormalStyle

	if m.selected {
		yesStyle = SelectedItemStyle
	} else {
		noStyle = SelectedItemStyle
	}

	b.WriteString("  ")
	b.WriteString(yesStyle.Render("[ Yes ]"))
	b.WriteString("    ")
	b.WriteString(noStyle.Render("[ No ]"))
	b.WriteString("\n\n")

	// Help
	helpText := []string{
		"←/→/h/l switch", "↑/↓/k/j switch", "y yes", "n no", "enter confirm", "q/esc cancel",
	}
	b.WriteString(HelpStyle.Render(strings.Join(helpText, " • ")))

	return BorderStyle.Render(b.String())
}

func (m ConfirmModel) GetResult() ConfirmResult {
	if m.cancelled {
		return ConfirmResult{
			Confirmed: false,
			Cancelled: true,
		}
	}

	return ConfirmResult{
		Confirmed: m.confirmed && m.selected,
		Cancelled: false,
	}
}
